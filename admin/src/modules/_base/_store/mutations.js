import Vue from 'vue'
import StoreMethodTypes from './private-types'
const { MutationTypes: Mutations } = StoreMethodTypes

const mutations = {

  [Mutations.SET_ALL] (state, records) {
    if (!records || records.length === 0) {
      console.warn('Most likely `records` arg should not be empty!')
    }
    Vue.set(state, 'records', records)
  },

  [Mutations.SET_ONE] (state, record) {
    if (!record) {
      throw new Error('An empty record cannot be put to store')
    }
    if (!record.id) {
      throw new Error('A record with no id cannot be put to store')
    }

    if (!record.attributes || Object.entries(record.attributes).length === 0) {
      throw new Error('An record with no attributes is invalid and cannot be put to store: ' + record.id)
    }

    const index = state.records.findIndex(r => r.id === record.id)
    if (index < 0) {
      state.gearspecs.records.push(record)
    } else {
      console.error('Record ' + record.id + ' already exists')
    }
  },

  [Mutations.UPDATE] (state, payload) {
    const { record, recordData } = payload
    Vue.set(record, 'attributes', { ...record.attributes, ...recordData.attributes })
  },

  [Mutations.DELETE] (state, id) {
    const index = state.records.findIndex(r => r.id === id)
    if (index >= 0) {
      state.records.splice(index, 1)
    } else {
      console.error('Could not delete record ' + id)
    }
  },

  [Mutations.SET_LIST_FILTER] (state, payload) {
    const { field, values } = payload
    Vue.set(state.filter, field, values)
  },

  [Mutations.SET_LIST_FILTER_SORT_BY] (state, fieldName) {
    state.sorting.fieldName = fieldName
  },

  [Mutations.SET_LIST_FILTER_SORT_ASC] (state, isAscending) {
    state.sorting.ascending = !!isAscending
  }
}

export default mutations
