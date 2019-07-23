import actions from './actions'
import getters from './getters'
import mutations from './mutations'

const state = {
  records: [],
  showProjectsHaving: {
    'states': ['running', 'stopped', 'candidates'],
    'basedir': 'all',
    'stacks': 'all',
    'programs': 'all'
  }
}

export default {
  namespaced: true,
  state,
  actions,
  getters,
  mutations
}
