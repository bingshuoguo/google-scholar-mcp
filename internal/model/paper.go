package model

type Paper struct {
	Title         string `json:"title"`
	URL           string `json:"url,omitempty"`
	AuthorsText   string `json:"authors_text,omitempty"`
	Snippet       string `json:"snippet,omitempty"`
	Publication   string `json:"publication,omitempty"`
	Year          int    `json:"year,omitempty"`
	CitationCount int    `json:"citation_count,omitempty"`
	VersionCount  int    `json:"version_count,omitempty"`
	PDFURL        string `json:"pdf_url,omitempty"`
	Source        string `json:"source"`
}

type SearchResponse struct {
	Query        string  `json:"query"`
	ResultsCount int     `json:"results_count"`
	Results      []Paper `json:"results"`
}
