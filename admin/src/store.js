import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { getConfig as raxConfig } from 'retry-axios'
import HTTP from './http-common'

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
    gears: [],
    connectionStatus: {
      networkError: null,
      remainingRetries: 5
    }
  },
  getters: {
    projectByName: (state) => (projectName) => {
      return state.projects.find(p => p.name === projectName)
    }
  },
  actions: {
    retryHTTPRequest (url) {

    },
    loadProjects ({ commit }) {
      try {
        HTTP.get(
          'projects',
          {
            crossDomain: true,
            raxConfig: {
              // You can detect when a retry is happening, and figure out how many
              // retry attempts have been made
              onRetryAttempt: (err) => {
                const cfg = raxConfig(err)
                commit('SET_NETWORK_ERROR', err.message)
                commit('SET_REMAINING_RETRIES', cfg.retry - cfg.currentRetryAttempt)
              }
            }
          }
        ).catch((error, config) => {
          // handle error
          // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
          console.log('rejected', error)
          // if (error.message === 'Network Error') {
          //   commit('SET_NETWORK_ERROR', error.message)
          // }
        })
          .then(r => r ? r.data.data : null)
          .then((projects) => {
            if (projects) {
              commit('SET_PROJECTS', projects)
            }
          })
      } catch (e) {
        console.log(e)
      }
    },
    loadStacks ({ commit }) {
      HTTP.get(
        'stacks',
        { crossDomain: true }
      ).catch((error) => {
        // handle error
        // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        console.log('rejected', error)
      })
        .then(r => r ? r.data.data : null)
        .then((stacks) => {
          if (stacks) {
            commit('SET_STACKS', stacks)
          }
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
    },
    SET_NETWORK_ERROR (state, message) {
      state.connectionStatus.networkError = message
    },
    CLEAR_NETWORK_ERROR (state) {
      state.connectionStatus.networkError = ''
    },
    SET_REMAINING_RETRIES (state, remainingRetries) {
      state.connectionStatus.remainingRetries = remainingRetries
    }
  }

})
