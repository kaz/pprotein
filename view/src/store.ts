import Vuex from "vuex";

type Status = "ok" | "fail" | "pending";
type EntryInfo = {
	Status: Status;
	Message: string;
	Entry: {
		ID: string;
		Datetime: Date;
		URL: string;
		Duration: number;
	};
};

export default new Vuex.Store({
	state: {
		profiles: [] as EntryInfo[],
	},
	mutations: {
		setProfiles(state, profiles: { [key: string]: EntryInfo; }) {
			state.profiles = Object.values(profiles)
				.map(e => {
					e.Entry.Datetime = new Date(e.Entry.Datetime);
					return e;
				})
				.sort((a, b) => b.Entry.Datetime.getTime() - a.Entry.Datetime.getTime());
		}
	},
	actions: {
		async updateProfiles({ commit }) {
			const resp = await fetch("/api/pprof/profiles");

			if (resp instanceof Error) {
				return alert(resp);
			}
			if (!resp.ok) {
				return alert("HTTP Error: " + await resp.text());
			}

			commit("setProfiles", await resp.json());
		}
	},
});
