import BaseActions from '../../_base/_store/actions'
import { UNSUPPORTED_ACTION } from '../../_helpers'
import api from '../_api'

import StackMethodTypes from './private-types'
const { ActionTypes: Actions } = StackMethodTypes

const OverrideActions = {
  // [Actions.LOAD_ALL]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.LOAD_ONE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.CREATE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.UPDATE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.DELETE]: ({ commit }, payload) => UNSUPPORTED_ACTION()
}

export default { ...BaseActions(api), ...OverrideActions }
