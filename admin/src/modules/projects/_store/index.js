import actions from './actions'
import getters from './getters'
import mutations from './mutations'

const state = {
  records: [],
  filter: {
    'states': ['running', 'stopped', 'candidates'],
    'basedir': 'all',
    'stacks': 'all',
    'programs': 'all'
  },
  sorting: {
    fieldName: 'id',
    ascending: true
  }
}

export default {
  namespaced: false,
  state,
  getters,
  actions,
  mutations
}
