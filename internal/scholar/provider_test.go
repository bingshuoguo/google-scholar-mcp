package scholar

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
)

func TestProviderSearchBuildsExpectedQuery(t *testing.T) {
	t.Parallel()

	var gotPath string
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, readFixture(t, "scholar_search_page.html"))
	}))
	defer server.Close()

	provider := newTestProvider(server.URL)
	response, err := provider.Search(context.Background(), KeywordSearchRequest{
		Query:      "graph neural networks",
		NumResults: 2,
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if gotPath != "/scholar" {
		t.Fatalf("unexpected request path: %q", gotPath)
	}
	if gotQuery != "hl=en&q=graph+neural+networks" && gotQuery != "q=graph+neural+networks&hl=en" {
		t.Fatalf("unexpected request query: %q", gotQuery)
	}
	if response.ResultsCount != 2 {
		t.Fatalf("unexpected result count: %d", response.ResultsCount)
	}
}

func TestProviderAdvancedSearchBuildsExpectedQuery(t *testing.T) {
	t.Parallel()

	var queryValues map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryValues = map[string]string{
			"q":       r.URL.Query().Get("q"),
			"as_auth": r.URL.Query().Get("as_auth"),
			"as_ylo":  r.URL.Query().Get("as_ylo"),
			"as_yhi":  r.URL.Query().Get("as_yhi"),
			"hl":      r.URL.Query().Get("hl"),
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, readFixture(t, "scholar_search_page.html"))
	}))
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.AdvancedSearch(context.Background(), AdvancedSearchRequest{
		Query:      "transformer interpretability",
		Author:     "Yoshua Bengio",
		StartYear:  2020,
		EndYear:    2024,
		NumResults: 1,
	})
	if err != nil {
		t.Fatalf("AdvancedSearch returned error: %v", err)
	}

	if queryValues["q"] != "transformer interpretability" {
		t.Fatalf("unexpected q value: %q", queryValues["q"])
	}
	if queryValues["as_auth"] != "Yoshua Bengio" {
		t.Fatalf("unexpected as_auth value: %q", queryValues["as_auth"])
	}
	if queryValues["as_ylo"] != "2020" || queryValues["as_yhi"] != "2024" {
		t.Fatalf("unexpected year values: %v", queryValues)
	}
	if queryValues["hl"] != "en" {
		t.Fatalf("unexpected hl value: %q", queryValues["hl"])
	}
}

func newTestProvider(baseURL string) *ScholarHTMLProvider {
	cfg := config.Config{
		Transport:      "stdio",
		BaseURL:        baseURL,
		Timeout:        5 * time.Second,
		MaxResults:     10,
		RateLimitRPS:   10,
		UserAgent:      "test-agent",
		AcceptLanguage: "en-US,en;q=0.9",
		EnableAuthor:   true,
		LogLevel:       slog.LevelError,
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewProvider(cfg, NewClient(cfg, logger))
}
