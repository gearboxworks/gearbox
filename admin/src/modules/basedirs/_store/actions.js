import BaseActions from '../../_base/_store/actions'
// import { UNSUPPORTED_ACTION } from '../../_helpers'
import api from '../_api'

import StoreMethodTypes from './private-types'
const { ActionTypes: Actions } = StoreMethodTypes

const OverrideActions = {
  // [Actions.LOAD_ALL]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  // [Actions.LOAD_ONE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  // [Actions.CREATE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  // [Actions.UPDATE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  // [Actions.DELETE]: ({ commit }, payload) => UNSUPPORTED_ACTION()
  [Actions.CHECK_DIRECTORY]: ({ commit }, dir) => api.checkDirectory(dir),
  // [Actions.CREATE_DIRECTORY]: ({ commit }, dir) => api.createDirectory(dir),
  [Actions.OPEN_DIRECTORY]: ({ commit }, dir) => api.openDirectory(dir)
}

export default { ...BaseActions(api), ...OverrideActions }
