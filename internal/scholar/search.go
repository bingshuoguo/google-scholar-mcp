package scholar

import (
	"context"
	"net/url"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
	"github.com/bingshuoguo/google-scholar-mcp/internal/model"
)

type ScholarHTMLProvider struct {
	cfg    config.Config
	client *Client
}

func NewProvider(cfg config.Config, client *Client) *ScholarHTMLProvider {
	return &ScholarHTMLProvider{
		cfg:    cfg,
		client: client,
	}
}

func (p *ScholarHTMLProvider) Search(ctx context.Context, req KeywordSearchRequest) (*model.SearchResponse, error) {
	validated, err := req.Validate(p.cfg.MaxResults)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("q", validated.Query)
	query.Set("hl", "en")

	response, err := p.client.Get(ctx, "/scholar", query)
	if err != nil {
		return nil, err
	}

	results, err := parseSearchResults(string(response.Body), validated.NumResults)
	if err != nil {
		return nil, err
	}

	return &model.SearchResponse{
		Query:        validated.Query,
		ResultsCount: len(results),
		Results:      results,
	}, nil
}

func (p *ScholarHTMLProvider) AdvancedSearch(ctx context.Context, req AdvancedSearchRequest) (*model.SearchResponse, error) {
	validated, err := req.Validate(p.cfg.MaxResults)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("q", validated.Query)
	query.Set("hl", "en")
	if validated.Author != "" {
		query.Set("as_auth", validated.Author)
	}
	if validated.StartYear != 0 {
		query.Set("as_ylo", itoa(validated.StartYear))
	}
	if validated.EndYear != 0 {
		query.Set("as_yhi", itoa(validated.EndYear))
	}

	response, err := p.client.Get(ctx, "/scholar", query)
	if err != nil {
		return nil, err
	}

	results, err := parseSearchResults(string(response.Body), validated.NumResults)
	if err != nil {
		return nil, err
	}

	return &model.SearchResponse{
		Query:        validated.Query,
		ResultsCount: len(results),
		Results:      results,
	}, nil
}
