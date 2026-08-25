import { readdir, readFile, stat } from "node:fs/promises";
import { dirname, extname, join, relative, resolve } from "node:path";

export const HTTP_METHODS = new Set([
  "get",
  "put",
  "post",
  "delete",
  "patch",
  "options",
  "head",
]);

export const LIFECYCLE_FIELDS = [
  "status",
  "authority",
  "audience",
  "last verified",
  "owner",
  "superseded by",
];

export const LIFECYCLE_STATUSES = new Set([
  "current",
  "target",
  "superseded",
  "historical",
  "generated",
]);

export const NOT_MCP_TOOLS = new Set([
  "healthCheck",
  "beansMcp",
  "espressoMcp",
]);

export const PORTAL_MOUNTS = new Set([
  "/api/beans",
  "/api/espresso",
  "/api/overview",
  "/account",
  "/settings",
  "/oauth/callback",
  "/404",
  "/contact",
  "/start",
  "/products/beans",
  "/products/espresso",
  "/products/cortado",
]);

export const STALE_PUBLIC_PATHS = [
  "/howtos/",
  "/howtos'",
  "](/howtos",
  "](/introduction)",
  "/start/troubleshooting",
  "/guides/cross-product-workflows",
];

/** Public-surface terms that violate AGENTS.md documentation boundary. */
export const FORBIDDEN_PUBLIC_PATTERNS = [
  { id: "cupboard", re: /\bcupboard\b/i, hint: "persistence architecture" },
  { id: "beansack", re: /\bbeansack\b/i, hint: "persistence architecture" },
  { id: "sips-table", re: /\bsips\b[^\n]{0,40}\b(table|schema|column)/i, hint: "sips table internals" },
  { id: "table-sips", re: /\b(table|schema)\b[^\n]{0,40}\bsips\b/i, hint: "sips table internals" },
  { id: "hnsw", re: /\bhnsw\b/i, hint: "retrieval internals" },
  {
    id: "embedding-vector",
    re: /\bembedding\s+vector\b|\bvector\(\d+\)|\bvector_cosine/i,
    hint: "embeddings internals",
  },
  { id: "embedder", re: /\bembedder(_model|_base_url|_api_key)?\b/i, hint: "embedder configuration" },
  { id: "x-api-key", re: /\bx-api-key\b/i, hint: "private backend header as public contract" },
  { id: "same_as", re: /\bsame_as\b/, hint: "relation storage values" },
  { id: "derived_from", re: /\bderived_from\b/, hint: "relation storage values" },
  { id: "create-table", re: /\bcreate\s+table\b/i, hint: "SQL/schema internals" },
  { id: "materialized-view", re: /\bmaterialized\s+view\b/i, hint: "SQL/schema internals" },
  { id: "foreign-key", re: /\bforeign\s+keys?\b/i, hint: "SQL/schema internals" },
  { id: "pg-trgm", re: /\bpg_trgm\b/i, hint: "SQL/schema internals" },
  { id: "ef_construction", re: /\bef_construction\b/i, hint: "HNSW internals" },
];

export const REQUIRED_GENERATED_SLUGS = [
  "/start",
  "/products/beans",
  "/products/espresso",
  "/guides/mcp-ai-agents",
  "/guides/api-conventions",
];

/** True when generated indexes mention the portal path (URL, .md publish, or llms.txt entry). */
export function generatedOutputHasSlug(blob, slug) {
  const text = String(blob).toLowerCase();
  const path = slug.startsWith("/") ? slug.toLowerCase() : `/${slug.toLowerCase()}`;
  const trimmed = path.replace(/^\//, "");
  const needles = [
    path,
    `${path}.md`,
    `${path}/`,
    `${trimmed}.md`,
    `](${path})`,
    `](${path}.md)`,
  ];
  return needles.some((n) => text.includes(n));
}

export const MAX_LIFECYCLE_AGE_DAYS = 180;

export async function walk(dir, filter) {
  const out = [];
  let entries;
  try {
    entries = await readdir(dir, { withFileTypes: true });
  } catch {
    return out;
  }
  for (const ent of entries) {
    const p = join(dir, ent.name);
    if (ent.isDirectory()) {
      if (ent.name === "node_modules" || ent.name === ".git") continue;
      out.push(...(await walk(p, filter)));
    } else if (!filter || filter(p, ent.name)) {
      out.push(p);
    }
  }
  return out;
}

export async function readText(path) {
  return readFile(path, "utf8");
}

export async function readJson(path) {
  return JSON.parse(await readText(path));
}

export async function exists(path) {
  try {
    await stat(path);
    return true;
  } catch {
    return false;
  }
}

export function rel(root, path) {
  return relative(root, path).split("\\").join("/");
}

export function issue(check, file, message) {
  return { check, file, message };
}

export function templateKey(path) {
  return path
    .replace(/\{[^}]+\}/g, "{p}")
    .replace(/:[A-Za-z0-9_]+/g, "{p}")
    .split("/")
    .map((seg) => (/^[A-Z][A-Z0-9_]{2,}$/.test(seg) ? "{p}" : seg))
    .join("/");
}

export function resolveRef(doc, node) {
  if (!node || typeof node !== "object") return node;
  if (typeof node.$ref === "string" && node.$ref.startsWith("#/")) {
    let cur = doc;
    for (const part of node.$ref.slice(2).split("/")) {
      cur = cur?.[part];
    }
    return cur;
  }
  return node;
}

export function loadOperations(oas) {
  const ops = [];
  const paths = oas?.paths ?? {};
  for (const [path, item] of Object.entries(paths)) {
    if (!item || typeof item !== "object") continue;
    for (const [method, op] of Object.entries(item)) {
      if (!HTTP_METHODS.has(method.toLowerCase())) continue;
      ops.push({
        method: method.toUpperCase(),
        path,
        operationId: op?.operationId ?? null,
        parameters: op?.parameters ?? [],
        pathItemParameters: item.parameters ?? [],
      });
    }
  }
  return ops;
}

export function mcpExportedIds(oas) {
  const ids = new Set();
  const visit = (node) => {
    if (!node || typeof node !== "object") return;
    const ops = node?.["x-zuplo-route"]?.handler?.options?.operations;
    if (Array.isArray(ops)) {
      for (const entry of ops) {
        if (typeof entry?.id === "string") ids.add(entry.id);
      }
    }
    if (Array.isArray(node)) {
      for (const v of node) visit(v);
      return;
    }
    for (const v of Object.values(node)) visit(v);
  };
  visit(oas?.paths ?? {});
  return ids;
}

export function queryParamNames(oas, operation) {
  const names = new Set();
  const list = [
    ...(operation.pathItemParameters ?? []),
    ...(operation.parameters ?? []),
  ];
  for (const raw of list) {
    const p = resolveRef(oas, raw);
    if (p?.in === "query" && typeof p.name === "string") names.add(p.name);
  }
  return names;
}

export function findOperation(ops, method, requestPath) {
  const m = method.toUpperCase();
  const want = templateKey(requestPath);
  return ops.find((op) => op.method === m && templateKey(op.path) === want);
}

export function extractJsonValues(text) {
  const found = [];
  for (let i = 0; i < text.length; i++) {
    if (text[i] !== "{") continue;
    let depth = 0;
    let inStr = false;
    let esc = false;
    for (let j = i; j < text.length; j++) {
      const c = text[j];
      if (inStr) {
        if (esc) esc = false;
        else if (c === "\\") esc = true;
        else if (c === '"') inStr = false;
        continue;
      }
      if (c === '"') inStr = true;
      else if (c === "{") depth++;
      else if (c === "}") {
        depth--;
        if (depth === 0) {
          const slice = text.slice(i, j + 1);
          try {
            found.push(JSON.parse(slice));
          } catch {
            /* not json */
          }
          i = j;
          break;
        }
      }
    }
  }
  return found;
}

export function codeFences(text) {
  const fences = [];
  const re = /```([^\n`]*)\n([\s\S]*?)```/g;
  let m;
  while ((m = re.exec(text))) {
    fences.push({ lang: (m[1] || "").trim().split(/\s+/)[0].toLowerCase(), body: m[2] });
  }
  const re2 = /~~~([^\n]*)\n([\s\S]*?)~~~/g;
  while ((m = re2.exec(text))) {
    fences.push({ lang: (m[1] || "").trim().split(/\s+/)[0].toLowerCase(), body: m[2] });
  }
  return fences;
}

export function scanForbidden(text) {
  const hits = [];
  for (const rule of FORBIDDEN_PUBLIC_PATTERNS) {
    if (rule.re.test(text)) hits.push(rule);
  }
  return hits;
}

export function parseLifecycle(markdown) {
  const fields = {};
  const tableRe =
    /^\|\s*([^|]+?)\s*\|\s*([^|]+?)\s*\|/gm;
  let m;
  while ((m = tableRe.exec(markdown))) {
    const key = m[1].trim().toLowerCase().replace(/\*\*/g, "");
    const value = m[2].trim().replace(/\*\*/g, "");
    if (key === "field" || key === "---") continue;
    fields[key] = value;
  }
  const boldRe = /\*\*([^*]+):\*\*\s*([^\n]+)/g;
  while ((m = boldRe.exec(markdown))) {
    const key = m[1].trim().toLowerCase();
    if (!fields[key]) fields[key] = m[2].trim();
  }
  if (fields["owner role"] && !fields.owner) fields.owner = fields["owner role"];
  if (fields["superseded-by"] && !fields["superseded by"]) {
    fields["superseded by"] = fields["superseded-by"];
  }
  return fields;
}

export function parseIsoDate(value) {
  const m = String(value || "").match(/(\d{4}-\d{2}-\d{2})/);
  if (!m) return null;
  const t = Date.parse(m[1]);
  return Number.isNaN(t) ? null : t;
}

export function daysSince(isoMs, nowMs) {
  return (nowMs - isoMs) / 86400000;
}

export function stripCodeFences(text) {
  return text.replace(/```[\s\S]*?```/g, " ").replace(/~~~[\s\S]*?~~~/g, " ");
}

export function mdLinks(text) {
  const links = [];
  const re = /\[([^\]]*)\]\(([^)]+)\)/g;
  let m;
  while ((m = re.exec(text))) {
    links.push({ label: m[1], href: m[2].trim() });
  }
  return links;
}

export function pageFileCandidates(pagesRoot, fileId) {
  const base = join(pagesRoot, fileId);
  return [`${base}.mdx`, `${base}.md`, base];
}

export function collectNavFiles(configText) {
  const files = [];
  const re = /file:\s*["']([^"']+)["']/g;
  let m;
  while ((m = re.exec(configText))) files.push(m[1]);
  return [...new Set(files)];
}

export function gatewayPathFromBruno(url, productPrefix) {
  let path = String(url || "");
  path = path.replace(/\{\{[^}]+\}\}/g, "");
  try {
    if (/^https?:\/\//i.test(path)) {
      path = new URL(path).pathname;
    }
  } catch {
    /* keep */
  }
  const q = path.indexOf("?");
  if (q >= 0) path = path.slice(0, q);
  path = path.trim();
  if (!path.startsWith("/")) path = `/${path}`;
  if (path === productPrefix || path.startsWith(`${productPrefix}/`)) return path;
  if (path === "/") return productPrefix;
  return `${productPrefix}${path}`;
}

export function brunoQueryNames(yaml) {
  const names = [];
  const block = yaml.split(/\n(?=examples:|docs:|settings:)/)[0] ?? yaml;
  const re =
    /-\s*name:\s*([A-Za-z0-9_]+)[\s\S]*?type:\s*(query|path)/g;
  let m;
  while ((m = re.exec(block))) {
    const chunk = m[0];
    if (/disabled:\s*true/.test(chunk) && m[2] === "query") {
      names.push({ name: m[1], in: m[2], disabled: true });
    } else {
      names.push({ name: m[1], in: m[2], disabled: false });
    }
  }
  return names;
}

export function brunoMethodUrl(yaml) {
  const method = yaml.match(/^\s*method:\s*(\w+)/m)?.[1];
  const url = yaml.match(/^\s*url:\s*(.+)$/m)?.[1]?.trim().replace(/^["']|["']$/g, "");
  return { method, url };
}

export function extractApiCalls(text) {
  const calls = [];
  const pathRe =
    /\b(GET|PUT|POST|DELETE|PATCH|HEAD|OPTIONS)\s+(\/(?:beans|espresso)\/[A-Za-z0-9_{}:./?-]*)/gi;
  let m;
  while ((m = pathRe.exec(text))) {
    const [path, query] = splitPathQuery(m[2]);
    calls.push({ method: m[1].toUpperCase(), path, query });
  }
  const curlPath =
    /(?:\$BASE_URL|https:\/\/api\.cafecito\.tech)(\/(?:beans|espresso)\/[A-Za-z0-9_{}:./?-]*)/gi;
  while ((m = curlPath.exec(text))) {
    const [path, query] = splitPathQuery(m[1]);
    const already = calls.some((c) => c.path === path && c.query.join("&") === query.join("&"));
    if (!already) calls.push({ method: "GET", path, query });
  }
  return calls;
}

export function splitPathQuery(raw) {
  const q = raw.indexOf("?");
  if (q < 0) return [raw.replace(/\/$/, "") || raw, []];
  const path = raw.slice(0, q).replace(/\/$/, "") || "/";
  const query = [];
  const qs = raw.slice(q + 1);
  for (const part of qs.split("&")) {
    if (!part) continue;
    query.push(decodeURIComponent(part.split("=")[0]));
  }
  return [path, query];
}

export function dataUrlEncodeNames(text) {
  const names = [];
  const re = /--data-urlencode\s+"([^=]+)=/g;
  let m;
  while ((m = re.exec(text))) names.push(m[1]);
  return names;
}

export function dirnameOf(file) {
  return dirname(file);
}

export { join, resolve, extname };
