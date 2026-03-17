package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
	"github.com/bingshuoguo/google-scholar-mcp/internal/model"
	"github.com/bingshuoguo/google-scholar-mcp/internal/scholar"
)

type Server struct {
	server   *mcp.Server
	logger   *slog.Logger
	provider scholar.Provider
	cfg      config.Config
	version  string
}

type keywordSearchInput struct {
	Query      string `json:"query"`
	NumResults int    `json:"num_results,omitempty"`
}

type advancedSearchInput struct {
	Query      string `json:"query"`
	Author     string `json:"author,omitempty"`
	StartYear  int    `json:"start_year,omitempty"`
	EndYear    int    `json:"end_year,omitempty"`
	YearRange  []int  `json:"year_range,omitempty"`
	NumResults int    `json:"num_results,omitempty"`
}

type authorInfoInput struct {
	AuthorName string `json:"author_name"`
}

func New(cfg config.Config, logger *slog.Logger, version string) *Server {
	impl := &mcp.Implementation{
		Name:       "google-scholar-mcp",
		Title:      "Google Scholar MCP",
		Version:    version,
		WebsiteURL: "https://github.com/bingshuoguo/google-scholar-mcp",
	}
	server := mcp.NewServer(impl, nil)
	client := scholar.NewClient(cfg, logger)
	provider := scholar.NewProvider(cfg, client)

	s := &Server{
		server:   server,
		logger:   logger,
		provider: provider,
		cfg:      cfg,
		version:  version,
	}
	s.registerTools()
	s.registerResources()
	s.registerPrompts()
	return s
}

func (s *Server) Run(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

func (s *Server) registerTools() {
	readOnly := true
	notDestructive := false
	openWorld := true

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "search_google_scholar_key_words",
		Title:       "Search Google Scholar By Keywords",
		Description: "Search Google Scholar by keyword query and return structured paper metadata.",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Search Google Scholar By Keywords",
			ReadOnlyHint:    true,
			DestructiveHint: &notDestructive,
			IdempotentHint:  true,
			OpenWorldHint:   &openWorld,
		},
	}, s.handleKeywordSearch)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "search_google_scholar_advanced",
		Title:       "Advanced Google Scholar Search",
		Description: "Search Google Scholar with author and year filters. Supports both start_year/end_year and legacy year_range.",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Advanced Google Scholar Search",
			ReadOnlyHint:    readOnly,
			DestructiveHint: &notDestructive,
			IdempotentHint:  true,
			OpenWorldHint:   &openWorld,
		},
	}, s.handleAdvancedSearch)

	if s.cfg.EnableAuthor {
		mcp.AddTool(s.server, &mcp.Tool{
			Name:        "get_author_info",
			Title:       "Get Google Scholar Author Info",
			Description: "Find a Google Scholar author profile by name and return structured author metadata. This is best-effort because Google Scholar author search is fragile.",
			Annotations: &mcp.ToolAnnotations{
				Title:           "Get Google Scholar Author Info",
				ReadOnlyHint:    readOnly,
				DestructiveHint: &notDestructive,
				IdempotentHint:  true,
				OpenWorldHint:   &openWorld,
			},
		}, s.handleAuthorInfo)
	}
}

func (s *Server) handleKeywordSearch(ctx context.Context, req *mcp.CallToolRequest, input keywordSearchInput) (*mcp.CallToolResult, model.SearchResponse, error) {
	start := time.Now()
	response, err := s.provider.Search(ctx, scholar.KeywordSearchRequest{
		Query:      input.Query,
		NumResults: input.NumResults,
	})
	return s.finishSearchTool("search_google_scholar_key_words", start, response, err)
}

func (s *Server) handleAdvancedSearch(ctx context.Context, req *mcp.CallToolRequest, input advancedSearchInput) (*mcp.CallToolResult, model.SearchResponse, error) {
	start := time.Now()

	startYear, endYear := input.StartYear, input.EndYear
	if len(input.YearRange) == 2 {
		if startYear == 0 {
			startYear = input.YearRange[0]
		}
		if endYear == 0 {
			endYear = input.YearRange[1]
		}
	}

	response, err := s.provider.AdvancedSearch(ctx, scholar.AdvancedSearchRequest{
		Query:      input.Query,
		Author:     input.Author,
		StartYear:  startYear,
		EndYear:    endYear,
		NumResults: input.NumResults,
	})
	return s.finishSearchTool("search_google_scholar_advanced", start, response, err)
}

func (s *Server) handleAuthorInfo(ctx context.Context, req *mcp.CallToolRequest, input authorInfoInput) (*mcp.CallToolResult, model.AuthorProfile, error) {
	start := time.Now()
	response, err := s.provider.GetAuthor(ctx, scholar.AuthorRequest{AuthorName: input.AuthorName})
	if err != nil {
		s.logTool("get_author_info", time.Since(start), 0, err)
		return errorResult(scholar.ClassifyMessage(err)), model.AuthorProfile{}, nil
	}

	s.logTool("get_author_info", time.Since(start), len(response.Publications), nil)
	summary := fmt.Sprintf("Found author profile for %q with %d publications.", response.Name, len(response.Publications))
	return textResult(summary), *response, nil
}

func (s *Server) finishSearchTool(name string, start time.Time, response *model.SearchResponse, err error) (*mcp.CallToolResult, model.SearchResponse, error) {
	if err != nil {
		s.logTool(name, time.Since(start), 0, err)
		return errorResult(scholar.ClassifyMessage(err)), model.SearchResponse{}, nil
	}

	s.logTool(name, time.Since(start), response.ResultsCount, nil)
	summary := fmt.Sprintf("Found %d Google Scholar results for %q.", response.ResultsCount, response.Query)
	return textResult(summary), *response, nil
}

func (s *Server) logTool(name string, duration time.Duration, results int, err error) {
	attrs := []any{
		"tool", name,
		"duration_ms", duration.Milliseconds(),
		"result_count", results,
	}
	if err != nil {
		attrs = append(attrs, "error", err, "error_class", errorClass(err))
		s.logger.Warn("tool call finished", attrs...)
		return
	}
	s.logger.Info("tool call finished", attrs...)
}

func errorClass(err error) string {
	switch {
	case errors.Is(err, scholar.ErrInvalidInput):
		return "invalid_input"
	case errors.Is(err, scholar.ErrNoResults):
		return "no_results"
	case errors.Is(err, scholar.ErrTimeout):
		return "timeout"
	case errors.Is(err, scholar.ErrUpstreamBlocked):
		return "upstream_blocked"
	case errors.Is(err, scholar.ErrUpstreamUnavailable):
		return "upstream_unavailable"
	case errors.Is(err, scholar.ErrParseFailed):
		return "parse_failed"
	default:
		return "unknown"
	}
}

func textResult(summary string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: summary},
		},
	}
}

func errorResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: message},
		},
	}
}
