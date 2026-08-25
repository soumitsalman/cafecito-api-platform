import {
  REQUIRED_GENERATED_SLUGS,
  exists,
  generatedOutputHasSlug,
  issue,
  readText,
  rel,
  scanForbidden,
  walk,
} from "./verify-doc-lib.mjs";

export async function checkGenerated(ctx) {
  const issues = [];
  const check = "generated";
  const searchRoots = ctx.generatedDirs ?? [];
  const files = [];
  for (const dir of searchRoots) {
    if (!(await exists(dir))) continue;
    files.push(
      ...(await walk(dir, (p, name) => {
        const n = name.toLowerCase();
        return n === "llms.txt" || n === "llms-full.txt" || n.endsWith(".md");
      })),
    );
  }

  if (files.length === 0) {
    if (ctx.requireGenerated) {
      issues.push(
        issue(
          check,
          "generated-output",
          "required generated Markdown/llms.txt not found (build the portal, then pass --require-generated)",
        ),
      );
    }
    return { issues, skipped: files.length === 0 && !ctx.requireGenerated };
  }

  const combined = [];
  for (const file of files) {
    const text = await readText(file);
    const loc = rel(ctx.root, file);
    combined.push(text);
    for (const rule of scanForbidden(text)) {
      issues.push(
        issue(check, loc, `forbidden public term in generated output (${rule.id})`),
      );
    }
  }

  const blob = combined.join("\n");
  const slugs = ctx.requiredSlugs ?? REQUIRED_GENERATED_SLUGS;
  for (const slug of slugs) {
    if (!generatedOutputHasSlug(blob, slug)) {
      issues.push(
        issue(check, "llms.txt", `generated output missing required page ${slug}`),
      );
    }
  }
  return { issues, skipped: false };
}
