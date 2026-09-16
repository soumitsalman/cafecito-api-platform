# Automated Collection Policy

**Status:** Draft for counsel and factual review; not approved for production reliance.

**Effective Date:** To be set on approval and publication.

Cafecito uses automated systems to discover and process public publisher pages, RSS or Atom feeds, public records, and other source material for search, classification, source linking, and related Services.

## Collection principles

Before this draft becomes effective, Cafecito will validate and implement the following collection controls:

- access public material and source-provided feeds;
- avoid intentionally bypassing authentication, paywalls, CAPTCHAs, or other access controls;
- use reasonable request rates and avoid unreasonable source load;
- retain canonical source links and available source attribution;
- review source terms, machine-readable directives, publisher requests, and license restrictions;
- avoid returning article bodies obtained solely by scraping ordinary publisher pages;
- treat feed delivery as distinct from permission to redistribute content; and
- provide a direct publisher request and opt-out channel.

Robots directives and other technical signals inform access policy but do not independently grant copyright or redistribution rights. Before publication, this page must be completed with the actual crawler user-agent string, crawler information URL, monitored contact, robots behavior, request-rate policy, caching/retention behavior, and whether source text is used for RAG or training.

## What automated processing may do

Automated processing may extract, process, and retain factual metadata and text to classify content, generate or normalize categories, regions, entities, sentiments and other attributes, identify related coverage, and support search. Text can remain stored for these service functions even when it is restricted from API output. Beans may return content obtained through a source feed, official public record, platform self-post, or another non-page-scrape ingestion path when available, subject to the [Third-Party Content and Attribution Policy](/company/content-rights-policy).

Cafecito generally returns media source URLs rather than copying third-party images, video, charts, or other embedded media. A URL does not grant a license to display or copy the underlying media.

## Publisher controls

A publisher or authorized representative may ask Cafecito to stop or limit collection, suppress full content, correct source information, remove specified material, or discuss a license through [Publisher Requests and Licensing](/company/publisher-requests).

Include the affected domain, feed, paths or URLs, the requested action, and information showing authority to act for the source. Cafecito may verify the request and apply URL-, path-, feed-, or domain-level restrictions.

## Security and vulnerability research

This policy does not authorize access to non-public systems, credential use, security testing, circumvention, or interference with a source. Report a security concern through [Contact](/contact) without exploiting or retaining unnecessary data.
