import NamespacedBaseMutations from '../../_base/_store/mutations'
import moduleConfig from '../config'
import api from '../_api'
// import { GearspecMutations as Mutations } from './method-names'

export default {
  ...NamespacedBaseMutations(api, moduleConfig.namespace)
}
