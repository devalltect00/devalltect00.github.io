import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "jsdom",
    // Isolated threads avoid Windows fork-handshake stalls without weakening per-file test isolation.
    pool: "threads",
    include: ["tests/**/*.test.ts"],
  },
});
