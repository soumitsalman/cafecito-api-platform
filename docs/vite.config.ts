import { createRequire } from "node:module";
import { defineConfig } from "zudoku/vite";

const require = createRequire(import.meta.url);

function tryResolve(id: string): string | null {
  try {
    return require.resolve(id);
  } catch {
    return null;
  }
}

// #region agent log
fetch("http://127.0.0.1:7457/ingest/9819a78f-6d80-4364-8f29-6bfca7332d3e", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "X-Debug-Session-Id": "8363eb",
  },
  body: JSON.stringify({
    sessionId: "8363eb",
    runId: "post-fix",
    hypothesisId: "A",
    location: "docs/vite.config.ts",
    message: "vite config loaded via zudoku/vite",
    data: {
      importer: "zudoku/vite",
      resolvedVite: tryResolve("vite"),
      resolvedZudokuVite: tryResolve("zudoku/vite"),
      cwd: process.cwd(),
    },
    timestamp: Date.now(),
  }),
}).catch(() => {});
// #endregion

export default defineConfig({
  server: {
    watch: {
      usePolling: true,
      interval: 1000,
    },
  },
});
