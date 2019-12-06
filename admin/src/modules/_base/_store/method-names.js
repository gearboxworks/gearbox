import { namespaceValues } from '../../_helpers'

export const BaseGetters = {
  FIND_BY: 'FIND_BY',
  LIST_OPTIONS: 'LIST_OPTIONS',
  LIST_FILTERED_BY: 'LIST_FILTERED_BY',
  LIST_FILTERED: 'LIST_FILTERED'
}

export const BaseActions = {
  LOAD_ALL: 'LOAD_ALL',
  LOAD_ONE: 'LOAD_ONE',
  CREATE: 'CREATE',
  UPDATE: 'UPDATE',
  DELETE: 'DELETE',
  SET_LIST_FILTER: 'SET_LIST_FILTER',
  SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
}

export const BaseMutations = {
  SET_ALL: 'SET_ALL',
  SET_ONE: 'SET_ONE',
  UPDATE: 'UPDATE',
  DELETE: 'DELETE',
  SET_LIST_FILTER: 'SET_LIST_FILTER',
  SET_LIST_FILTER_SORT_BY: 'SET_LIST_FILTER_SORT_BY',
  SET_LIST_FILTER_SORT_ASC: 'SET_LIST_FILTER_SORT_ASC'
}

export default {
  Getters: BaseGetters,
  Actions: BaseActions,
  Mutations: BaseMutations
}

/**
 * @type {{LIST_OPTIONS: string, LIST_FILTERED: string, LIST_FILTERED_BY: string, FIND_BY: string}}
 */
export function NamespacedBaseGetters (namespace) {
  return namespaceValues(BaseGetters, namespace)
}

/**
 * @type {{DELETE: string, LOAD_ALL: string, CREATE: string, LOAD_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export function NamespacedBaseActions (namespace) {
  return namespaceValues(BaseActions, namespace)
}

/**
 * @type {{DELETE: string, SET_ONE: string, SET_LIST_FILTER: string, SET_LIST_FILTER_SORT_ASC: string, UPDATE: string, SET_ALL: string, SET_LIST_FILTER_SORT_BY: string}}
 */
export function NamespacedBaseMutations (namespace) {
  return namespaceValues(BaseMutations, namespace)
}
