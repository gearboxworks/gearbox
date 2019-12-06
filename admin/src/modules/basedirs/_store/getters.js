import NamespacedBaseGetters from '../../_base/_store/getters'
import moduleConfig from '../config'

import { BasedirGetters as Getters } from './method-names'

export default {
  ...NamespacedBaseGetters(moduleConfig.namespace),

  [Getters.HAS_EXTRA_BASEDIRS]: (state, getters, rootState, rootGetters) => () => state.records.length > 1
}
