import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig(({ mode }) => {
  console.log(`Vite mode: ${mode}`);

  let proxy = {};
  if (mode === "development") {
    proxy = {
      "^/api(?:/|$)": {
        target: "http://localhost:8081",
        changeOrigin: true,
        secure: false,
      },
    };
  }

  return {
    plugins: [react(), tailwindcss()],
    base: process.env.BASE_URL || "/",
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      proxy,
    },
  };
});
