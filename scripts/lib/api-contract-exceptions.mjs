export const REQUIRED_EXCEPTION_FIELDS = [
  "id",
  "product",
  "kind",
  "intentional_difference",
  "owned_by",
  "reason",
  "reviewed_on",
  "review_due",
];

export const EXCEPTION_KINDS = new Set([
  "gateway_only_path",
  "backend_only_path",
  "path_rewrite",
  "health",
  "mcp",
  "docs",
  "prefix_mapping",
  "security_scheme",
  "extra_response_status",
  "parameter",
  "envelope",
]);

export function validateExceptions(doc) {
  const errors = [];
  if (!doc || typeof doc !== "object" || !Array.isArray(doc.exceptions)) {
    return ["exceptions file must be an object with an exceptions array"];
  }
  const ids = new Set();
  for (const [i, ex] of doc.exceptions.entries()) {
    const loc = `exceptions[${i}]`;
    for (const field of REQUIRED_EXCEPTION_FIELDS) {
      if (typeof ex?.[field] !== "string" || ex[field].trim() === "") {
        errors.push(`${loc}.${field} must be a non-empty string`);
      }
    }
    if (ex?.id) {
      if (ids.has(ex.id)) errors.push(`${loc}.id is duplicated: ${ex.id}`);
      ids.add(ex.id);
    }
    if (ex?.kind && !EXCEPTION_KINDS.has(ex.kind)) {
      errors.push(`${loc}.kind must be one of: ${[...EXCEPTION_KINDS].join(", ")}`);
    }
    if (ex?.reviewed_on && Number.isNaN(Date.parse(ex.reviewed_on))) {
      errors.push(`${loc}.reviewed_on must be an ISO date`);
    }
    if (ex?.review_due && Number.isNaN(Date.parse(ex.review_due))) {
      errors.push(`${loc}.review_due must be an ISO date`);
    }
  }
  return errors;
}

export function expiredExceptions(exceptions, now = new Date()) {
  return exceptions.filter((ex) => {
    const due = Date.parse(ex.review_due);
    return !Number.isNaN(due) && due < now.getTime();
  });
}

export function pathRewrites(exceptions, product) {
  return exceptions.filter((ex) => ex.product === product && ex.kind === "path_rewrite");
}

function fieldMatch(expected, actual) {
  if (expected == null || expected === "") return true;
  return expected === actual;
}

export function exceptionCovers(ex, mismatch) {
  if (ex.product !== mismatch.product) return false;
  if (ex.code && ex.code !== mismatch.code) return false;
  if (ex.method && mismatch.method && ex.method.toUpperCase() !== mismatch.method) return false;
  if (!fieldMatch(ex.backend_path, mismatch.backend_path)) return false;
  if (!fieldMatch(ex.gateway_path, mismatch.gateway_path)) return false;
  if (!fieldMatch(ex.operation_id, mismatch.operationId)) return false;
  if (!fieldMatch(ex.parameter, mismatch.parameter)) return false;
  if (ex.status && mismatch.status && String(ex.status) !== String(mismatch.status)) return false;

  if (ex.kind === "gateway_only_path") {
    return mismatch.code === "path.extra_on_gateway";
  }
  if (ex.kind === "backend_only_path") {
    return mismatch.code === "path.missing_from_gateway";
  }
  if (ex.kind === "path_rewrite") {
    return mismatch.code === "path.rewrite";
  }
  if (ex.kind === "mcp") {
    return String(mismatch.code).startsWith("mcp.");
  }
  if (ex.kind === "health") {
    return mismatch.gateway_path?.endsWith("/health") || mismatch.backend_path === "/health";
  }
  if (ex.kind === "security_scheme") {
    return String(mismatch.code).startsWith("security.");
  }
  if (ex.kind === "extra_response_status") {
    return mismatch.code === "response.extra_status";
  }
  if (ex.kind === "parameter") {
    return String(mismatch.code).startsWith("param.");
  }
  if (ex.kind === "envelope") {
    return String(mismatch.code).startsWith("envelope.");
  }
  if (ex.code) return true;
  return false;
}

export function consumeExceptions(exceptions, mismatches) {
  const used = new Set();
  const remaining = [];
  for (const mismatch of mismatches) {
    const hit = exceptions.find((ex) => exceptionCovers(ex, mismatch));
    if (hit) used.add(hit.id);
    else remaining.push(mismatch);
  }
  const unused = exceptions.filter((ex) => ex.kind !== "path_rewrite" && !used.has(ex.id));
  return { remaining, unused, usedIds: [...used] };
}
