import BaseGetters from '../../_base/_store/getters'
import { UNSUPPORTED_GETTER } from '../../_helpers'

import StoreMethodTypes from './private-types'
const { GetterTypes: Getters } = StoreMethodTypes

const OverrideGetters = {

  [Getters.FIND_BY]: (state, getters, rootState, rootGetters) => (fieldName, fieldValue) => {
    /**
     * manipulate arguments
     */
    const results = BaseGetters[Getters.FIND_BY](state)(fieldName, fieldValue)
    /**
     * manipulate results
     */
    return results
  },

  [Getters.DEMO_GETTER]: (state, getters, rootState, rootGetters) => () => {
    return 'This is the result from DEMO_GETTER.'
  },

  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => () => BaseGetters[Getters.LIST_OPTIONS](state, getters, rootState, rootGetters)('someAttributeName')
}

export default { ...BaseGetters, ...OverrideGetters }
