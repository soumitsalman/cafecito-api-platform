import { join } from "node:path";
import { exists, issue, readText, rel } from "./verify-doc-lib.mjs";

const POSITIONING = {
  beans:
    "Beans is a publisher-content API for news, blogs, financial and earnings reports, litigation and lawsuits, official statements, research, technical documents, and related coverage context.",
  espresso:
    "Espresso is a market and business intelligence API for discovering market actions, signals, and tracing concrete evidence.",
};

const SITE_METADATA = {
  description:
    "Cafecito developer docs for Beans news and publisher-content APIs, Espresso market, business, event, and news intelligence APIs, and their MCP servers.",
  keywords: [
    "news API",
    "publisher content API",
    "article search API",
    "earnings reports API",
    "litigation monitoring",
    "market intelligence API",
    "business intelligence API",
    "event intelligence API",
    "news intelligence API",
    "MCP servers",
    "AI agents",
  ],
};

const TARGETS = [
  {
    file: "docs/pages/start/overview.mdx",
    required: [POSITIONING.beans, POSITIONING.espresso, POSITIONING.boundary],
  },
  {
    file: "docs/pages/api-overview.mdx",
    required: [POSITIONING.beans, POSITIONING.espresso, POSITIONING.boundary],
  },
  {
    file: "docs/pages/products/beans/overview.mdx",
    required: [POSITIONING.beans],
  },
  {
    file: "docs/pages/products/espresso/overview.mdx",
    required: [POSITIONING.espresso, POSITIONING.boundary],
  },
  {
    file: "docs/pages/products/espresso/migration.mdx",
    required: ["Espresso is a market, business, event, and news intelligence API", POSITIONING.boundary],
  },
  {
    file: "docs/pages/guides/mcp-ai-agents.mdx",
    required: [POSITIONING.espresso, POSITIONING.boundary],
  },
  {
    file: "docs/README.md",
    required: [POSITIONING.beans, POSITIONING.espresso, POSITIONING.boundary],
  },
  {
    file: "docs/zudoku.config.tsx",
    required: [
      SITE_METADATA.description,
      "canonicalUrlOrigin",
      'robots: "index,follow"',
      'creator: "Cafecito"',
      'publisher: "Cafecito"',
      ...SITE_METADATA.keywords,
    ],
  },
  {
    file: "config/beans.oas.json",
    required: [POSITIONING.beans],
  },
  {
    file: "config/espresso.oas.json",
    required: [POSITIONING.espresso, POSITIONING.boundary],
  },
];

export async function checkPositioning(ctx) {
  const issues = [];
  const check = "positioning";

  for (const target of TARGETS) {
    let path = join(ctx.root, target.file);
    if (!(await exists(path)) && target.file.startsWith("docs/")) {
      path = join(ctx.root, target.file.slice("docs/".length));
    }
    if (!(await exists(path))) {
      issues.push(issue(check, target.file, "positioning entry point is missing"));
      continue;
    }
    const text = await readText(path);
    for (const needle of target.required) {
      if (!text.includes(needle)) {
        issues.push(issue(check, target.file, `missing canonical positioning: ${needle}`));
      }
    }
  }

  const robotsPath = join(ctx.root, "docs/public/robots.txt");
  if (!(await exists(robotsPath))) {
    issues.push(issue(check, rel(ctx.root, robotsPath), "robots.txt source is missing"));
  } else {
    const robots = await readText(robotsPath);
    for (const needle of ["User-agent: *", "Allow: /", "Sitemap: https://developer.cafecito.tech/sitemap.xml"]) {
      if (!robots.includes(needle)) {
        issues.push(issue(check, rel(ctx.root, robotsPath), `robots.txt is missing ${needle}`));
      }
    }
  }

  return issues;
}
