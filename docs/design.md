# Design Spec: Google Scholar MCP Server Go Rewrite

## Goal
Build a Go-based MCP server that replaces the previous implementation, preserves the existing Google Scholar search capabilities, and provides a production-oriented foundation for maintainability, testing, and client compatibility.

## Background
The previous implementation exposed three tools:

1. `search_google_scholar_key_words`
2. `search_google_scholar_advanced`
3. `get_author_info`

The previous implementation had these characteristics:

- MCP integration and business logic are mixed together in a single script.
- Search tools scrape Google Scholar HTML directly.
- Author lookup relied on an external library without a direct Go equivalent in this repository.
- There is no test suite, no package structure, and limited error handling.

The Go rewrite must preserve the externally useful behavior while improving structure, reliability, observability, and testability.

## Success Criteria
The rewrite is successful if all of the following are true:

1. A local MCP client can start the Go server over `stdio`.
2. The Go server exposes tool definitions that are compatible with current usage patterns.
3. Keyword search and advanced search return structured, parseable results.
4. The codebase has clear separation between MCP transport, domain logic, and HTTP scraping logic.
5. Critical parsing logic is covered by automated tests using fixture HTML.
6. The implementation does not write logs to `stdout` while running over `stdio`.

## Requirements

### Functional Requirements

1. The server must expose the following tools in V1:
   - `search_google_scholar_key_words`
   - `search_google_scholar_advanced`
   - `get_author_info`
2. The server must support local MCP clients over `stdio`.
3. The keyword search tool must accept:
   - `query`
   - `num_results`
4. The advanced search tool must accept:
   - `query`
   - `author`
   - `start_year`
   - `end_year`
   - `num_results`
5. The author tool must accept:
   - `author_name`
6. Search tools must return structured paper metadata, not free-form text only.
7. The author tool must return structured author metadata, including basic publication summaries when available.
8. The server must validate tool inputs before making upstream requests.
9. The server must return actionable error messages for invalid input, upstream blocking, timeout, parse failure, and empty results.

### Engineering Requirements

1. The implementation must use Go.
2. The implementation must use the official MCP Go SDK.
3. The codebase must separate:
   - MCP server wiring
   - configuration
   - Google Scholar provider logic
   - HTML parsing
   - shared domain models
4. The server must support timeouts, retry policy, and rate limiting for upstream requests.
5. The implementation must be testable without live Google Scholar traffic for the main parsing paths.
6. Logging must use `stderr` or file output only.
7. The project must build as a single executable binary.

### Compatibility Requirements

1. Existing tool names must remain available in V1 for compatibility.
2. The default transport in V1 must be `stdio`.
3. Existing client workflows that call the legacy tool names must continue to work after switching to the Go binary.

## Non-Requirements (Out of Scope)

- V1 will not implement MCP `Resources`.
- V1 will not implement MCP `Prompts`.
- V1 will not provide a remote HTTP deployment mode.
- V1 will not download paper full text or PDFs automatically.
- V1 will not solve anti-bot challenges such as CAPTCHA.
- V1 will not guarantee high-volume crawling or bulk export.
- V1 will not guarantee stable access if Google Scholar changes its HTML structure.
- V1 will not implement user authentication.
- V1 will not persist search history or user state.

## External Constraints

1. Google Scholar does not provide a stable public API for this use case.
2. HTML scraping is inherently brittle and may break when Google changes markup.
3. The previous author lookup path relied on a dependency that has no direct Go equivalent in this repository.
4. MCP `stdio` transport requires strict separation between protocol traffic and logs.

## Data Access Strategy

### V1 Strategy

V1 will use direct HTML retrieval from Google Scholar as the primary data source.

Rationale:

- It avoids introducing a paid runtime dependency for basic local MCP usage.
- It preserves the current repository's operating model.
- It keeps the main control points inside the Go codebase, including parsing, error handling, and throttling.

### Third-Party Providers

Paid API or proxy-backed providers are explicitly optional in this design.

Examples:

- structured SERP APIs
- proxy and anti-bot scraping services

These providers may be added later behind the existing provider interfaces, but V1 must not depend on them for core functionality.

### Operating Envelope

The self-scraping implementation is intended for:

- low-volume local usage
- interactive tool calls
- best-effort retrieval

It is not intended for:

- bulk crawling
- high-concurrency harvesting
- guaranteed stability across Google Scholar markup changes

## Data Extraction Boundaries

The implementation is expected to reliably target these result-page fields when present:

- paper title
- paper link
- author line / publication line
- Scholar snippet text
- publication year when detectable
- citation count when detectable
- version count when detectable
- PDF or resource link when present

Important boundary:

- The Scholar snippet must be treated as a result-page excerpt, not a guaranteed full paper abstract.
- Full abstracts usually require visiting the publisher or repository landing page and are out of scope for V1.

## Architecture

### High-Level Architecture

The Go implementation will use a layered architecture:

1. `cmd/google-scholar-mcp/main.go`
   - Process entrypoint
   - Config loading
   - Logger initialization
   - Dependency wiring
   - MCP server startup

2. `internal/mcpserver`
   - Tool registration
   - Tool input schema definitions
   - Translation between MCP requests and domain services
   - Structured tool responses

3. `internal/scholar`
   - Provider interfaces
   - Google Scholar request construction
   - HTTP client logic
   - Retry, timeout, and rate limiting
   - HTML parsing integration
   - Domain error classification

4. `internal/model`
   - Shared request and response models

5. `testdata`
   - Static HTML fixtures used by parser tests

### Proposed Directory Layout

```text
cmd/google-scholar-mcp/
  main.go

internal/config/
  config.go

internal/mcpserver/
  server.go
  tools_search.go
  tools_author.go
  result.go

internal/scholar/
  provider.go
  request.go
  search.go
  author.go
  client.go
  parser_search.go
  parser_author.go
  errors.go
  rate_limit.go
  useragent.go

internal/model/
  paper.go
  author.go

testdata/
  scholar_search_page.html
  scholar_author_search_page.html
  scholar_author_profile_page.html
```

### Component Responsibilities

`main.go`

- Parse flags and environment variables.
- Initialize a logger that writes to `stderr`.
- Construct the HTTP client and Google Scholar provider.
- Create the MCP server and run it over `stdio`.

`internal/mcpserver`

- Register tools with the official Go SDK.
- Define request schemas and field documentation.
- Perform lightweight validation and coercion.
- Convert service results into MCP tool responses.

`internal/scholar`

- Encapsulate all upstream interactions.
- Build Google Scholar URLs using `url.Values`.
- Execute HTTP requests with common headers and limits.
- Parse HTML into domain models.
- Map failures into typed domain errors.

`internal/model`

- Provide stable internal response shapes independent of HTML structure.

## Transport Design

### V1 Transport

V1 will support only `stdio`.

Rationale:

- It matches the current local MCP client use case.
- It minimizes deployment complexity.
- It aligns with the simplest official build-server flow.
- It keeps the rewrite focused on behavior and reliability instead of remote serving concerns.

### Future Transport

Streamable HTTP may be added in V2, but it is explicitly out of scope for this spec.

## Tool Specification

### Tool 1: `search_google_scholar_key_words`

#### Purpose
Search Google Scholar by general query terms.

#### Input

```json
{
  "query": "graph neural networks",
  "num_results": 5
}
```

#### Input Rules

- `query` is required and must be non-empty after trimming whitespace.
- `num_results` is optional.
- Default `num_results` is `5`.
- Maximum `num_results` in V1 is `10`.

#### Output

```json
{
  "query": "graph neural networks",
  "results_count": 2,
  "results": [
    {
      "title": "A Survey on Graph Neural Networks",
      "url": "https://example.org/paper",
      "authors_text": "Z. Wu, S. Pan, ...",
      "snippet": "Graph neural networks have emerged...",
      "publication": "IEEE Transactions ...",
      "year": 2021,
      "citation_count": 1234,
      "pdf_url": "https://example.org/paper.pdf",
      "source": "google_scholar"
    }
  ]
}
```

#### Behavior

- Build a Scholar search URL from the query.
- Fetch the first result page only in V1.
- Parse up to `num_results` entries.
- Return an empty result set if the page loads successfully but no matches are found.
- Return Scholar snippet text when available, but do not claim it is the full abstract.

### Tool 2: `search_google_scholar_advanced`

#### Purpose
Search Google Scholar with author and year filters.

#### Input

```json
{
  "query": "transformer interpretability",
  "author": "Yoshua Bengio",
  "start_year": 2020,
  "end_year": 2024,
  "num_results": 5
}
```

#### Input Rules

- `query` is required.
- `author` is optional.
- `start_year` and `end_year` are optional.
- If one year bound is set, the other may be omitted.
- If both are set, `start_year` must be less than or equal to `end_year`.
- `num_results` follows the same rules as the keyword search tool.

#### Behavior

- Build a Scholar URL using query parameters:
  - `q`
  - `as_auth`
  - `as_ylo`
  - `as_yhi`
- Reuse the same search result parser as the keyword search tool.

### Tool 3: `get_author_info`

#### Purpose
Return structured information about a Google Scholar author profile.

#### Input

```json
{
  "author_name": "Geoffrey Hinton"
}
```

#### Input Rules

- `author_name` is required and must be non-empty after trimming.

#### Output

```json
{
  "name": "Geoffrey Hinton",
  "affiliation": "University of Toronto",
  "interests": ["machine learning", "neural networks"],
  "cited_by": 912345,
  "scholar_id": "abc123",
  "profile_url": "https://scholar.google.com/...",
  "publications": [
    {
      "title": "Reducing the dimensionality of data with neural networks",
      "year": 2006,
      "citations": 12345
    }
  ]
}
```

#### Author Tool Strategy

The author tool will be implemented through a provider abstraction because the previous author lookup dependency is not directly portable.

V1 decision:

1. The external tool name remains `get_author_info`.
2. The internal implementation uses an `AuthorProvider` interface.
3. The default V1 provider is a pure Go HTML-based author provider.
4. The provider is allowed to return partial metadata if the profile page is reachable but some optional fields cannot be parsed.
5. The tool must clearly surface upstream blocking and parse errors.
6. Name-based author resolution is best-effort in V1.
7. A future `author_id`-based lookup tool is preferred for stronger stability.

Rationale:

- This keeps the rewrite fully in Go.
- It avoids embedding another runtime in the Go execution path.
- It makes it possible to swap in another provider later without changing the MCP layer.

## Domain Model

### Paper

```go
type Paper struct {
    Title         string
    URL           string
    AuthorsText   string
    Snippet       string
    Publication   string
    Year          int
    CitationCount int
    PDFURL        string
    Source        string
}
```

### Publication

```go
type Publication struct {
    Title     string
    Year      int
    Citations int
}
```

### AuthorProfile

```go
type AuthorProfile struct {
    Name         string
    Affiliation  string
    Interests    []string
    CitedBy      int
    ScholarID    string
    ProfileURL   string
    Publications []Publication
}
```

## Provider Design

### Interfaces

```go
type SearchProvider interface {
    Search(ctx context.Context, req KeywordSearchRequest) ([]model.Paper, error)
    AdvancedSearch(ctx context.Context, req AdvancedSearchRequest) ([]model.Paper, error)
}

type AuthorProvider interface {
    GetAuthor(ctx context.Context, req AuthorRequest) (*model.AuthorProfile, error)
}
```

### Concrete Implementation

`ScholarHTMLProvider` will implement both interfaces in V1.

Responsibilities:

- Build URLs for Scholar search and author profile navigation.
- Execute upstream requests.
- Parse search pages and author pages.
- Return typed domain models.

Future provider adapters may use third-party APIs or proxy-backed retrieval without changing the MCP tool layer.

## HTTP Client Design

### Request Policy

- Use one shared `http.Client`.
- Apply a default timeout of 15 seconds.
- Use a custom `Transport` with connection reuse enabled.
- Set a browser-like `User-Agent`.
- Set `Accept-Language` to improve deterministic parsing behavior where possible.

### Retry Policy

- Retry at most 2 times for:
  - `429`
  - `500`
  - `502`
  - `503`
  - `504`
- Do not retry client-side validation errors.
- Use exponential backoff with jitter.

### Rate Limiting

- Apply a process-wide rate limiter.
- Default target rate: `0.5` requests per second.
- Allow a small burst of `1`.

### Concurrency

- V1 will keep low concurrency by default.
- Tool calls may execute concurrently, but upstream requests must still pass through the limiter.

## Parsing Design

### Search Page Parser

The search page parser must extract:

- title
- title link
- author/publication line
- snippet
- year when detectable
- citation count when detectable
- PDF link when present

Parser design rules:

1. Parsing must tolerate missing optional nodes.
2. Missing optional data must not fail the whole result item.
3. Parsing must fail only when the page itself is unusable or the expected result container cannot be interpreted at all.

### Author Search and Profile Parser

The author flow may require two steps:

1. Search for an author profile by name.
2. Resolve the best candidate profile page.
3. Parse the profile page fields.

The parser must extract when available:

- profile name
- affiliation
- interests
- cited-by total
- scholar profile id
- top publications

Selection rule for multiple author candidates:

1. Exact display-name match after normalization wins.
2. Otherwise select the first visible candidate and mark the result as best-effort.

## Future Tool Candidates

The following tools are good candidates after V1 because their underlying data is often visible from Scholar result or author pages:

1. `get_cited_by_papers`
   - Input: citation or result identifier
   - Output: papers that cite the target paper

2. `get_paper_versions`
   - Input: cluster or result identifier
   - Output: alternate versions of the same paper

3. `get_related_papers`
   - Input: result identifier
   - Output: related-paper results when the Scholar page exposes that navigation

4. `get_author_by_id`
   - Input: `author_id`
   - Output: author profile data with stronger stability than name search

5. `get_author_publications`
   - Input: `author_id`
   - Output: paginated publication list from an author profile

6. `get_author_metrics`
   - Input: `author_id`
   - Output: cited-by total, h-index, and i10-index when present

7. `get_author_coauthors`
   - Input: `author_id`
   - Output: visible coauthor list from the author profile

8. `fetch_paper_landing_page`
   - Input: paper URL
   - Output: landing-page metadata such as resolved title, canonical URL, and abstract when the destination site is parseable

9. `fetch_paper_abstract`
   - Input: paper URL or DOI
   - Output: best-effort full abstract from the publisher or repository page

The first seven tools build naturally on top of Scholar parsing. The last two cross the boundary from Scholar into publisher-specific parsing and should remain explicitly best-effort.

## MCP Response Design

Each tool response should contain:

- a structured payload
- a concise human-readable summary string

The structured payload is the source of truth.

The summary string exists to improve readability in clients that display plain text prominently.

Example summary:

- `"Found 5 Google Scholar results for 'graph neural networks'."`
- `"Found author profile for 'Geoffrey Hinton' with 5 publications."`

## Error Model

The domain layer must define typed errors:

```go
var (
    ErrInvalidInput        = errors.New("invalid input")
    ErrNoResults           = errors.New("no results")
    ErrTimeout             = errors.New("request timeout")
    ErrUpstreamBlocked     = errors.New("upstream blocked")
    ErrUpstreamUnavailable = errors.New("upstream unavailable")
    ErrParseFailed         = errors.New("parse failed")
)
```

Behavioral rules:

1. Invalid input must be caught before upstream I/O.
2. Empty result sets must not be treated as internal server failures.
3. Upstream challenge pages or obvious anti-bot pages should map to `ErrUpstreamBlocked`.
4. HTML structure changes that break extraction should map to `ErrParseFailed`.
5. MCP tool responses must convert internal errors into actionable user-facing messages.

## Logging and Observability

### Logging Rules

1. Never log to `stdout` while using `stdio`.
2. Default logger target is `stderr`.
3. Log one structured event per tool call with:
   - tool name
   - duration
   - result count
   - error class
4. Do not log full raw HTML responses.
5. Avoid logging raw research queries at debug level unless explicitly enabled.

### Metrics

Metrics are out of scope for V1, but the code should keep a clean seam for adding them later.

## Configuration

The binary must support configuration via environment variables and optional flags.

### Required Configuration Surface

- `MCP_TRANSPORT`
- `SCHOLAR_BASE_URL`
- `SCHOLAR_TIMEOUT`
- `SCHOLAR_MAX_RESULTS`
- `SCHOLAR_RATE_LIMIT_RPS`
- `SCHOLAR_USER_AGENT`
- `SCHOLAR_ENABLE_AUTHOR_TOOL`
- `LOG_LEVEL`

### Defaults

- `MCP_TRANSPORT=stdio`
- `SCHOLAR_BASE_URL=https://scholar.google.com`
- `SCHOLAR_TIMEOUT=15s`
- `SCHOLAR_MAX_RESULTS=10`
- `SCHOLAR_RATE_LIMIT_RPS=0.5`
- `SCHOLAR_ENABLE_AUTHOR_TOOL=true`
- `LOG_LEVEL=info`

## Security and Reliability Considerations

1. The server must bound `num_results` to a small maximum to reduce load and anti-bot risk.
2. The server must use timeouts on all external requests.
3. The server must not execute arbitrary shell commands.
4. The server must not trust upstream HTML structure without validation.
5. The server must fail closed on malformed inputs rather than attempt unsafe coercion.

## Implementation Plan

### Milestone 1: Foundation

Deliverables:

- Go module initialization
- project directory scaffold
- config loader
- logger setup
- minimal MCP server boot over `stdio`
- one non-network smoke tool for local verification

Exit criteria:

- Binary starts under an MCP client without protocol corruption.

### Milestone 2: Keyword Search

Deliverables:

- URL builder for keyword search
- shared HTTP client
- search result parser
- `search_google_scholar_key_words` tool
- parser fixture tests

Exit criteria:

- Tool returns structured search results using fixture-backed tests and local smoke verification.

### Milestone 3: Advanced Search

Deliverables:

- advanced query builder
- input validation for year bounds
- `search_google_scholar_advanced` tool

Exit criteria:

- Tool reuses the shared parser and returns filtered structured results.

### Milestone 4: Author Info

Deliverables:

- author search flow
- author profile parser
- `get_author_info` tool
- best-effort handling for partial profile pages

Exit criteria:

- Tool returns structured author metadata when the profile page is available.

### Milestone 5: Hardening

Deliverables:

- typed domain errors
- retry policy
- rate limiting
- better summaries
- end-to-end MCP tests
- documentation updates

Exit criteria:

- All V1 acceptance criteria are met.

## Test Plan

### Unit Tests

- query parameter construction
- request validation
- search page parser
- author page parser
- error classification logic

### Integration Tests

- `httptest.Server` for simulated Scholar responses
- provider behavior over mocked upstream HTML

### End-to-End Tests

- start the binary or server instance locally
- invoke tools via MCP client transport
- assert structured responses

### Live Tests

Live Google Scholar traffic tests are optional and must not run in default CI.

If implemented, they must:

- be guarded by an explicit environment variable
- use very low request counts
- tolerate flakiness due to upstream anti-bot behavior

## Migration Plan

1. Keep the current Go implementation isolated from compatibility cleanup until the binary reaches feature parity for the required tools.
2. Introduce the Go implementation in a separate directory or branch of the codebase.
3. Validate the Go binary with MCP Inspector and at least one real local client.
4. Switch client configuration from the old entrypoint to the Go binary only after acceptance criteria are met.
5. Remove deprecated implementation artifacts after a short overlap period.

## Acceptance Criteria

- [ ] The repository contains a Go binary entrypoint for the MCP server.
- [ ] The Go server starts successfully over `stdio`.
- [ ] The server exposes `search_google_scholar_key_words`.
- [ ] The server exposes `search_google_scholar_advanced`.
- [ ] The server exposes `get_author_info`.
- [ ] Tool names are compatible with the previous server names.
- [ ] Search tool inputs are validated before outbound requests.
- [ ] Search results are returned as structured JSON-compatible objects.
- [ ] Author info is returned as a structured JSON-compatible object.
- [ ] Upstream timeout, blocked, no-result, and parse-failure cases are distinguished.
- [ ] Logs do not write to `stdout`.
- [ ] Parser behavior is covered by automated tests using fixture HTML.
- [ ] The project builds into a single executable binary.
- [ ] A local MCP client can successfully call at least the keyword search tool end to end.

## Open Questions

1. Should V1 expose only the legacy tool names, or expose both legacy names and new canonical names?
2. Should `get_author_info` be enabled by default from day one, or guarded behind a feature flag until author parsing is stable?
3. Do we need a small in-memory TTL cache in V1, or can we defer caching to V2?
4. Should V2 add a provider toggle so deployments can opt into proxy-backed or paid retrieval without changing tool names?
5. Should future author-centric tools standardize on `author_id` as the primary identifier while keeping name search as a convenience path?

## Decision Summary

This spec makes the following concrete decisions:

1. The rewrite will be pure Go.
2. The implementation will use the official MCP Go SDK.
3. V1 will support only `stdio`.
4. V1 will preserve the current tool names.
5. V1 will use a provider abstraction for Google Scholar access.
6. The default provider in V1 will be pure Go HTML scraping.
7. V1 will not depend on a paid third-party Scholar API.
8. Scholar result-page snippet text is treated as an excerpt, not a guaranteed full abstract.
9. The implementation will prioritize testable parsing and low-risk local integration over feature expansion.
