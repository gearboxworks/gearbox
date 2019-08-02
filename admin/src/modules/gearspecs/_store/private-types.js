import BaseTypes from '../../_base/_store/private-types'
import moduleConfig from '../config'
/**
 * Note, you should avoid overriding method names that already exist in BaseTypes
 */
const ExtraGetters = {
  FIND_COMPATIBLE_SERVICE: 'FIND_COMPATIBLE_SERVICE',
  DEFAULT_GEARSPEC_SERVICE: 'DEFAULT_GEARSPEC_SERVICE',
  GEARSPEC_SERVICES: 'GEARSPEC_SERVICES',
  GEARSPEC_SERVICE_VERSIONS_GROUPED_BY_PROGRAM: 'GEARSPEC_SERVICE_VERSIONS_GROUPED_BY_PROGRAM'
}

const ExtraActions = {}

const ExtraMutations = {}

export default {
  Namespace: moduleConfig.namespace,
  GetterTypes: { ...BaseTypes.GetterTypes, ...ExtraGetters },
  ActionTypes: { ...BaseTypes.ActionTypes, ...ExtraActions },
  MutationTypes: { ...BaseTypes.MutationTypes, ...ExtraMutations }
}
