import NamespacedBaseGetters from '../../_base/_store/getters'
import moduleConfig from '../config'
// import { StackGetters as Getters } from './method-names'

export default {
  ...NamespacedBaseGetters(moduleConfig.namespace)
}
