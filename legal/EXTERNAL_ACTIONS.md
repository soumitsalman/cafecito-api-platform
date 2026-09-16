# External Legal and Operational Actions

**Prepared:** September 11, 2026

These items cannot be completed by repository text alone. Do not mark an item complete without the listed evidence.

## P0 — before relying on the new policies or broad full-content output

| # | Owner | Action | Completion evidence |
| --- | --- | --- | --- |
| 1 | U.S. copyright/internet counsel | Review Cafecito's actual collection, RSS ingestion, scraping, enrichment, body retention, storage, and every `full_content` route. Provide a written rule for unlicensed metadata, excerpts, summaries, embeddings/RAG, and full text. | Written product-risk memorandum approved by management |
| 2 | Product + counsel | Decide whether to disable all unknown-source full text immediately or maintain a counsel-approved allowlist during transition. | Signed decision and enforced source policy |
| 2A | Operations + engineering | Produce an Article-level provenance and body-retention inventory for every source class, including feeds, official records, platform self-posts, manual/legacy rows, page scrapes, processing caches, backups, logs, and fixtures. | Reviewed inventory, data-flow diagram, and retention/purge owner |
| 2B | Release owner + counsel | Publish only approved final policies with a confirmed effective date; notify existing users of material changes and collect renewed affirmative assent where counsel requires it. | Live-page capture, version archive, delivery record, and acceptance report |
| 3 | Counsel | Review and approve the public Terms, Privacy Policy, AUP, Content Rights Policy, Copyright Policy, Publisher Request Policy, Corrections Policy, and Automated Collection Policy. | Dated approval and approved policy versions |
| 4 | Corporate counsel/manager | Verify the LLC's exact legal name, active status, assumed-name filings, public address, contracting authority, and governing-law/venue choice. | State records and manager approval |
| 5 | Operations | Create monitored role addresses for legal, privacy, copyright, publishers, security, and abuse; add a business telephone number where required. | Working inboxes, routing, owner, backup owner, test messages |
| 6 | Copyright counsel | Determine which § 512 safe harbors could apply to Cafecito's own activity and user-directed functions. If advisable, register a designated DMCA agent and publish the exact registered details. | Copyright Office registration and matching public contact |
| 6A | Copyright counsel + operations | Before relying on any § 512 safe harbor, adopt, inform users of, and reasonably implement a repeat-infringer policy; accommodate standard technical measures; operate notice, user-notification, counter-notice, restoration, and escalation workflows; publish the registered agent's exact contact details including required telephone information; and calendar the designation renewal cycle. | Approved SOPs, public policy, training record, test cases, registration and renewal calendar |
| 7 | Publisher relations + counsel | Inventory every active feed/source and classify it as licensed, public domain/open license, permission pending, or unknown. | Source inventory with evidence links and reviewer/date |
| 8 | Publisher relations + counsel | Obtain full-display/API redistribution permission for every source whose full content remains available, including feeds, public records, platform self-posts, manual/legacy ingestion, and any other non-page-scrape source. Do not rely on delivery availability alone. | Signed agreement or verified public-domain/open-license basis and completed rights schedule per source |
| 9 | Engineering | Implement the content-rights registry and deny full text by default unless rights permit it. Treat unknown, NULL, legacy, revoked, expired, and source-unmatched records as restricted. Cover REST, MCP, private routes, all formats, collection/detail/related routes, caches, logs, and exports. | Automated policy tests and production verification |
| 10 | Engineering | Replace browser-like automated request identity with a dedicated Cafecito crawler user agent, publish the identity/contact page, and implement auditable robots/source-policy and rate controls. | Production request evidence and policy tests |
| 11 | Product/engineering | Implement affirmative clickwrap at signup or API-key creation and store acceptance evidence by policy version. | Acceptance UI, records, exportable audit report |
| 12 | Privacy counsel + engineering | Complete the data inventory, retention schedule, cookie/analytics review, subprocessors list, privacy-request workflow, notice-at-collection analysis, opt-out preference-signal/GPC analysis, privacy appeal workflow, and state applicability analysis. Reconcile the public policy to reality. | Signed data map, schedule, request SOP, GPC/appeal decision, and revised policy if needed |
| 12A | Privacy counsel | Determine whether Cafecito is a data broker in California or another state and whether California DROP or other registration, deletion, or consumer-right duties apply before broad operation. | Written analysis and registrations/workflows if required |
| 12B | Washington privacy counsel | Assess Washington My Health My Data for publisher content, health-related inference, queries, and customer use before broad operation. | Written applicability analysis and separate policy/consent flow if required |
| 13 | Security + counsel | Establish an incident-response plan and current breach-notification matrix for all U.S. states, D.C., Puerto Rico, and U.S. territories. | Approved plan, contacts, tabletop record |
| 14 | Billing/product counsel | Review checkout, renewal, cancellation, refunds, pricing disclosures, taxes, and commercial email before paid subscriptions or marketing. | Approved screenshots/flows and compliance memo |

## P1 — before licensed full display at scale

| # | Owner | Action | Completion evidence |
| --- | --- | --- | --- |
| 15 | Content counsel | Finalize a Master Publisher Content License Agreement and publisher order form/rights schedule. | Approved templates |
| 16 | Publisher relations | Identify and exclude wire-service, syndicated, freelance, image, video, chart, and other material the publisher cannot sublicense. | Contract schedules and enforcement rules |
| 17 | Engineering | Add machine-readable response rights: scope, attribution, retention, redistribution, training, expiration, withdrawal, and license reference. | Public contract, tests, and live response |
| 18 | Engineering + operations | Implement urgent suppression, license-expiry cutoff, cache purge, downstream notification, correction/retraction propagation, and legal holds. | Runbook exercise and audit logs |
| 19 | Finance + publisher relations | Define royalty units, retries/cache-hit rules, reporting, invoices, audit support, and revenue recognition. | Reconciled sample statement |
| 20 | Editorial/legal | Approve heightened-risk rules for crime, health, finance, minors, elections, sanctions, allegations, identity matching, corrections, and retractions. | Editorial standard and review escalation |
| 21 | Insurance broker + counsel | Obtain and compare media liability, Tech E&O, cyber/privacy, and IP coverage; review exclusions for scraping, copyright, defamation, and prior acts. | Bound policies or signed risk acceptance |
| 22 | Accessibility owner | Audit the portal and account/API-key flows against WCAG 2.2 AA; fix barriers and add an accessibility channel. | Automated and manual audit with remediation evidence |

## P2 — enterprise and international expansion

| # | Owner | Action | Completion evidence |
| --- | --- | --- | --- |
| 23 | Commercial counsel | Finalize enterprise MSA, order form, SLA, security exhibit, support terms, and procurement playbook. | Approved templates |
| 24 | Privacy counsel | Finalize the DPA, subprocessor schedule, international transfer mechanism, data-subject assistance, audit terms, and deletion certification. | Executable DPA package |
| 25 | International counsel | Assess GDPR, UK GDPR, EU database rights, text-and-data-mining rules, local publisher rights, privacy laws, and required representatives before targeting non-U.S. markets. | Country/region launch memo |
| 26 | Export/sanctions owner | Adopt a proportionate sanctions/export program for customers, payments, and restricted jurisdictions. | Policy, screening process, escalation records |
| 27 | Tax counsel/accounting | Determine sales/VAT/GST and marketplace obligations for API subscriptions and publisher royalties. | Nexus matrix and configured tax process |

## Publisher outreach sequence

1. Verify publisher ownership and the rights holder for each property/feed.
2. Send a concise use description and the completed rights schedule.
3. Ask the publisher to identify syndicated/media exclusions.
4. Negotiate search, summary/RAG, and full-display rights separately.
5. Confirm customer display, API distribution, retention, attribution, territory, termination, and derived-data rules.
6. Obtain signatures before enabling full text.
7. Enter the signed terms into the rights registry and test enforcement.
8. Calendar renewal, reporting, audit, and deletion dates.

## Counsel engagement brief

Ask counsel to deliver:

1. product copyright/access risk memorandum;
2. Master Publisher Content License and rights schedule;
3. API/customer Terms review;
4. privacy/state-data-broker/My Health My Data analysis;
5. DMCA safe-harbor and designated-agent advice;
6. crawler/source-Terms and cease-and-desist playbook;
7. corrections, defamation, publicity, and high-risk editorial rules;
8. enterprise MSA, DPA, and security exhibit;
9. subscription/marketing compliance review; and
10. insurance-coverage review.

Provide counsel with:

- this repository's `legal/` package;
- the supplied `tavily_news_content_distribution_legal_framework.md`;
- current source/feed inventory and source terms;
- sample API responses for every content-bearing route and format;
- current retention, logging, caching, deletion, and subprocessor facts;
- account, checkout, API-key, policy-link, and cancellation screenshots; and
- representative publisher complaints or permissions, if any.
