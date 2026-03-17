package scholar

import (
	"context"

	"github.com/bingshuoguo/google-scholar-mcp/internal/model"
)

type SearchProvider interface {
	Search(ctx context.Context, req KeywordSearchRequest) (*model.SearchResponse, error)
	AdvancedSearch(ctx context.Context, req AdvancedSearchRequest) (*model.SearchResponse, error)
}

type AuthorProvider interface {
	GetAuthor(ctx context.Context, req AuthorRequest) (*model.AuthorProfile, error)
}

type Provider interface {
	SearchProvider
	AuthorProvider
}
