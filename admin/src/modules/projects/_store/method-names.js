import { NamespacedBaseGetters, NamespacedBaseActions, NamespacedBaseMutations } from '../../_base/_store/method-names'
import moduleConfig from '../config'

/**
 * @type {{LIST_OPTIONS: string, LIST_FILTERED: string, LIST_FILTERED_BY: string, FIND_BY: string}}
 */
export const ProjectGetters = {
  ...NamespacedBaseGetters(moduleConfig.namespace),
  /**
   * The following are added by destructuring the result of the above function call
   */
  // FIND_BY: 'FIND_BY',
  // LIST_OPTIONS: 'LIST_OPTIONS',
  // LIST_FILTERED_BY: 'LIST_FILTERED_BY',
  // LIST_FILTERED: 'LIST_FILTERED',
  /**
   *  The following getters are specific to projects
   */
  PROJECT_STACK_ITEM_INDEX_BY: 'projects/PROJECT_STACK_ITEM_INDEX_BY',
  SERVICES_GROUPED_BY_GEARSPEC_ROLE: 'projects/SERVICES_GROUPED_BY_GEARSPEC_ROLE',
  GEARS_GROUPED_BY_STACK: 'projects/GEARS_GROUPED_BY_STACK',
  UNUSED_STACKS: 'projects/UNUSED_STACKS'
}
/**
 * @type {{DELETE: string, LOAD_ALL: string, CREATE: string, LOAD_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const ProjectActions = {
  ...NamespacedBaseActions(moduleConfig.namespace),
  // LOAD_ALL: 'LOAD_ALL',
  // LOAD_ONE: 'LOAD_ONE',
  // CREATE: 'CREATE',
  // UPDATE: 'UPDATE',
  // DELETE: 'DELETE',
  // SET_LIST_FILTER: 'SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
  LOAD_ALL_DETAILS: 'projects/LOAD_ALL_DETAILS',
  ADD_STACK: 'projects/ADD_STACK',
  REMOVE_STACK: 'projects/REMOVE_STACK',
  UPDATE_NOTES: 'projects/UPDATE_NOTES',
  UPDATE_HOSTNAME: 'projects/UPDATE_HOSTNAME',
  UPDATE_STATE: 'projects/UPDATE_STATE',
  CHANGE_GEAR: 'projects/CHANGE_GEAR'
}

/**
 * @type {{DELETE: string, SET_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_ALL: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export const ProjectMutations = {
  ...NamespacedBaseMutations(moduleConfig.namespace),
  // SET_ALL: 'SET_ALL',
  // SET_ONE: 'SET_ONE',
  // UPDATE: 'UPDATE',
  // DELETE: 'DELETE',
  // SET_LIST_FILTER: 'SET_LIST_FILTER',
  // SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  // SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
  SET_STACK: 'projects/SET_STACK',
  UPDATE_HOSTNAME: 'projects/UPDATE_HOSTNAME',
  UPDATE_STATE: 'projects/UPDATE_STATE',
  UPDATE_NOTES: 'projects/UPDATE_NOTES',
  ADD_STACK: 'projects/ADD_STACK',
  REMOVE_STACK: 'projects/REMOVE_STACK',
  CHANGE_GEAR: 'projects/CHANGE_GEAR'
}

export default {
  Namespace: moduleConfig.namespace,
  Getters: ProjectGetters,
  Actions: ProjectActions,
  Mutations: ProjectMutations
}
