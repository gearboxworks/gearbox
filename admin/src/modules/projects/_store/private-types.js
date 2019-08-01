import BaseTypes from '../../_base/_store/private-types'
import moduleConfig from '../config'

export const ExtraGetters = {
  PROJECT_STACK_ITEM_INDEX_BY: 'PROJECT_STACK_ITEM_INDEX_BY',
  SERVICES_GROUPED_BY_GEARSPEC_ROLE: 'SERVICES_GROUPED_BY_GEARSPEC_ROLE',
  GEARS_GROUPED_BY_STACK: 'GEARS_GROUPED_BY_STACK',
  UNUSED_STACKS: 'UNUSED_STACKS'
}

export const ExtraActions = {
  LOAD_ALL_DETAILS: 'LOAD_ALL_DETAILS',
  ADD_STACK: 'ADD_STACK',
  REMOVE_STACK: 'REMOVE_STACK',
  UPDATE_NOTES: 'UPDATE_NOTES',
  UPDATE_HOSTNAME: 'UPDATE_HOSTNAME',
  UPDATE_STATE: 'UPDATE_STATE',
  CHANGE_GEAR: 'CHANGE_GEAR'
}

export const ExtraMutations = {
  SET_STACK: 'SET_STACK',
  UPDATE_HOSTNAME: 'UPDATE_HOSTNAME',
  UPDATE_STATE: 'UPDATE_STATE',
  UPDATE_NOTES: 'UPDATE_NOTES',
  ADD_STACK: 'ADD_STACK',
  REMOVE_STACK: 'REMOVE_STACK',
  CHANGE_GEAR: 'CHANGE_GEAR'
}

export default {
  Namespace: moduleConfig.namespace,
  GetterTypes: { ...BaseTypes.GetterTypes, ...ExtraGetters },
  ActionTypes: { ...BaseTypes.ActionTypes, ...ExtraActions },
  MutationTypes: { ...BaseTypes.MutationTypes, ...ExtraMutations }
}
