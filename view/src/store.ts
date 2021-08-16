import Vuex from "vuex";

type Status = "ok" | "fail" | "pending";
export type EntryInfo = {
  Status: Status;
  Message: string;
  Entry: {
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
    remote: {} as { [key: string]: EntryInfo[] },
  },
  mutations: {
    setStoreData(state, { endpoint, entries }) {
      state.remote[endpoint] = Object.values(entries as EntryInfo[])
        .map((e) => {
          e.Entry.Datetime = new Date(e.Entry.Datetime);
          return e;
        })
        .sort(
          (a, b) => b.Entry.Datetime.getTime() - a.Entry.Datetime.getTime()
        );
    },
  },
  actions: {
    async syncStoreData({ commit }, { endpoint }) {
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
