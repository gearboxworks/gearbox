import NamespacedBaseActions from '../../_base/_store/actions'
// import { FORCED_DELAY } from '../../_helpers'
import moduleConfig from '../config'
import api from '../_api'
import { BasedirActions as Actions } from './method-names'

export default {
  ...NamespacedBaseActions(api, moduleConfig.namespace),
  [Actions.CHECK_DIRECTORY]: ({ commit }, dir) => api.checkDirectory(dir),
  [Actions.CREATE_DIRECTORY]: ({ commit }, dir) => api.createDirectory(dir),
  [Actions.OPEN_DIRECTORY]: ({ commit }, dir) => api.openDirectory(dir)
}
