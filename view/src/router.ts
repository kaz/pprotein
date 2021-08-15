import { createRouter, createWebHashHistory } from "vue-router";

import PProfList from "./components/PProfList.vue";
import PProfEntry from "./components/PProfEntry.vue";
import HttpLogList from "./components/HttpLogList.vue";
import HttpLogEntry from "./components/HttpLogEntry.vue";
import SlowLogEntry from "./components/SlowLogEntry.vue";
import SlowLogList from "./components/SlowLogList.vue";

export default createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: "/",
			redirect: "/pprof/"
		},
		{
			path: "/pprof/",
			component: PProfList,
			meta: {
				title: "PProf",
			},
		},
		{
			path: "/pprof/:id/",
			component: PProfEntry,
			meta: {
				title: "PProf {{id}}",
			},
		},
		{
			path: "/httplog/",
			component: HttpLogList,
			meta: {
				title: "HTTP Log",
			},
		},
		{
			path: "/httplog/:id/",
			component: HttpLogEntry,
			meta: {
				title: "HTTP Log {{id}}",
			},
		},
		{
			path: "/slowlog/",
			component: SlowLogList,
			meta: {
				title: "Slow Log",
			},
		},
		{
			path: "/slowlog/:id/",
			component: SlowLogEntry,
			meta: {
				title: "Slow Log {{id}}",
			},
		},
	],
});
