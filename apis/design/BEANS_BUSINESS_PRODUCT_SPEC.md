# Beans business and product specification: source packs for B2B briefing builders

| Field | Value |
| --- | --- |
| Decision date | 2026-10-30 |
| Status | Conditional discovery specification; commercial demand unvalidated |
| Decision | Test a narrow offer with the existing API before funding new retrieval features |
| Initial target | Small technical teams building English-language B2B technology briefings or cited assistants for a defined publisher list |
| Product boundary | Publisher information and coverage context; structured Events and interpreted Signals belong to Espresso |
| Authority | Internal business proposal, not a public contract, price list, source license, or implementation authorization |
| Owner | Product owner; technical and rights owners must be assigned before a pilot |

## 1. Decision and sufficient conjecture

**Decision:** pursue a bounded evaluation of **Beans Source Packs**, a proposed way to package existing article, source, and Story retrieval for small briefing-product builders. Do not build another general news API, market MCP as unique, or compete solely against NewsAPI.org's expensive entry plan. Low-cost APIs, feed readers, static JSON feeds, and self-hosted pipelines already serve much of this job.

**Conjecture:** for a developer or technical founder maintaining daily or weekly B2B technology briefings, repeatedly finding the same reporting, checking source freshness, and repairing feed ingestion may consume enough effort to justify paying for a small, explicitly scoped source pack. They currently use RSS scripts, a reader API, general news APIs, or manually curated links. We expect a source list accepted against their actual workflow, compact citable article/Story responses, and visible coverage limits to shorten integration and repeated source review. A proposed $49/month offer is worth testing only if it solves something their existing free Beans allowance and cheaper substitutes do not.

This is a workflow and packaging hypothesis, not established whitespace. No customer interviews, benchmark against a customer's reference set, live coverage audit, paid pilots, purchase commitments, or retention measurements were performed in this research. The evidence supports the next small experiment; it does not support a paid launch or a revenue forecast.

The immediate route is **GTM first using existing retrieval**, conditional on rights and coverage readiness. Minimum enhancements become justified only when observations identify a repeatable missing outcome. Source-specific rights enforcement is already an unresolved release dependency, not an optional differentiator.

## 2. Method, boundaries, and evidence standard

The supplied *Product and Business Idea Development Manual* determines the sequence: operating fit, observed behavior, segment alternatives, objective and source of volume, target economics, positioning, complete offer, price, channels, experiments, and evaluation. Its examples and numerical thresholds are methodology, not evidence. Its activities do not authorize outreach, purchases, legal changes, implementation, or publishing.

Statements below are labeled **Observed**, **Estimate**, **Assumption**, or **Decision**. External links were read through Exa on 2026-10-30. Local capability claims derive from inspected source and documentation, not production tests. The [evidence record](research/2026-10-30/beans/EVIDENCE_AND_RESEARCH.md) preserves research scope, source quality, contradictions, and next tests; raw captures are beside it.

Operating constraints are a software subscription with bounded support, no physical fulfillment, no new marketplace, no publisher-rights promises beyond documented permission, and no editorial service whose labor grows with every customer's daily brief. Do not promise exhaustive global coverage, trading latency, regulatory completeness, litigation docket access, fact verification from source counts, or unrestricted full text.

**Observed readiness limit:** the [internal legal readiness plan](../../legal/LEGAL_READINESS_PLAN.md) says all products remain free, paid billing and commercial agreements are not operational, and monetization stays disabled pending corporate prerequisites, rights, insurance, and assent controls. It also identifies unresolved fail-closed source-specific full-content enforcement. Metadata and summaries are not assumed automatically cleared either. This specification proposes future commercial tests, conditional on those gates. It grants no downstream rights.

## 3. Existing product and actual gaps

The code and current gateway/portal have already moved beyond older route proposals. The historical [route proposal](BEANS_API_ROUTE_PROPOSAL.md) and parts of the [market report](NEWS_AND_BLOG_API_MARKET_REPORT.md) describe offset/bare-array behavior and propagation-era assumptions; they are research inputs, not the current baseline.

| Capability | Observed current behavior and source | Business implication |
| --- | --- | --- |
| Article search | [`router/routes.go`](../beans/router/routes.go), [`router/params.go`](../beans/router/params.go), and [gateway OAS](../../config/beans.oas.json) expose topic relevance, exact IDs/URLs, dates and source/domain/byline/language/taxonomy filters | Already supports bounded source retrieval; do not propose basic search as new |
| Multiple source types | Fourteen request-filterable types include news, blog, official statement, earnings/financial reports, lawsuit, research paper, technical documentation and others; `post` is response-only | Can test specialist public-information workflows; accepted enum values do not prove useful corpus depth |
| Feeds | Latest, trending, news-only shortcuts and recent top-headlines routes | Chronology and attention are distinct; attention is not importance, truth, or a forecast |
| Detail and optional body | Article UUID detail, optional `full_content`; [`responses.go`](../beans/router/responses.go) returns publisher URL and available metadata | Citation handoff exists; full-content availability and downstream permission remain separate |
| Stories | Collection, detail, and paginated member Articles; a selected Article may link to a Story | Related coverage can reduce repeated reading; not proof that every duplicate is grouped or every member independently verifies a fact |
| Similar and mentions | Similar Articles and separate external Mention observations | Already present. Mentions are not a complete social listening product; similarity is not guaranteed Story membership |
| Discovery | Source collection/detail and categories/entities/regions/sentiments | A customer can inspect sources and accepted filter values instead of guessing |
| Pagination and transport | JSON `{data,pagination,meta}`, opaque cursors, page-only `num_results`, default 20/max 100; hosted MCP described in [current portal](../../docs/pages/products/beans/overview.mdx) | Agent integration is an existing capability; API-response `meta.as_of` is not publisher collection freshness |
| Current access | Free tier 100 requests/minute and 50,000 requests/month per authenticated user; Bearer except health | Small briefing workloads may fit entirely within free allowance, weakening a volume-based paid upsell |

Inspected query implementations in [`db/queries.go`](../beans/db/queries.go) include actual Article, Source, similarity, Mention and Story retrieval. Registration is not the only evidence. No live data, deployment parity, latency, source rights, classification accuracy, Story quality, or freshness guarantee was verified. The current code and guides also differ on when a nonzero `score_threshold` is required with `q`; evaluation examples should provide an explicit threshold until that behavior is verified. This research does not change the contract.

**Gap diagnosis:** market evidence points toward coverage confidence, content-use certainty, and workflow maintenance, not absence of Article/Story/MCP primitives. Current public Source responses lack a demonstrated per-source coverage/freshness service commitment. A Source record does not establish that it is actively collected. A filterable `lawsuit` or `earnings_report` type is not a guarantee of primary-document coverage. These are proof and delivery gaps to measure before declaring feature gaps.

## 4. Observed problem episodes and competing explanations

| Observation | What it supports | What it does not support |
| --- | --- | --- |
| Kapyn's operator describes a 2026-10-16 RSS/LLM loop in which 79 of 86 inserted stories were older than its seven-day window; deletion caused reprocessing and inflated insertion metrics ([operator account](https://kapyn.app/blog/llm-pipeline-re-summarising-deleted-articles)) | Repeated ingestion and misleading freshness metrics cause real operational waste in one small briefing product | Not a representative frequency, verified cost saving, or request to buy Beans; operator fixed it locally |
| Brian Baldock describes an agent consuming a daily digest as roughly 700 KB HTML; a per-story JSON feed and an explicit license solved the handoff ([2026-10-21 account](https://blog.brianbaldock.net/posts/aigregator-for-agents/)) | Compact story-level data, original-source links, and permissions can improve a real agent workflow | Counterevidence to needing a paid API or MCP: two static files were enough |
| HN's simonw asks whether search APIs allow users to “store and resyndicate results” ([2026-10-05 discussion comment](https://news.ycombinator.com/item?id=49965369)) | An adjacent technical buyer treats use permissions as an adoption criterion | Search APIs generally are not identical to news APIs; this is not permission advice or validated demand for Beans |
| In an RSS discussion, xyzsparetimexyz writes “1. Have a RSS feed” and “2. Don't get rid of it” ([2026-10-03 comment](https://news.ycombinator.com/item?id=49949882)) | Dependable access matters to some users | Consumer RSS frustration does not establish B2B budget or willingness to outsource |

The last30days engine returned 14 Reddit threads, 24 HN items and 24 GitHub items, but explicitly judged evidence thin. Most returned Reddit activity concerns consumer RSS readers, social-feed access or app promotion. Unrelated n8n support and repository issues were rejected for the business case. The number of returned items is not the number of independent relevant demand observations.

Alternative explanations remain strong: builders enjoy maintaining their own pipelines; RSS is sufficient; the bottleneck is editorial judgment rather than retrieval; consumers expect free news; rights uncertainty cannot be solved economically; cheap incumbents already provide enough. The selected segment must beat these explanations behaviorally.

## 5. Candidate applications and target selection

Scores are ordinal judgments, not measured market shares. H/M/L indicate relative attractiveness under current evidence.

| Application and buyer | Pain/frequency | Reach and ability to pay | Fit/competition | Decision |
| --- | --- | --- | --- | --- |
| English-language B2B technology briefing/cited-assistant builders with 20-100 named sources and daily/weekly delivery | Maintenance and repeated coverage observed; recurring | Public project authors can be identified; budget unproven | Strong current fit, very crowded including free specialist feeds | First evaluation cohort; require existing repeated workflow and a technical budget holder |
| General consumer news apps/indie aggregators | Frequent, but often hobby use | Reachable, weak payment incentive | Very high commodity competition; free Beans may suffice | Free acquisition channel, not initial paid target |
| Healthcare/insurance, policy or litigation intelligence vendors | Potentially high-value missing coverage | Buyer access and rights/completeness harder | Structured impact and policy judgments overlap Espresso; no corpus audit | Park until representative workflows and source readiness are established |
| Financial trading and real-time market systems | High urgency | Potential spend, demanding procurement | finlight and enterprise feeds strong; no measured latency/SLA | Exclude |
| Large commercial news redistribution or training corpora | High-volume and licensing-heavy | Enterprise sales and legal costs | Current rights readiness fails | Exclude |

**Initial cohort eligibility:** one to five builders, an operating daily/weekly B2B technology brief or cited assistant, an agreed list of 20-100 English sources, a named maintainer/payer, willingness to supply a reference sample and current workaround, and no requirement for full-text republication, proprietary feeds, guaranteed completeness, or real-time trading. Firm size is a recruiting convenience; the defining need is an existing repeated source-selection workflow.

**User:** developer/information-product maintainer. **Buyer and budget holder:** founder or product lead. **Beneficiary:** the briefing reader or assistant user. **Approver/blocker:** rights/security reviewer or the same founder in a small business. Time saved is not treated as cash saved until the buyer explains what work it displaces.

Technology is the first test context because concrete operator accounts were found, not because it is less competitive. If a qualified cohort cannot be recruited or better served than its static feeds, stop this niche; do not silently relabel the target as all B2B industries.

## 6. Alternatives, pricing, and the actual competitive standard

Observed public prices below were fetched on 2026-10-30, monthly billing unless stated. They are technical-access prices and do not establish identical coverage or downstream rights. Competitors' uptime and quality claims were not benchmarked.

| Alternative | Current public offer | Strength and limitation for this job |
| --- | --- | --- |
| [TheNewsAPI](https://www.thenewsapi.com/pricing) | Free 100/day × 3 articles; Basic $19/month, 2,500/day × 25 articles | Direct low-price disconfirmation. General news metadata/search; full-text and derivative-use rights must be separately verified |
| [Currents](https://currentsapi.services/en/product/price) | Free 250/day × 20; Builder $69/month, 2,500/day and 75,000/month × 50 | Free commercial attribution/link-out use is explicitly offered. [FAQ](https://currentsapi.services/en/faq) excludes end-user AI summaries, permanent stores and republication by default; scope those uses separately |
| [finlight](https://finlight.me/pricing) | $29/10,000 requests with real-time data, ticker/entities and one webhook; $99/50,000 with sentiment/history | Strong financial vertical alternative, already has push. Pro Light history is inconsistent between feature list and comparison table; do not promise archive scope without confirmation |
| [NewsAPI.org](https://newsapi.org/pricing) | $449/250,000 requests for production; free 100/day is dev-only, 24h delay | Well-known incumbent, but an incomplete price benchmark. Full article content unavailable even on paid plans |
| [NewsAPI.ai](https://www.newsapi.ai/plans) | Visible 5K plan $90/month; recent Article search 1 token, recent Event search 5; archive costs by searched year | Article bodies, duplicate detection and event clustering already available. Tokens are not interchangeable with requests; full tier selection and rights require confirmation |
| [Perigon Data APIs](https://perigon.io/products/pricing/apis) | Eligible-startup Basic $250/10,000; Plus $550/50,000; full commercial API plans from $24,000/year | Semantic search, Stories, source discovery, MCP and custom sources already exist. Basic 1h delay; startup license eligibility and endpoints constrain comparability |
| [Perigon custom APIs + Monitors](https://perigon.io/products/pricing/custom) | Individual $25/1,000 personal; startup $250/10,000; commercial from $1,250/25,000 | Different customizable product ladder, includes monitors. Do not select the $24k API floor as Perigon's universal cheapest commercial price; written quote needed to reconcile use case and entitlement |
| [Newsflash](https://newsflash.sh/docs) | Free 1,000/day; $29/month premium 50,000/day, 5y archive, REST/MCP/SSE, deduped events/articles | Direct contrary offer to cheap agent-ready Stories. Vendor claims, not adoption/coverage proof. Its source-count-derived confidence is not demonstrated independent verification |
| [Inoreader](https://www.inoreader.com/pricing/feature/api_access) | Pro $9.99 monthly or $7.50/month annual; reader automation and limited private API access | Excellent cheap substitute for a small known feed list. [Developer portal](https://www.inoreader.com/developers/) requires scoping API use/quotes beyond eligible applications; reader subscription is not unrestricted resale access |
| DIY RSS/static JSON/manual curation | No mandatory data-provider subscription; engineering and editorial effort vary | Strongest substitute for 5-20 sources and daily batches. AIgregator demonstrates that compact licensed discovery can be served statically; no invented $/hour saving is assumed |

Perigon's [vector search documentation](https://perigon.io/docs/api/vector-search) already describes metadata filters and reprint exclusion; its [source documentation](https://perigon.io/docs/api/sources) supports source groups and custom source requests. Therefore semantic filtering, clustering, source selection, and MCP are capabilities to prove economically, not unique inventions.

### Comparable workload rather than headline request quotas

**Estimate A: bounded editorial evaluation.** Ten topics, two polls/day, one collection page/poll: 10 × 2 × 30 = 600 collection calls/month. Two selected Article details/poll add 1,200, and 200 Source/Story/member calls give **2,000 calls/month**. At 20 results/page that is at most 12,000 collection rows before repeat retrieval and overlap; it is not 12,000 unique or useful Articles. Ten hourly topics instead produce 7,200 collection calls/month before details. A 50,000 free allowance may already cover both.

**Estimate B: 10,000 calls/month**, distributed at approximately 334/day, for a product with a larger source/topic set. Effective subscription cost per 1,000 actual calls is TheNewsAPI $1.90; Currents Builder $6.90 (free daily cap insufficient at even load); finlight Pro Light $2.90; NewsAPI.org Business $44.90; eligible Perigon Basic $25.00 with 1h delay. Newsflash's advertised free 1,000/day tier already covers this request count with 30-day history, if the source slice and rights fit; its $29 Premium/$2.90 per 1,000 actual calls is an optional archive/higher-usage tier, not a necessary cost for this workload. NewsAPI.ai's visible 5K/$90 plan covers only 5,000 simple recent Article searches; 10,000 tokens require a confirmed larger tier, so no invented $180 quotation is used. At 2,000 calls the same plan floors dominate, and several free offers may suffice.

This is an access-cost comparison, not an equal-rights/equal-recall ranking. Each page, detail, empty search, and repeat poll can consume another request. Some providers return 25 versus 50 versus 100 rows, but an upper result limit is not useful yield. Record **total spend per 1,000 unique, relevant, permitted items** during the pilot, including detail follow-ups, discarded duplicates, retrieval/LLM costs and required licensing. If a tier's source slice or rights do not qualify, mark it ineligible rather than comparing price alone.

## 7. Positioning, objective, and source of volume

**Positioning hypothesis:** for small teams delivering a cited B2B technology brief from an explicit publisher list, Beans provides a compact source-retrieval backend with inspectable coverage and related reporting, reducing feed-maintenance and repeat-reading work. The reason to believe must be a dated evaluation on the team's source list, not the number of sources in a catalog or a claimed AI advantage.

**Price role:** affordable enough to test and predict, with a small premium only for a demonstrated complete offer. A $49 offer is more expensive than several substitutes. It must earn that difference with useful coverage and delivery proof. “Cheaper than NewsAPI.org” alone is rejected.

**Category:** outsourced source discovery and retrieval for recurring B2B briefings. The primary source of volume is **earn share** from an existing paid feed/news API or from a workflow already consuming engineering time. Converting a team doing no repeated briefing today is demand creation and should be tracked as a different experiment. Free hobby signups are not counted as paid share capture.

**Assumed objective:** decide within 90 days whether this can support a founder-operated subscription; aim for 10 qualified paying accounts after readiness gates, not a large TAM claim. An illustrative 12-month operating goal of 50 accounts × $49 = **$2,450 MRR / $29,400 annual run rate** is a planning scenario, not market evidence, revenue achieved, or revenue recognized during the acquisition year. If the owner's required income is materially higher, this offer needs expansion economics or should be parked.

## 8. Smallest complete offer and requirements

Journey: a builder supplies its current source list and sample brief; checks a sample of existing Beans results; sees missing/late sources and allowed-use limits; completes one cited brief; integrates existing REST/MCP; repeats against a reference set; evaluates a disclosed future price after commercial gates; can export its own query configuration and stop. API delivery does not transfer publisher ownership or guarantee a right to retain the results.

| Priority / requirement | Existing versus proposed | Acceptance and customer outcome |
| --- | --- | --- |
| P0: scoped coverage assessment | Proposed delivery artifact using existing search/source routes, initially manual | Every requested domain is labeled observed-active, not observed in sample, unsupported, or unverified; dated method/window; no Source-catalog-only success claims |
| P0: allowed-use scope | Existing warnings are insufficient; readiness dependency | Rights owner approves collection, processing, output fields, retention and customer use; unknown rights fail closed; no paid distribution before gates |
| P0: one useful cited brief | Existing search, detail, Stories, Sources | Five qualified teams complete a real task; retain original links and distinguish related coverage from corroboration; timestamp interpretation correct |
| P0: measured coverage/quality | New evaluation, not a new route | Buyer supplies 100 reference items across at least 20 active sources and 14 days; target ≥85% recall within agreed daily briefing deadline, ≥90% top-20 relevance, ≤10% duplicate read effort, and ≥95% working links; report denominators and source exclusions |
| P0: cost/reliability characterization | New measurement using existing service | Measure 14-day response success and latency, per-source publication-to-availability, useful yield and support time; proposed test threshold ≥99% successful scheduled jobs after bounded retries, no contractual SLA implied |
| P1: source-pack setup recipe | GTM/prototype artifact; current filter routes suffice | A maintainer completes first filtered retrieval and a cited output in ≤60 minutes without custom per-account code; record assistance |
| P1: coverage/freshness disclosure | Initially dated CSV/report, not claimed present API field | Customer can tell collection time from publication time and API response time. If ≥3 teams need programmatic monitoring, scope a later public coverage/status interface with explicit semantics |
| P2: incremental update/replay control | Not an existing durable change feed; deferred | Only design if repeat polling/late items demonstrably create material misses or cost; test bounded overlap and client ID tracking first |
| P2: push/webhooks, SDKs, entity resolution or new enrichment | Proposed only; deferred | Do not build for parity. Require ≥3 qualified teams blocked on same outcome and evidence existing alternatives are inadequate |

Thresholds are **decisions for a bounded experiment**, not established benchmarks or guaranteed performance. A 100-item reference set can be biased and only supports the agreed slice. Measure disagreements with two reviewers on 30 overlap/relatedness cases; report Story false merges and misses separately. Do not count syndicated copies as independent confirmation. Small-source or sparse-day samples can be inconclusive.

**No enhancement branch:** if source rights and the existing endpoints pass the reference evaluation and onboarding test, launch a narrow GTM experiment with current capabilities; do not add a route. The sellable value is the complete scoped offer, its proof, and fit. If teams find free Beans sufficient and reject payment, retain a free developer acquisition product or park monetization. Do not degrade the frozen free allowance to force an upgrade.

**Enhancement branch:** if missing active sources cause failure, investigate permissible source acquisition and its cost before API work. If lack of freshness visibility blocks ≥3 qualified teams, propose a separately reviewed coverage/freshness interface. If rights cannot be granted within viable cost, stop that use case. Every later public API change must follow repository contract/docs requirements; this document implements none.

## 9. Price, entitlement, and unit economics hypotheses

**Current fact:** all products are free; 100/min and 50,000/month remains the public free policy. There is no operational paid plan, paid support, SLA, or commercial content license.

**Proposed offer to test after readiness:** $49/month for one production application/source pack, up to 100,000 product requests/month, a dated coverage report and self-service integration recipe, capped support, no overage by default, and a hard stop/explicit upgrade choice. Billing scope across Beans and Espresso must be settled before quoting: the current free allowance is per user, and this model must not silently turn it into separate per-product free quotas. Trial/entitlement language, cancellation, taxes, refunds, payment failure and customer retention permissions need approved commercial terms. Higher request allowance alone is unlikely to sell at the observed workload.

Do not promise adding arbitrary sources on demand in $49. Scope any new source list separately before acceptance; if rights or maintenance costs exceed the unit model, decline the source or revise the target. No training/full-text resale rights included by implication. Before readiness, only a nonbinding price/offer reaction can be collected; it is weaker than purchase evidence.

**Estimate, USD per account per month:** payment fee modeled at 3% + $0.30 = $1.77; founder support valued at $50/hour; infrastructure, collection, extraction and rights allocations below are unmeasured assumptions. Fixed shared ingestion should be allocated transparently across a conservative paid base, not disappear from margin.

| Scenario | Allocated hosting/collection/query cost | Rights/data allocation | Support | Payment | Contribution before fixed costs |
| --- | --- | --- | --- | --- | --- |
| Low cost | $4 | $0 only if expressly cleared | 0.10h = $5 | $1.77 | $38.23 / 78.0% |
| Base | $8 | $5 | 0.25h = $12.50 | $1.77 | $21.73 / 44.3% |
| High cost | $18 | $15 | 0.50h = $25 | $1.77 | -$10.77 / -22.0% |

At assumed shared fixed operating cost $300/month, break-even is 8/14/not attainable accounts in the low/base/high cases, before founder salary, sales, legal costs, taxes and development. At base contribution, six-month acquisition payback caps fully loaded CAC at about $130. Founder time at $50/hour permits roughly 2.6 hours acquisition/onboarding per account within that cap. These constraints may make a $49 source-pack service unattractive; measure them before selling custom attention as software.

Fifty base-case accounts yield $1,086.50 monthly contribution before fixed costs, far below $2,450 gross MRR. Gross revenue is not founder income. The most important sensitivities are rights costs, shared ingestion allocation at low customer count, and support. Pricing tests should compare a $29 basic-access concept, the $49 scoped offer, and a $99 higher-service concept qualitatively with matched buyers; do not call a small preference exercise an A/B conversion result.

## 10. Attainable market and route to market

No verified registry or deduplicated directory of eligible paying accounts was established. Do not substitute a global news-API market report. The initial **serviceable population** is restricted to teams with named English sources, permissible use, no enterprise SLA, and repeated B2B briefing demand. Count accounts, not developers, repositories, subscribers, or feeds.

**Recruitment model, assumptions for 90 days:** 120 individually qualified accounts identified, 30 agree to a relevant discussion (25%), 15 share a source/reference workflow (50%), 10 finish evaluation (67%), and 4 accept and pay $49 when payment is permissible (40%). This yields $196 MRR, not $2,450. Low: 60 × 20% × 50% × 50% × 25% = 0.75 expected accounts, rounded to a 0-1 planning outcome. High: 200 × 30% × 60% × 70% × 50% = 12.6, or 12-13 accounts. All conversion inputs are assumptions, no outcomes achieved. If legal gates remain closed, replace payment with recorded conditional offers and keep the commercial conclusion inconclusive.

The 50-account twelve-month scenario needs a materially larger reachable prospect pool, referrals, stronger conversion, expansion or a partner channel. At the base funnel's 3.33% end-to-end conversion it needs around 1,500 eligible accounts. That population and acquisition capacity have not been established, so the objective is not yet supported.

| Stage | Channel/asset and requested behavior | Proof or blocker |
| --- | --- | --- |
| Discovery | A public worked example: one named source list, duplicate coverage and missing-source disclosures; targeted documentation on briefing retrieval | Source-specific results and rights must be accurate; generic “AI news API” SEO is saturated |
| Evaluation | Invite technical maintainers of existing briefing projects to supply last week's source list and output | Need real recurrent job and payer, not interest in a free tool |
| Approval | Explain allowable output, retention, limitations and proposed full price before integration | Rights/security and commercial terms may block; founder cannot approve publisher rights by assertion |
| Purchase | Self-service subscription only after corporate and rights readiness | Billing does not exist today; evaluate willingness conditionally meanwhile |
| Implementation | Existing REST/MCP query recipe, explicit cursor handling, source-filtered sample and selected details | ≤60m assisted-to-unassisted onboarding target; source gaps exposed early |
| Retention | Four scheduled weekly uses, comparison against prior method, then renewal opportunity | Repeated useful retrieval and payment, not call volume alone |

Outreach and publication are proposed next activities, not actions taken or authorized by this document. Respect community promotion rules: the current r/rss research includes complaints about repeated app promotion, so a link-dump channel is a poor default. Start with known project maintainers and narrowly useful examples only when outreach is separately authorized. An integration partner must earn reduced customer setup/support work; no partner has committed to distribution.

## 11. Experiments, launch gates, and stop decisions

Proposed caps exclude remediation of unresolved legal readiness; that cost must be scoped independently. Product owner records approached/qualified/participated/completed/paid counts and reasons for every rejection.

| Timing and uncertainty | Experiment, participants and cap | Pass / fail / inconclusive decision |
| --- | --- | --- |
| Days 1-15: real problem and reachable payer | Observe 8 qualified maintainers reconstruct their last two briefs, inspect source list and current tools; ≤12 founder hours/$600 imputed | Proceed if ≥5 show repeated relevant retrieval/maintenance friction and ≥3 identify a budget holder/current spend or measurable opportunity cost. Otherwise revise/park target; compliments do not pass |
| Days 1-30: rights and feasible source slice | Rights owner and technical owner audit every requested source/use in one 20-source slice; measure existing outputs for 14 days; ≤16 technical hours/$800, external rights work separately scoped | Proceed only for expressly cleared use and stated acceptance thresholds. Rights unresolved = no commercial launch; source deficits without viable acquisition = stop this slice |
| Days 16-45: minimum complete workflow | 5 eligible teams complete a real brief with current API and their prior method; paired task baseline and blind relevance review; ≤15 support hours/$750 plus metered costs | ≥4 finish, median task time ≥30% lower, ≥3 use it again without founder assistance. If task-time saving is absent but quality improves, require buyer evidence quality is worth paying for; otherwise fail |
| Days 31-60: price and delivery | Present the same scoped $49 complete offer to 10 qualified evaluators; ≤$150 cash variable cost and ≤2.6 acquisition/onboarding hours per acquired account | After readiness, ≥4 pay undiscounted for a month and base contribution ≥$20/account. Before readiness, nonbinding acceptance is only a discovery signal; payment hypothesis remains unresolved |
| Days 61-90: recurrence and channel | Observe 4-week use and the first eligible renewal; recruit a second cohort through the same channel | ≥3 of first 4 paid accounts renew without custom work, ≥3 complete weekly jobs without rescue, support ≤15m/account/month, measured contribution positive. Fewer than 4 buyers is insufficient to extrapolate retention |

**Stop or park:** source rights cannot support the customer use; 2 independently recruited cohorts prefer their free/cheap substitute after a fair comparison; most use cases fit free Beans and nobody pays for the added offer; source additions require recurring editorial labor; measured high-cost economics persist; the reachable qualified account list cannot support the income goal. Do not turn missing demand into an uncontrolled feature roadmap.

**Expand only after** recurring value, renewal and margin are observed and the next segment has its own rights/source evaluation. A healthcare, finance or policy label is not an expansion plan by itself.

## 12. Decision record and next review

**Decision:** conditional pursue, narrowly, with existing capabilities. The most plausible opportunity is a repeatable source-pack experience with proof, not a novel general retrieval capability. Low-price competitors and static feeds substantially weaken the original commodity pricing thesis. Existing Beans may need only GTM and evaluation assets for the initial task; rights readiness remains mandatory, and measured source/freshness deficiencies may justify later enhancements.

**Next review:** after 8 workflow observations and the first cleared 14-day coverage evaluation, or 2026-11-30, whichever occurs first. Product owner must then choose continue / revise target / GTM-only / enhancement proposal / park, recording evidence, counterevidence and costs. No implementation, billing, outreach, public pricing, or source-rights change is approved by this specification.
