# Espresso business and product specification

| Field | Value |
| --- | --- |
| Decision date | 2026-10-30 |
| Status | Internal research specification; conditional opportunity, not a validated business or shipped roadmap |
| Proposed offer | Evidence-qualified market Signals for software that produces recurring company and sector briefs |
| Initial target | Small B2B competitive-intelligence software teams producing weekly, human-reviewed briefs from named company watchlists |
| First decision | Test the existing API and sell the workflow only after demonstrating coverage and review-time value; do not build a broad intelligence platform |
| Evidence | [Dated register](research/2026-10-30/espresso/EVIDENCE_REGISTER.md), preserved Exa and last30days outputs |
| Method | Attached *Product and Business Idea Development Manual*: boundaries → behavior → applications → alternatives → objective/source of volume → target/economics → position → offer → price → distribution → experiments → evaluation |

## 1. Business decision and its limits

**Pursue a bounded evaluation of an evidence-review job, not a general cheap intelligence API launch.** Espresso already distinguishes concrete Events from synthesized Signals and exposes supporting Events and evidence. The most plausible near-term product is a reusable input to someone else's weekly market brief: a proposed implication, its concrete support, and visible reasons to accept, qualify or reject it. The paying outcome would be less human time checking whether a synthesized statement is defensible, rather than more articles retrieved.

The market does not lack affordable event feeds, citations, MCP, monitoring or policy APIs. [Newsflash](https://newsflash.sh/docs) advertises $29/month event feeds; [Parallel](https://docs.parallel.ai/getting-started/pricing) offers monitoring from $3/1,000 checks; [FiscalNote](https://fiscalnote.com/policynote-api-pricing) now advertises a $300/month self-serve PolicyNote API. [Seerly's October engineering account](https://seerly.app/resources/engineering/building-a-content-intelligence-loop) describes an evidence-backed weekly competitor product. These are direct reasons to reject generic whitespace and uniqueness claims.

**Confidence:** high that Espresso's current route composition can support an evaluation; moderate that dates, evidence qualification and coverage visibility matter to builders; low that the proposed target will switch or pay. No interviews, observed customer workflow sessions, product trials, paid pilots or willingness-to-pay tests occurred. Recent social research yielded little independent buyer evidence. This specification justifies the next experiment, not a substantial development commitment.

**Commercial state:** all Cafecito products currently remain free. Existing internal [legal readiness requirements](../../legal/LEGAL_READINESS_PLAN.md) record monetization prerequisites and unresolved source-specific rights controls. Proposed paid offers below are hypotheses for later testing, conditional on completing those prerequisites. A public URL, a summary, a model output or a free API key is not proof of permission for downstream client display.

## 2. Operating fit and business intent

Fit: digitally delivered API/MCP subscription; no physical fulfillment, two-sided market, travel or bespoke advisory practice. Reuse current retrieval and Event/Signal/evidence relationships. A customer continues owning its watchlist, editorial review, delivery interface and decisions. Espresso assumes bounded data-normalization and evidence-quality responsibilities only after they are proven deliverable.

Proposed planning objective (not confirmed by the owner): **10 retained paying accounts at $149/month within six months of commercial readiness**, or $1,490 MRR, with at least 60% contribution after variable delivery and support. This is a learning-stage objective, not a founder-income forecast. At that margin it generates only $894/month before fixed costs, acquisition and founder compensation. A sustainable-income goal would require more accounts, greater demonstrated value per account, or a different model.

Constraints assumed for planning: at most 60 founder hours and $750 incremental cash for the first 30-day discovery/benchmark phase; at most 30 minutes recurring support per production account per month; no paid data source without explicit rights and a cost model. Resource availability has not been confirmed. If these assumptions are unrealistic, reduce the experiment before building.

Disqualifiers: needs for guaranteed all-market coverage; trading-grade latency or predictive returns; legal applicability decisions; unlicensed redistributable full text; broad multi-jurisdiction obligation tracking; unlimited human analysis; or a buyer who wants only a finished dashboard rather than an API.

## 3. Current capability baseline, verified in the repository

Authority is the current runtime plus gateway specification, not old design proposals. No live service/data-quality audit was performed. Backend paths below acquire `/espresso` at the public gateway.

| Customer job | Implemented surface | What is established | What is not established |
| --- | --- | --- | --- |
| Retrieve developments | `GET /events`, `GET /events/{id}` | Event search/detail; semantic `q`, normalized company/person/product/category/region/event-type filters, tags, source filters and date bounds | A complete sector corpus, occurrence-date coverage or verified facts |
| Retrieve interpretation | `GET /signals`, `GET /signals/{id}` | Separate Signal search/detail; optional impact, confidence, outlook and other upstream fields can be projected | Calibrated confidence, prediction accuracy or a fixed complete interpretation schema |
| Inspect support | `GET /signals/{id}/events`; `GET /events/{id}/evidence`; `GET /events/{id}/signals` | Linked support/context and related interpretations; detail links and relationship counts | Claim-level entailment, source independence, primary-source preference, universal support completeness |
| Resolve filters/provenance | `/sources`, `/sources/{id}`, `/entities`, `/regions`, `/event-types`, `/tags` | Source metadata and discoverable filter vocabulary | Canonical company IDs, complete company aliases, authoritative jurisdiction mapping or source-health status |
| Integrate with agents | Gateway `/espresso/mcp`; JSON, YAML, TOON | Hosted MCP mapping and alternate representations of the logical payload | Measured token/cost advantage, end-to-end agent accuracy or production reliability |
| Read incrementally | Opaque cursor pagination; `limit` default 20, max 100 | Collection `{data,pagination,meta}`; detail `{data}`; successful empty lists | A change-data feed, revisions, deletions, webhook delivery or exactly-once monitoring |

Inspected sources: [`router/routes.go`](../espresso/router/routes.go), [`router/params.go`](../espresso/router/params.go), [`router/responses.go`](../espresso/router/responses.go), [`db/types.go`](../espresso/db/types.go), [`db/queries.go`](../espresso/db/queries.go), [`config/espresso.oas.json`](../../config/espresso.oas.json), [portal overview](../../docs/pages/products/espresso/overview.mdx), [workflows](../../docs/pages/products/espresso/workflows.mdx), [migration](../../docs/pages/products/espresso/migration.mdx), and [MCP guide](../../docs/pages/guides/mcp-ai-agents.mdx).

Important boundaries:

- `from` and `to` currently bound record `created_at`. They do not mean publication time, occurrence time, a comment deadline, a rule's effective date or a forecast horizon. The current request binder accepts date-only bounds.
- `meta.as_of` is generated when the response is assembled. It is not evidence that every upstream source was checked then, and a cached response can retain that earlier value.
- Stable record fields are `id`, `kind`, `created_at`, `tags`. Conditional/extension fields are not a guarantee that every record has source, confidence, forecast, impact or summary values.
- An evidence relationship provides related records. It does not prove that the source supports every phrase in the Signal or that three publishers are independent witnesses.
- Discovery returns filter vocabulary, not a canonical entity registry. `companies=apple` by itself is not a tested entity-disambiguation product.
- No public confidence endpoint, watchlist manager, alert webhook or point-in-time correction feed is part of the current public contract. Internal preview operations are excluded from this product promise.
- Current published free policy remains **100 requests/minute and 50,000 requests/month per authenticated user**, one key across Cafecito products. Public REST uses Bearer authentication except health. Shared usage matters when a client also uses Beans.

The superseded route proposal's old data counts cannot establish present coverage. Current migration-page statements about synthesized Signals and token efficiency are claims to benchmark; they are not comparative proof.

## 4. Problem hypothesis and candidate applications

Conjecture: **For a small CI software team preparing weekly briefs for named company watchlists, assembling source-backed implications from company announcements and industry/policy coverage creates repeated review work. The team currently combines news/search APIs, model synthesis and manual checks. Espresso may reduce that review work by serving interpretations with inspectable supporting developments and explicit uncertainty, motivating weekly reuse and eventually payment for qualified evidence.**

Roles: user is the engineer maintaining collection/synthesis and the analyst reviewing claims; beneficiary is the recipient of a reliable brief; buyer is founder, product lead or research lead; approver may be the client's data-rights/security owner; budget competes with search APIs, analyst time and existing monitoring subscriptions. The analyst's saved time becomes monetary value only if it reduces a real cost or frees capacity the buyer values.

Observed support is narrow. [Seerly](https://seerly.app/resources/engineering/building-a-content-intelligence-loop) documents real engineering decisions around publication dates, grouping and incomplete coverage. It also shows that a capable substitute already exists. A [Valyu tutorial](https://dev.to/valyuai/how-to-build-a-competitive-intelligence-agent-for-public-companies-op6) and a reader's peer-sampling critique support the importance of defensible source selection, but promotional tutorials are not customer research. A current [Reddit post title](https://www.reddit.com/r/SaaS/comments/1xgi83x/building_an_ai_competitor_research_tool_taught_me/) describes a reliability issue; its unavailable body prevents a stronger claim. Historical [medical-device monitoring noise](https://github.com/RASAAS/docmcp-knowledge/issues/13) is an analogy, not evidence for this segment's demand.

| Application/segment | Recurring job and present workaround | Fit and competitive pressure | Decision |
| --- | --- | --- | --- |
| B2B CI-tool builders with weekly named-company briefs | Turn developments into defensible implications; API/search + model + analyst review | Closest current Event→Signal→evidence fit; reachable technical trial; strong low-cost substitutes | First experiment, confidence low on purchase |
| Boutique research consultancies | Deliver recurring client market notes; reading and spreadsheets | User pain plausible; many want a finished service, not developer integration | Secondary only if technically capable and repeatable |
| Regulatory/compliance platforms | Track exact obligations, lifecycle and affected jurisdictions | Higher stakes; lacks stable effective dates/jurisdiction/lifecycle; FiscalNote and official APIs strong | Exclude from first release |
| Quant/retail trading systems | React to market news, forecast movements | Latency, licensing, backtesting and calibrated accuracy unproved | Exclude |
| Demand forecasting/local events | Predict store/site demand from scheduled events | Wrong data/date/geography requirements; PredictHQ-type workflows distinct | Exclude |
| Generic agent news retrieval | Answer what's new with citations | Beans/search APIs already address retrieval; low differentiation | Use Beans; do not duplicate |

These rankings are qualitative judgments based on fit and evidence, not measured demand scores. **No industrial vertical is selected yet:** this research did not establish a particular sector's buyer pain, source coverage or rights. First recruitment should target teams already briefing 20-50 named technology companies in English, then select a vertical only if their real watchlists and available evidence converge. This is a recruitable workflow boundary, not a claim of a proven technology-sector opportunity.

## 5. Alternatives, pricing and the actual gap

Published offers checked 2026-10-30, USD unless specified. Prices are vendor statements, not purchases or validated delivery quality.

| Alternative | Published entry/units | Strength for this job | Scope/rights caveat |
| --- | --- | --- | --- |
| [Newsflash](https://newsflash.sh/docs) | Premium $29/month; 50,000 requests/day; five-year history advertised | Deduplicated Events, underlying articles, sources, MCP/SSE | Not a claim-level implication audit; source-count confidence is heuristic; customer display rights require checking |
| [Parallel](https://docs.parallel.ai/getting-started/pricing) | Monitor lite $3/1,000 checks, base $10/1,000; Extract $1/1,000 URLs; Task core $25/1,000 successful runs | Custom-source monitoring, deduplication, structured cited research; substantial build-versus-buy threat | Different steps bill separately; terms/source-use permission and quality need evaluation |
| [GDELT Cloud](https://gdeltcloud.com/pricing) | Builder $59/1,500 QU; Watch $149/3,000 QU; Analyst $399/10,000 QU | Structured Events/Stories, REST/MCP, included scheduled monitors | Ordinary reads 1 QU; scheduled monitor checks use 0 QU; listed display rights internal/prototypes. Distinct from free GDELT Project |
| [FiscalNote PolicyNote API](https://fiscalnote.com/policynote-api-pricing) | Free 4,000 returned results/month; Self-Serve $300/20,000 results/month; month-to-month | Broad legislative/regulatory/local data and MCP; established policy specialization | Each returned result bills; 20-result search costs 20 results, not one unit; dataset/use rights still require review |
| [Perigon API packages](https://perigon.io/products/pricing/apis) | Startup Basic $250/10k requests; Plus $550/50k; commercial from $24k/year | Articles/Stories, structured enrichment, search, MCP | Startup eligibility and limited licenses; Basic article scope and delay differ from Plus |
| [Perigon configurable offers](https://perigon.io/products/pricing/custom) | Individual from $25, Startup $250, Commercial $1,250/month; commercial example includes 25k requests/10 monitors | Bundled monitors/API and commercial scope | Different advertised packaging from API page; obtain scoped quote, do not assume identical entitlements |
| [Direct Federal Register API](https://www.federalregister.gov/reader-aids/developer-resources/rest-api) | Keyless JSON/CSV API | Official US federal record retrieval | Customer must add relevance, lifecycle interpretation, integrations and review |
| Search API + model + spreadsheet/Notion | Actual cost varies; no single comparable list price | Fully customizable, already familiar; direct access to company's own pages | Customer maintains schema, dates, entity matching and verification |
| Existing monitoring/briefing products | [Seerly](https://seerly.app/resources/engineering/building-a-content-intelligence-loop) already describes weekly evidence-backed interpretation | A finished workflow can be better than an API for nondevelopers | Do not infer API price or license from an engineering article |

**Normalized workload example, estimate, not a feature-equivalence benchmark:** 25 company monitors checked daily for 30 days = 750 checks. Assume 200 changed pages/month, 50 structured research runs and 250 human-reviewed Signal candidates. Parallel Monitor lite costs `750 × $3/1,000 = $2.25`; base costs $7.50. Adding 200 Extract URLs = $0.20 and 50 core Task runs = $1.25, giving **$3.70-$8.95/month in these API line items** before integration, model work outside those runs, rights, storage and review. This is a credible low-cost substitute, not an expensive enterprise straw man.

For the same scale, 750 Espresso collection calls + 250 Signal details + 250 supporting-Event calls + 500 Event/evidence follow-ups = **1,750 HTTP calls/month**, assuming one page per step. Pagination, retries and discovery add calls. This fits today's free allowance unless shared usage elsewhere exhausts it. A $149 paid hypothesis cannot be justified by request access that is currently free. It must earn payment through demonstrated evidence qualification, source coverage, support/reliability terms or another genuinely additional complete offer.

GDELT Builder's scheduled-monitors allowance covers only three daily monitors; Watch covers ten. Analyst's 30 monitors can accommodate this 25-target scenario within its stated target limits, at $399/month, with a separate QU pool for ad hoc reads. Newsflash's $29 subscription accommodates 1,750 monthly reads well below its daily cap, subject to actual data/rights fit. FiscalNote at 1,750 requests returning 20 results each would consume 35,000 results, beyond the self-serve allowance; 1,750 total returned results fit its free tier. Counting billing units correctly changes the comparison.

**Gap classification:** a potential trust/proof and fragmentation gap: eliminate part of the customer's repeated evidence audit across synthesized implications. Price and developer access are packaging tools, not defensible gaps on their own. The gap may persist because entailment, stale sources, syndicated coverage, ambiguous dates and rights are costly to solve. It may also be adequately solved already by a customer's current stack. The next test must distinguish those explanations.

## 6. Target, category and source of volume

Category: **software inputs for recurring, human-reviewed competitive market briefs**. Initial target must satisfy all of:

1. Already produces or is shipping weekly company/sector briefs for a fixed 20-50-company watchlist.
2. Has an engineer who can integrate REST/MCP and an identified human reviewer.
3. Can show a recent real brief, sources and time spent checking implications.
4. Uses source-permitted English information; accepts bounded US-first evaluation and missing-coverage disclosure.
5. Buyer can redirect an existing data/research-tool budget or name a specific delivery cost reduced.

Explicit exclusions: users merely exploring agents, passive newsletter readers, regulated legal/compliance determinations, people seeking investment recommendations, and teams requiring global source completeness or external redistribution before source-use terms are resolved.

Primary source of volume is **earn share** from existing search/API plus synthesis and review workflows within that category. Replacing a spreadsheet can count as share capture in this broad workflow category; a user who previously produced no brief is demand creation and must be modeled separately. Existing free Cafecito accounts are not assumed to convert, and Beans customers are not counted as a second independent demand pool without evidence.

Behavior chain: qualified target supplies watchlist and prior brief → integrates one evaluation → creates a supported interpretation faster → repeats for four weekly cycles → buyer accepts a scoped paid offer after readiness → renews. Measure the first unsupported transition; API calls and signups do not substitute for retained usefulness.

## 7. Complete product and positioning

Positioning hypothesis: **For CI-tool builders who must check the support behind weekly market interpretations, Espresso provides structured Signal candidates linked to concrete Events and inspectable evidence, helping reviewers decide what is defensible before it reaches a client brief.** Proof must be a same-watchlist, same-week comparison of review effort and retained factual support, not a claim that synthesis predicts markets.

Product distinction from Beans: Beans supplies relevant source material and article context. Espresso's unit of value is **a reviewable interpretation plus its support**, and the measurable outcome is reduced review work on that interpretation. Buying one does not require buying the other. Do not relabel article retrieval or related-story counts as verified reasoning.

Smallest complete experience:

| Journey stage | Experience | Proof/measurement |
| --- | --- | --- |
| Trigger | Existing weekly brief is due | Prior brief and real watchlist establish a baseline |
| Evaluate | Reproduce one historical week with existing Event/Signal/evidence calls | Coverage, unsupported claims and review minutes compared with current stack |
| First useful outcome | Reviewer accepts/rejects candidate implications and opens their support | A supported brief-ready claim, not a populated dashboard |
| Repeat | Customer schedules its own reads for four weeks; stores returned IDs to avoid repeats | Weekly useful candidates and coverage gaps; no promise of a revision stream |
| Purchase | Optional future month-to-month qualified-evidence offer after readiness | Explicit price, allowed use, scope and no invented SLA |
| Support/exit | Bounded integration help, cancellation and export of customer-owned review annotations | Support cost and renewal reasons, with source rights preserved |

Initial source scope must be declared from an actual audit: named company announcements and selected industry/public-agency sources relevant to the enrolled watchlist. Report which sources are covered, unavailable or unverified. No new public completeness claim until measured. Policy information may appear as business context with original-source links; a reviewer determines applicability.

## 8. Enhancement decision and requirements

**No implementation is authorized by this document.** Current routes are sufficient to test the central workflow. First inspect actual outputs for missing support and dates. Use an offline evaluation ledger for the benchmark; do not add watchlists, webhooks, new forecasts, a dashboard or a generalized confidence score before evidence identifies a blocker.

| ID/priority | Requirement and customer outcome | Current vs proposed | Acceptance gate before release |
| --- | --- | --- | --- |
| R1 / evaluate now | Trace selected Signal → supporting Event → original source where available | Current traversal; completeness unverified | On a blinded 200-candidate sample, every accepted material claim has supporting source context; at least 95% of reviewer-accepted material claims directly supported. Report denominator and all failures; never call this universal accuracy |
| R2 / conditional must | Distinguish observation from interpretation and identify support for each material clause | Relationships exist; claim-level support qualification proposed | Reviewer can identify each material claim's support or see explicit insufficient-evidence status. Unsupported candidates excluded from ready-to-use set |
| R3 / conditional must | Explain time: source publication, first observation, event time if known, interpretation version | Only record creation bounds stable today | All new qualified records carry observation time and explicit known/unknown publication/event times. No creation timestamp silently relabeled as publication/effective date |
| R4 / conditional must | Make bounded source/watchlist coverage visible so silence is interpretable | Source catalog exists; health/coverage reporting proposed | Each enrolled company's expected sources marked checked, unavailable or out of scope for the evaluated period. Empty results never equated with no business change |
| R5 / conditional must | Preserve corrections/withdrawals for previously served qualified Signals | No public change feed guarantee | During a seeded correction test, a customer can identify superseded/withdrawn interpretation and amended support. Define version and retention semantics before selling persistent reuse |
| R6 / conditional must | Avoid counting syndicated reports as independent corroboration | Related evidence exists; independence not guaranteed | Sample all multi-source qualified candidates; same announcement reprints counted as one origin, with uncertain independence disclosed |
| R7 / conditional must | Scope rights for the actual customer use | Readiness register/control work unresolved | Owner-approved source-by-use permissions cover collection, processing, retention and distribution of each output type. Unknown permissions cannot be sold as an unrestricted commercial feed |
| R8 / later | Canonical IDs/aliases for selected watchlist when names cause misses or false matches | Normalized exact strings and discovery exist | Only build after measured entity errors prevent R1/coverage gates; assess precision/recall on agreed names and hard negatives |
| R9 / defer | Push alerts, user-specific relevance, point-in-time audit export, packaged connectors | Not part of current public surface | Require repeated paid-customer request and a scoped retention/rights/economics case |

Proposed thresholds are decision rules, not measured results. A 200-record evaluation only supports the sampled workload; report uncertainty and recurrence of failure modes. Assess recall against an independently assembled source/watchlist reference, not the API's own returned universe. Two reviewers should adjudicate at least 40 disputed/high-impact items to reveal label ambiguity. Neither source counts nor an LLM confidence label establish causal support. Forecasts remain clearly labeled interpretations; exclude return predictions and untested probabilities from the qualified offer.

A proposed future qualified record would include a material claim, interpretation status, supporting Event IDs, source URLs and permitted quotations/excerpts where rights allow, publication/observation times with unknown states, coverage scope and a version/correction status. These are **proposed product requirements**, not existing response fields or new routes. Prefer enriching the existing composition over a parallel API. Any later public contract change must follow the repository's full contract/documentation process; this internal research makes no such changes.

**GTM-only branch:** if existing records already meet R1-R7 for the selected workload and reviewers save at least 30% time, keep the current API. Improve adoption through a concrete task example, evaluation dataset and transparent coverage/use terms. If a paid offer has no additional value over free access, continue free adoption or park monetization rather than force a paywall rationale.

**Enhancement branch:** if buyer interest and source rights pass but a specific quality requirement blocks task success, implement only that requirement after a separate approved scope. A general evidence graph, global jurisdiction model or universal monitoring layer is not the first release.

## 9. Market model, price hypotheses and economics

There is no verified count of eligible buying accounts. Public developer directories, GitHub projects and social posts are recruitment leads, not a market census. Do not quote a global market-research CAGR as attainable revenue. Build a deduplicated prospect list of CI-software vendors, qualify named workflow owners, and count separate buying accounts only once.

Illustrative six-month funnel, **all inputs assumptions**:

| Scenario | Distinct prospects reached | Qualified with real weekly job | Complete evaluation | Accept future paid offer | Pay/retain to month 6 | MRR at $149 | Annualized run rate, not booked revenue |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Low | 100 | 20 | 10 | 2 | 1 | $149 | $1,788 |
| Base | 200 | 60 | 30 | 12 | 10 | $1,490 | $17,880 |
| High | 300 | 120 | 60 | 30 | 24 | $3,576 | $42,912 |

Base requires 30 completed evaluations across six months, roughly five/month, which may exceed founder capacity; validate onboarding effort. This is an attainable-funnel scenario, not a measured TAM/SAM. The serviceable pool exists only if source scope and permissions fit the buyer. Most sensitive inputs are qualified prospect access, repeated usefulness and acceptance of a premium over cheap monitoring. Twelve subscriptions purchased over time do not produce twelve full annual subscriptions of realized first-year revenue.

Pricing test, after readiness: **$149/month per integrating organization for a bounded evidence-qualified pack and repeatable review workflow**, starting with one agreed company watchlist and a proposed cap of 250 qualified Signal candidates/month. No overage billing or automated charge on trial exhaustion; larger scope requires a separate offer. Proposed 30-day paid pilot at the same price tests the complete recurring offer, not a discounted one-off consulting project. Until monetization is enabled, use unpaid evaluations and nonbinding price commitments clearly described as such. Do not collect payment now.

The current free allowance remains available and unchanged. The paid qualification pack is not operational, and neither source coverage nor client-display permissions can be assumed simply by assigning it a name. If qualification requires recurring manual analyst work beyond the cap, this is a service business and must be repriced or stopped under the operating constraints.

Illustrative monthly variable costs per account, USD; all are assumptions to replace with measured logs:

| Cost input | Efficient | Base | Stress |
| --- | ---: | ---: | ---: |
| Incremental retrieval/processing/inference | $4 | $12 | $35 |
| Source rights/data allocation | $0 | $10 | $40 |
| Support/review at $40/hour | $10 (15 min) | $20 (30 min) | $60 (90 min) |
| Payment processing estimate | $4.77 | $4.77 | $4.77 |
| Variable total | $18.77 | $46.77 | $139.77 |
| Contribution at $149 | $130.23 / 87.4% | $102.23 / 68.6% | $9.23 / 6.2% |

Processing estimate assumes 3% + $0.30 per monthly charge; not a provider quote. Data rights can be far above the assumed stress amount, making the offer unviable. Separate shared fixed platform costs, legal/entity/insurance costs, research/acquisition, refunds/tax and founder compensation. A 60% contribution target permits at most `$149 × 40% = $59.60` variable costs. At 250 qualified candidates that is under $0.2384 per candidate including support and payment costs. If 100 candidates require two minutes human review each, that alone is $133.33 at $40/hour: fail or change delivery, not silently assume automation.

At base contribution $102.23/month, a hypothetical $250 acquisition cost pays back in 2.45 retained months. Neither CAC nor retention is known, so do not infer lifetime value. At 30 minutes evaluation onboarding per completed trial plus follow-ups, count failed evaluations in acquisition cost as well. The $149 offer must win on whole-workflow value; raw retrieval cost leadership is not a credible claim against the cited alternatives.

## 10. Route to market and proof

First route: narrowly recruited developer evaluations, then evidence-led integration content. No outbound messages, public posts, partner commitments or ads were sent in this research.

| Stage | Planned action | Buyer/gatekeeper | Measured bottleneck |
| --- | --- | --- | --- |
| Discovery | Identify independent software teams already selling/prototyping weekly named-company briefs, through their public product/docs and relevant builder communities | Founder/engineer; exclude dormant/demo-only projects | Number of real target accounts reachable, not search traffic |
| Evaluation | Offer same-week blinded comparison using their current stack and Espresso, with source/rights scope visible | Engineer + human analyst | Time to first supported useful interpretation and integration effort |
| Approval | Show measured savings, missing coverage and allowed-use contract | Buyer/data owner | Whether value reaches the budget holder; commercial readiness |
| Purchase | Explicit month-to-month scoped offer at tested price | Founder/product/research lead | Paid or nonbinding commitment as appropriate; reason for rejection |
| Implementation | Minimal REST/MCP call recipe and IDs carried into their existing workflow | Engineer | Can integrate within two hours without bespoke analysis |
| Continued use | Four weekly outputs and subsequent renewal opportunity | Analyst/user plus buyer | Useful repeat candidates, support burden, renewal and churn reason |

Proof asset hypothesis: a transparent side-by-side example of one weak Signal, its support, a justified qualification/rejection and a resulting accepted brief sentence. Label historical/source-permitted samples and synthetic examples accurately. A comparison must include low-cost competitors and the user's own stack; show failures as well as wins. Do not say “outperforms,” “causal,” “accurate forecasts,” “complete,” “real-time,” or “audit-ready” without task-specific substantiation.

Possible later channels: developer templates for workflow tools, MCP directories and partnerships with CI-tool builders. Each requires its own audience and incentive test. A directory listing is awareness, not adoption. A partner only matters if it reaches qualified buyers and has a reason to distribute a component rather than build it. Social research did not prove r/SaaS or RegTech communities provide repeatable access to this exact buyer.

## 11. Experiments and 30/60/90-day decisions

Clock starts when an owner begins the experiment; deadlines are relative, not claims that actions were completed. Each gate can stop work early. Numbers are management decision thresholds, not statistically validated population estimates.

| Experiment | Participants/method | Metric, threshold and duration | Cost cap | Decision |
| --- | --- | --- | --- | --- |
| X1 / days 1-14: real job and access | Approach 30 independently qualified accounts; recruit 8 interviews, at least 5 providing a recent brief + source workflow, including incumbent-satisfied teams | At least 5 report a repeated implication-verification job and can identify its owner, current tool and actual effort; at least 3 accept a task trial | 16 hours + $100 | Fewer than 3 trials after 30 qualified approaches: revise access/message once; do not build features to cure an access failure |
| X2 / by day 30: current API coverage and rights | Three representative watchlists, 20-50 companies each; independent source checklist; 200 candidate claim sample over two reporting weeks | At least 80% of agreed in-scope meaningful developments found; at least 95% of accepted claims supported; explicit missing-source states; no material support fabricated | 24 hours + $250 | Misses are classified: coverage, dates, evidence, entity or reviewer disagreement. If critical source rights unavailable, stop that scope. If coverage fails in 2/3 lists, revise niche or park |
| X3 / days 15-45: task superiority | Three teams compare matched weekly tasks with current stack/low-cost alternative; counterbalance order, preserve acceptance criteria and reviewer notes | Median review minutes per accepted useful Signal falls ≥30%, with no drop in accepted-support precision; ≥2 teams reuse for four weeks; integration ≤2 hours/team | 20 hours + $400 | Passing existing API: GTM-only. One fixable blocker plus repeat use: scope minimal enhancement. No practical advantage: park |
| X4 / days 31-60: price/complete offer | Offer six qualified teams exact $149 monthly terms, including scope and alternative costs | After commercial readiness, ≥3 start paid pilot, ≥2 renew for second month; before readiness, record only nonbinding intent, not revenue | 15 hours + $250 incremental | If nobody commits and raw API is enough, maintain free access or revise value. No paid launch without readiness |
| X5 / days 61-90: reproducibility and retention | Up to 6 paying/evaluating organizations outside a single accommodating partner | At least 4 use on 3 of 4 weekly cycles; at least 3 renew where due; contribution ≥60%; support ≤30 min/account/month | 20 hours + $300 incremental | Continue bounded growth only if quality, repeat value, economics and readiness all pass. Document rejection/dropout causes |

Do not add all later caps to the initial $750 budget automatically; each subsequent phase requires a fresh owner commitment after passing earlier gates. Three teams cannot establish prevalence or general quality. Record approached, qualified, interviewed, started, completed, paid, renewed and dropped out separately; report rejected hypotheses and losses. If fewer than six qualified teams reach pricing, label conversion findings inconclusive rather than presenting a small denominator as success.

A suggested quality-control trial compares at least 40 disputed/high-impact records using two reviewers. If adjudication reveals that “supported implication” is inherently ambiguous, narrow output to attributable observations and clearly conditional interpretations. If that becomes indistinguishable from Newsflash/Parallel plus a small prompt, stop paid qualification rather than inflate the product story.

## 12. Decision log and expansion rules

| Decision | Rationale | Revisit condition | Proposed owner |
| --- | --- | --- | --- |
| Reject broad low-price API differentiation | Multiple inexpensive event/monitoring/policy competitors already exist | Only revisit if a specific rights/coverage/workflow benchmark shows a credible cost-value advantage | Product owner |
| Select recurring named-company brief builders for evaluation | Closest fit to current Signal-support composition; technically recruitable | X1 cannot recruit actual workflow owners or X2 reveals inadequate source fit | Product/research owner |
| Leave vertical open pending real watchlists | Evidence does not justify inventing a vertical demand story | Three or more target workflows converge on a deliverable source scope | Product/research owner |
| Use existing API before enhancement | Current routes can test value without new platform work | X3 identifies a repeated blocker after buyer/rights gates pass | Engineering/product owner |
| Treat claim support, time, coverage and corrections as potential paid differentiators | Current routes do not guarantee these quality outcomes | Existing records pass, or buyers value a different gap | Engineering/product owner |
| Gate payment and customer display on readiness | Recorded local policy and unresolved source-specific use permissions | Owner verifies prerequisites and allowed-use offer | Business owner |
| Expand only after repeated value and positive contribution | Signups and trial interest do not prove a business | Quality gates, second/third renewal opportunities and a supported next-segment case | Business owner |

The next required result is a real watchlist benchmark and buyer behavior, not another broad list of market features. If no meaningful enhancement is necessary, the product decision is a focused developer GTM experiment using current capabilities. If no switch/purchase reason survives the alternatives comparison, the correct decision is to park paid positioning while preserving the useful free API.
