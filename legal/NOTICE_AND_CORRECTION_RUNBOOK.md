# Notice, Publisher Request, Privacy, and Correction Runbook

**Status:** Internal operational draft  
**Owner:** Assign before publication  
**Counsel escalation:** Assign primary and backup

## 1. Intake

Monitor the published inbox continuously during business days. Create one case per request and record:

- case ID and category;
- date/time and delivery method;
- requester name, organization, authority, and contact;
- affected domains, feeds, URLs, Article IDs, routes, customers, and output formats;
- requested action and stated legal basis;
- supplied evidence;
- urgency and safety factors;
- interim action;
- assigned owner and counsel;
- decisions, reasons, notices, and closure date; and
- rights-registry policy versions before and after action.

Never place credentials or unnecessary sensitive personal information in the case log.

## 2. Triage

| Category | Examples | Immediate owner |
| --- | --- | --- |
| Copyright/DMCA | Full text, image, takedown, counter-notice | Legal/copyright |
| Publisher opt-out/license | Stop crawl, suppress content, terms, commercial rights | Publisher relations + legal |
| Correction/retraction | Wrong source, date, identity, label, stale/retracted content | Editorial/product |
| Privacy | Access, deletion, correction, sensitive/public personal data | Privacy |
| Defamation/safety | Allegations, minors, doxxing, threats, court order | Legal + safety |
| Trademark/endorsement | Logo, source branding, false affiliation | Legal |
| Security/abuse | Credential leak, exploit, attack | Security |
| Subpoena/government | Legal process or authority request | Counsel only |

Escalate immediately when a request involves imminent physical harm, a minor, exposed credentials, active exploitation, a court order, law-enforcement demand, threatened litigation, a major publisher, or widespread sensitive data.

## 3. Interim restriction

When credible harm or rights risk may continue, apply the narrowest effective interim control while review proceeds:

1. block full text;
2. block copied media;
3. suppress the URL or record;
4. stop collection for the path/feed/domain;
5. disable a customer's affected use;
6. preserve evidence under access restriction; and
7. notify engineering of cache/export/downstream implications.

An interim restriction is not an admission. Record the reason, approver, scope, start time, and re-review time.

## 4. DMCA notice review

Confirm the notice substantially includes:

- signature;
- identified copyrighted work(s);
- identified material and enough information to locate it;
- sender contact information;
- good-faith unauthorized-use statement; and
- accuracy and perjury/authority statement.

If incomplete, ask for missing information without unnecessarily delaying a justified interim restriction. If facially valid, route to counsel and act expeditiously as required. Determine which § 512 category, if any, applies; do not assume safe harbor covers Cafecito's own copying.

Record customer/user notification, repeat-infringer implications, removal/suppression scope, and source-policy changes.

## 5. Counter-notice

Confirm:

- signature;
- removed material and prior location;
- mistake/misidentification statement under penalty of perjury;
- name, address, and telephone;
- required federal-court jurisdiction consent; and
- acceptance-of-service statement.

Send to the original notice sender when legally appropriate. Calendar the statutory window. When § 512(g) applies, restoration generally occurs no sooner than 10 and no later than 14 business days after receipt unless the notice sender reports a filed court action. Counsel must approve restoration.

## 6. Publisher opt-out or rights restriction

Verify the requester represents the publisher or rights holder. Resolve domain/feed ownership. Determine whether the request affects access, full text, metadata, excerpts, summaries, embeddings/RAG, media, historical data, or customer output.

Apply the result to the rights registry. A source-level opt-out must control every REST/MCP route, format, cache, export, and future collection job. Record whether already-delivered customers require notice or deletion.

## 7. Correction or retraction

Compare Cafecito output to the canonical source and any publisher correction/retraction.

- Correct Cafecito metadata and derived labels directly when wrong.
- Do not silently rewrite publisher text; link to or reflect the publisher's correction.
- Preserve “alleged,” “reported,” “charged,” and similar status language.
- Disambiguate people and organizations.
- Invalidate affected summaries, clusters, rankings, Events, Signals, caches, and exports as applicable.
- Record the before/after value, source evidence, reason, reviewer, and time.
- Notify affected customers when contract, license, risk, or feasibility requires it.

## 8. Privacy request

Determine whether the data is account/usage data, Customer-controlled data, third-party published material, or a mixture. Verify identity and authorized-agent authority proportionately. Record jurisdiction, requested right, deadline, extensions, exceptions, searches performed, decision, response, deletion propagation, and appeal.

Do not remove accurate public-interest publisher material automatically. Consult privacy counsel for sensitive health data, minors, doxxing, criminal allegations, data-broker requests, and conflicts with journalism/public-record exceptions.

## 9. Government and legal process

Do not respond substantively without counsel. Validate jurisdiction, authority, scope, signature, service, confidentiality/gag restrictions, preservation duties, and ability to narrow or challenge. Disclose only what is legally required and record the production securely.

## 10. Closure

Before closing:

- requested and required actions are complete;
- rights registry and collection/output systems agree;
- caches, exports, and downstream notices are handled;
- retention or legal hold is documented;
- response was sent through an approved channel;
- appeal/counter-notice dates are calendared;
- systemic root cause and preventive change are recorded; and
- case access is limited to personnel with a need to know.

## 11. Readiness exercises

Quarterly until the process is mature, then at least annually:

- test a publisher domain-level full-text block;
- test a single-URL correction and cache invalidation;
- run a mock DMCA notice and counter-notice calendar;
- run a privacy deletion request through active systems and backups;
- verify all public contact links and inbox routing;
- export an audit record showing the exact rights policy applied; and
- review unresolved, repeated, and late cases with counsel.
