import BaseTypes from '../../_base/_store/private-types'
import moduleConfig from '../config'
/**
 * Note, you should avoid overriding method names that already exist in BaseTypes
 */
const ExtraGetters = {
  DEMO_GETTER: 'DEMO_GETTER'
}

const ExtraActions = {
  DEMO_ACTION: 'DEMO_ACTION'
}

const ExtraMutations = {
  DEMO_MUTATION: 'DEMO_MUTATION'
}

export default {
  Namespace: moduleConfig.namespace,
  GetterTypes: { ...BaseTypes.GetterTypes, ...ExtraGetters },
  ActionTypes: { ...BaseTypes.ActionTypes, ...ExtraActions },
  MutationTypes: { ...BaseTypes.MutationTypes, ...ExtraMutations }
}
