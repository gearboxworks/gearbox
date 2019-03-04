import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { HTTP } from './http-common'

Vue.use(Vuex)
Vue.use(VueAxios, axios)

export default new Vuex.Store({
  /**
   * In strict mode, whenever Vuex state is mutated outside of mutation handlers, an error will be thrown.
   */
  strict: true,
  state: {
    projects: [],
    stacks: [],
    stack_members: [],
    gears: []
  },
  getters: {
    projectByName: (state) => (projectName) => {
      return state.projects.find(p => p.name === projectName)
    }
  },
  actions: {
    loadProjects ({ commit }) {
      HTTP.get(
        'projects',
        { crossDomain: true }
      ).catch((error) => {
        // handle error
        // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        console.log('rejected', error)
      }).then(r => r.data.data).then((projects) => {
        commit('SET_PROJECTS', projects)
      })
    },
    loadStacks ({ commit }) {
      HTTP.get(
        'stacks',
        { crossDomain: true }
      ).catch((error) => {
        // handle error
        // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        console.log('rejected', error)
      }).then(r => r.data.data).then((stacks) => {
        commit('SET_STACKS', stacks)
      })
    },
    loadGears ({ commit }) {
      // axios
      //   .get(
      //     'http://127.0.0.1:9999/gears',
      //     { crossDomain: true },
      //   )
      //   .catch((error) => {
      //     // handle error
      //     // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
      //   })
      //   .then(r => r.data.data)
      //   .then((gears) => {
      //     commit('SET_GEARS', gears);
      //   });
    },
    updateProject ({ commit }, payload) {
      const { projectName, project } = payload

      commit('UPDATE_PROJECT', { projectName, project })

      HTTP({
        method: 'post',
        url: 'project/' + projectName,
        data: project
      }).then(r => r.data).then((project) => {
        // move commit here
        // resolve()
      }).catch((error) => {
        console.log('rejected', error)
        // resolve();
      })
    }
  },
  mutations: {
    /**
     * Names of mutation functions should be all-caps -- that's "idiomatic Vue"
     */
    SET_PROJECTS (state, projects) {
      state.projects = projects
    },
    SET_STACKS (state, stacks) {
      state.stacks = stacks
    },
    SET_GEARS (state, gears) {
      state.gears = gears
    },
    UPDATE_PROJECT (state, args) {
      const { projectName, project } = args
      const p = this.getters.projectByName(projectName)

      p.name = project.name
      p.hostname = project.hostname
      p.group = project.group
      p.enabled = project.enabled
    }
  }

})
