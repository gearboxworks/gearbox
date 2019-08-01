import BaseGetters from '../../_base/_store/getters'
// import { UNSUPPORTED_GETTER } from '../../_helpers'
import StoreMethodTypes from './private-types'
const { GetterTypes: Getters } = StoreMethodTypes

const OverrideGetters = {

  // [Getters.DEMO_GETTER]: (state) => {
  //   return 'This is the result from DEMO_GETTER.'
  // },
  //
  // [Getters.LIST_FILTERED_BY]: (state) => (state) => (fieldName, allowedValues, allowedFields) => UNSUPPORTED_GETTER(),
  // [Getters.LIST_FILTERED]: (state, getters) => UNSUPPORTED_GETTER(),

  [Getters.HAS_EXTRA_BASEDIRS]: (state, getters, rootState, rootGetters) => () => state.records.length > 1,

  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => () => BaseGetters[Getters.LIST_OPTIONS](state, getters, rootState, rootGetters)('basedir')
}

export default { ...BaseGetters, ...OverrideGetters }
