import NamespacedBaseGetters from '../../_base/_store/getters'
import moduleConfig from '../config'

import { DemoGetters as Getters } from './method-names'

export default {
  ...NamespacedBaseGetters(moduleConfig.namespace),

  [Getters.DEMO_GETTER]: (state, getters, rootState, rootGetters) => () => {
    return 'This is the result from DEMO_GETTER.'
  }
}
