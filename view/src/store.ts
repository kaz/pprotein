import Vuex, { Store } from "vuex";

export type Entry = {
  Status: "ok" | "fail" | "pending";
  Message: string;
  Snapshot: SnapshotMeta & SnapshotTarget;
};

type SnapshotMeta = {
  Type: string;
  ID: string;
  Datetime: Date;
  GitRevision: string;
};
export type SnapshotTarget = {
  GroupId: string;
  Label: string;
  URL: string;
  Duration: number;
};

const state = {
  endpoints: ["pprof", "httplog", "slowlog"],
  entries: {} as { [key: string]: Entry[] },
};

const syncPlugin = (store: Store<typeof state>) => {
  const es = new EventSource("/api/event");
  es.addEventListener("message", ({ data }) => {
    store.dispatch("sync", { endpoint: data });
  });

  store.state.endpoints.forEach((endpoint) => {
    store.dispatch("sync", { endpoint });
  });
};

export default new Vuex.Store({
  state,
  plugins: [syncPlugin],
  mutations: {
    sync(state, { endpoint, entries }: { endpoint: string; entries: Entry[] }) {
      state.entries[endpoint] = entries
        .map((e) => {
          e.Snapshot.Datetime = new Date(e.Snapshot.Datetime);
          return e;
        })
        .sort(
          (a, b) =>
            b.Snapshot.Datetime.getTime() - a.Snapshot.Datetime.getTime()
        );
    },
  },
  actions: {
    async sync({ commit }, { endpoint }: { endpoint: string }) {
      try {
        const resp = await fetch(`/api/${endpoint}`);
        if (!resp.ok) {
          return alert("HTTP Error: " + (await resp.text()));
        }

        commit("sync", {
          endpoint,
          entries: await resp.json(),
        });
      } catch (e) {
        return alert(e);
      }
    },
  },
});
