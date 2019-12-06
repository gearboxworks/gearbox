import NamespacedBaseActions from '../../_base/_store/actions'
import { FORCED_DELAY } from '../../_helpers'
import moduleConfig from '../config'
import api from '../_api'

import { DemoActions as Actions, DemoMutations as Mutations } from './method-names'

export default {
  ...NamespacedBaseActions(api, moduleConfig.namespace),

  [Actions.DEMO_ACTION]: ({ commit }, payload) => {
    console.warn('Called DEMO_ACTION. Returning a promise to call SOME_MUTATION in ' + FORCED_DELAY + 'ms with this payload:', payload)
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit(Mutations.DEMO_MUTATION, payload)
        resolve()
      }, FORCED_DELAY)
    })
  }
}
