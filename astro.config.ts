import sitemap from "@astrojs/sitemap";
import { defineConfig } from "astro/config";

export default defineConfig({
  site: "https://devalltect00.github.io",
  output: "static",
  trailingSlash: "always",
  integrations: [sitemap()],
});
