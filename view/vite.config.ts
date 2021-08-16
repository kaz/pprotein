import vue from "@vitejs/plugin-vue";
import { UserConfig } from "vite";

export default {
	plugins: [vue()],
	proxy: {
		"/api": "http://localhost:9000",
	},
} as UserConfig;
