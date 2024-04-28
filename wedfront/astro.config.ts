import tailwind from "@astrojs/tailwind";
import icon from "astro-icon";
import { defineConfig } from "astro/config";
import { loadEnv } from "vite";
import node from "@astrojs/node";
const {
  WEDFRONT_BACKEND_URL,
  WEDFRONT_BACKEND_PUBLIC_URL
} = loadEnv(process.env.NODE_ENV || "development", process.cwd(), "");
interface Config {
  backend_url: string;
  backend_public_url: string;
}
const config: Config = {
  backend_url: WEDFRONT_BACKEND_URL || "http://127.0.0.1:1378",
  backend_public_url: WEDFRONT_BACKEND_PUBLIC_URL || "http://127.0.0.1:1378"
};
export { config };


// https://astro.build/config
export default defineConfig({
  site: "https://parham-alvani.github.com/wedding",
  output: 'server',
  integrations: [tailwind(), icon()],
  adapter: node({
    mode: "middleware"
  })
});