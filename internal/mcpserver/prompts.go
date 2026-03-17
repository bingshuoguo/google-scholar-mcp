package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerPrompts() {
	s.server.AddPrompt(&mcp.Prompt{
		Name:        searchPromptName,
		Title:       "Scholar Literature Scan",
		Description: "Guide the model through a keyword-oriented Google Scholar literature scan.",
		Arguments: []*mcp.PromptArgument{
			{Name: "query", Description: "Primary Google Scholar keyword query", Required: true},
			{Name: "focus", Description: "Optional research focus or evaluation angle", Required: false},
		},
	}, s.handleLiteratureScanPrompt)

	if s.cfg.EnableAuthor {
		s.server.AddPrompt(&mcp.Prompt{
			Name:        authorPromptName,
			Title:       "Scholar Author Brief",
			Description: "Guide the model through an author-centric research brief using Google Scholar author metadata.",
			Arguments: []*mcp.PromptArgument{
				{Name: "author_name", Description: "Author name to look up in Google Scholar", Required: true},
				{Name: "research_goal", Description: "Optional goal for the author brief", Required: false},
			},
		}, s.handleAuthorBriefPrompt)
	}
}

func (s *Server) handleLiteratureScanPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	query, err := requiredPromptArgument(req, "query")
	if err != nil {
		return nil, err
	}
	focus := optionalPromptArgument(req, "focus")

	return &mcp.GetPromptResult{
		Description: "Workflow prompt for a Scholar literature scan",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.EmbeddedResource{
					Resource: &mcp.ResourceContents{
						URI:      serverToolsURI,
						MIMEType: "text/markdown",
						Text:     s.toolsResourceText(),
					},
				},
			},
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: literatureScanPromptText(query, focus)},
			},
		},
	}, nil
}

func (s *Server) handleAuthorBriefPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	authorName, err := requiredPromptArgument(req, "author_name")
	if err != nil {
		return nil, err
	}
	researchGoal := optionalPromptArgument(req, "research_goal")

	return &mcp.GetPromptResult{
		Description: "Workflow prompt for an author-centric Google Scholar brief",
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.EmbeddedResource{
					Resource: &mcp.ResourceContents{
						URI:      serverLimitationsURI,
						MIMEType: "text/markdown",
						Text:     limitationsResourceText(),
					},
				},
			},
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: authorBriefPromptText(authorName, researchGoal)},
			},
		},
	}, nil
}

func requiredPromptArgument(req *mcp.GetPromptRequest, name string) (string, error) {
	value := strings.TrimSpace(optionalPromptArgument(req, name))
	if value == "" {
		return "", fmt.Errorf("missing required prompt argument %q", name)
	}
	return value, nil
}

func optionalPromptArgument(req *mcp.GetPromptRequest, name string) string {
	if req == nil || req.Params == nil || req.Params.Arguments == nil {
		return ""
	}
	return req.Params.Arguments[name]
}

func literatureScanPromptText(query, focus string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Run a literature scan for the Google Scholar query %q.\n\n", query)
	b.WriteString("Use search_google_scholar_key_words first. If the result set looks too broad or needs author/year narrowing, follow up with search_google_scholar_advanced.\n\n")
	if strings.TrimSpace(focus) != "" {
		fmt.Fprintf(&b, "Focus the analysis on: %s.\n\n", focus)
	}
	b.WriteString("In the final answer:\n")
	b.WriteString("- group results by theme or method when possible\n")
	b.WriteString("- mention year, citation_count, version_count, source, and any obvious PDF availability\n")
	b.WriteString("- state clearly that Scholar snippets are result-page excerpts, not guaranteed full abstracts\n")
	b.WriteString("- call out missing coverage or low-confidence inferences explicitly\n")
	return b.String()
}

func authorBriefPromptText(authorName, researchGoal string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Prepare a Google Scholar author brief for %q.\n\n", authorName)
	b.WriteString("Use get_author_info first. If the profile looks ambiguous or sparse, say so rather than over-claiming.\n\n")
	if strings.TrimSpace(researchGoal) != "" {
		fmt.Fprintf(&b, "Research goal: %s.\n\n", researchGoal)
	}
	b.WriteString("In the final answer:\n")
	b.WriteString("- summarize affiliation, interests, cited_by, and notable publication patterns\n")
	b.WriteString("- separate profile facts from your own interpretation\n")
	b.WriteString("- explicitly note that author lookup is best-effort and may miss or mis-rank profiles\n")
	return b.String()
}
