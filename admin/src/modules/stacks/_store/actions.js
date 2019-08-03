import NamespacedBaseActions from '../../_base/_store/actions'
import moduleConfig from '../config'
import api from '../_api'
// import { StackActions as Actions } from './method-names'

export default {
  ...NamespacedBaseActions(api, moduleConfig.namespace)
}
