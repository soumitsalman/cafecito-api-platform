import {
  applyRewrite,
  effectiveSecurity,
  envelopeShape,
  mcpCatalog,
  mcpEnabledTools,
  operationParameters,
  operationsFromPaths,
  publicServerUrl,
  responseStatuses,
  securitySchemeNames,
} from "./api-contract-spec.mjs";

function opKey(method, path) {
  return `${method} ${path}`;
}

function jsonEqual(a, b) {
  return JSON.stringify(a) === JSON.stringify(b);
}

function setDiff(a, b) {
  const bs = new Set(b);
  return a.filter((x) => !bs.has(x));
}

export function compareProduct(product, swagger, gateway, rewrites) {
  const mismatches = [];
  const push = (m) => mismatches.push({ product: product.name, ...m });

  const backendOps = operationsFromPaths(swagger.paths).map((op) => ({
    ...op,
    gatewayPath: applyRewrite(op.path, product.prefix, rewrites),
  }));
  const gatewayOps = operationsFromPaths(gateway.paths);

  const backendByGw = new Map(backendOps.map((op) => [opKey(op.method, op.gatewayPath), op]));
  const gatewayByPath = new Map(gatewayOps.map((op) => [opKey(op.method, op.path), op]));

  for (const [key, op] of backendByGw) {
    if (!gatewayByPath.has(key)) {
      push({
        code: "path.missing_from_gateway",
        method: op.method,
        backend_path: op.path,
        gateway_path: op.gatewayPath,
        operationId: op.operationId,
        detail: key,
      });
    }
  }
  for (const [key, op] of gatewayByPath) {
    if (!backendByGw.has(key)) {
      push({
        code: "path.extra_on_gateway",
        method: op.method,
        backend_path: null,
        gateway_path: op.path,
        operationId: op.operationId,
        detail: key,
      });
    }
  }

  const gwSchemes = securitySchemeNames(gateway);
  const swSchemes = securitySchemeNames(swagger);
  if (!gwSchemes.includes("BearerAuth")) {
    push({
      code: "security.missing_scheme",
      detail: "gateway must declare BearerAuth",
    });
  }
  if (!swSchemes.includes("BackendAPIKey")) {
    push({
      code: "security.missing_scheme",
      detail: "swagger must declare BackendAPIKey (X-API-KEY)",
    });
  }

  const server = publicServerUrl(gateway);
  if (!server) {
    push({
      code: "security.missing_server",
      detail: "gateway must declare servers[0].url",
    });
  }

  for (const [key, backend] of backendByGw) {
    const gw = gatewayByPath.get(key);
    if (!gw) continue;

    if (backend.operationId && gw.operationId && backend.operationId !== gw.operationId) {
      push({
        code: "operation.id_mismatch",
        method: backend.method,
        backend_path: backend.path,
        gateway_path: gw.path,
        operationId: backend.operationId,
        detail: `${backend.operationId} vs ${gw.operationId}`,
      });
    }

    const bParams = operationParameters(swagger, backend);
    const gParams = operationParameters(gateway, gw);
    const bNames = [...bParams.keys()].sort();
    const gNames = [...gParams.keys()].sort();
    for (const name of setDiff(bNames, gNames)) {
      push({
        code: "param.missing_from_gateway",
        method: backend.method,
        backend_path: backend.path,
        gateway_path: gw.path,
        operationId: backend.operationId,
        parameter: name,
      });
    }
    for (const name of setDiff(gNames, bNames)) {
      push({
        code: "param.extra_on_gateway",
        method: backend.method,
        backend_path: backend.path,
        gateway_path: gw.path,
        operationId: gw.operationId,
        parameter: name,
      });
    }
    for (const name of bNames) {
      const bp = bParams.get(name);
      const gp = gParams.get(name);
      if (!gp) continue;
      if (bp.required !== gp.required) {
        push({
          code: "param.required",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: `swagger required=${bp.required} gateway required=${gp.required}`,
        });
      }
      if (bp.enum && gp.enum && !jsonEqual(bp.enum, gp.enum)) {
        push({
          code: "param.enum",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: `swagger ${JSON.stringify(bp.enum)} gateway ${JSON.stringify(gp.enum)}`,
        });
      }
      if (bp.enum && !gp.enum) {
        push({
          code: "param.enum",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: "gateway missing enum",
        });
      }
      if (!bp.enum && gp.enum) {
        push({
          code: "param.enum",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: "swagger missing enum",
        });
      }
      if (bp.default !== undefined && gp.default !== undefined && !jsonEqual(bp.default, gp.default)) {
        push({
          code: "param.default",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: `swagger ${JSON.stringify(bp.default)} gateway ${JSON.stringify(gp.default)}`,
        });
      }
      if (bp.default !== undefined && gp.default === undefined) {
        push({
          code: "param.default",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: "gateway missing default",
        });
      }
      if (bp.default === undefined && gp.default !== undefined) {
        push({
          code: "param.default",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: "swagger missing default",
        });
      }
      if (bp.type && gp.type && bp.type !== gp.type) {
        push({
          code: "param.type",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          parameter: name,
          detail: `${bp.type} vs ${gp.type}`,
        });
      }
    }

    const bStatus = responseStatuses(swagger, backend);
    const gStatus = responseStatuses(gateway, gw);
    for (const status of setDiff(bStatus, gStatus)) {
      push({
        code: "response.missing_status",
        method: backend.method,
        backend_path: backend.path,
        gateway_path: gw.path,
        operationId: backend.operationId,
        status,
        detail: status,
      });
    }
    for (const status of setDiff(gStatus, bStatus)) {
      push({
        code: "response.extra_status",
        method: backend.method,
        backend_path: backend.path,
        gateway_path: gw.path,
        operationId: gw.operationId,
        status,
        detail: status,
      });
    }

    const bEnv = envelopeShape(swagger, backend);
    const gEnv = envelopeShape(gateway, gw);
    if (bEnv && gEnv) {
      if (!jsonEqual(bEnv.required, gEnv.required)) {
        push({
          code: "envelope.required",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          detail: `swagger ${JSON.stringify(bEnv.required)} gateway ${JSON.stringify(gEnv.required)}`,
        });
      }
      if (!jsonEqual(bEnv.properties, gEnv.properties)) {
        push({
          code: "envelope.properties",
          method: backend.method,
          backend_path: backend.path,
          gateway_path: gw.path,
          operationId: backend.operationId,
          detail: `swagger ${JSON.stringify(bEnv.properties)} gateway ${JSON.stringify(gEnv.properties)}`,
        });
      }
      if (bEnv.pagination && gEnv.pagination) {
        if (!jsonEqual(bEnv.pagination.properties, gEnv.pagination.properties)) {
          push({
            code: "envelope.pagination_properties",
            method: backend.method,
            backend_path: backend.path,
            gateway_path: gw.path,
            operationId: backend.operationId,
            detail: `swagger ${JSON.stringify(bEnv.pagination.properties)} gateway ${JSON.stringify(gEnv.pagination.properties)}`,
          });
        }
        if (!jsonEqual(bEnv.pagination.required, gEnv.pagination.required)) {
          push({
            code: "envelope.pagination_required",
            method: backend.method,
            backend_path: backend.path,
            gateway_path: gw.path,
            operationId: backend.operationId,
            detail: `swagger ${JSON.stringify(bEnv.pagination.required)} gateway ${JSON.stringify(gEnv.pagination.required)}`,
          });
        }
      }
    }

    const bSec = effectiveSecurity(swagger, backend);
    const gSec = effectiveSecurity(gateway, gw);
    const bAuth = bSec.includes("BackendAPIKey");
    const gAuth = gSec.includes("BearerAuth");
    if (bAuth !== gAuth) {
      push({
        code: "security.requirement",
        method: backend.method,
        backend_path: backend.path,
        gateway_path: gw.path,
        operationId: backend.operationId,
        detail: `swagger auth=${bAuth} gateway auth=${gAuth}`,
      });
    }
  }

  const catalog = mcpCatalog(gateway);
  const enabled = mcpEnabledTools(gateway);
  const gwIds = new Set(gatewayOps.map((op) => op.operationId).filter(Boolean));
  const swIds = new Set(backendOps.map((op) => op.operationId).filter(Boolean));

  for (const id of catalog.tools) {
    if (!gwIds.has(id)) {
      push({
        code: "mcp.unknown_operation",
        operationId: id,
        detail: `MCP catalog id ${id} is not a gateway operationId`,
      });
    }
    if (!swIds.has(id)) {
      push({
        code: "mcp.missing_from_backend",
        operationId: id,
        detail: `MCP catalog id ${id} is not a backend operationId`,
      });
    }
  }
  for (const name of enabled) {
    if (catalog.tools.length && !catalog.tools.includes(name)) {
      push({
        code: "mcp.enabled_not_exported",
        operationId: name,
        detail: `enabled MCP tool ${name} missing from mcpServerHandler operations`,
      });
    }
  }

  return {
    backendOperations: backendOps.length,
    gatewayOperations: gatewayOps.length,
    mcpTools: catalog.tools,
    server: server,
    mismatches,
  };
}
