package scholar

import "testing"

func TestParseAuthorSearchCandidates(t *testing.T) {
	t.Parallel()

	html := readFixture(t, "scholar_author_search_page.html")
	candidates, err := parseAuthorSearchCandidates(html, "https://scholar.google.com")
	if err != nil {
		t.Fatalf("parseAuthorSearchCandidates returned error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(candidates))
	}

	best := chooseAuthorCandidate("Geoffrey Hinton", candidates)
	if best == nil {
		t.Fatal("expected a best candidate")
	}
	if best.ScholarID != "abc123" {
		t.Fatalf("unexpected scholar id: %q", best.ScholarID)
	}
}

func TestParseAuthorProfile(t *testing.T) {
	t.Parallel()

	html := readFixture(t, "scholar_author_profile_page.html")
	profile, err := parseAuthorProfile(html, "https://scholar.google.com")
	if err != nil {
		t.Fatalf("parseAuthorProfile returned error: %v", err)
	}
	if profile.Name != "Geoffrey Hinton" {
		t.Fatalf("unexpected profile name: %q", profile.Name)
	}
	if profile.CitedBy != 912345 {
		t.Fatalf("unexpected cited_by: %d", profile.CitedBy)
	}
	if len(profile.Publications) != 2 {
		t.Fatalf("unexpected publications count: %d", len(profile.Publications))
	}
	if profile.Publications[0].Year != 2006 {
		t.Fatalf("unexpected first publication year: %d", profile.Publications[0].Year)
	}
}
