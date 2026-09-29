import * as esbuild from "esbuild";
import { copyFile, mkdir } from "node:fs/promises";

try {
  process.loadEnvFile("../.env");
} catch {
  // sem .env, usa as variáveis do ambiente
}

const watch = process.argv.includes("--watch");

await mkdir("dist", { recursive: true });
await copyFile("public/index.html", "dist/index.html");

const options = {
  entryPoints: ["src/main.tsx"],
  bundle: true,
  outdir: "dist",
  format: "esm",
  sourcemap: watch,
  minify: !watch,
  define: {
    "process.env.DISCORD_CLIENT_ID": JSON.stringify(process.env.DISCORD_CLIENT_ID ?? ""),
    "process.env.NODE_ENV": JSON.stringify(watch ? "development" : "production"),
  },
  logLevel: "info",
};

if (watch) {
  const ctx = await esbuild.context(options);
  await ctx.watch();
} else {
  await esbuild.build(options);
}
