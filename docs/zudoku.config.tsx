import type { ZudokuConfig } from "zudoku";
import AccountPage from "./src/AccountPage";

const serverUrl =
  process.env.ZUDOKU_PUBLIC_GATEWAY_URL || import.meta.env.ZUPLO_SERVER_URL;
/**
 * Developer Portal Configuration
 * For more information, see:
 * https://zuplo.com/docs/dev-portal/zudoku/configuration/overview
 */
const config: ZudokuConfig = {
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
    title: "Cafecito Developer Portal",
    description: "Developer documentation for Cafecito APIs & MCPs. Project Cafecito offers AI tools to automate business tasks like news, PR and medical billing.",
    keywords: ["AI automation", "business tools", "billing software", "PR management", "news api", "news mcp", "market intelligence", "tech startup"],
    favicon: "/cafecito-dark.png",
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
          type: "doc",
          file: "introduction",
        },
        {
          type: "category",
          label: "Getting Started",
          icon: "puzzle",
          items: [
            {
              type: "doc",
              file: "howtos/api-keys",
            },
            {
              type: "doc",
              file: "howtos/beans-howto",
            },
            {
              type: "doc",
              file: "howtos/espresso-howto",
            },
            {
              type: "doc",
              file: "howtos/cortado-howto",
            },
          ],
        },
        {
          type: "doc",
          file: "howtos/mcp-howto",
        },
        {
          type: "doc",
          file: "pricing",
        },
        {
          type: "doc",
          file: "contact",
        },
        {
          type: "category",
          label: "Company & Policies",
          icon: "building",
          items: [
            {
              type: "link",
              to: "https://cafecito.tech",
              label: "Cafecito Website",
            },
            {
              type: "doc",
              file: "company/about-us",
            },
            {
              type: "doc",
              file: "company/privacy-policy",
            },
            {
              type: "doc",
              file: "company/terms-of-use",
            },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "API Reference",
      link: {
        type: "doc",
        file: "api-overview",
        path: "/api/overview",
        label: "Overview",
      },
      items: [
        {
          type: "doc",
          file: "api-overview",
          path: "/api/overview",
          label: "Overview",
        },
        {
          type: "link",
          to: "/api/beans",
          label: "Beans API",
        },
        {
          type: "link",
          to: "/api/espresso",
          label: "Espresso API",
        },
      ],
    },
    {
      type: "category",
      label: "Account",
      items: [
        {
          type: "custom-page",
          path: "/account",
          element: <AccountPage />,
        },
      ],
    },
  ],
  redirects: [{ from: "/", to: "/introduction" }],
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
  authentication: {
    type: "clerk",
    clerkPubKey: process.env.ZUDOKU_PUBLIC_CLERK_PUBLISHABLE_KEY,
    jwtTemplateName: process.env.ZUDOKU_PUBLIC_CLERK_JWT_TEMPLATE_NAME,
  },
  apiKeys: {
    enabled: true,
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
