import { describe, it } from "@zuplo/test";
import { expect } from "chai";

const base =
  process.env.ZUPLO_TEST_ENDPOINT ??
  process.env.TEST_ENDPOINT ??
  "http://localhost:9000";

describe("health stays public", () => {
  it("GET /beans/health without Authorization is not 401", async () => {
    const res = await fetch(`${base}/beans/health`);
    expect(res.status).to.not.equal(401);
  });

  it("GET /espresso/health without Authorization is not 401", async () => {
    const res = await fetch(`${base}/espresso/health`);
    expect(res.status).to.not.equal(401);
  });
});

describe("REST and MCP require Bearer", () => {
  const protectedPaths: Array<{ method: string; path: string }> = [
    { method: "GET", path: "/beans/articles" },
    { method: "POST", path: "/beans/mcp" },
    { method: "GET", path: "/espresso/events" },
    { method: "POST", path: "/espresso/mcp" },
  ];

  for (const { method, path } of protectedPaths) {
    it(`${method} ${path} missing credentials returns 401`, async () => {
      const res = await fetch(`${base}${path}`, { method });
      expect(res.status).to.equal(401);
    });

    it(`${method} ${path} invalid credentials returns 401`, async () => {
      const res = await fetch(`${base}${path}`, {
        method,
        headers: { Authorization: "Bearer not-a-valid-key" },
      });
      expect(res.status).to.equal(401);
    });
  }
});
