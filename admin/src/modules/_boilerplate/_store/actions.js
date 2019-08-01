import BaseActions from '../../_base/_store/actions'
import { UNSUPPORTED_ACTION, FORCED_DELAY } from '../../_helpers'
import api from '../_api'

import StoreMethodTypes from './private-types'
const { ActionTypes: Actions, MutationTypes: Mutations } = StoreMethodTypes

const OverrideActions = {

  [Actions.DEMO_ACTION]: ({ commit }, payload) => {
    console.warn('Called DEMO_ACTION. Returning a promise to call SOME_MUTATION in ' + FORCED_DELAY + 'ms with this payload:', payload)
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit(Mutations.DEMO_MUTATION, payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  [Actions.LOAD_ALL]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.LOAD_ONE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.CREATE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.UPDATE]: ({ commit }, payload) => UNSUPPORTED_ACTION(),
  [Actions.DELETE]: ({ commit }, payload) => UNSUPPORTED_ACTION()

}

export default { ...BaseActions(api), ...OverrideActions }
