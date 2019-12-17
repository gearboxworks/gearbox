import { NamespacedBaseGetters, NamespacedBaseActions, NamespacedBaseMutations } from '../../_base/_store/method-names'
import moduleConfig from '../config'

/**
 * @type {{LIST_OPTIONS: string, LIST_FILTERED: string, LIST_FILTERED_BY: string, FIND_BY: string}}
 */
export const GearspecGetters = {
  ...NamespacedBaseGetters(moduleConfig.namespace),
  /**
   * The following are added by destructuring the result of the above function call
   */
  // FIND_BY: 'gearspecs/FIND_BY',
  // LIST_OPTIONS: 'gearspecs/LIST_OPTIONS',
  // LIST_FILTERED_BY: 'gearspecs/LIST_FILTERED_BY',
  // LIST_FILTERED: 'gearspecs/LIST_FILTERED',
  /**
   *  The following getters are specific to gearspecs
   */
  STACK: 'gearspecs/STACK',
  SERVICES: 'gearspecs/SERVICES',
  FIND_COMPATIBLE_SERVICE: 'gearspecs/FIND_COMPATIBLE_SERVICE',
  DEFAULT_SERVICE: 'gearspecs/DEFAULT_SERVICE',
  SERVICE_VERSIONS_GROUPED_BY_PROGRAM: 'gearspecs/SERVICE_VERSIONS_GROUPED_BY_PROGRAM'
}

/**
 * @type {{DELETE: string, LOAD_ALL: string, CREATE: string, LOAD_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const GearspecActions = {
  ...NamespacedBaseActions(moduleConfig.namespace)
  // LOAD_ALL: 'LOAD_ALL',
  // LOAD_ONE: 'LOAD_ONE',
  // CREATE: 'CREATE',
  // UPDATE: 'UPDATE',
  // DELETE: 'DELETE',
  // SET_LIST_FILTER: 'SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
}

/**
 * @type {{DELETE: string, SET_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_ALL: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const GearspecMutations = {
  ...NamespacedBaseMutations(moduleConfig.namespace)
  // SET_ALL: 'gearspecs/SET_ALL',
  // SET_ONE: 'gearspecs/SET_ONE',
  // UPDATE: 'gearspecs/UPDATE',
  // DELETE: 'gearspecs/DELETE',
  // SET_LIST_FILTER: 'gearspecs/SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'gearspecs/SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'gearspecs/SET_LIST_FILTER_SORT_ASC'
}

export default {
  Namespace: moduleConfig.namespace,
  Getters: GearspecGetters,
  Actions: GearspecActions,
  Mutations: GearspecMutations
}
