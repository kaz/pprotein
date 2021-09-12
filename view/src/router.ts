import { createRouter, createWebHashHistory } from "vue-router";
import GroupConfig from "./components/GroupConfig.vue";
import GroupEntry from "./components/GroupEntry.vue";
import GroupIndex from "./components/GroupIndex.vue";
import GroupList from "./components/GroupList.vue";
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
      redirect: "/group/",
    },
    {
      path: "/group/",
      component: GroupList,
      meta: {
        title: "group",
      },
    },
    {
      path: "/group/config/",
      component: GroupConfig,
      meta: {
        title: "group:config",
      },
    },
    {
      path: "/group/:gid/",
      component: GroupEntry,
      children: [
        {
          path: "index/",
          component: GroupIndex,
          meta: {
            title: "group:{{gid}}",
          },
        },
        {
          path: "pprof/:id/",
          component: PProfEntry,
          meta: {
            title: "pprof:{{id}} | group:{{gid}}",
          },
        },
        {
          path: "httplog/:id/",
          component: HttpLogEntry,
          meta: {
            title: "httplog:{{id}} | group:{{gid}}",
          },
        },
        {
          path: "slowlog/:id/",
          component: SlowLogEntry,
          meta: {
            title: "slowlog:{{id}} | group:{{gid}}",
          },
        },
      ],
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
        title: "pprof:{{id}}",
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
        title: "httplog:config",
      },
    },
    {
      path: "/httplog/:id/",
      component: HttpLogEntry,
      meta: {
        title: "httplog:{{id}}",
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
        title: "slowlog:{{id}}",
      },
    },
  ],
});
