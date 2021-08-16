import HttpLogEntry from "./components/HttpLogEntry.vue";
import HttpLogList from "./components/HttpLogList.vue";
import Index from "./components/Index.vue";
import PProfEntry from "./components/PProfEntry.vue";
import PProfList from "./components/PProfList.vue";
import SlowLogEntry from "./components/SlowLogEntry.vue";
import SlowLogList from "./components/SlowLogList.vue";
import { createRouter, createWebHashHistory } from "vue-router";

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: Index,
      meta: {
        title: "top",
      },
    },
    {
      path: "/pprof/",
      component: PProfList,
      meta: {
        title: "pprof",
      },
    },
    {
      path: "/pprof/:id/",
      component: PProfEntry,
      meta: {
        title: "pprof - {{id}}",
      },
    },
    {
      path: "/httplog/",
      component: HttpLogList,
      meta: {
        title: "httplog",
      },
    },
    {
      path: "/httplog/:id/",
      component: HttpLogEntry,
      meta: {
        title: "httplog - {{id}}",
      },
    },
    {
      path: "/slowlog/",
      component: SlowLogList,
      meta: {
        title: "slowlog",
      },
    },
    {
      path: "/slowlog/:id/",
      component: SlowLogEntry,
      meta: {
        title: "slowlog - {{id}}",
      },
    },
  ],
});
