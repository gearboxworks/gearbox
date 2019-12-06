import actions from './actions'
import getters from './getters'
import mutations from './mutations'

const state = {
  records: [],
  filter: {},
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
