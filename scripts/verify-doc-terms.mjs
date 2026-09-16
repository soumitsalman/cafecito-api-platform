import {
  STALE_PUBLIC_PATHS,
  exists,
  issue,
  readText,
  rel,
  scanForbidden,
} from "./verify-doc-lib.mjs";

export async function checkTerms(ctx) {
  const issues = [];
  const check = "terms";
  for (const file of ctx.publicFiles) {
    const text = await readText(file);
    const loc = rel(ctx.root, file);
    for (const rule of scanForbidden(text)) {
      issues.push(
        issue(check, loc, `forbidden public term (${rule.id}): ${rule.hint}`),
      );
    }
    for (const stale of STALE_PUBLIC_PATHS) {
      if (text.includes(stale)) {
        issues.push(issue(check, loc, `stale public path or term: ${stale}`));
      }
    }
  }
  if (ctx.zudokuPath && await exists(ctx.zudokuPath)) {
    const navigation = await readText(ctx.zudokuPath);
    if (/file:\s*["'](?:\.\.\/)?legal\//i.test(navigation) || /to:\s*["']\/legal\//i.test(navigation)) {
      issues.push(issue(check, rel(ctx.root, ctx.zudokuPath), "internal legal directory exposed in portal navigation"));
    }
  }
  return issues;
}
