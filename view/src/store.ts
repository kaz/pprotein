import Vuex from "vuex";

export type Entry = {
  Status: "ok" | "fail" | "pending";
  Message: string;
  Snapshot: {
    ID: string;
    Datetime: Date;
    URL: string;
    Duration: number;
  };
};

type AddRequest = {
  endpoint: string;
  URL: string;
  Duration: number;
};

export default new Vuex.Store({
  state: {
    remote: {} as { [key: string]: Entry[] },
  },
  mutations: {
    setStoreData(
      state,
      { endpoint, entries }: { endpoint: string; entries: Entry[] }
    ) {
      state.remote[endpoint] = entries
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
    async syncStoreData({ commit }, { endpoint }: { endpoint: string }) {
      try {
        const resp = await fetch(`/api/${endpoint}`);
        if (!resp.ok) {
          return alert("HTTP Error: " + (await resp.text()));
        }

        commit("setStoreData", {
          endpoint,
          entries: await resp.json(),
        });
      } catch (e) {
        return alert(e);
      }
    },
    async postStoreData(_, { endpoint, URL, Duration }: AddRequest) {
      const resp = await fetch(`/api/${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ URL, Duration }),
      });

      if (!resp.ok) {
        alert(await resp.text());
      }
    },
  },
});
