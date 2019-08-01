import Vue from 'vue'
import Vuex from 'vuex'

import projects from './modules/projects/_store'
import stacks from './modules/stacks/_store'
import services from './modules/services/_store'
import basedirs from './modules/basedirs/_store'
import gearspecs from './modules/gearspecs/_store'

Vue.use(Vuex)

export default new Vuex.Store({
  /**
   * In strict mode, whenever Vuex state is mutated outside of mutation handlers, an error will be thrown.
   */
  strict: true,
  state: {
    connectionStatus: {
      networkError: null,
      remainingRetries: 5
    }
  },
  modules: {
    basedirs,
    gearspecs,
    projects,
    services,
    stacks
  },
  getters: {},
  actions: {},
  mutations: {
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
