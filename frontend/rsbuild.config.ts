import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";
import { pluginHtmlMinifierTerser } from "rsbuild-plugin-html-minifier-terser";

export default defineConfig({
  plugins: [pluginReact(), pluginHtmlMinifierTerser()],
  source: {
    entry: { index: "./src/main.tsx" },
    define: {
      "process.env.DEV_AUTH": JSON.stringify(process.env.DEV_AUTH ?? ""),
    },
  },
  html: {
    template: "./index.html",
  },
  output: {
    distPath: { root: process.env.DIST_DIR ?? "dist" },
    cleanDistPath: true,
  },
});
