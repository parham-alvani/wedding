import node from "@astrojs/node";
import { defineConfig } from "astro/config";
import icon from "astro-icon";

const backendUrl = process.env.WEDFRONT_BACKEND_URL || "http://127.0.0.1:1378";

interface Config {
  backend_url: string;
}
const config: Config = {
  backend_url: backendUrl,
};

export { config };

// https://astro.build/config
export default defineConfig({
  site: process.env.WEDFRONT_SITE_URL || "https://parham-alvani.github.com/wedding",
  output: "server",
  integrations: [icon()],
  adapter: node({
    mode: "standalone",
  }),
  vite: {
    build: {
      rollupOptions: {
        external: ["@astrojs/compiler-rs"],
      },
    },
  },
});
