# Tools

This server exposes three read-only Scholar tools.

## `search_google_scholar_key_words`

Purpose:
- Search Google Scholar by a plain keyword query.

Arguments:
- `query` required
- `num_results` optional

Returns:
- `SearchResponse`
- fields: `query`, `results_count`, `results[]`

Result fields:
- `title`
- `url`
- `authors_text`
- `snippet`
- `publication`
- `year`
- `citation_count`
- `version_count`
- `pdf_url`
- `source`

## `search_google_scholar_advanced`

Purpose:
- Search Google Scholar with optional author and year filters.

Arguments:
- `query` required
- `author` optional
- `start_year` optional
- `end_year` optional
- `year_range` optional legacy compatibility field
- `num_results` optional

Returns:
- `SearchResponse`

## `get_author_info`

Purpose:
- Find a best-effort Google Scholar author profile by name.

Arguments:
- `author_name` required

Returns:
- `AuthorProfile`
- fields: `name`, `affiliation`, `interests`, `cited_by`, `citedby`, `scholar_id`, `profile_url`, `publications[]`

## Notes

- All tools are read-only.
- Search snippets are result-page excerpts, not guaranteed full abstracts.
- `get_author_info` is best-effort and depends on Scholar author search quality.
