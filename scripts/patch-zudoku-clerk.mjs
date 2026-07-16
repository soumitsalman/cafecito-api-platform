import { readFile, writeFile } from "node:fs/promises";
import path from "node:path";

const providerPath = path.resolve(
  "node_modules/zudoku/src/lib/authentication/providers/clerk.tsx",
);

const originalClerkType = `type Clerk = {
  session: ClerkSession | null | undefined;
  load: () => Promise<void>;
  signOut: (opts: { redirectUrl: string }) => Promise<void>;
  redirectToSignIn: (opts: ClerkRedirectOptions) => Promise<void>;
  redirectToSignUp: (opts: ClerkRedirectOptions) => Promise<void>;
};`;

const patchedClerkType = `type ClerkUIConstructor = unknown;

type Clerk = {
  session: ClerkSession | null | undefined;
  load: (options?: {
    ui?: { ClerkUI: ClerkUIConstructor };
  }) => Promise<void>;
  signOut: (opts: { redirectUrl: string }) => Promise<void>;
  redirectToSignIn: (opts: ClerkRedirectOptions) => Promise<void>;
  redirectToSignUp: (opts: ClerkRedirectOptions) => Promise<void>;
};`;

const originalLoader = `  clerkPromise = new Promise<void>((resolve, reject) => {
    const frontendApiUrl = getClerkFrontendApi(publishableKey);

    const script = document.createElement("script");
    script.src = \`https://\${frontendApiUrl}/npm/@clerk/clerk-js@6/dist/clerk.browser.js\`;
    script.async = true;
    script.crossOrigin = "anonymous";
    script.dataset.clerkPublishableKey = publishableKey;
    script.onload = () => resolve();
    script.onerror = () => {
      clerkPromise = undefined;
      reject(new Error("Failed to load Clerk from CDN"));
    };
    document.head.appendChild(script);
  }).then(async () => {
    const clerk = (window as { Clerk?: Clerk }).Clerk;
    if (!clerk) {
      throw new Error("Clerk script loaded but window.Clerk is not available");
    }
    await clerk.load();
    return clerk;
  });`;

const patchedLoader = `  clerkPromise = new Promise<void>((resolve, reject) => {
    const frontendApiUrl = getClerkFrontendApi(publishableKey);
    const uiScript = document.createElement("script");
    uiScript.src = \`https://\${frontendApiUrl}/npm/@clerk/ui@1/dist/ui.browser.js\`;
    uiScript.async = true;
    uiScript.crossOrigin = "anonymous";
    uiScript.onload = () => {
      const script = document.createElement("script");
      script.src = \`https://\${frontendApiUrl}/npm/@clerk/clerk-js@6/dist/clerk.browser.js\`;
      script.async = true;
      script.crossOrigin = "anonymous";
      script.dataset.clerkPublishableKey = publishableKey;
      script.onload = () => resolve();
      script.onerror = () => {
        clerkPromise = undefined;
        reject(new Error("Failed to load Clerk from CDN"));
      };
      document.head.appendChild(script);
    };
    uiScript.onerror = () => {
      clerkPromise = undefined;
      reject(new Error("Failed to load Clerk UI from CDN"));
    };

    document.head.appendChild(uiScript);
  }).then(async () => {
    const clerk = (window as { Clerk?: Clerk }).Clerk;
    if (!clerk) {
      throw new Error("Clerk script loaded but window.Clerk is not available");
    }
    const clerkUI = (
      window as { __internal_ClerkUICtor?: ClerkUIConstructor }
    ).__internal_ClerkUICtor;
    if (!clerkUI) {
      throw new Error("Clerk UI script loaded but ClerkUI is not available");
    }
    await clerk.load({ ui: { ClerkUI: clerkUI } });
    return clerk;
  });`;

const source = await readFile(providerPath, "utf8");

if (
  source.includes("@clerk/ui@1/dist/ui.browser.js") &&
  source.includes("__internal_ClerkUICtor")
) {
  console.log("Zudoku Clerk UI compatibility patch already applied.");
} else {
  if (!source.includes(originalClerkType) || !source.includes(originalLoader)) {
    throw new Error(
      "Unsupported Zudoku Clerk provider source; update the compatibility patch.",
    );
  }

  await writeFile(
    providerPath,
    source
      .replace(originalClerkType, patchedClerkType)
      .replace(originalLoader, patchedLoader),
  );
  console.log("Applied Zudoku Clerk UI compatibility patch.");
}
