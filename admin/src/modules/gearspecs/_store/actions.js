import NamespacedBaseActions from '../../_base/_store/actions'
import moduleConfig from '../config'
import api from '../_api'
// import { GearspecActions as Actions } from './method-names'

export default {
  ...NamespacedBaseActions(api, moduleConfig.namespace)
}
