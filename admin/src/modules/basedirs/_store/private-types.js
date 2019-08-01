import BaseTypes from '../../_base/_store/private-types'
import moduleConfig from '../config'
/**
 * Note, you should avoid overriding method names that already exist in BaseTypes
 */
const ExtraGetters = {
  HAS_EXTRA_BASEDIRS: 'HAS_EXTRA_BASEDIRS'
}

const ExtraActions = {
  CHECK_DIRECTORY: 'CHECK_DIRECTORY',
  // CREATE_DIRECTORY: 'CREATE_DIRECTORY',
  OPEN_DIRECTORY: 'OPEN_DIRECTORY'
}

const ExtraMutations = {}

export default {
  Namespace: moduleConfig.namespace,
  GetterTypes: { ...BaseTypes.GetterTypes, ...ExtraGetters },
  ActionTypes: { ...BaseTypes.ActionTypes, ...ExtraActions },
  MutationTypes: { ...BaseTypes.MutationTypes, ...ExtraMutations }
}
