import { NamespacedBaseGetters, NamespacedBaseActions, NamespacedBaseMutations } from '../../_base/_store/method-names'
import moduleConfig from '../config'

/**
 * @type {{LIST_OPTIONS: string, LIST_FILTERED: string, LIST_FILTERED_BY: string, FIND_BY: string}}
 */
export const ServiceGetters = {
  ...NamespacedBaseGetters(moduleConfig.namespace),
  /**
   * The following are added by destructuring the result of the above function call
   */
  // FIND_BY: 'FIND_BY',
  // LIST_OPTIONS: 'LIST_OPTIONS',
  // LIST_FILTERED_BY: 'LIST_FILTERED_BY',
  // LIST_FILTERED: 'LIST_FILTERED',
  /**
   *  The following getters are specific to services
   */
  LIST_PROGRAM_OPTIONS: 'services/LIST_PROGRAM_OPTIONS'
}

/**
 * @type {{DELETE: string, LOAD_ALL: string, CREATE: string, LOAD_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const ServiceActions = {
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
export const ServiceMutations = {
  ...NamespacedBaseMutations(moduleConfig.namespace)
  // SET_ALL: 'SET_ALL',
  // SET_ONE: 'SET_ONE',
  // UPDATE: 'UPDATE',
  // DELETE: 'DELETE',
  // SET_LIST_FILTER: 'SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
}

export default {
  Namespace: moduleConfig.namespace,
  Getters: ServiceGetters,
  Actions: ServiceActions,
  Mutations: ServiceMutations
}
