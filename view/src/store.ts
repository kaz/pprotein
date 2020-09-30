import Vuex from "vuex";

type ProfileStatus = "ok" | "fail" | "pending";
type ProfileInfo = {
	Status: ProfileStatus;
	Message: string;
	Profile: {
		ID: string;
		Datetime: Date;
		URL: string;
		Duration: number;
	};
};

export default new Vuex.Store({
	state: {
		profiles: [] as ProfileInfo[],
	},
	mutations: {
		setProfiles(state, profiles: { [key: string]: ProfileInfo; }) {
			state.profiles = Object.values(profiles).map(e => {
				e.Profile.Datetime = new Date(e.Profile.Datetime);
				return e;
			});
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
