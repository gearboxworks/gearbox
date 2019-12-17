import Vue from 'vue'
import Vuex from 'vuex'

import projects from './modules/projects/_store'
import stacks from './modules/stacks/_store'
import services from './modules/services/_store'
import basedirs from './modules/basedirs/_store'
import gearspecs from './modules/gearspecs/_store'

Vue.use(Vuex)

export default new Vuex.Store({
  strict: true,
  state: {},
  modules: {
    basedirs,
    gearspecs,
    projects,
    services,
    stacks
  },
  getters: {},
  actions: {},
  mutations: {}
})
