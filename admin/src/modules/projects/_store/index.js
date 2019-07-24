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
  },
  sortBy: 'project_title',
  sortOrder: true
}

export default {
  namespaced: true,
  state,
  actions,
  getters,
  mutations
}
