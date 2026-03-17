package mcpserver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerResources() {
	s.server.AddResource(&mcp.Resource{
		Name:        "server_overview",
		Title:       "Google Scholar Server Overview",
		Description: "High-level description of the server, version, transport, and enabled capabilities.",
		MIMEType:    "text/markdown",
		URI:         serverOverviewURI,
	}, s.handleOverviewResource)

	s.server.AddResource(&mcp.Resource{
		Name:        "tool_catalog",
		Title:       "Google Scholar Tool Catalog",
		Description: "Catalog of tools, arguments, and return shapes exposed by this MCP server.",
		MIMEType:    "text/markdown",
		URI:         serverToolsURI,
	}, s.handleToolsResource)

	s.server.AddResource(&mcp.Resource{
		Name:        "runtime_config",
		Title:       "Google Scholar Runtime Config",
		Description: "Current runtime configuration values relevant to Scholar scraping behavior.",
		MIMEType:    "text/markdown",
		URI:         serverConfigURI,
	}, s.handleConfigResource)

	s.server.AddResource(&mcp.Resource{
		Name:        "limitations",
		Title:       "Google Scholar Limitations",
		Description: "Operational boundaries and important caveats for this server.",
		MIMEType:    "text/markdown",
		URI:         serverLimitationsURI,
	}, s.handleLimitationsResource)

	s.server.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "search_guide",
		Title:       "Topic Search Guide",
		Description: "Topic-specific guidance for how to search and interpret Google Scholar results through this server.",
		MIMEType:    "text/markdown",
		URITemplate: searchGuideTemplate,
	}, s.handleSearchGuideResource)
}

func (s *Server) handleOverviewResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return markdownResource(req.Params.URI, s.overviewResourceText()), nil
}

func (s *Server) handleToolsResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return markdownResource(req.Params.URI, s.toolsResourceText()), nil
}

func (s *Server) handleConfigResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return markdownResource(req.Params.URI, s.configResourceText()), nil
}

func (s *Server) handleLimitationsResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return markdownResource(req.Params.URI, limitationsResourceText()), nil
}

func (s *Server) handleSearchGuideResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	topic, err := topicFromGuideURI(req.Params.URI)
	if err != nil {
		return nil, err
	}
	return markdownResource(req.Params.URI, searchGuideResourceText(topic)), nil
}

func markdownResource(uri, text string) *mcp.ReadResourceResult {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "text/markdown",
				Text:     text,
			},
		},
	}
}

func topicFromGuideURI(uri string) (string, error) {
	rawTopic := strings.TrimPrefix(uri, searchGuidePrefix)
	if rawTopic == uri || rawTopic == "" {
		return "", fmt.Errorf("expected resource URI beginning with %q", searchGuidePrefix)
	}

	topic, err := url.PathUnescape(rawTopic)
	if err != nil {
		return "", fmt.Errorf("decode topic from resource URI: %w", err)
	}
	if strings.TrimSpace(topic) == "" {
		return "", fmt.Errorf("search guide topic must not be empty")
	}
	return topic, nil
}
