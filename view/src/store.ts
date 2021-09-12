import { createStore, Store } from "vuex";

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
  groups: [] as string[],
  entries: {} as { [key: string]: Entry },
};

const syncPlugin = (store: Store<typeof state>) => {
  const es = new EventSource("/api/event");
  es.addEventListener("message", ({ data }) => {
    store.commit("saveEntry", JSON.parse(data));
  });

  store.state.endpoints.forEach((endpoint) => {
    store.dispatch("sync", { endpoint });
  });
};

export default createStore({
  state,
  plugins: [syncPlugin],
  mutations: {
    saveEntry(state, entry: Entry) {
      entry.Snapshot.Datetime = new Date(entry.Snapshot.Datetime);
      state.entries[entry.Snapshot.ID] = entry;

      if (
        entry.Snapshot.GroupId &&
        !state.groups.includes(entry.Snapshot.GroupId)
      ) {
        state.groups.push(entry.Snapshot.GroupId);
        state.groups.sort((a, b) => b.localeCompare(a));
      }
    },
  },
  actions: {
    async sync({ commit }, { endpoint }: { endpoint: string }) {
      try {
        const resp = await fetch(`/api/${endpoint}`);
        if (!resp.ok) {
          return alert("HTTP Error: " + (await resp.text()));
        }

        const entries = (await resp.json()) as Entry[];
        entries.forEach((entry) => commit("saveEntry", entry));
      } catch (e) {
        return alert(e);
      }
    },
  },
  getters: {
    entriesByType: (state) => (snapshotType: string) => {
      return Object.values(state.entries)
        .filter((e) => e.Snapshot.Type == snapshotType)
        .sort(
          (a, b) =>
            b.Snapshot.Datetime.getTime() - a.Snapshot.Datetime.getTime()
        );
    },
    entriesByGroup: (state) => (groupId: string) => {
      return Object.values(state.entries)
        .filter((e) => e.Snapshot.GroupId == groupId)
        .sort((a, b) => a.Snapshot.Label.localeCompare(b.Snapshot.Label));
    },
    availableEntriesByGroup: (_, getters) => (groupId: string) => {
      return getters
        .entriesByGroup(groupId)
        .filter((e: Entry) => e.Status == "ok");
    },
  },
});
