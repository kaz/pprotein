import { createRouter, createWebHashHistory } from "vue-router";
import All from "./components/All.vue";
import AllConfig from "./components/AllConfig.vue";
import HttpLogConfig from "./components/HttpLogConfig.vue";
import HttpLogEntry from "./components/HttpLogEntry.vue";
import HttpLogList from "./components/HttpLogList.vue";
import PProfEntry from "./components/PProfEntry.vue";
import PProfList from "./components/PProfList.vue";
import SlowLogEntry from "./components/SlowLogEntry.vue";
import SlowLogList from "./components/SlowLogList.vue";

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      redirect: "/all/",
    },
    {
      path: "/all/",
      component: All,
      meta: {
        title: "all",
      },
    },
    {
      path: "/all/config/",
      component: AllConfig,
      meta: {
        title: "all - config",
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
      path: "/httplog/config/",
      component: HttpLogConfig,
      meta: {
        title: "httplog - config",
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
