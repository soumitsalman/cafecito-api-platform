import type { ZudokuConfig } from "zudoku";
import AccountPage from "./src/AccountPage";

const serverUrl =
  process.env.ZUDOKU_PUBLIC_GATEWAY_URL || import.meta.env.ZUPLO_SERVER_URL;

const clerkPubKey =
  process.env.ZUDOKU_PUBLIC_CLERK_PUBLISHABLE_KEY ||
  process.env.CLERK_PUBLISHABLE_KEY;
const clerkJwtTemplateName =
  process.env.ZUDOKU_PUBLIC_CLERK_JWT_TEMPLATE_NAME ||
  process.env.CLERK_JWT_TEMPLATE_NAME ||
  "dev-portal";
/**
 * Developer Portal Configuration
 * For more information, see:
 * https://zuplo.com/docs/dev-portal/zudoku/configuration/overview
 */
const config: ZudokuConfig = {
  canonicalUrlOrigin: "https://developer.cafecito.tech",
  sitemap: {
    siteUrl: "https://developer.cafecito.tech",
    changefreq: "weekly",
    exclude: ["/account", "/settings", "/oauth/callback", "/404"],
  },
  site: {
    title: "Cafecito Developer Portal",
    logo: {
      src: {
        light: "/cafecito-light-v2.png",
        dark: "/cafecito-dark-v2.png",
      },
      alt: "Project Cafecito Logo",
      width: "60px",
    },    
  },
  search: {
    type: "pagefind",
    maxSubResults: 3,
    ranking: {
      termFrequency: 0.7,
      pageLength: 0.5,
      termSimilarity: 1.2,
      termSaturation: 1.1,
    },
  },
  docs: {    
    publishMarkdown: true,
    llms: {
      llmsTxt: true,
      llmsTxtFull: true,
    },
  },
  metadata: {
    title: "%s | Cafecito Developer Portal",
    description:
      "Developer docs for Cafecito REST APIs and MCP servers. Beans is a news and publisher-content API. Espresso is a market, business, event, and news intelligence API. One API key; Cortado and Latte are future products.",
    applicationName: "Cafecito Developer Portal",
    logo: "https://developer.cafecito.tech/cafecito-banner.png",
    favicon: "/cafecito-dark.png",
    keywords: [
      "news API",
      "blog API",
      "publisher content API",
      "article search API",
      "earnings reports API",
      "financial reports API",
      "lawsuit monitoring",
      "litigation news",
      "press releases",
      "official statements",
      "research papers",
      "technical documents",
      "story clustering",
      "news MCP",
      "market intelligence API",
      "business intelligence API",
      "event intelligence API",
      "news intelligence API",
      "event monitoring",
      "company monitoring",
      "impact",
      "outlook",
      "evidence",
      "provenance",
      "GDELT comparison",
      "Perigon comparison",
      "MCP tools",
      "AI research workflows",
    ],
  },
  theme: {
    light: {
      background: "#fafafa",
      foreground: "#0b0b0b",
      card: "#fffaf0",
      cardForeground: "#0b0b0b",
      popover: "#fff7ea",
      popoverForeground: "#0b0b0b",
      primary: "#A66A2A",           // caramel
      primaryForeground: "#ffffff",
      secondary: "#f3e6d0",         // latte milk
      secondaryForeground: "#0b0b0b",
      muted: "#f6f2ee",
      mutedForeground: "#726250",
      accent: "#D4A574",            // light coffee glaze
      accentForeground: "#0b0b0b",
      destructive: "#ef4444",
      destructiveForeground: "#ffffff",
      border: "#eadac8",
      input: "#fff2e6",
      ring: "#c47a3b",
      radius: "0.5rem",
    },
    dark: {
      background: "#141414",
      foreground: "#f8fafc",
      card: "#171212",              // dark espresso
      cardForeground: "#f8fafc",
      popover: "#1a1513",
      popoverForeground: "#f8fafc",
      primary: "#D4A574",           // caramel highlight on dark
      primaryForeground: "#0b0b0b",
      secondary: "#2a2320",         // steamed milk shadow
      secondaryForeground: "#f8fafc",
      muted: "#1f1b1a",
      mutedForeground: "#9aa0a6",
      accent: "#33241b",            // deep roast
      accentForeground: "#f8fafc",
      destructive: "#ef4444",
      destructiveForeground: "#f8fafc",
      border: "#2b2017",
      input: "#1e1a18",
      ring: "#f59e0b",
      radius: "0.5rem",
    },
    customCss: {
      // Badge: Live (green across both modes)
      ".badge-live": {
        "background-color": "#6b8e23",  // earthy olive-green
        color: "#ffffff",
      },

      // Button: Get API Key (swaps primary/accent on hover)
      ".btn-with-link": {
        "background-color": "var(--primary)",
        color: "var(--primary-foreground)",
        transition: "all 200ms ease",
        "border-radius": "var(--radius)",
      },
      ".btn-with-link:hover": {
        "background-color": "var(--accent)",
        color: "var(--accent-foreground)",
      },

      // Dark mode: swap accent → primary on hover for stronger contrast
      "@media (prefers-color-scheme: dark)": {
        ".btn-with-link:hover": {
          "background-color": "#A66A2A",  // light mode's primary
          color: "#ffffff",
        },
      },
    },
  },
  navigation: [
    {
      type: "category",
      label: "Documentation",
      items: [
        {
          type: "category",
          label: "Getting Started",
          icon: "rocket",
          collapsible: true,
          collapsed: false,
          link: {
            type: "doc",
            file: "start/overview",
            path: "/start",
          },
          items: [
            { type: "doc", file: "start/overview", path: "/start" },
            { type: "doc", file: "start/api-keys", path: "/start/api-keys" },
            { type: "doc", file: "start/first-api-call", path: "/start/first-api-call" },
          ],
        },
        {
          type: "category",
          label: "Products",
          icon: "boxes",
          collapsible: true,
          collapsed: false,
          items: [
            {
              type: "category",
              label: "Beans",
              collapsible: true,
              collapsed: true,
              link: { type: "doc", file: "products/beans/overview", path: "/products/beans" },
              items: [
                { type: "doc", file: "products/beans/overview", path: "/products/beans" },
                { type: "doc", file: "products/beans/scenarios", path: "/products/beans/scenarios" },
                { type: "doc", file: "products/beans/migration", path: "/products/beans/migration" },
                { type: "link", to: "/api/beans", label: "API reference" },
              ],
            },
            {
              type: "category",
              label: "Espresso",
              collapsible: true,
              collapsed: true,
              link: { type: "doc", file: "products/espresso/overview", path: "/products/espresso" },
              items: [
                { type: "doc", file: "products/espresso/overview", path: "/products/espresso" },
                { type: "doc", file: "products/espresso/workflows", path: "/products/espresso/workflows" },
                { type: "doc", file: "products/espresso/migration", path: "/products/espresso/migration" },
                { type: "link", to: "/api/espresso", label: "API reference" },
              ],
            },
            {
              type: "category",
              label: "Cortado",
              collapsible: true,
              collapsed: true,
              link: { type: "doc", file: "products/cortado/overview", path: "/products/cortado" },
              items: [
                { type: "doc", file: "products/cortado/overview", path: "/products/cortado" },
              ],
            },
          ],
        },
        {
          type: "category",
          label: "Build with Cafecito",
          icon: "bot",
          collapsible: true,
          collapsed: true,
          items: [
            { type: "doc", file: "guides/mcp-ai-agents", path: "/guides/mcp-ai-agents" },
            { type: "doc", file: "guides/cross-product-workflow", path: "/guides/cross-product-workflow" },
            { type: "doc", file: "guides/api-conventions", path: "/guides/api-conventions" },
            { type: "doc", file: "guides/client-patterns", path: "/guides/client-patterns" },
            { type: "doc", file: "guides/troubleshooting", path: "/guides/troubleshooting" },
            { type: "doc", file: "guides/pricing-limits", path: "/guides/pricing-limits" },
          ],
        },
        { type: "doc", file: "contact" },
        {
          type: "category",
          label: "Company & Policies",
          icon: "building",
          items: [
            { type: "link", to: "https://cafecito.tech", label: "Cafecito Website" },
            { type: "doc", file: "company/about-us" },
            { type: "doc", file: "company/privacy-policy" },
            { type: "doc", file: "company/terms-of-use" },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Reference",
      link: {
        type: "doc",
        file: "api-overview",
        path: "/api/overview",
      },
      items: [
        { type: "doc", file: "api-overview", path: "/api/overview" },
        { type: "link", to: "/api/beans", label: "Beans API" },
        { type: "link", to: "/api/espresso", label: "Espresso API" },
      ],
    },
    {
      type: "category",
      label: "Account",
      items: [
        { type: "custom-page", path: "/account", element: <AccountPage /> },
      ],
    },
  ],
  redirects: [
    { from: "/", to: "/start" },
    { from: "/introduction", to: "/start" },
    { from: "/howtos/api-keys", to: "/start/api-keys" },
    { from: "/howtos/beans-howto", to: "/products/beans" },
    { from: "/howtos/espresso-howto", to: "/products/espresso" },
    { from: "/howtos/espresso-scenarios", to: "/products/espresso/workflows" },
    { from: "/howtos/espresso-migration", to: "/products/espresso/migration" },
    { from: "/howtos/cortado-howto", to: "/products/cortado" },
    { from: "/howtos/mcp-howto", to: "/guides/mcp-ai-agents" },
    { from: "/pricing", to: "/guides/pricing-limits" },
    { from: "/guides/cross-product-workflows", to: "/guides/cross-product-workflow" },
    { from: "/start/troubleshooting", to: "/guides/troubleshooting" },
  ],
  apis: [
    {
      type: "file",
      input: "../config/beans.oas.json",
      path: "/api/beans",
    },
    {
      type: "file",
      input: "../config/espresso.oas.json",
      path: "/api/espresso",
    },
  ],
  authentication: clerkPubKey
    ? {
        type: "clerk",
        clerkPubKey,
        jwtTemplateName: clerkJwtTemplateName,
      }
    : undefined,
  apiKeys: {
    enabled: Boolean(clerkPubKey),
    createKey: async ({ apiKey, context, auth }) => {
      const createApiKeyRequest = new Request(
        serverUrl + "/v1/developer/api-key",
        {
          method: "POST",
          body: JSON.stringify({
            ...apiKey,
            email: auth.profile?.email,
            metadata: {
              userId: auth.profile?.sub,
              name: auth.profile?.name,
              subscription_plan: auth.profile?.subscription_plan,
              subscription_status: auth.profile?.subscription_status,
            },
          }),
          headers: {
            "Content-Type": "application/json",
          },
        },
      );

      const response = await fetch(
        await context.signRequest(createApiKeyRequest),
      );

      if (!response.ok) {
        throw new Error("Could not create API Key");
      }

      return true;
    },
  },
};

export default config;
