import { createRouter, createWebHashHistory } from "vue-router";

import PProf from "./components/PProf.vue";
import PProfInternal from "./components/PProfInternal.vue";

export default createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: "/",
			redirect: "/pprof"
		},
		{
			path: "/pprof",
			component: PProf,
			meta: {
				title: "PProf",
			},
		},
		{
			path: "/pprof/:id",
			component: PProfInternal,
			meta: {
				title: "Profile {{id}}",
			},
		}
	],
});
