package scholar

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSearchResults(t *testing.T) {
	t.Parallel()

	html := readFixture(t, "scholar_search_page.html")
	results, err := parseSearchResults(html, 5)
	if err != nil {
		t.Fatalf("parseSearchResults returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	first := results[0]
	if first.Title != "A Survey on Graph Neural Networks" {
		t.Fatalf("unexpected first title: %q", first.Title)
	}
	if first.CitationCount != 1234 {
		t.Fatalf("unexpected citation count: %d", first.CitationCount)
	}
	if first.VersionCount != 7 {
		t.Fatalf("unexpected version count: %d", first.VersionCount)
	}
	if first.Year != 2021 {
		t.Fatalf("unexpected year: %d", first.Year)
	}
	if first.PDFURL == "" {
		t.Fatal("expected pdf url to be parsed")
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("..", "..", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}
