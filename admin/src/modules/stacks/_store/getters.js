import BaseGetters from '../../_base/_store/getters'

import StackMethodTypes from './private-types'
const { GetterTypes: Getters } = StackMethodTypes

const OverrideGetters = {
  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => () => BaseGetters[Getters.LIST_OPTIONS](state, getters, rootState, rootGetters)('basedir')
}

export default { ...BaseGetters, ...OverrideGetters }
