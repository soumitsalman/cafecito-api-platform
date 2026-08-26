export const HTTP_METHODS = new Set([
  "get",
  "put",
  "post",
  "delete",
  "patch",
  "options",
  "head",
  "trace",
]);

export function withPrefix(path, prefix) {
  if (path === prefix || path.startsWith(`${prefix}/`)) return path;
  if (path === "/") return prefix;
  return `${prefix}${path}`;
}

export function applyRewrite(backendPath, prefix, rewrites) {
  const match = rewrites.find((ex) => ex.backend_path === backendPath);
  if (match?.gateway_path) return match.gateway_path;
  return withPrefix(backendPath, prefix);
}

export function resolveRef(doc, node, seen = new Set()) {
  if (!node || typeof node !== "object") return node;
  if (!node.$ref) return node;
  const ref = node.$ref;
  if (seen.has(ref)) return node;
  seen.add(ref);
  const parts = ref.replace(/^#\//, "").split("/");
  let cur = doc;
  for (const p of parts) {
    cur = cur?.[p];
    if (cur == null) return node;
  }
  const merged = { ...cur, ...Object.fromEntries(Object.entries(node).filter(([k]) => k !== "$ref")) };
  return resolveRef(doc, merged, seen);
}

function mergeAllOf(doc, schema) {
  const resolved = resolveRef(doc, schema);
  if (!resolved?.allOf) return resolved;
  const out = { type: "object", properties: {}, required: [] };
  for (const part of resolved.allOf) {
    const s = mergeAllOf(doc, part);
    Object.assign(out.properties, s?.properties ?? {});
    out.required.push(...(s?.required ?? []));
    if (s?.type) out.type = s.type;
  }
  Object.assign(out.properties, resolved.properties ?? {});
  out.required.push(...(resolved.required ?? []));
  out.required = [...new Set(out.required)];
  return out;
}

export function operationsFromPaths(paths = {}) {
  const ops = [];
  for (const [path, item] of Object.entries(paths)) {
    if (!item || typeof item !== "object") continue;
    for (const [method, op] of Object.entries(item)) {
      if (!HTTP_METHODS.has(method.toLowerCase())) continue;
      ops.push({
        method: method.toUpperCase(),
        path,
        operationId: op?.operationId ?? null,
        spec: op,
        pathItem: item,
      });
    }
  }
  return ops;
}

function paramLocationRequired(param) {
  if (typeof param.required === "boolean") return param.required;
  return param.in === "path";
}

export function normalizeParameter(doc, raw) {
  const param = resolveRef(doc, raw);
  const schema = resolveRef(doc, param.schema ?? {});
  const type = schema.type ?? param.type ?? null;
  const enumValues = schema.enum ?? param.enum ?? null;
  const defaultValue = schema.default ?? param.default;
  return {
    name: param.name,
    in: param.in,
    required: paramLocationRequired(param),
    type: Array.isArray(type) ? type.filter((t) => t !== "null").join("|") : type,
    enum: Array.isArray(enumValues) ? [...enumValues].sort() : null,
    default: defaultValue,
  };
}

export function operationParameters(doc, op) {
  const list = [...(op.pathItem?.parameters ?? []), ...(op.spec?.parameters ?? [])];
  const byKey = new Map();
  for (const raw of list) {
    const p = normalizeParameter(doc, raw);
    if (!p.name || !p.in) continue;
    byKey.set(`${p.in}:${p.name}`, p);
  }
  return byKey;
}

export function responseStatuses(doc, op) {
  const responses = op.spec?.responses ?? {};
  return Object.keys(responses)
    .filter((k) => k !== "default")
    .sort();
}

export function jsonContentSchema(doc, response) {
  const resolved = resolveRef(doc, response);
  if (resolved?.schema) return mergeAllOf(doc, resolved.schema);
  const content = resolved?.content ?? {};
  const json = content["application/json"] ?? Object.values(content)[0];
  if (!json) return null;
  return mergeAllOf(doc, json.schema);
}

export function envelopeShape(doc, op) {
  const ok = op.spec?.responses?.["200"] ?? op.spec?.responses?.["201"];
  if (!ok) return null;
  const schema = jsonContentSchema(doc, ok);
  if (!schema) return null;
  const properties = Object.keys(schema.properties ?? {}).sort();
  const required = [...(schema.required ?? [])].sort();
  let pagination = null;
  if (schema.properties?.pagination) {
    const p = mergeAllOf(doc, schema.properties.pagination);
    pagination = {
      properties: Object.keys(p?.properties ?? {}).sort(),
      required: [...(p?.required ?? [])].sort(),
    };
  }
  return { properties, required, pagination };
}

export function effectiveSecurity(doc, op) {
  const opSec = op.spec?.security;
  const rootSec = doc.security;
  const chosen = opSec !== undefined ? opSec : rootSec;
  if (!chosen || chosen.length === 0) return [];
  return chosen.flatMap((req) => Object.keys(req ?? {}));
}

export function securitySchemeNames(doc) {
  const oas = Object.keys(doc.components?.securitySchemes ?? {});
  const swagger = Object.keys(doc.securityDefinitions ?? {});
  return [...new Set([...oas, ...swagger])].sort();
}

export function mcpCatalog(doc) {
  const tools = [];
  const servers = [];
  for (const [path, item] of Object.entries(doc.paths ?? {})) {
    for (const [method, op] of Object.entries(item ?? {})) {
      if (!HTTP_METHODS.has(method.toLowerCase())) continue;
      const handler = op?.["x-zuplo-route"]?.handler;
      if (handler?.export !== "mcpServerHandler") continue;
      servers.push({ path, method: method.toUpperCase(), operationId: op.operationId });
      for (const entry of handler.options?.operations ?? []) {
        if (entry?.id) tools.push(entry.id);
      }
    }
  }
  return { tools: [...new Set(tools)].sort(), servers };
}

export function mcpEnabledTools(doc) {
  const tools = [];
  for (const item of Object.values(doc.paths ?? {})) {
    for (const [method, op] of Object.entries(item ?? {})) {
      if (!HTTP_METHODS.has(method.toLowerCase())) continue;
      const mcp = op?.["x-zuplo-route"]?.mcp;
      if (mcp?.enabled === false) continue;
      if (mcp?.type === "tool" || mcp?.name) {
        tools.push(mcp.name ?? op.operationId);
      }
    }
  }
  return [...new Set(tools)].sort();
}

export function publicServerUrl(doc) {
  return doc.servers?.[0]?.url ?? null;
}
