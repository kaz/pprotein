import { createRouter, createWebHashHistory } from "vue-router";
import EntryList from "./components/EntryList.vue";
import GroupEntry from "./components/GroupEntry.vue";
import GroupIndex from "./components/GroupIndex.vue";
import GroupList from "./components/GroupList.vue";
import HttpLogEntry from "./components/HttpLogEntry.vue";
import PProfEntry from "./components/PProfEntry.vue";
import SettingList from "./components/SettingList.vue";
import SlowLogEntry from "./components/SlowLogEntry.vue";
import MemoEntry from "./components/MemoEntry.vue";

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
        {
          path: "memo/:id/",
          component: MemoEntry,
          meta: {
            title: "memo:{{id}} | group:{{gid}}",
          },
        },
      ],
    },
    {
      path: "/pprof/",
      component: EntryList,
      meta: {
        title: "pprof",
      },
      props: {
        endpoint: "pprof",
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
      component: EntryList,
      meta: {
        title: "httplog",
      },
      props: {
        endpoint: "httplog",
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
      component: EntryList,
      meta: {
        title: "slowlog",
      },
      props: {
        endpoint: "slowlog",
      },
    },
    {
      path: "/slowlog/:id/",
      component: SlowLogEntry,
      meta: {
        title: "slowlog:{{id}}",
      },
    },
    {
      path: "/setting/",
      component: SettingList,
      meta: {
        title: "setting",
      },
    },
  ],
});
