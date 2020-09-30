import { UserConfig } from "vite";

export default {
	proxy: {
		"/api": "http://localhost:9000",
	},
} as UserConfig;
