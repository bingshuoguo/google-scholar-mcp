package scholar

import (
	"fmt"
	"strings"
)

const defaultNumResults = 5

type KeywordSearchRequest struct {
	Query      string
	NumResults int
}

type AdvancedSearchRequest struct {
	Query      string
	Author     string
	StartYear  int
	EndYear    int
	NumResults int
}

type AuthorRequest struct {
	AuthorName string
}

func (r KeywordSearchRequest) Validate(maxResults int) (KeywordSearchRequest, error) {
	r.Query = strings.TrimSpace(r.Query)
	if r.Query == "" {
		return r, wrap(ErrInvalidInput, "query is required", nil)
	}

	if r.NumResults == 0 {
		r.NumResults = defaultNumResults
	}
	if r.NumResults < 1 || r.NumResults > maxResults {
		return r, wrap(ErrInvalidInput, fmt.Sprintf("num_results must be between 1 and %d", maxResults), nil)
	}

	return r, nil
}

func (r AdvancedSearchRequest) Validate(maxResults int) (AdvancedSearchRequest, error) {
	validated, err := KeywordSearchRequest{Query: r.Query, NumResults: r.NumResults}.Validate(maxResults)
	if err != nil {
		return r, err
	}
	r.Query = validated.Query
	r.NumResults = validated.NumResults
	r.Author = strings.TrimSpace(r.Author)

	if r.StartYear != 0 && r.StartYear < 1000 {
		return r, wrap(ErrInvalidInput, "start_year must be a valid year", nil)
	}
	if r.EndYear != 0 && r.EndYear < 1000 {
		return r, wrap(ErrInvalidInput, "end_year must be a valid year", nil)
	}
	if r.StartYear != 0 && r.EndYear != 0 && r.StartYear > r.EndYear {
		return r, wrap(ErrInvalidInput, "start_year must be less than or equal to end_year", nil)
	}

	return r, nil
}

func (r AuthorRequest) Validate() (AuthorRequest, error) {
	r.AuthorName = strings.TrimSpace(r.AuthorName)
	if r.AuthorName == "" {
		return r, wrap(ErrInvalidInput, "author_name is required", nil)
	}
	return r, nil
}
