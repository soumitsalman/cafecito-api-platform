# Cafecito API Legal Readiness Plan

**Prepared:** September 11, 2026  
**Primary jurisdiction:** United States, with Washington as the operator's home state  
**Operator:** Strategic Implementation Advising, LLC, doing business through Project Cafecito

> This is a product and documentation readiness assessment, not legal advice. U.S. copyright/internet, privacy, and commercial counsel must confirm the analysis and approve customer-facing terms before production reliance.

## Executive decision

Cafecito should treat full-text output as a licensed capability, not as a default consequence of receiving content through RSS or being able to fetch a public page.

The safest operating rule is:

| Source-rights state | Metadata and canonical URL | Limited excerpt or summary | Full-content API | Persistent body storage | Training |
| --- | --- | --- | --- | --- | --- |
| Unknown or unlicensed | Case-by-case access policy | Counsel-approved policy only | No | No, except transient processing approved by counsel | No |
| Search/discovery permission | Yes, as licensed | As licensed | No | As licensed | No |
| Summary/RAG permission | Yes | Yes, within license | No | Only as licensed | No unless separate grant |
| Full-display permission | Yes | Yes | Yes, within license and customer rights | Within retention terms | No unless separate grant |
| Public domain or government work | Confirm status and any exceptions | Yes | Yes if status confirmed | Per policy | Per policy |

RSS or Atom availability is not a full-display license. Feed terms, publisher site terms, copyright notices, syndication arrangements, and third-party media exclusions still apply.

## Current implementation findings

### What is already risk-reducing

- PyCoffeeMaker marks bodies obtained by page scraping as restricted.
- Beans database selection suppresses restricted bodies from API output.
- Beans omits content by default and returns it only when a caller requests `full_content=true`.
- Article output includes a canonical URL; Source metadata is available when the publisher join succeeds and may be absent.
- Public docs already tell customers that full content is not guaranteed and to cite the canonical URL.

### Material gaps and completed remediations

| Area | Audit status | Risk | Completed or required disposition |
| --- | --- | --- | --- |
| Source-delivered full text | Unrestricted body candidates include RSS feeds, official filings, platform self-posts, manually inserted rows, and legacy rows; acquisition provenance is not retained per Article | Delivery method does not itself establish redistribution rights or a complete rights chain | Source-by-source license/status review; default unknown sources to no full text |
| Scraped-body retention | Marked scraped bodies are suppressed at query time but can remain in processing cache and the main database; the observed default cleanup window is six months | Copying and storage risk remains even when public output is suppressed | Inventory and minimize retention, define purge/backup/legal-hold rules, and obtain counsel approval before relying on the policy |
| Rights enforcement | One nullable restricted-content flag distinguishes only some scraped bodies; NULL or unknown rows fail open and the flag does not express provenance or permission | It cannot represent crawl, summary, full display, retention, RAG, media, termination, and customer sublicensing rights | Implement a fail-closed rights registry and output policy |
| Full-content routes | `full_content=true` is available on nine public REST routes, seven MCP tools, and an internal unique-Article route; public collections can return up to 100 bodies per request | A single request can deliver many full bodies | Gate every content-producing route; consider detail-only full text and a lower body limit |
| API response | No license or permitted-use metadata | Customers cannot determine retention/display/training rights | Add machine-readable rights metadata before licensed full-display launch |
| OpenAPI license scope | Before this work, the API specs described the API as MIT; all four current Beans/Espresso gateway and backend specs now omit a data/content license | A software license can be confused with rights in API data and publisher content | Completed in the API contracts; keep the repository software license separately scoped |
| Terms | The prior page was informal and materially incomplete; a comprehensive draft is now published in the portal source | Weak notice, rights allocation, remedies, and downstream restrictions if the draft is not approved and assented to | Counsel review plus affirmative clickwrap remain required |
| Privacy | The prior page omitted published-content data, legal uses, request handling, and state disclosures; a comprehensive draft is now published in the portal source | Transparency and state-law risk if statements do not match operations | Data inventory, retention validation, and counsel review remain required |
| Assent | No clear clickwrap or versioned acceptance record found | Terms may be harder to enforce | Implement affirmative acceptance at signup and API-key creation |
| Copyright process | A public notice/counter-notice policy is now drafted | Publication alone does not create safe-harbor eligibility | Counsel must assess eligibility and register an agent if appropriate |
| Publisher channel | A public opt-out/licensing policy and internal runbook are now drafted | Requests will not be effective unless intake and runtime enforcement exist | Establish monitored intake and enforcement |
| Crawler identity | Collector has a Cafecito user agent, but ordinary page/feed requests currently use a browser-like user agent; no robots enforcement was found | Transparency, contract, and access-law risk | Use a dedicated disclosed crawler identity and implement auditable robots/source-policy checks |
| Corrections | A public correction/retraction policy and internal runbook are now drafted | Defamation, false-light, accuracy, and stale-content risk remains without operations | Establish monitored intake, decision ownership, and propagation controls |
| Company identity | About, Contact, Terms, and Privacy now identify the operator and contact | Incorrect or unverified entity/contact facts can undermine notices and contracts | Verify entity, address, contact, authority, and trade-name status before reliance |
| OpenAPI legal link | All four current Beans/Espresso gateway and backend specs link to the canonical Terms | Future annotation/spec changes could drift | Completed in the API contracts; preserve with the contract cascade |
| Licenses | A rights schedule exists, but no executed publisher master license was found | Full-text rights cannot be proven | Negotiate licenses using the supplied schedule and a counsel-approved agreement |
| Privacy and content operations | No verified retention schedule, request log, subprocessors schedule, breach playbook, or production body-provenance inventory was found in scope | Public claims may not match operations and stored test fixtures may themselves redistribute source text | Complete operational inventory, counsel review, and a test-fixture/license audit |
| Subscription law | Subscription events exist, but checkout disclosures/cancellation/renewal compliance were not assessed | State auto-renewal and consumer-protection risk | Review purchase and cancellation flows before paid rollout |

## Comprehensive U.S. legal and operational control inventory

This is the exhaustive control inventory for the product facts currently known. Counsel should add state-, source-, customer-, and industry-specific obligations as scope changes.

### 1. Entity, authority, and public identity

- Confirm Strategic Implementation Advising, LLC is active and in good standing.
- Confirm Project Cafecito is properly registered or used as an assumed/trade name where required.
- Use the LLC's exact legal name on Terms, invoices, publisher licenses, privacy notices, and DMCA registration.
- Verify the public mailing address, legal email, privacy email, copyright contact, publisher contact, security contact, and telephone number.
- Confirm who has authority to accept customer and publisher contracts.
- Maintain tax registrations, sales-tax analysis, and required business licenses.
- Preserve signed contracts, policy versions, acceptance records, notices, and corporate approvals.

### 2. Source access and crawler governance

- Inventory every source, feed, API, public record repository, and page-scraping target.
- Record source Terms of Service, feed terms, API terms, robots directives, copyright notices, paywall/login status, cease-and-desist history, and jurisdiction.
- Prohibit bypass of authentication, paywalls, CAPTCHAs, session controls, and technical access measures.
- Publish a stable crawler user agent, information page, and monitored contact.
- Implement auditable robots and per-domain rate controls.
- Provide URL-, path-, feed-, and domain-level blocking.
- Record access decisions and re-review on material source-term changes.
- Obtain counsel guidance for post-block access, cease-and-desist letters, contract claims, trespass/tortious-interference theories, CFAA risk, and DMCA § 1201.
- Treat paid crawl access separately from copying and redistribution permission.

### 3. Copyright and content licensing

For every source, determine rights to:

- access and crawl;
- reproduce technical copies;
- parse and extract;
- convert HTML/feed content to Markdown, text, JSON, YAML, or TOON;
- normalize, deduplicate, classify, and cluster;
- cache transiently;
- store persistently or archive;
- index and search;
- quote or excerpt;
- summarize;
- create embeddings and use content for RAG;
- return complete text;
- distribute through REST, MCP, UI, feeds, newsletters, or exports;
- allow customer display or sublicensing;
- translate;
- analyze and retain derived data after termination;
- use for training or fine-tuning; and
- use names, logos, trademarks, images, video, charts, and other media.

Also:

- Preserve copyright-management information; review DMCA § 1202 before removing bylines, notices, or metadata.
- Exclude wire-service, syndicated, freelance, embedded, and other third-party material the publisher cannot sublicense.
- Distinguish U.S. federal government works from state/local government, contractor, third-party, and otherwise copyrighted material.
- Maintain public-domain and open-license evidence rather than inferring status from accessibility.
- Ensure customer rights never exceed Cafecito's source rights.
- Define termination, purge, correction, withdrawal, legal-hold, and surviving-derived-data rules.
- Measure and report royalty events consistently with publisher agreements.

### 4. Rights-aware product controls

- Default unknown sources to no full-text output.
- Enforce rights at URL and domain level on every route that can return content.
- Separate search, excerpt, summary/RAG, full display, archive, and training permissions.
- Limit full text to the smallest necessary response surface and page size.
- Apply license expiration, territory, product, customer, and retention rules at request time.
- Return rights metadata that states scope, attribution, retention, redistribution, training, and deletion obligations.
- Block or separately license images, video, charts, logos, and syndicated material.
- Support urgent suppression, revocation, cache invalidation, downstream notification, and audit logging.
- Test that restricted bodies never leak through alternate routes, formats, MCP tools, errors, logs, exports, or caches.
- Preserve canonical source URL, publisher, author, date, license, and correction status.

### 5. Customer contracting and assent

- Publish attorney-approved API Terms, AUP, Privacy Policy, Content Rights Policy, Copyright Policy, Publisher Request Policy, Corrections Policy, and Automated Collection Policy.
- Implement affirmative clickwrap at signup or API-key creation; do not rely only on footer links or continued use.
- Store user/account, organization, policy version, timestamp, IP or equivalent evidence, and acceptance event.
- Re-consent when a material change requires it; do not silently expand data use.
- Align documentation, plan descriptions, checkout, invoices, order forms, and marketing.
- Define permitted use, prohibited use, credential security, usage limits, suspension, termination, fees, taxes, renewals, refunds, confidentiality, IP, warranties, disclaimers, liability, indemnity, export/sanctions, governing law, and notices.
- Use an enterprise MSA/order form/SLA/security exhibit where customers need negotiated commitments.
- Ensure customer terms incorporate source-specific restrictions and deletion duties.
- Prohibit consumer-report/high-impact-decision use unless separately designed and reviewed.

### 6. Copyright, publisher, and correction operations

- Have monitored intake channels and trained owners.
- Register and maintain a DMCA agent if counsel determines Cafecito should seek § 512 safe-harbor protection.
- Publish registered-agent details exactly and renew registration when required.
- Maintain valid-notice, counter-notice, repeat-infringer, restoration, and legal-escalation procedures.
- Keep a separate publisher opt-out and licensing process.
- Triage privacy, defamation, safety, trademark, correction, retraction, and copyright complaints differently.
- Log requester, authority, URLs/IDs, dates, evidence, interim action, decision, reason, notifications, and closure.
- Support litigation holds without silently defeating required deletion.
- Maintain correction/retraction history and downstream notification where required.

### 7. Editorial, tort, trademark, and media controls

- Preserve allegation language, attribution, dates, and source provenance.
- Prohibit fabricated quotations and misleading headline generation.
- Disambiguate people and organizations before merging stories or claims.
- Apply heightened review to crime, minors, health, finance, elections, sanctions, and reputational allegations.
- Provide correction, appeal, and urgent safety channels.
- Assess defamation, false light, public disclosure of private facts, right of publicity, negligence, and state-law misappropriation.
- Do not imply publisher endorsement.
- License trademarks/logos separately; follow brand rules.
- Avoid copying images, video, charts, or embedded media without explicit rights.
- Review Section 230 only as a claim-specific defense; it does not resolve intellectual-property claims or Cafecito's own content.

### 8. Privacy and data protection

- Inventory account data, API queries, submitted URLs, logs, cookies, billing metadata, support data, rights requests, and personal information in publisher material.
- Map sources, purposes, legal roles, recipients, locations, retention, deletion, security, and customer access.
- Validate that the Privacy Policy matches actual practices.
- Implement access, correction, deletion, portability, opt-out, limitation, appeal, authorized-agent, and identity-verification procedures where applicable.
- Assess every applicable state comprehensive privacy law and threshold annually.
- Assess data-broker status and registrations, including California's Delete Act and DROP requirements.
- Assess Washington My Health My Data if the service collects or infers consumer health data within scope; publish a separate policy and obtain required consent if applicable.
- Assess biometric, genetic, precise-location, criminal-record, minor, and other sensitive-data laws.
- Assess COPPA and maintain the 18+ account rule.
- Execute service-provider/processor contracts and a DPA; maintain a subprocessor schedule and change-notice process.
- Assess GDPR/UK GDPR and international transfer mechanisms before targeting or monitoring covered individuals.
- Adopt a documented retention schedule, deletion propagation, backup handling, and legal-hold procedure.
- Implement data minimization, access controls, encryption, vulnerability management, vendor diligence, incident response, and security training.
- Maintain a 50-state plus territory breach-notification matrix and response counsel.
- Evaluate whether the FTC Health Breach Notification Rule, HIPAA, GLBA, FCRA, DPPA, VPPA, FERPA, or other sectoral law applies to a product or customer use.

### 9. Consumer protection, billing, and marketing

- Substantiate accuracy, freshness, coverage, “real-time,” “licensed,” “complete,” and publisher-authorization claims.
- Disclose that full content is source- and rights-dependent and that generated/derived fields may be wrong.
- Review checkout, automatic renewal, cancellation, refund, pricing, tax, and free-trial disclosures under federal and state law.
- Comply with CAN-SPAM for commercial email, including sender identity, postal address, and opt-out handling.
- Obtain consent for marketing channels where required.
- Avoid dark patterns and misleading policy-change notices.
- Maintain sanctions/export screening proportionate to customers and use.
- Review portal accessibility and target WCAG 2.2 AA as an operational standard; provide an accessibility contact.
- Review state unfair/deceptive-practices and publicity/privacy claims before new features or marketing.

### 10. Insurance, governance, and ongoing review

- Obtain quotes for media liability, Tech E&O, cyber/privacy, and appropriate IP coverage; review exclusions and retroactive dates.
- Assign executive owners for legal, privacy, security, publisher relations, and corrections.
- Maintain a source/license inventory, risk register, incident register, and renewal calendar.
- Require legal/product review before new source classes, new full-content surfaces, training uses, international launch, or high-risk fields.
- Audit source rights, access behavior, customer use, policy accuracy, and takedown performance at least annually.
- Budget for publisher royalties, legal review, DMCA administration, enterprise redlines, and disputes.

## Release gates

### Gate A — publish revised policies

Before treating the public drafts as approved:

- attorney reviews Terms, Privacy, AUP, Content Rights, Copyright, publisher request, corrections, and crawler language;
- legal entity, address, email, telephone, and governing-law/venue choices are verified;
- liability cap, indemnity, confidentiality, age, subscription, and privacy statements are approved;
- actual retention, subprocessors, cookies/analytics, billing, and marketing practices are confirmed.

### Gate B — continue any unlicensed full-content output

Do not continue broad full-content output merely because it came from RSS. Counsel must approve a source-class policy or the source must have a license or confirmed public-domain/open-license basis. Unknown sources should return metadata, canonical URL, and counsel-approved limited output only.

### Gate C — licensed full display

Before enabling licensed full text:

- signed rights schedule covers format conversion, API distribution, customer display, storage, territory, media exclusions, and termination;
- the rights registry controls every content-producing route;
- response metadata communicates downstream limits;
- purge, withdrawal, corrections, usage reporting, and royalty measurement are tested;
- customer Terms/order forms grant no broader rights.

## Repository artifacts produced

### Public drafts

- `docs/pages/company/terms-of-use.md`
- `docs/pages/company/privacy-policy.md`
- `docs/pages/company/acceptable-use-policy.md`
- `docs/pages/company/content-rights-policy.md`
- `docs/pages/company/copyright-policy.md`
- `docs/pages/company/publisher-requests.md`
- `docs/pages/company/corrections-policy.md`
- `docs/pages/company/automated-collection-policy.md`
- updated About, Contact, API-key, Beans, API-conventions, navigation, and OpenAPI legal metadata

### Internal working documents

- `legal/EXTERNAL_ACTIONS.md`
- `legal/CONTENT_RIGHTS_REGISTER_TEMPLATE.md`
- `legal/PUBLISHER_RIGHTS_SCHEDULE_TEMPLATE.md`
- `legal/DATA_PROCESSING_ADDENDUM_TEMPLATE.md`
- `legal/NOTICE_AND_CORRECTION_RUNBOOK.md`

## Primary authorities and official guidance

- [17 U.S.C. § 106 — exclusive rights](https://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title17-section106)
- [17 U.S.C. § 107 — fair use](https://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title17-section107)
- [U.S. Copyright Office — fair-use guidance](https://www.copyright.gov/fair-use/)
- [17 U.S.C. § 512 — service-provider limitations](https://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title17-section512)
- [U.S. Copyright Office — Section 512 safe harbors and notice-and-takedown](https://www.copyright.gov/512/)
- [18 U.S.C. § 1030 — computer fraud and abuse](https://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title18-section1030)
- [47 U.S.C. § 230 — online-platform liability provision](https://uscode.house.gov/view.xhtml?req=granuleid:USC-prelim-title47-section230)
- [U.S. Copyright Office — Section 1201 anti-circumvention](https://www.copyright.gov/dmca/)
- [17 U.S.C. Chapter 12 — §§ 1201 and 1202](https://uscode.house.gov/view.xhtml?edition=prelim&path=%2Fprelim%40title17%2Fchapter12)
- [Ninth Circuit — hiQ Labs v. LinkedIn (2022)](https://cdn.ca9.uscourts.gov/datastore/opinions/2022/04/18/17-16783.pdf)
- [Supreme Court — Andy Warhol Foundation v. Goldsmith (2023)](https://www.supremecourt.gov/opinions/22pdf/21-869_87ad.pdf)
- [California Attorney General — CCPA](https://www.oag.ca.gov/privacy/ccpa)
- [California Privacy Protection Agency — Delete Act regulations](https://cppa.ca.gov/regulations/)
- [Washington My Health My Data Act, RCW 19.373](https://app.leg.wa.gov/RCW/default.aspx?cite=19.373)
- [FTC — Protecting Personal Information](https://www.ftc.gov/business-guidance/resources/protecting-personal-information-guide-business)
- [FTC — Data Breach Response](https://www.ftc.gov/business-guidance/resources/data-breach-response-guide-business)
- [FTC — COPPA compliance](https://www.ftc.gov/business-guidance/resources/complying-coppa-frequently-asked-questions)
- [FTC — CAN-SPAM compliance](https://www.ftc.gov/business-guidance/resources/can-spam-act-compliance-guide-business)
- [15 U.S.C. § 1681b — permissible purposes for consumer reports](https://uscode.house.gov/view.xhtml?edition=prelim&req=granuleid%3AUSC-prelim-title15-section1681b)
- [DOJ — Web accessibility and the ADA](https://www.ada.gov/resources/web-guidance/)
- [OFAC — Framework for Compliance Commitments](https://ofac.treasury.gov/recent-actions/20190502_33)
- [15 U.S.C. § 7001 — electronic records and signatures](https://uscode.house.gov/view.xhtml?edition=2023&num=0&req=granuleid%3AUSC-2023-title15-section7001)
