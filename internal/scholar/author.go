package scholar

import (
	"context"
	"net/url"
	"strconv"

	"github.com/bingshuoguo/google-scholar-mcp/internal/model"
)

func (p *ScholarHTMLProvider) GetAuthor(ctx context.Context, req AuthorRequest) (*model.AuthorProfile, error) {
	validated, err := req.Validate()
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("view_op", "search_authors")
	query.Set("mauth", validated.AuthorName)
	query.Set("hl", "en")

	searchResponse, err := p.client.Get(ctx, "/citations", query)
	if err != nil {
		return nil, err
	}

	candidates, err := parseAuthorSearchCandidates(string(searchResponse.Body), p.cfg.BaseURL)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, wrap(ErrNoResults, "no Google Scholar author profile found", nil)
	}

	candidate := chooseAuthorCandidate(validated.AuthorName, candidates)
	if candidate == nil || candidate.ProfileURL == "" {
		return nil, wrap(ErrParseFailed, "failed to resolve Google Scholar author profile URL", nil)
	}

	profileURL, err := url.Parse(candidate.ProfileURL)
	if err != nil {
		return nil, wrap(ErrParseFailed, "invalid Google Scholar author profile URL", err)
	}

	profileResponse, err := p.client.Get(ctx, profileURL.Path, profileURL.Query())
	if err != nil {
		return nil, err
	}

	profile, err := parseAuthorProfile(string(profileResponse.Body), p.cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	if profile.ScholarID == "" {
		profile.ScholarID = candidate.ScholarID
	}
	if profile.ProfileURL == "" {
		profile.ProfileURL = candidate.ProfileURL
	}
	if profile.Affiliation == "" {
		profile.Affiliation = candidate.Affiliation
	}
	if len(profile.Interests) == 0 {
		profile.Interests = candidate.Interests
	}
	if profile.CitedBy == 0 {
		profile.CitedBy = candidate.CitedBy
		profile.CitedByLegacy = candidate.CitedBy
	}

	return profile, nil
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
