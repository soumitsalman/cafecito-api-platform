#!/usr/bin/env node
/**
 * Fail-closed Swagger 2 vs gateway OpenAPI 3 contract cascade.
 *
 * Normalizes backend paths with /beans or /espresso (and reviewed rewrites),
 * then compares operations, parameters, defaults, enums, required flags,
 * response statuses, envelopes, security, and MCP catalogs.
 *
 * Portal/example checks belong in scripts/verify-documentation.mjs.
 *
 * Usage:
 *   node scripts/verify-api-contract-cascade.mjs
 *   node scripts/verify-api-contract-cascade.mjs --strict
 */

import { readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { compareProduct } from "./lib/api-contract-compare.mjs";
import {
  consumeExceptions,
  expiredExceptions,
  pathRewrites,
  validateExceptions,
} from "./lib/api-contract-exceptions.mjs";

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const EXCEPTIONS_REL = "config/api-contract-exceptions.json";
const STRICT = process.argv.includes("--strict");

const PRODUCTS = [
  {
    name: "beans",
    prefix: "/beans",
    swagger: "apis/beans/docs/swagger.json",
    gateway: "config/beans.oas.json",
  },
  {
    name: "espresso",
    prefix: "/espresso",
    swagger: "apis/espresso/docs/swagger.json",
    gateway: "config/espresso.oas.json",
  },
];

async function readJson(relPath) {
  return JSON.parse(await readFile(resolve(ROOT, relPath), "utf8"));
}

async function main() {
  const exceptionsDoc = await readJson(EXCEPTIONS_REL);
  const exceptionErrors = validateExceptions(exceptionsDoc);
  if (exceptionErrors.length > 0) {
    console.error(`Invalid ${EXCEPTIONS_REL}:`);
    for (const err of exceptionErrors) console.error(`  - ${err}`);
    process.exit(1);
  }

  const expired = expiredExceptions(exceptionsDoc.exceptions);
  const reports = [];
  const allRemaining = [];
  const usedIds = new Set();

  for (const product of PRODUCTS) {
    const swagger = await readJson(product.swagger);
    const gateway = await readJson(product.gateway);
    const rewrites = pathRewrites(exceptionsDoc.exceptions, product.name);
    const result = compareProduct(product, swagger, gateway, rewrites);
    const { remaining, usedIds: consumed } = consumeExceptions(
      exceptionsDoc.exceptions,
      result.mismatches,
    );
    for (const id of consumed) usedIds.add(id);
    reports.push({
      product: product.name,
      prefix: product.prefix,
      backendOperations: result.backendOperations,
      gatewayOperations: result.gatewayOperations,
      mcpTools: result.mcpTools,
      server: result.server,
      mismatchCount: remaining.length,
      mismatches: remaining,
    });
    allRemaining.push(...remaining);
  }

  const inventory = {
    generated: true,
    scaffold: false,
    strict: STRICT,
    exceptionsFile: EXCEPTIONS_REL,
    remainingUnexplained: allRemaining.length,
    products: reports,
  };
  console.log(JSON.stringify(inventory, null, 2));

  if (allRemaining.length) {
    console.error(`\nUnexplained contract diffs (${allRemaining.length}):`);
    for (const m of allRemaining) {
      console.error(
        `  [${m.product}] ${m.code} ${m.method ?? ""} ${m.gateway_path ?? m.backend_path ?? ""} ${m.parameter ?? m.status ?? m.operationId ?? ""} ${m.detail ?? ""}`.trim(),
      );
    }
  } else {
    console.error("\nNo unexplained Swagger-to-gateway diffs.");
  }

  const failReasons = [];
  if (allRemaining.length) failReasons.push(`${allRemaining.length} unexplained diffs`);
  if (expired.length) failReasons.push(`${expired.length} expired exceptions`);
  const unused = exceptionsDoc.exceptions.filter(
    (ex) => ex.kind !== "path_rewrite" && !usedIds.has(ex.id),
  );
  if (unused.length) {
    console.error(`Unused exceptions: ${unused.map((e) => e.id).join(", ")}`);
    failReasons.push(`${unused.length} unused exceptions`);
  }
  if (expired.length) {
    console.error(`Expired exceptions: ${expired.map((e) => e.id).join(", ")}`);
  }

  if (STRICT && failReasons.length) {
    console.error(`verify-api-contract-cascade failed: ${failReasons.join("; ")}`);
    process.exit(1);
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
