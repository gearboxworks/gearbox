import BaseTypes from '../../_base/_store/private-types'
import moduleConfig from '../config'
/**
 * Note, you should avoid overriding method names that already exist in BaseTypes
 */
const ExtraGetters = {
  FIND_COMPATIBLE_SERVICE: 'FIND_COMPATIBLE_SERVICE',
  DEFAULT_SERVICE_FOR_GEARSPEC: 'DEFAULT_SERVICE_FOR_GEARSPEC',
  GEARSPEC_SERVICES: 'GEARSPEC_SERVICES'
}

const ExtraActions = {}
const ExtraMutations = {}

export default {
  Namespace: moduleConfig.namespace,
  GetterTypes: { ...BaseTypes.GetterTypes, ...ExtraGetters },
  ActionTypes: { ...BaseTypes.ActionTypes, ...ExtraActions },
  MutationTypes: { ...BaseTypes.MutationTypes, ...ExtraMutations }
}
