# Prompts

The server exposes MCP prompts for clients that support prompt discovery and execution.

## `scholar_literature_scan`

Purpose:
- Guide the model through a keyword-oriented Scholar literature scan.

Arguments:
- `query` required
- `focus` optional

Expected behavior:
- start with `search_google_scholar_key_words`
- optionally refine with `search_google_scholar_advanced`
- summarize themes, years, citation counts, and coverage gaps

## `scholar_author_brief`

Purpose:
- Guide the model through an author-centric Scholar research brief.

Arguments:
- `author_name` required
- `research_goal` optional

Expected behavior:
- start with `get_author_info`
- separate profile facts from interpretation
- explicitly mention author-search ambiguity when relevant
