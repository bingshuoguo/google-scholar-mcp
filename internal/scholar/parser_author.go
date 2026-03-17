package scholar

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bingshuoguo/google-scholar-mcp/internal/model"
)

type authorCandidate struct {
	Name        string
	Affiliation string
	Interests   []string
	CitedBy     int
	ScholarID   string
	ProfileURL  string
}

func parseAuthorSearchCandidates(html, baseURL string) ([]authorCandidate, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, wrap(ErrParseFailed, "parse author search page HTML", err)
	}

	users := doc.Find("div.gsc_1usr")
	if users.Length() == 0 {
		if doc.Find("#gs_captcha_ccl, form#captcha-form").Length() > 0 {
			return nil, wrap(ErrUpstreamBlocked, "Google Scholar returned an author challenge page", nil)
		}
		return []authorCandidate{}, nil
	}

	candidates := make([]authorCandidate, 0, users.Length())
	users.Each(func(_ int, selection *goquery.Selection) {
		nameNode := selection.Find("h3.gsc_1usr_name a").First()
		name := cleanText(nameNode.Text())
		href, _ := nameNode.Attr("href")
		profileURL, scholarID := resolveScholarURL(baseURL, href)

		candidate := authorCandidate{
			Name:        name,
			Affiliation: cleanText(selection.Find(".gsc_1usr_aff").First().Text()),
			CitedBy:     parseTrailingInt(selection.Find(".gsc_1usr_cby").First().Text()),
			ScholarID:   scholarID,
			ProfileURL:  profileURL,
		}

		selection.Find(".gsc_1usr_int a").Each(func(_ int, interest *goquery.Selection) {
			if text := cleanText(interest.Text()); text != "" {
				candidate.Interests = append(candidate.Interests, text)
			}
		})

		candidates = append(candidates, candidate)
	})

	return candidates, nil
}

func chooseAuthorCandidate(name string, candidates []authorCandidate) *authorCandidate {
	if len(candidates) == 0 {
		return nil
	}

	normalizedTarget := normalizeName(name)
	for _, candidate := range candidates {
		if normalizeName(candidate.Name) == normalizedTarget {
			c := candidate
			return &c
		}
	}

	candidate := candidates[0]
	return &candidate
}

func parseAuthorProfile(html, baseURL string) (*model.AuthorProfile, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, wrap(ErrParseFailed, "parse author profile HTML", err)
	}

	if doc.Find("#gs_captcha_ccl, form#captcha-form").Length() > 0 {
		return nil, wrap(ErrUpstreamBlocked, "Google Scholar returned an author profile challenge page", nil)
	}

	name := cleanText(doc.Find("#gsc_prf_in").First().Text())
	if name == "" {
		return nil, wrap(ErrParseFailed, "author profile name not found", nil)
	}

	profile := &model.AuthorProfile{
		Name:        name,
		Affiliation: cleanText(doc.Find("#gsc_prf_i .gsc_prf_il").First().Text()),
	}

	doc.Find("#gsc_prf_int a, .gsc_prf_inta").Each(func(_ int, interest *goquery.Selection) {
		if text := cleanText(interest.Text()); text != "" {
			profile.Interests = append(profile.Interests, text)
		}
	})

	stats := doc.Find("#gsc_rsb_st td.gsc_rsb_std")
	if stats.Length() > 0 {
		profile.CitedBy = parseTrailingInt(stats.First().Text())
		profile.CitedByLegacy = profile.CitedBy
	}

	if href, ok := doc.Find("a[href*=\"/citations?user=\"]").First().Attr("href"); ok {
		profile.ProfileURL, profile.ScholarID = resolveScholarURL(baseURL, href)
	}

	doc.Find("tr.gsc_a_tr").Each(func(_ int, row *goquery.Selection) {
		title := cleanText(row.Find(".gsc_a_at").First().Text())
		if title == "" {
			return
		}

		yearText := cleanText(row.Find(".gsc_a_y span, .gsc_a_y").First().Text())
		citationsText := cleanText(row.Find(".gsc_a_c a, .gsc_a_c").First().Text())

		profile.Publications = append(profile.Publications, model.Publication{
			Title:     title,
			Year:      parseTrailingInt(yearText),
			Citations: parseTrailingInt(citationsText),
		})
	})

	return profile, nil
}

func parseTrailingInt(value string) int {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return 0
	}
	last := fields[len(fields)-1]
	number, err := strconv.Atoi(strings.TrimSpace(last))
	if err == nil {
		return number
	}
	return parseFirstInt(citedByPattern, value)
}

func normalizeName(value string) string {
	value = strings.ToLower(cleanText(value))
	value = strings.ReplaceAll(value, ".", "")
	return value
}

func resolveScholarURL(baseURL, href string) (string, string) {
	if href == "" {
		return "", ""
	}

	base, _ := url.Parse(baseURL)
	ref, err := url.Parse(href)
	if err != nil {
		return "", ""
	}
	resolved := base.ResolveReference(ref)
	return resolved.String(), resolved.Query().Get("user")
}
