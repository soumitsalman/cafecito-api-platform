import {
  NOT_MCP_TOOLS,
  brunoMethodUrl,
  brunoQueryNames,
  codeFences,
  dataUrlEncodeNames,
  extractApiCalls,
  extractJsonValues,
  findOperation,
  gatewayPathFromBruno,
  issue,
  mcpExportedIds,
  queryParamNames,
  readText,
  rel,
  walk,
} from "./verify-doc-lib.mjs";

const SKIP_PATHS = new Set(["/beans/mcp", "/espresso/mcp"]);

function envelopeIssues(obj, product, file, check) {
  const issues = [];
  if (!obj || typeof obj !== "object" || Array.isArray(obj)) return issues;
  if (!obj.pagination || typeof obj.pagination !== "object") return issues;
  const p = obj.pagination;
  for (const field of ["limit", "num_results", "next_cursor"]) {
    if (!(field in p)) {
      issues.push(
        issue(check, file, `collection pagination missing ${field}`),
      );
    }
  }
  if (product === "beans" && "cursor" in p) {
    issues.push(
      issue(
        check,
        file,
        "Beans collection JSON must not serialize pagination.cursor",
      ),
    );
  }
  if (product === "espresso" && "cursor" in p) {
    issues.push(
      issue(
        check,
        file,
        "Espresso collection JSON must not serialize pagination.cursor",
      ),
    );
  }
  return issues;
}

function productOfPath(path, fallback) {
  if (path.startsWith("/beans")) return "beans";
  if (path.startsWith("/espresso")) return "espresso";
  return fallback;
}

export async function checkExamples(ctx) {
  const issues = [];
  const check = "examples";
  const { beansOps, espressoOps, beansOas, espressoOas } = ctx;
  const beansMcp = mcpExportedIds(beansOas);
  const espressoMcp = mcpExportedIds(espressoOas);

  for (const file of ctx.mdxFiles) {
    const text = await readText(file);
    const loc = rel(ctx.root, file);
    const fallback = loc.includes("/espresso/") ? "espresso" : loc.includes("/beans/") ? "beans" : null;
    const fences = codeFences(text);
    const chunks = fences.length ? fences.map((f) => f.body) : [text];
    for (const chunk of chunks) {
      const split = chunk.split(/(?=^\s*curl\b)/m);
      const units = split.filter((s) => s.trim()).length
        ? split.filter((s) => s.trim())
        : [chunk];
      for (const unit of units) {
      const calls = extractApiCalls(unit);
      const extraQs = dataUrlEncodeNames(unit);
      for (const call of calls) {
        if (SKIP_PATHS.has(call.path)) continue;
        const product = productOfPath(call.path, fallback);
        const ops = product === "espresso" ? espressoOps : beansOps;
        const oas = product === "espresso" ? espressoOas : beansOas;
        const op = findOperation(ops, call.method, call.path);
        if (!op) {
          issues.push(
            issue(
              check,
              loc,
              `unknown ${call.method} ${call.path} (not in gateway OpenAPI)`,
            ),
          );
          continue;
        }
        const allowed = queryParamNames(oas, op);
        for (const q of [...call.query, ...extraQs]) {
          if (!q || q === "YOUR_API_KEY") continue;
          if (!allowed.has(q)) {
            issues.push(
              issue(
                check,
                loc,
                `query parameter ${q} is not accepted on ${call.method} ${call.path}`,
              ),
            );
          }
        }
      }
      }
    }

    if (loc.includes("mcp-ai-agents")) {
      const tableTools = [...text.matchAll(/^\|\s*`([a-z][a-zA-Z0-9]+)`\s*\|/gm)].map(
        (m) => m[1],
      );
      for (const tool of tableTools) {
        if (NOT_MCP_TOOLS.has(tool)) {
          issues.push(issue(check, loc, `${tool} must not be documented as an MCP tool`));
          continue;
        }
        const ok = beansMcp.has(tool) || espressoMcp.has(tool);
        if (!ok) {
          issues.push(issue(check, loc, `MCP tool ${tool} is not exported in gateway OAS`));
        }
      }
    }

    const jsonLangs = new Set(["json", "jsonc", ""]);
    for (const fence of fences) {
      if (fence.lang && !jsonLangs.has(fence.lang) && fence.lang !== "json") continue;
      if (fence.lang && !["json", "jsonc"].includes(fence.lang)) continue;
      for (const obj of extractJsonValues(fence.body)) {
        issues.push(
          ...envelopeIssues(obj, fallback ?? productOfPath("", fallback), loc, check),
        );
      }
    }
  }

  for (const { dir, prefix, product, ops, oas } of ctx.brunoCollections) {
    const files = await walk(dir, (p, name) => name.endsWith(".yml") && name !== "folder.yml" && name !== "opencollection.yml" && !name.includes("environments"));
    for (const file of files) {
      const yaml = await readText(file);
      const loc = rel(ctx.root, file);
      if (!/^\s*method:/m.test(yaml) && !/^http:/m.test(yaml)) continue;
      const { method, url } = brunoMethodUrl(yaml);
      if (!method || !url) continue;
      let path = gatewayPathFromBruno(url, prefix);
      if (path === "/beans/top-headlines") path = "/beans/articles/top-headlines";
      if (SKIP_PATHS.has(path)) continue;
      const op = findOperation(ops, method, path);
      if (!op) {
        issues.push(issue(check, loc, `Bruno ${method} ${path} is not in gateway OpenAPI`));
        continue;
      }
      const allowed = queryParamNames(oas, op);
      const documentsBadRequest = /^\s*status:\s*400\b/m.test(yaml);
      for (const q of brunoQueryNames(yaml)) {
        if (q.in !== "query") continue;
        if (!allowed.has(q.name) && !documentsBadRequest) {
          issues.push(
            issue(check, loc, `Bruno query ${q.name} is not accepted on ${method} ${path}`),
          );
        }
      }
      for (const obj of extractJsonValues(yaml)) {
        issues.push(...envelopeIssues(obj, product, loc, check));
      }
    }
  }

  return issues;
}
