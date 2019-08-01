import StoreMethodTypes from './private-types'
const { ActionTypes: Actions, MutationTypes: Mutations } = StoreMethodTypes

function actions (api) {
  return {
    [Actions.LOAD_ALL] ({ commit }) {
      return api.fetchAll()
        .then(records => {
          if (records) {
            commit(Mutations.SET_ALL, records)
          } else {
            throw new Error('Could not load all ' + api.apiEndpoint)
          }
        })
    },

    [Actions.LOAD_ONE] ({ commit }, id) {
      return api.fetchOne(id)
        .then(record => {
          if (record) {
            commit(Mutations.SET_ONE, record)
          } else {
            throw new Error('Could not load one ' + api.recordType)
          }
        })
    },

    [Actions.CREATE] ({ commit }, recordData) {
      console.log('Base Actions.CREATE', recordData)

      return api.create(recordData)
        .then(record => {
          if (record) {
            commit(Mutations.SET_ONE, record)
          } else {
            throw new Error('Could not create one ' + api.recordType)
          }
        })
    },

    [Actions.UPDATE] ({ commit }, payload) {
      const { record, recordData } = payload

      console.log('Base Actions.UPDATE', record, recordData)

      return api.update(record, recordData)
        .then(result => {
        /**
         * on success, commit the new project to the projects store before resolving
         */
          if (result) {
            commit(Mutations.UPDATE, payload)
          } else {
            throw new Error('Could not update ' + api.recordType + ': ' + record.id)
          }
          return (result)
        })
    },

    [Actions.DELETE] ({ commit }, id) {
      return api.delete(id)
        .then(result => {
        /**
         * on success, commit the new project to the projects store before resolving
         */
          if (result) {
            commit(Mutations.DELETE, id)
          } else {
            throw new Error('Could not delete record ' + id)
          }
          return (result)
        })
    },

    [Actions.SET_LIST_FILTER] ({ commit }, payload) {
      commit(Mutations.SET_LIST_FILTER, payload)
    },

    [Actions.SET_LIST_FILTER_SORT_BY] ({ commit }, fieldName) {
      commit(Mutations.SET_LIST_FILTER_SORT_BY, fieldName)
    },

    [Actions.SET_LIST_FILTER_SORT_ASC] ({ commit }, isAscending) {
      commit(Mutations.SET_LIST_FILTER_SORT_ASC, isAscending)
    }
  }
}

export default actions
