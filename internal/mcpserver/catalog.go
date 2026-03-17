package mcpserver

import (
	"fmt"
	"strings"
)

const (
	serverOverviewURI    = "scholar://server/overview"
	serverToolsURI       = "scholar://server/tools"
	serverConfigURI      = "scholar://server/config"
	serverLimitationsURI = "scholar://server/limitations"
	searchGuidePrefix    = "scholar://search-guide/"
	searchGuideTemplate  = "scholar://search-guide/{topic}"
	searchPromptName     = "scholar_literature_scan"
	authorPromptName     = "scholar_author_brief"
)

type toolDescriptor struct {
	Name        string
	Summary     string
	Args        []string
	ReturnShape string
}

func (s *Server) toolDescriptors() []toolDescriptor {
	tools := []toolDescriptor{
		{
			Name:        "search_google_scholar_key_words",
			Summary:     "Search Google Scholar by keyword query and return structured paper metadata.",
			Args:        []string{"query (required)", "num_results (optional)"},
			ReturnShape: "SearchResponse with query, results_count, and results[].",
		},
		{
			Name:        "search_google_scholar_advanced",
			Summary:     "Search Google Scholar with author and year filters.",
			Args:        []string{"query (required)", "author", "start_year", "end_year", "year_range (legacy)", "num_results"},
			ReturnShape: "SearchResponse with query, results_count, and results[].",
		},
	}

	if s.cfg.EnableAuthor {
		tools = append(tools, toolDescriptor{
			Name:        "get_author_info",
			Summary:     "Find a Google Scholar author profile by name and return structured author metadata.",
			Args:        []string{"author_name (required)"},
			ReturnShape: "AuthorProfile with affiliation, interests, cited_by, scholar_id, and publications[].",
		})
	}

	return tools
}

func (s *Server) overviewResourceText() string {
	return fmt.Sprintf(`# Google Scholar MCP

Google Scholar MCP is a Go-based MCP server that exposes Google Scholar search and author lookup as local stdio tools.

- Version: %s
- Transport: %s
- Base URL: %s
- Author tool enabled: %t
- Tool count: %d

This server exposes tools, resources, and prompts so clients can both call Scholar operations directly and inspect the server's operating envelope before they do so.
`, s.version, s.cfg.Transport, s.cfg.BaseURL, s.cfg.EnableAuthor, len(s.toolDescriptors()))
}

func (s *Server) toolsResourceText() string {
	var b strings.Builder
	b.WriteString("# Tool Catalog\n\n")
	b.WriteString("The server currently exposes these tools:\n")
	for _, tool := range s.toolDescriptors() {
		b.WriteString("\n## ")
		b.WriteString(tool.Name)
		b.WriteString("\n")
		b.WriteString(tool.Summary)
		b.WriteString("\n\nArguments:\n")
		for _, arg := range tool.Args {
			b.WriteString("- ")
			b.WriteString(arg)
			b.WriteString("\n")
		}
		b.WriteString("\nReturns:\n- ")
		b.WriteString(tool.ReturnShape)
		b.WriteString("\n")
	}
	return b.String()
}

func (s *Server) configResourceText() string {
	return fmt.Sprintf(`# Runtime Configuration

- MCP_TRANSPORT: %s
- SCHOLAR_BASE_URL: %s
- SCHOLAR_TIMEOUT: %s
- SCHOLAR_MAX_RESULTS: %d
- SCHOLAR_RATE_LIMIT_RPS: %.2f
- SCHOLAR_ACCEPT_LANGUAGE: %s
- SCHOLAR_ENABLE_AUTHOR_TOOL: %t
- LOG_LEVEL: %s

The server keeps logs on stderr so stdio protocol traffic stays clean.
`, s.cfg.Transport, s.cfg.BaseURL, s.cfg.Timeout, s.cfg.MaxResults, s.cfg.RateLimitRPS, s.cfg.AcceptLanguage, s.cfg.EnableAuthor, s.cfg.LogLevel.String())
}

func limitationsResourceText() string {
	return `# Limitations and Usage Notes

- Google Scholar does not provide a stable public API for this workflow.
- Search result snippets are not guaranteed to be full paper abstracts.
- HTML scraping can break when Google Scholar changes markup or introduces anti-bot measures.
- The server is intended for low-volume local usage, not bulk harvesting.
- Full-text download, PDF crawling, and CAPTCHA solving are out of scope.
`
}

func searchGuideResourceText(topic string) string {
	return fmt.Sprintf(`# Search Guide: %s

Recommended workflow for this topic:

1. Start with search_google_scholar_key_words using a concise topic query.
2. Inspect year, citation_count, version_count, source, and snippet quality before drawing conclusions.
3. If the topic is too broad, refine with a narrower keyword phrase or move to search_google_scholar_advanced for author and year filters.
4. Treat Scholar snippets as result-page excerpts, not authoritative abstracts.
5. Call out coverage gaps explicitly if the results appear sparse or noisy.
`, topic)
}
