import { extname } from "node:path";
import {
  collectNavFiles,
  exists,
  issue,
  pageFileCandidates,
  readText,
  rel,
  walk,
} from "./verify-doc-lib.mjs";

export async function checkInventory(ctx) {
  const issues = [];
  const check = "inventory";
  const configText = await readText(ctx.zudokuPath);
  const navFiles = collectNavFiles(configText);
  if (navFiles.length === 0) {
    issues.push(issue(check, rel(ctx.root, ctx.zudokuPath), "no navigation file: entries"));
  }

  for (const fileId of navFiles) {
    const candidates = pageFileCandidates(ctx.pagesDir, fileId);
    let found = false;
    for (const c of candidates) {
      if (await exists(c)) {
        found = true;
        break;
      }
    }
    if (!found) {
      issues.push(
        issue(
          check,
          rel(ctx.root, ctx.zudokuPath),
          `nav file ${fileId} has no matching docs/pages source`,
        ),
      );
    }
  }

  const pages = await walk(
    ctx.pagesDir,
    (p, name) => [".md", ".mdx"].includes(extname(name)),
  );
  const navSet = new Set(navFiles);
  for (const page of pages) {
    const loc = rel(ctx.pagesDir, page).replace(/\.(mdx|md)$/, "");
    if (!navSet.has(loc)) {
      issues.push(
        issue(check, rel(ctx.root, page), "page is not listed in zudoku navigation"),
      );
    }
  }
  return issues;
}
