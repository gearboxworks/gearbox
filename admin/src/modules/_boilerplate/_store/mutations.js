import NamespacedBaseMutations from '../../_base/_store/mutations'
import moduleConfig from '../config'
import { DemoMutations as Mutations } from './method-names'

export default {
  ...NamespacedBaseMutations(moduleConfig.namespace),

  [Mutations.DEMO_MUTATION] (state, payload) {
    console.error('Called DEMO_MUTATION mutation:', payload)
  }
}
