import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  console.log(`Vite mode: ${mode}`);

  let proxy = {};
  if (mode === "development") {
    proxy = {
      "^/api(?:/|$)": {
        target: "http://localhost:8082", // your Gin backend
        changeOrigin: true,
        secure: false,
      },
    };
  }

  return {
    plugins: [react(), tailwindcss()],
    base: process.env.BASE_URL || "/", // Ensures correct asset paths
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
