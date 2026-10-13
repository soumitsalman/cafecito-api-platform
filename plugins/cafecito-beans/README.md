# Cafecito Beans for Cursor and Grok Build

Discover what publishers have reported across news, blogs, financial and earnings reports, litigation, official statements, research, and technical documents. Beans connects Article reading with related coverage, cross-publisher Stories, Source details, and external Mentions. The hosted [Beans MCP server](https://developer.cafecito.tech/api/beans/mcp) supplies the tools; no local server or SDK is required.

## Connect in Cursor

1. Create a Cafecito key at [API Keys](https://developer.cafecito.tech/settings/api-keys).
2. For local use before marketplace publication, copy this plugin directory to `~/.cursor/plugins/local/cafecito-beans` and reload Cursor. Once listed, install `cafecito-beans` from Cursor's marketplace. See Cursor's [plugin instructions](https://cursor.com/docs/plugins).
3. Configure the required `CAFECITO_API_KEY` variable in Cursor's plugin settings. Enter only the key value; the plugin supplies `Bearer` in the request header.
4. Enable the `cafecito-beans` MCP server in Cursor if needed.

## Connect in Grok Build

1. Create a Cafecito key at [API Keys](https://developer.cafecito.tech/settings/api-keys).
2. Set `CAFECITO_API_KEY` in the environment that starts Grok Build. Enter only the key value; the plugin adds `Bearer` to the request header.
3. Install and enable `cafecito-beans` from the Grok Build marketplace after its listing is approved. For local use before publication, run `grok plugin install ./plugins/cafecito-beans --trust` and `grok plugin enable cafecito-beans` from this repository's root.
4. Use `grok inspect` to confirm that Grok loaded the Beans skills and MCP server.

Grok Build reads the three `skills/` folders and the hosted MCP server declared in `.mcp.json`. No local server, OAuth client, or API client package is needed. See [Grok Build's plugin](https://docs.x.ai/build/features/skills-plugins-marketplaces) and [MCP](https://docs.x.ai/build/features/mcp-servers) documentation.

For Bash, enter the key without putting it in shell history, then start Grok from that shell:

```bash
read -r -s -p "Cafecito API key: " CAFECITO_API_KEY
export CAFECITO_API_KEY
printf '\n'
grok
```

The key must be available each time Grok starts. Grok expands `${CAFECITO_API_KEY}` in the Authorization header; the Cursor `variables` field is not needed for this setup. Beans tool requests go to `https://api.cafecito.tech/beans/mcp` with the API key and the requested search parameters or resource IDs.

This repository also supplies `.grok-plugin/marketplace.json`. To browse it as a local marketplace, run these commands from the repository root:

```bash
grok plugin marketplace add .
grok plugin install cafecito-beans --trust
grok plugin enable cafecito-beans
```

Refresh the Plugins tab or start a new session after installation. For public listing in xAI's catalog, use the [Grok marketplace submission content](GROK-MARKETPLACE.md).

Keep the key out of prompts and repository files. Missing or invalid keys return HTTP 401. REST and MCP calls share the same [usage limits](https://developer.cafecito.tech/guides/pricing-limits).

## Ask Beans

| Question | Skill |
| --- | --- |
| "What have publishers reported about battery storage this week? Show sources, dates, and links." | `beans-research` |
| "Which articles are trending about climate policy, and which are only the latest?" | `beans-research` |
| "Find a story about semiconductor exports and show which publishers covered it." | `beans-story-coverage` |
| "Show recent earnings reports from this publisher." | `beans-source-monitoring` |
| "Find lawsuits and official statements about this company since the start of the month." | `beans-source-monitoring` |

The skills guide tool selection and follow-up calls. Beans collections use cursor pagination and include `meta.as_of` for snapshot freshness. Coverage and update timing vary by source. Beans provides publisher material; Cafecito's Espresso product handles structured Events, Signals, impact, and outlook.

In Grok Build, invoke `/beans-research`, `/beans-story-coverage`, or `/beans-source-monitoring` to select a workflow explicitly. Each skill also describes when Grok should use it automatically.

## Filter by content type

Beans accepts these 14 `content_type` values in requests:

| Material | Values |
| --- | --- |
| News and blogs | `news`, `blog` |
| Financial updates | `earnings_report`, `financial_report` |
| Legal and regulatory | `lawsuit`, `enforcement_action`, `contract` |
| Corporate communications | `official_statement`, `press_release` |
| Research and technical | `research_paper`, `technical_documentation`, `whitepaper` |
| Audio and sites | `podcast`, `site` |

Use one value per request, or omit `content_type` for broad coverage. For multiple types, make separate calls. The filter is available on Article search, latest and trending Article feeds, similar Articles, Stories, and Story Articles. Top headlines is already news only and does not accept it. See `beans-source-monitoring` for agent guidance.

See the [Beans overview](https://developer.cafecito.tech/products/beans), [workflow examples](https://developer.cafecito.tech/products/beans/scenarios), and [MCP guide](https://developer.cafecito.tech/guides/mcp-ai-agents) for the current contract and all 16 exported tools.

The plugin source is MIT licensed under this repository's [license](../../LICENSE). Hosted API access and third-party content follow the Cafecito service terms and content policy linked in the repository's root README.
