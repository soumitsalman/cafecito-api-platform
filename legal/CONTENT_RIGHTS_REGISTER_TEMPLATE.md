# Content Rights Register Template

**Status:** Internal specification for counsel and engineering review  
**Default:** Unknown means deny full text, durable storage, embeddings/RAG, redistribution, and training.

The register must be authoritative for automated source and output decisions. A spreadsheet may support review, but runtime enforcement must consume a versioned policy record.

## Required source record

| Group | Field | Meaning |
| --- | --- | --- |
| Identity | `source_id` | Stable internal source identifier |
| Identity | `legal_name` | Rights-holder or publisher legal name |
| Identity | `domains` | Covered domains and subdomain rules |
| Identity | `source_references` | Covered public-information source references |
| Identity | `path_rules` | Included and excluded URL paths |
| Identity | `jurisdictions` | Relevant source/operator jurisdictions |
| Scope | `territories` | Territories where rights may be exercised |
| Scope | `languages` | Covered content and output languages |
| Scope | `products` | Covered Cafecito products |
| Scope | `customer_tiers` | Covered customer plans or negotiated accounts |
| Access | `crawl` | deny, allow, allow_with_conditions, unknown |
| Access | `feed_access` | deny, allow, allow_with_conditions, unknown |
| Access | `source_terms_urls` | Applicable site, feed, API, and publisher terms |
| Access | `source_terms_version_or_hash` | Evidence of the terms version reviewed |
| Access | `source_terms_reviewed_at` | Last legal/rights review time |
| Access | `robots_policy` | Current directives and last check |
| Access | `authentication_required` | Whether login, subscription, cookie, or credential is required |
| Access | `technical_controls` | Paywall, CAPTCHA, blocks, or other controls |
| Access | `request_rate` | Contractual or source-requested rate |
| Access | `cease_and_desist_history` | Notices, objections, and resulting restrictions |
| Rights | `metadata` | URL/title/byline/date/source use |
| Rights | `excerpt` | Excerpt allowed and any length/placement rule |
| Rights | `summary` | Generated summary permission |
| Rights | `full_display` | Complete text in API/UI |
| Rights | `format_conversion` | HTML/feed to Markdown/text/JSON/YAML/TOON |
| Rights | `indexing` | Search indexing |
| Rights | `temporary_processing` | Transient copies for enrichment |
| Rights | `persistent_storage` | Stored source text |
| Rights | `embeddings` | Vector representation permission |
| Rights | `rag` | Use as inference context |
| Rights | `training` | Model training/fine-tuning |
| Rights | `translation` | Translation permission |
| Rights | `analytics` | Classifications, trends, and derived data |
| Channels | `rest`, `mcp`, `portal`, `feeds`, `exports` | Permitted delivery surfaces |
| Downstream | `customer_display` | Internal, single display, public display, or denied |
| Downstream | `redistribution` | Sublicense/resale/republishing scope |
| Downstream | `customer_retention` | Maximum customer retention |
| Downstream | `customer_training` | Whether customers may train models |
| Attribution | `required_fields` | URL, publisher, author, date, notice, link text |
| Attribution | `brand_rules_url` | Approved marks and usage rules |
| Exclusions | `syndicated`, `wire`, `images`, `video`, `charts`, `embeds` | Material not covered |
| Lifecycle | `license_status` | draft, active, suspended, expired, terminated, or revoked |
| Lifecycle | `effective_at`, `expires_at`, `revoked_at` | Rights lifecycle dates |
| Lifecycle | `termination_notice` | Notice period |
| Lifecycle | `purge_deadline` | Required deletion deadline |
| Lifecycle | `derived_data_survives` | Which derived fields may remain |
| Lifecycle | `correction_obligation` | Refresh/withdrawal requirements |
| Evidence | `basis` | signed_license, direct_permission, open_license, public_domain, counsel_approved_fair_use, unknown |
| Evidence | `open_license_id`, `open_license_version` | Exact open-license evidence, when applicable |
| Evidence | `agreement_id` | Contract/permission reference |
| Evidence | `evidence_location` | Controlled link to signed evidence |
| Evidence | `approved_by`, `approved_at` | Accountable reviewer and time |
| Evidence | `next_review_at` | Mandatory re-review date |
| Enforcement | `blocked` | Emergency deny switch |
| Enforcement | `legal_hold` | Preservation override that does not itself allow display |
| Enforcement | `policy_version` | Version used for audit decisions |

## Example policy record

```yaml
source_id: SOURCE_UUID
legal_name: Example Publisher, Inc.
domains:
  - example.com
feeds:
  - https://example.com/source-reference
territories: [US]
languages: [en]
products: [beans]
customer_tiers: [licensed_full_display]
path_rules:
  include: [/news/]
  exclude: [/wire/, /photos/, /video/]
access:
  crawl: allow_with_conditions
  feed_access: allow
  source_terms_urls: [https://example.com/terms]
  source_terms_version_or_hash: TERMS_HASH
  source_terms_reviewed_at: 2026-09-11T00:00:00Z
  request_rate: 30_per_minute
  cease_and_desist_history: []
rights:
  metadata: allow
  excerpt:
    decision: allow
    max_characters: 300
  summary: allow
  full_display: allow_with_conditions
  format_conversion: allow
  indexing: allow
  temporary_processing: allow
  persistent_storage:
    decision: allow
    max_age: 24h
  embeddings: deny
  rag: deny
  training: deny
channels:
  rest: allow
  mcp: allow
  portal: deny
  feeds: deny
  exports: deny
downstream:
  customer_display: single_end_user_display
  redistribution: deny
  customer_retention: 24h
  customer_training: deny
attribution:
  required_fields: [publisher, author, canonical_url, published_at]
exclusions:
  syndicated: true
  wire: true
  images: true
  video: true
  charts: true
  embeds: true
lifecycle:
  license_status: active
  effective_at: 2026-09-11T00:00:00Z
  expires_at: 2027-09-10T23:59:59Z
  revoked_at: null
  purge_deadline: 24h_after_termination
  derived_data_survives: [canonical_url, factual_metadata]
evidence:
  basis: signed_license
  open_license_id: null
  open_license_version: null
  agreement_id: PUBLISHER-2026-001
  evidence_location: CONTROLLED_CONTRACT_URL
  approved_by: LEGAL_APPROVER
  approved_at: 2026-09-11T00:00:00Z
  next_review_at: 2027-08-01T00:00:00Z
enforcement:
  blocked: false
  legal_hold: false
  policy_version: 1
```

## Request decision

For each record and request:

1. Resolve canonical URL, source, feed, path, and content components.
2. Apply emergency block, publisher opt-out, court order, and license termination first.
3. Deny access that would require authentication or circumvention unless a contract expressly authorizes it.
4. Confirm the requested channel, territory, customer tier, use, and time are within scope.
5. Apply the most restrictive applicable component rule, including syndicated and media exclusions.
6. Return only the fields permitted by the effective policy.
7. Attach downstream rights metadata.
8. Record policy version, decision, reason, request/account, source, output class, and time.
9. For unknown or conflicting facts, deny full text and escalate.

## Response rights object

Before full-display launch, add a stable public object similar to:

```json
{
  "rights": {
    "scope": "full_display",
    "basis": "licensed",
    "attribution_required": true,
    "retention_seconds": 86400,
    "redistribution_allowed": false,
    "training_allowed": false,
    "media_included": false,
    "expires_at": "2027-09-10T23:59:59Z",
    "policy_url": "https://cafecito.tech/docs/third-party-content-policy/"
  }
}
```

Do not expose private contract terms, agreement IDs, prices, legal analysis, or internal evidence locations.

## Required enforcement tests

- restricted third-party bodies are absent from all REST and MCP routes and formats;
- unknown sources never return content;
- license expiration cuts off content without deployment;
- excluded paths and syndicated items override domain permission;
- media exclusions prevent copied media even when text is licensed;
- collection, detail, similar, Story Article, feed, and export routes apply identical policy;
- cache hits apply current policy rather than historical policy;
- customer tier and channel restrictions are enforced;
- attribution and rights fields appear when required;
- publisher opt-out and emergency block take effect immediately;
- purge completes within the contractual deadline; and
- audit records identify the exact policy version used.
