import vue from "@vitejs/plugin-vue";
import { UserConfig } from "vite";

export default {
  plugins: [vue()],
  server: {
    proxy: {
      "/api": "http://localhost:9000",
    },
  },
} as UserConfig;
