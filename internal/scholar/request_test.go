package scholar

import "testing"

func TestKeywordSearchRequestValidate(t *testing.T) {
	t.Parallel()

	req, err := KeywordSearchRequest{Query: "  graph neural networks  "}.Validate(10)
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if req.Query != "graph neural networks" {
		t.Fatalf("unexpected normalized query: %q", req.Query)
	}
	if req.NumResults != 5 {
		t.Fatalf("unexpected default num_results: %d", req.NumResults)
	}
}

func TestKeywordSearchRequestValidateRange(t *testing.T) {
	t.Parallel()

	_, err := KeywordSearchRequest{Query: "test", NumResults: 11}.Validate(10)
	if err == nil {
		t.Fatal("expected range validation error")
	}
}

func TestAdvancedSearchRequestValidateYearRange(t *testing.T) {
	t.Parallel()

	_, err := AdvancedSearchRequest{
		Query:      "transformer interpretability",
		StartYear:  2025,
		EndYear:    2020,
		NumResults: 5,
	}.Validate(10)
	if err == nil {
		t.Fatal("expected year range validation error")
	}
}

func TestAuthorRequestValidate(t *testing.T) {
	t.Parallel()

	req, err := AuthorRequest{AuthorName: "  Geoffrey Hinton  "}.Validate()
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if req.AuthorName != "Geoffrey Hinton" {
		t.Fatalf("unexpected normalized author name: %q", req.AuthorName)
	}
}
