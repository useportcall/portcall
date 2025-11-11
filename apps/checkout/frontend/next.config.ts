import type { NextConfig } from "next";

const isDev = process.env.NODE_ENV !== "production";

const nextConfig: NextConfig = {
  output: "export",
  distDir: "dist",
  async rewrites() {
    if (!isDev) return [];
    return [
      {
        source: "/api/:path*", // Incoming request path
        destination: "http://localhost:8010/api/:path*", // Destination URL (e.g., your backend API)
      },
    ];
  },
};

export default nextConfig;
