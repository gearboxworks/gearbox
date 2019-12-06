import { NamespacedBaseGetters, NamespacedBaseActions, NamespacedBaseMutations } from '../../_base/_store/method-names'
import moduleConfig from '../config'

/**
 * @type {{LIST_OPTIONS: string, LIST_FILTERED: string, LIST_FILTERED_BY: string, FIND_BY: string}}
 */
export const DemoGetters = {
  ...NamespacedBaseGetters(moduleConfig.namespace),
  /**
   * The following are added by destructuring the result of the above function call
   */
  // FIND_BY: 'basedirs/FIND_BY',
  // LIST_OPTIONS: 'basedirs/LIST_OPTIONS',
  // LIST_FILTERED_BY: 'basedirs/LIST_FILTERED_BY',
  // LIST_FILTERED: 'basedirs/LIST_FILTERED',
  /**
   *  The following getters are specific to basedirs
   */
  DEMO_GETTER: 'items/DEMO_GETTER'
}

/**
 * @type {{DELETE: string, LOAD_ALL: string, CREATE: string, LOAD_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const DemoActions = {
  ...NamespacedBaseActions(moduleConfig.namespace),
  // LOAD_ALL: 'LOAD_ALL',
  // LOAD_ONE: 'LOAD_ONE',
  // CREATE: 'CREATE',
  // UPDATE: 'UPDATE',
  // DELETE: 'DELETE',
  // SET_LIST_FILTER: 'SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
  /**
   *
   */
  DEMO_ACTION: 'items/DEMO_ACTION'
}

/**
 * @type {{DELETE: string, SET_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_ALL: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const DemoMutations = {
  ...NamespacedBaseMutations(moduleConfig.namespace),
  // SET_ALL: 'basedirs/SET_ALL',
  // SET_ONE: 'basedirs/SET_ONE',
  // UPDATE: 'basedirs/UPDATE',
  // DELETE: 'basedirs/DELETE',
  // SET_LIST_FILTER: 'basedirs/SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'basedirs/SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'basedirs/SET_LIST_FILTER_SORT_ASC'
  /**
   *
   */
  DEMO_MUTATION: 'items/DEMO_MUTATION'
}

export default {
  Namespace: moduleConfig.namespace,
  Getters: DemoGetters,
  Actions: DemoActions,
  Mutations: DemoMutations
}
