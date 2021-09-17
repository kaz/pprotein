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

type SettingRecord = {
  key: string;
  value: string;
};

export type Config = { Type: string } & Omit<SnapshotTarget, "GroupId">;

const state = {
  endpoints: ["pprof", "httplog", "slowlog"],
  groups: [] as string[],
  entries: {} as { [key: string]: Entry },

  settingKeys: ["group", "repository", "httplog"],
  settings: {} as { [key: string]: SettingRecord },
};

const syncEntriesPlugin = (store: Store<typeof state>) => {
  const es = new EventSource("/api/event");
  es.addEventListener("message", ({ data }) => {
    store.commit("saveEntry", JSON.parse(data));
  });

  store.state.endpoints.forEach((endpoint) => {
    store.dispatch("fetchEntries", { endpoint });
  });
};

const syncSettingsPlugin = (store: Store<typeof state>) => {
  store.state.settingKeys.forEach(async (key) => {
    const resp = await fetch(`/api/setting/${key}`);
    if (!resp.ok) {
      return alert("HTTP Error: " + (await resp.text()));
    }

    store.commit("saveSetting", {
      key,
      value: await resp.text(),
    } as SettingRecord);
  });
};

export default createStore({
  state,
  plugins: [syncEntriesPlugin, syncSettingsPlugin],
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
    saveSetting(state, record: SettingRecord) {
      state.settings[record.key] = record;
    },
  },
  actions: {
    async fetchEntries({ commit }, { endpoint }: { endpoint: string }) {
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
    async updateSetting({ commit }, { key, value }: SettingRecord) {
      try {
        const resp = await fetch(`/api/setting/${key}`, {
          method: "POST",
          body: value,
        });
        if (!resp.ok) {
          return alert("HTTP Error: " + (await resp.text()));
        }

        commit("saveSetting", { key, value } as SettingRecord);
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
    groupConfig: (state): Config[] => {
      return JSON.parse(state.settings["group"]?.value || "[]");
    },
  },
});
