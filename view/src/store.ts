import { createStore, Store } from "vuex";

export type StatusText = "ok" | "fail" | "pending";

export interface Entry {
  Status: StatusText;
  Message: string;
  Snapshot: SnapshotMeta & SnapshotTarget;
}

interface SnapshotMeta {
  Type: string;
  ID: string;
  Datetime: Date;
  Repository?: RepositoryInfo;
}
export interface SnapshotTarget {
  GroupId: string;
  Label: string;
  URL: string;
  Duration: number;
}

export interface RepositoryInfo {
  Ref: string;
  Hash: string;
  Author: string;
  Message: string;
  Remote: string;
}

interface SettingRecord {
  key: string;
  value: string;
}

export interface Config extends Omit<SnapshotTarget, "GroupId"> {
  Type: string;
}

const state = {
  endpoints: ["memo", "pprof", "httplog", "slowlog"],
  groups: [] as string[],
  entries: {} as { [key: string]: Entry },

  settingKeys: ["group/targets", "httplog/config", "slowlog/config"],
  settings: {} as { [key: string]: SettingRecord },
};

const syncEntriesPlugin = (store: Store<typeof state>) => {
  const es = new EventSource("/api/event");
  es.addEventListener("message", ({ data }) => {
    const entry = JSON.parse(data) as Entry;
    store.commit("saveEntry", entry);
  });

  store.state.endpoints.forEach((endpoint) => {
    store.dispatch("fetchEntries", { endpoint });
  });
};

const syncSettingsPlugin = (store: Store<typeof state>) => {
  store.state.settingKeys.forEach(async (key) => {
    const resp = await fetch(`/api/${key}`);
    if (!resp.ok) {
      return alert(
        `http error: status=${resp.status}, message=${await resp.text()}`
      );
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
    async fetchEntries(store, { endpoint }: { endpoint: string }) {
      try {
        const resp = await fetch(`/api/${endpoint}`);
        if (!resp.ok) {
          return alert(
            `http error: status=${resp.status}, message=${await resp.text()}`
          );
        }

        const entries = (await resp.json()) as Entry[];
        entries.forEach((entry) => {
          store.commit("saveEntry", entry);
        });
      } catch (e) {
        return alert(e);
      }
    },
    async updateSetting(store, { key, value }: SettingRecord) {
      try {
        const resp = await fetch(`/api/${key}`, {
          method: "POST",
          body: value,
        });
        if (!resp.ok) {
          return alert(
            `http error: status=${resp.status}, message=${await resp.text()}`
          );
        }

        store.commit("saveSetting", { key, value } as SettingRecord);
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
        .sort((a, b) => {
          const ai = state.endpoints.indexOf(a.Snapshot.Type);
          const bi = state.endpoints.indexOf(b.Snapshot.Type);
          return ai == bi
            ? a.Snapshot.Label.localeCompare(b.Snapshot.Label)
            : ai - bi;
        });
    },
    availableEntriesByGroup: (_, getters) => (groupId: string) => {
      return getters
        .entriesByGroup(groupId)
        .filter((e: Entry) => e.Status == "ok");
    },
  },
});
