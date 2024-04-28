import tailwind from "@astrojs/tailwind";
import icon from "astro-icon";
import { defineConfig } from "astro/config";
import { loadEnv } from "vite";

const {
  WEDFRONT_BACKEND_URL,
  WEDFRONT_BACKEND_PUBLIC_URL,
} = loadEnv(process.env.NODE_ENV, process.cwd(), "");

WEDFRONT_BACKEND_URL = WEDFRONT_BACKEND_URL || "http://127.0.0.1:1378"
WEDFRONT_BACKEND_PUBLIC_URL = WEDFRONT_BACKEND_PUBLIC_URL || "http://127.0.0.1:1378"

export {
  WEDFRONT_BACKEND_URL as backend_url,
  WEDFRONT_BACKEND_PUBLIC_URL as backend_public_url,
}

// https://astro.build/config
export default defineConfig({
  site: "https://parham-alvani.github.com/wedding",
  output: 'server',
  integrations: [tailwind(), icon()],
});
