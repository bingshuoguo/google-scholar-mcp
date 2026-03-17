package model

type Publication struct {
	Title     string `json:"title"`
	Year      int    `json:"year,omitempty"`
	Citations int    `json:"citations,omitempty"`
}

type AuthorProfile struct {
	Name          string        `json:"name"`
	Affiliation   string        `json:"affiliation,omitempty"`
	Interests     []string      `json:"interests,omitempty"`
	CitedBy       int           `json:"cited_by,omitempty"`
	CitedByLegacy int           `json:"citedby,omitempty"`
	ScholarID     string        `json:"scholar_id,omitempty"`
	ProfileURL    string        `json:"profile_url,omitempty"`
	Publications  []Publication `json:"publications,omitempty"`
}
