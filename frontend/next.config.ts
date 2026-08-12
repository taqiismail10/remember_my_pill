import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Playwright uses this loopback origin when exercising the local dev server.
  allowedDevOrigins: ["127.0.0.1"],
};

export default nextConfig;
