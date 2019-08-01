import BaseGetters from '../../_base/_store/getters'
import { UNSUPPORTED_GETTER } from '../../_helpers'

import StoreMethodTypes from './private-types'
const { GetterTypes: Getters } = StoreMethodTypes

const OverrideGetters = {

  // [Getters.FIND_BY]: (state) => (fieldName, fieldValue) => {
  //   /**
  //    * manipulate arguments
  //    */
  //   const results = BaseGetters[Getters.FIND_BY](state)(fieldName, fieldValue)
  //   /**
  //    * manipulate results
  //    */
  //   return results
  // },

  // [Getters.DEMO_GETTER]: (state) => {
  //   return 'This is the result from DEMO_GETTER.'
  // },

  // [Getters.LIST_FILTERED_BY]: (state) => (state) => (fieldName, allowedValues) => UNSUPPORTED_GETTER(),
  // [Getters.LIST_FILTERED]: (state, getters) => UNSUPPORTED_GETTER(),

  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => () => BaseGetters[Getters.LIST_OPTIONS](state, getters, rootState, rootGetters)('role')

}

export default { ...BaseGetters, ...OverrideGetters }
