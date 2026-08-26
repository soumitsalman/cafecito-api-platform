#!/usr/bin/env node
/**
 * Documentation contract verifier (Wave 4/5).
 *
 * Checks:
 *   examples   MDX/Bruno paths, query params, collection envelopes, MCP tools
 *   terms      forbidden public internals and stale portal paths
 *   lifecycle  apis/design metadata (status, authority, audience, last verified, owner, superseded-by)
 *   links      relative and portal markdown links
 *   inventory  docs/pages vs zudoku navigation
 *   generated  llms.txt / published Markdown required slugs and forbidden terms
 *
 * Usage (local):
 *   node scripts/verify-documentation.mjs
 *   node scripts/verify-documentation.mjs --self-test
 *
 * After a portal build, generated output exists under docs/dist (or docs/build).
 * Local default skips the generated check when those dirs are empty.
 *
 * CI (after the docs build, from the repo root):
 *   (cd docs && npm ci && npm run build)
 *   node scripts/verify-documentation.mjs --require-generated
 *
 * --require-generated fails closed if llms.txt/Markdown is missing or omits
 * required portal slugs (/start, /products/beans, /products/espresso,
 * /guides/mcp-ai-agents, /guides/api-conventions).
 *
 * ci_governance wires: npm run verify:docs
 */

import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { checkExamples } from "./verify-doc-examples.mjs";
import { checkGenerated } from "./verify-doc-generated.mjs";
import { checkInventory } from "./verify-doc-inventory.mjs";
import { checkLifecycle } from "./verify-doc-lifecycle.mjs";
import { checkLinks } from "./verify-doc-links.mjs";
import { checkPositioning } from "./verify-doc-positioning.mjs";
import { checkTerms } from "./verify-doc-terms.mjs";
import {
  PORTAL_MOUNTS,
  exists,
  loadOperations,
  readJson,
  readText,
  walk,
} from "./verify-doc-lib.mjs";

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const SELF_TEST = process.argv.includes("--self-test");
const REQUIRE_GENERATED = process.argv.includes("--require-generated");

async function liveContext() {
  const beansOas = await readJson(join(ROOT, "config/beans.oas.json"));
  const espressoOas = await readJson(join(ROOT, "config/espresso.oas.json"));
  const pagesDir = join(ROOT, "docs/pages");
  const zudokuPath = join(ROOT, "docs/zudoku.config.tsx");
  const mdxFiles = await walk(pagesDir, (p, n) => n.endsWith(".mdx") || n.endsWith(".md"));
  const zudokuText = await readText(zudokuPath);
  const allowedHrefs = new Set(PORTAL_MOUNTS);
  for (const m of zudokuText.matchAll(/path:\s*["']([^"']+)["']/g)) {
    allowedHrefs.add(m[1]);
  }
  for (const m of zudokuText.matchAll(/to:\s*["']([^"']+)["']/g)) {
    allowedHrefs.add(m[1]);
  }
  const publicFiles = [...mdxFiles];

  const generatedDirs = [
    join(ROOT, "docs/dist"),
    join(ROOT, "docs/dist/client"),
    join(ROOT, "docs/build"),
    join(ROOT, "docs/.zudoku"),
    join(ROOT, "dist"),
  ];
  if (process.env.ZUDOKU_OUTPUT_DIR) {
    generatedDirs.unshift(resolve(ROOT, process.env.ZUDOKU_OUTPUT_DIR));
  }

  return {
    root: ROOT,
    pagesDir,
    designDir: join(ROOT, "apis/design"),
    zudokuPath,
    mdxFiles,
    publicFiles,
    linkFiles: [
      ...mdxFiles,
      join(ROOT, "docs/README.md"),
      join(ROOT, "docs/CONTRIBUTING.md"),
      join(ROOT, "apis/design/README.md"),
    ],
    allowedHrefs,
    beansOas,
    espressoOas,
    beansOps: loadOperations(beansOas),
    espressoOps: loadOperations(espressoOas),
    brunoCollections: [
      {
        dir: join(ROOT, "apis/beans/bruno"),
        prefix: "/beans",
        product: "beans",
        ops: loadOperations(beansOas),
        oas: beansOas,
      },
      {
        dir: join(ROOT, "apis/espresso/bruno"),
        prefix: "/espresso",
        product: "espresso",
        ops: loadOperations(espressoOas),
        oas: espressoOas,
      },
    ],
    generatedDirs,
    requireGenerated: REQUIRE_GENERATED,
    nowMs: Date.now(),
  };
}

async function runAll(ctx) {
  const examples = await checkExamples(ctx);
  const terms = await checkTerms(ctx);
  const lifecycle = await checkLifecycle(ctx);
  const links = await checkLinks(ctx);
  const inventory = await checkInventory(ctx);
  const positioning = await checkPositioning(ctx);
  const generated = await checkGenerated(ctx);
  return {
    examples,
    terms,
    lifecycle,
    links,
    inventory,
    positioning,
    generated: generated.issues,
    generatedSkipped: generated.skipped,
  };
}

function printIssues(groups) {
  let n = 0;
  for (const [name, list] of Object.entries(groups)) {
    if (name === "generatedSkipped") continue;
    if (!list?.length) continue;
    console.error(`\n[${name}] ${list.length} issue(s)`);
    for (const i of list) {
      n++;
      console.error(`  - ${i.file}: ${i.message}`);
    }
  }
  return n;
}

async function selfTest() {
  const fixtureRoot = join(ROOT, "scripts/fixtures/docs");
  const spec = await readJson(join(fixtureRoot, "cases.json"));
  const cases = spec.cases;
  const beansOas = await readJson(join(ROOT, "config/beans.oas.json"));
  const espressoOas = await readJson(join(ROOT, "config/espresso.oas.json"));
  let failed = 0;

  for (const c of cases) {
    const dir = join(fixtureRoot, c.dir);
    const ctx = await fixtureContext(dir, c, beansOas, espressoOas);
    let list = [];
    if (c.check === "examples") list = await checkExamples(ctx);
    else if (c.check === "terms") list = await checkTerms(ctx);
    else if (c.check === "lifecycle") list = await checkLifecycle(ctx);
    else if (c.check === "links") list = await checkLinks(ctx);
    else if (c.check === "inventory") list = await checkInventory(ctx);
    else if (c.check === "positioning") list = await checkPositioning(ctx);
    else if (c.check === "generated") list = (await checkGenerated(ctx)).issues;
    const all = list;
    const matched = all.filter((i) => i.message.includes(c.expect));
    if (matched.length === 0) {
      failed++;
      console.error(
        `self-test FAIL ${c.id}: expected ${c.check} issue containing ${JSON.stringify(c.expect)}`,
      );
      if (all.length) {
        console.error("  got:");
        for (const i of all) console.error(`    [${i.check}] ${i.file}: ${i.message}`);
      } else {
        console.error("  got no issues");
      }
    } else {
      console.log(`self-test PASS ${c.id}`);
    }
  }
  return failed;
}

async function fixtureContext(dir, c, beansOas, espressoOas) {
  const pagesDir = join(dir, "pages");
  const designDir = join(dir, "design");
  const zudokuPath = join(dir, "zudoku.config.tsx");
  const mdxFiles = (await exists(pagesDir))
    ? await walk(pagesDir, (p, n) => n.endsWith(".mdx") || n.endsWith(".md"))
    : [];
  const brunoDir = join(dir, "bruno");
  const publicFiles = [...mdxFiles];
  if (await exists(join(dir, "public.mdx"))) publicFiles.push(join(dir, "public.mdx"));
  const linkFiles = [...mdxFiles];
  if (await exists(join(dir, "links.mdx"))) linkFiles.push(join(dir, "links.mdx"));

  const beansOps = loadOperations(beansOas);
  const espressoOps = loadOperations(espressoOas);

  return {
    root: dir,
    pagesDir: (await exists(pagesDir)) ? pagesDir : dir,
    designDir: (await exists(designDir)) ? designDir : dir,
    zudokuPath: (await exists(zudokuPath)) ? zudokuPath : join(dir, "zudoku.config.tsx"),
    mdxFiles: mdxFiles.length ? mdxFiles : await walk(dir, (p, n) => n.endsWith(".mdx")),
    publicFiles: publicFiles.length ? publicFiles : await walk(dir, (p, n) => n.endsWith(".mdx") || n.endsWith(".json")),
    linkFiles: linkFiles.length ? linkFiles : await walk(dir, (p, n) => n.endsWith(".mdx") || n.endsWith(".md")),
    allowedHrefs: new Set(PORTAL_MOUNTS),
    beansOas,
    espressoOas,
    beansOps,
    espressoOps,
    brunoCollections: (await exists(brunoDir))
      ? [
          {
            dir: brunoDir,
            prefix: c.prefix || "/beans",
            product: c.product || "beans",
            ops: beansOps,
            oas: beansOas,
          },
        ]
      : [],
    generatedDirs: [join(dir, "dist")],
    requireGenerated: Boolean(c.requireGenerated),
    nowMs: Date.parse("2026-08-25"),
  };
}

if (SELF_TEST) {
  const failed = await selfTest();
  process.exit(failed === 0 ? 0 : 1);
}

const ctx = await liveContext();
const results = await runAll(ctx);
const n = printIssues(results);
if (results.generatedSkipped) {
  console.log(
    "generated: skipped (no llms.txt/Markdown output; build docs or pass --require-generated)",
  );
}
if (n > 0) {
  console.error(`\nverify-documentation: ${n} issue(s)`);
  process.exit(1);
}

const fixtureFails = await selfTest();
if (fixtureFails > 0) {
  console.error(`\nverify-documentation: ${fixtureFails} negative fixture(s) did not fail as expected`);
  process.exit(1);
}

console.log("verify-documentation: ok");
console.log("  examples, terms, lifecycle, links, inventory passed");
console.log("  generated skipped or passed");
console.log("  negative fixtures proved each check fails");
