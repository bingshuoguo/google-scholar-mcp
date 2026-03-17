package scholar

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bingshuoguo/google-scholar-mcp/internal/model"
)

var (
	yearPattern    = regexp.MustCompile(`\b(?:19|20)\d{2}\b`)
	citedByPattern = regexp.MustCompile(`(?i)cited by (\d+)`)
	versionPattern = regexp.MustCompile(`(?i)all (\d+) versions`)
	spacePattern   = regexp.MustCompile(`\s+`)
)

func parseSearchResults(html string, limit int) ([]model.Paper, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, wrap(ErrParseFailed, "parse search page HTML", err)
	}

	containers := doc.Find("div.gs_r.gs_or.gs_scl")
	if containers.Length() == 0 {
		containers = doc.Find("div.gs_ri").ParentsFiltered("div")
	}
	if containers.Length() == 0 {
		containers = doc.Find("div.gs_ri")
	}
	if containers.Length() == 0 {
		if doc.Find("#gs_captcha_ccl, form#captcha-form").Length() > 0 {
			return nil, wrap(ErrUpstreamBlocked, "Google Scholar returned a challenge page", nil)
		}
		if doc.Find("#gs_res_ccl_mid, .gs_med").Length() > 0 {
			return []model.Paper{}, nil
		}
		return nil, wrap(ErrParseFailed, "no recognizable Google Scholar result containers found", nil)
	}

	results := make([]model.Paper, 0, min(limit, containers.Length()))
	containers.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if len(results) >= limit {
			return false
		}

		root := selection
		content := selection.Find("div.gs_ri").First()
		if content.Length() == 0 {
			content = selection
			if selection.HasClass("gs_ri") {
				root = selection.Parent()
			}
		}

		titleNode := content.Find("h3.gs_rt").First()
		title := cleanText(titleNode.Text())
		if title == "" {
			title = "Untitled Google Scholar result"
		}

		link, _ := titleNode.Find("a").Attr("href")
		authorsText := cleanText(content.Find("div.gs_a").First().Text())
		snippet := cleanText(content.Find("div.gs_rs, div.gs_rs_w").First().Text())
		publication := parsePublication(authorsText)
		paper := model.Paper{
			Title:       title,
			URL:         link,
			AuthorsText: authorsText,
			Snippet:     snippet,
			Publication: publication,
			Year:        parseFirstInt(yearPattern, authorsText),
			Source:      "google_scholar",
		}

		if pdfLink, ok := root.Find("div.gs_or_ggsm a").First().Attr("href"); ok {
			paper.PDFURL = pdfLink
		}

		content.Find("a").Each(func(_ int, anchor *goquery.Selection) {
			text := cleanText(anchor.Text())
			if paper.CitationCount == 0 {
				paper.CitationCount = parseFirstInt(citedByPattern, text)
			}
			if paper.VersionCount == 0 {
				paper.VersionCount = parseFirstInt(versionPattern, text)
			}
		})

		results = append(results, paper)
		return true
	})

	return results, nil
}

func cleanText(value string) string {
	value = strings.ReplaceAll(value, "\u00a0", " ")
	return spacePattern.ReplaceAllString(strings.TrimSpace(value), " ")
}

func parseFirstInt(pattern *regexp.Regexp, value string) int {
	match := pattern.FindStringSubmatch(value)
	if len(match) == 0 {
		return 0
	}
	token := match[0]
	if len(match) >= 2 {
		token = match[1]
	}
	number, err := strconv.Atoi(token)
	if err != nil {
		return 0
	}
	return number
}

func parsePublication(authorsText string) string {
	parts := strings.Split(authorsText, " - ")
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(strings.Join(parts[1:], " - "))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
