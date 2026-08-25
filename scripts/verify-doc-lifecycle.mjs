import {
  LIFECYCLE_FIELDS,
  LIFECYCLE_STATUSES,
  MAX_LIFECYCLE_AGE_DAYS,
  daysSince,
  issue,
  parseIsoDate,
  parseLifecycle,
  readText,
  rel,
  walk,
} from "./verify-doc-lib.mjs";

export async function checkLifecycle(ctx) {
  const issues = [];
  const check = "lifecycle";
  const now = ctx.nowMs ?? Date.now();
  const files = await walk(ctx.designDir, (p, name) => name.endsWith(".md"));
  for (const file of files) {
    const text = await readText(file);
    const loc = rel(ctx.root, file);
    const fields = parseLifecycle(text);
    const isIndex = /\/README\.md$/i.test(loc.replace(/\\/g, "/"));

    if (isIndex) {
      if (!fields["last verified"] && !/\*\*Last verified:\*\*/i.test(text)) {
        issues.push(issue(check, loc, "missing last verified"));
      }
      if (!fields.owner && !/\*\*Owner role:\*\*/i.test(text)) {
        issues.push(issue(check, loc, "missing owner"));
      }
      continue;
    }

    for (const key of LIFECYCLE_FIELDS) {
      const present = Boolean(fields[key] && String(fields[key]).trim());
      if (!present) {
        issues.push(issue(check, loc, `missing lifecycle field: ${key}`));
      }
    }

    const statusRaw = String(fields.status || "")
      .replace(/[*_`]/g, "")
      .trim()
      .toLowerCase()
      .replace(/[^a-z].*$/, "");
    if (fields.status && !LIFECYCLE_STATUSES.has(statusRaw)) {
      issues.push(
        issue(
          check,
          loc,
          `status must be one of ${[...LIFECYCLE_STATUSES].join(", ")}`,
        ),
      );
    }

    const verified = parseIsoDate(fields["last verified"]);
    if (fields["last verified"] && !verified) {
      issues.push(issue(check, loc, "last verified must include an ISO date"));
    } else if (verified && daysSince(verified, now) > MAX_LIFECYCLE_AGE_DAYS) {
      issues.push(
        issue(
          check,
          loc,
          `last verified is older than ${MAX_LIFECYCLE_AGE_DAYS} days`,
        ),
      );
    }

    if (statusRaw === "superseded") {
      const sup = String(fields["superseded by"] || "").toLowerCase();
      if (!sup || sup === "n/a" || sup === "none") {
        issues.push(
          issue(check, loc, "superseded documents require a superseded-by target"),
        );
      }
    }
  }
  return issues;
}
