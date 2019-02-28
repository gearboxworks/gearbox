import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';
import VueAxios from 'vue-axios';

Vue.use(Vuex);
Vue.use(VueAxios, axios);

export default new Vuex.Store({
  state: {
    projects: [],
    stacks: [],
    stack_members: [],
    gears: [],
  },
  getters: {

  },
  actions: {
    loadProjects({ commit }) {
      axios
        .get(
          'http://127.0.0.1:9999/projects',
          { crossDomain: true },
        )
        .catch((error) => {
          // handle error
          //alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        })
        .then(r => r.data.data)
        .then((projects) => {
          commit('SET_PROJECTS', projects);
        });
    },
    loadStacks({ commit }) {
      axios
        .get(
          'http://127.0.0.1:9999/stacks',
          { crossDomain: true },
        )
        .catch((error) => {
          // handle error
          // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        })
        .then(r => r.data.data)
        .then((stacks) => {
          commit('SET_STACKS', stacks);
        });
    },
    loadGears({ commit }) {
      axios
        .get(
          'http://127.0.0.1:9999/gears',
          { crossDomain: true },
        )
        .catch((error) => {
          // handle error
          // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        })
        .then(r => r.data.data)
        .then((gears) => {
          commit('SET_GEARS', gears);
        });
    },
  },
  mutations: {
    SET_PROJECTS(state, projects) {
      state.projects = projects;
    },
    SET_STACKS(state, stacks) {
      state.stacks = stacks;
    },
    SET_GEARS(state, gears) {
      state.gears = gears;
    },

  },

});
