import { createRouter, createWebHistory } from "vue-router";

import PProf from "./components/PProf.vue";
import PProfInternal from "./components/PProfInternal.vue";

export default createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: "/",
			redirect: "/pprof"
		},
		{
			path: "/pprof",
			component: PProf,
		},
		{
			path: "/pprof/:id",
			component: PProfInternal,
		}
	],
});
