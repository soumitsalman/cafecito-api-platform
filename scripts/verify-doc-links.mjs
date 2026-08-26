import { dirname, resolve } from "node:path";
import { exists, issue, mdLinks, readText, rel } from "./verify-doc-lib.mjs";

const SKIP_HREF = /^(https?:|mailto:|tel:|#)/i;

export async function checkLinks(ctx) {
  const issues = [];
  const check = "links";
  const pagesRoot = ctx.pagesDir;

  for (const file of ctx.linkFiles) {
    const text = await readText(file);
    const loc = rel(ctx.root, file);
    for (const { href } of mdLinks(text)) {
      const raw = href.split(/\s+/)[0].replace(/^<|>$/g, "");
      if (!raw || SKIP_HREF.test(raw)) continue;
      const noHash = raw.split("#")[0];
      if (!noHash) continue;

      if (noHash.startsWith("/")) {
        if (/\.(png|jpe?g|gif|svg|webp|ico|css|js)$/i.test(noHash)) continue;
        if (noHash.startsWith("/settings")) continue;
        const okMount = ctx.allowedHrefs?.has(noHash);
        const pageMdx = resolve(pagesRoot, noHash.replace(/^\//, "") + ".mdx");
        const pageMd = resolve(pagesRoot, noHash.replace(/^\//, "") + ".md");
        const pageIndex = resolve(pagesRoot, noHash.replace(/^\//, ""), "index.mdx");
        const mapped = mapPortalHref(noHash, pagesRoot);
        const found =
          okMount ||
          (await exists(pageMdx)) ||
          (await exists(pageMd)) ||
          (await exists(pageIndex)) ||
          (mapped && (await exists(mapped)));
        if (!found && !ctx.allowedHrefs?.has(noHash)) {
          issues.push(issue(check, loc, `broken portal path ${noHash}`));
        }
        continue;
      }

      const target = resolve(dirname(file), noHash);
      if (!(await exists(target))) {
        issues.push(issue(check, loc, `broken relative link ${noHash}`));
      }
    }
  }
  return issues;
}

function mapPortalHref(href, pagesRoot) {
  const aliases = {
    "/products/beans": "products/beans/overview.mdx",
    "/products/espresso": "products/espresso/overview.mdx",
    "/products/cortado": "products/cortado/overview.mdx",
    "/start": "start/overview.mdx",
    "/contact": "contact.mdx",
    "/api/overview": "api-overview.mdx",
  };
  const mapped = aliases[href];
  return mapped ? resolve(pagesRoot, mapped) : null;
}
