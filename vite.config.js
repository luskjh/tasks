import { defineConfig } from "vite";

export default defineConfig({
  build: {
    manifest: true,
    rolldownOptions: {
      input: "js/dashboard.js",
    },
  },
  root: "internal/static",
});
