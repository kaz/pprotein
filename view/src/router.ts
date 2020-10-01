import { createRouter, createWebHashHistory } from "vue-router";

import PProfList from "./components/PProfList.vue";
import PProfEntry from "./components/PProfEntry.vue";

export default createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: "/",
			redirect: "/pprof"
		},
		{
			path: "/pprof",
			component: PProfList,
			meta: {
				title: "PProf",
			},
		},
		{
			path: "/pprof/:id",
			component: PProfEntry,
			meta: {
				title: "PProf {{id}}",
			},
		},
	],
});
