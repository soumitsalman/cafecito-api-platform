# Beans listing for the xAI Grok Build marketplace

The repository's `.grok-plugin/marketplace.json` makes the Beans plugin available as a Cafecito marketplace source. Listing it in xAI's official catalog requires a separate PR to [xai-org/plugin-marketplace](https://github.com/xai-org/plugin-marketplace/blob/main/CONTRIBUTING.md).

## Catalog entry

Add the following object to the `plugins` array in xAI's `.grok-plugin/marketplace.json`. Replace `<OWNER>` and `<PUBLIC_REPO>` with the public repository containing this plugin. Replace `<FULL_PUBLIC_COMMIT_SHA>` with the full 40-character lowercase commit SHA that contains these files. These placeholders must be replaced before submission.

```json
{
  "name": "cafecito-beans",
  "version": "1.0.0",
  "description": "Discover publisher articles across news, reports, research, and official statements; explore stories, sources, related coverage, and mentions with Beans MCP.",
  "category": "research",
  "source": {
    "source": "url",
    "url": "https://github.com/<OWNER>/<PUBLIC_REPO>.git",
    "sha": "<FULL_PUBLIC_COMMIT_SHA>",
    "path": "plugins/cafecito-beans"
  },
  "homepage": "https://developer.cafecito.tech/products/beans",
  "keywords": [
    "cafecito",
    "cafecito beans",
    "cafecito beans mcp",
    "cafecito beans api"
  ],
  "domains": [
    "cafecito.tech",
    "developer.cafecito.tech",
    "api.cafecito.tech",
    "beans.cafecito.tech"
  ]
}
```

Keep discovery keywords specific to Cafecito. xAI uses keywords and domains to suggest plugins, so generic words such as `news`, `research`, and `api` should not be used as triggers.

## Submission description

Beans connects Grok Build to Cafecito's hosted, read-only public-information MCP server. It includes three skills for article research, cross-publisher story coverage, and publisher or document-type searches. The skills cover tool selection, follow-up calls, source attribution, freshness, and cursor pagination.

Supported request content types are `news`, `blog`, `earnings_report`, `financial_report`, `lawsuit`, `enforcement_action`, `contract`, `official_statement`, `press_release`, `research_paper`, `technical_documentation`, `whitepaper`, `podcast`, and `site`.

The MCP endpoint is `https://api.cafecito.tech/beans/mcp`. Users obtain their own API key from [Cafecito API Keys](https://developer.cafecito.tech/settings/api-keys) and set `CAFECITO_API_KEY` before starting Grok. The plugin sends the key as an Authorization Bearer header to that endpoint. It requires no local executable, install script, lifecycle hook, or OAuth client. The plugin source uses this repository's MIT license; hosted API access and publisher content follow Cafecito's service terms and content policy.

## Complete the upstream PR

After adding the catalog entry in the xAI marketplace checkout, regenerate its component index and validate it:

```bash
python3 scripts/generate-plugin-index.py
python3 scripts/validate-catalog.py
python3 scripts/generate-plugin-index.py --check
```

The generated entry should show `beans-research`, `beans-source-monitoring`, `beans-story-coverage`, and the `cafecito-beans` HTTP MCP server. Include the catalog change and generated index in the PR. Later plugin updates require a new pinned SHA and regenerated index.
