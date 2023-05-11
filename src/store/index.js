import { createStore } from 'vuex';

export default createStore({
	state: {},
	getters: {},
	mutations: {
		SET_TOKEN(state, token) {
			state.token = token;
		},
	},
	actions: {},
	modules: {},
});
