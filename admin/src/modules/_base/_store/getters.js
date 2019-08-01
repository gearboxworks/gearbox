import StoreMethodTypes from './private-types'
const { GetterTypes: Getters } = StoreMethodTypes

const getters = {

  [Getters.FIND_BY]: (state) => (fieldName, fieldValue) => {
    return (fieldName === 'id')
      ? state.records.find(p => p.id === fieldValue)
      : state.records.find(p => p.attributes[fieldName] === fieldValue)
  },

  [Getters.LIST_FILTERED_BY]: (state) => (fieldName, allowedValues) => {
    let valuesArray = Array.isArray(allowedValues) ? allowedValues : [allowedValues]
    // 'notes' and 'stack' are not included on purpose because simple comparison does not work on them
    let list = []

    if (fieldName === 'id') {
      list = state.records.filter(p => valuesArray.indexOf(p.id) !== -1)
    } else if (fieldName) {
      list = state.records.filter(p => valuesArray.indexOf(p.attributes[fieldName]) !== -1)
    } else {
      /**
       * return unfiltered
       */
      return state.records
    }
    // Project specific:
    //
    // else if (fieldName === 'stacks') {
    //   list = state.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.gearspec_id.indexOf(val) > -1)))
    // } else if (fieldName === 'programs') {
    //   list = state.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.service_id.split('/')[1].split(':')[0] === val)))
    // }

    return list
  },

  [Getters.LIST_FILTERED]: (state, getters) => () => {
    let list = state.records
    const sortAscending = !!state.sorting.ascending
    // const sortBy = state.sortBy
    for (const field in state.filter) {
      const values = state.filter[field]
      if (values === 'all') {
        continue
      }
      list = (list.filter(p => getters[Getters.LIST_FILTERED_BY](field, values).includes(p)))
    }
    /**
     * TODO implement sorting by attributes
     */
    return list.concat().sort((a, b) => a.id > b.id ? (sortAscending ? 1 : -1) : (a.id === b.id) ? 0 : (sortAscending ? -1 : 1))
  },

  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => (textAttribute) => {
    const options = []
    state.records.forEach((el, idx) => {
      options.push({
        value: el.id,
        text: (textAttribute && el.attributes[textAttribute]) ? el.attributes[textAttribute] : el.id
      })
    })
    return options.sort((a, b) => a.value > b.value ? 1 : (a.value === b.value) ? 0 : -1)
  }
}

export default getters
