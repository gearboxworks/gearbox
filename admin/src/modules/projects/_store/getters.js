import NamespacedBaseGetters from '../../_base/_store/getters'
import moduleConfig from '../config'

import { ProjectGetters as Getters } from './method-names'
import { GearspecGetters } from '../../gearspecs/_store/method-names'
import { ServiceGetters } from '../../services/_store/method-names'

import { programFromServiceId } from '../../_helpers'

const BaseGetters = NamespacedBaseGetters(moduleConfig.namespace)

export default {
  ...BaseGetters,

  [Getters.LIST_FILTERED_BY]: (state, getters, rootState, rootGetters) => (fieldName, allowedValues) => {
    const valuesArray = Array.isArray(allowedValues) ? allowedValues : [allowedValues]
    // 'notes' and 'stack' are not included on purpose because simple comparison does not work on them
    let projects = []

    if (fieldName === 'stacks') {
      projects = state.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.gearspec_id.indexOf(val) > -1)))
    } else if (fieldName === 'programs') {
      projects = state.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => programFromServiceId(s.service_id) === val)))
    } else {
      projects = BaseGetters[Getters.LIST_FILTERED_BY.replace(moduleConfig.namespace)](state, getters, rootState, rootGetters)(fieldName, allowedValues)
    }

    return projects
  },

  [Getters.LIST_FILTERED]: (state, getters, rootState, rootGetters) => () => {
    let list = state.records
    const sortAscending = !!state.sorting.ascending
    // const sortBy = state.sortBy
    for (const field in state.filter) {
      const values = state.filter[field]
      if (values === 'all') {
        continue
      }

      // Project specific:

      if (field === 'states') {
        if (values.length === 3) {
          continue
        } else {
          if (values.indexOf('running') > -1) {
            list = list.filter(p => getters[Getters.LIST_FILTERED_BY]('enabled', true).includes(p))
          }
          if (values.indexOf('stopped') > -1) {
            list = list.filter(p => getters[Getters.LIST_FILTERED_BY]('enabled', false).includes(p))
          }
        }
        continue
      }
      list = (list.filter(p => getters[Getters.LIST_FILTERED_BY](field, values).includes(p)))
    }
    /**
     * TODO implement sorting by attributes
     */
    return list.concat().sort((a, b) => a.id > b.id ? (sortAscending ? 1 : -1) : (a.id === b.id) ? 0 : (sortAscending ? -1 : 1))
  },

  [Getters.PROJECT_STACK_ITEM_INDEX_BY]: (state, getters, rootState, rootGetters) => (project, fieldName, fieldValue) => {
    let memberIndex = -1
    project.attributes.stack.find((m, idx) => {
      /**
       * fieldName can be "service_id" or "gearspec_id"
       */
      if (m[fieldName] === fieldValue) {
        memberIndex = idx
        return true
      }
      return false
    })
    return memberIndex
  },

  [Getters.GEARS_GROUPED_BY_STACK]: (state, getters, rootState, rootGetters) => (project) => {
    /**
     * Effectively this returns project's stacks: gearspecs (and their currently-selected services) grouped by stack (i.e. indexed by stack_id)
     */
    var result = {}
    const stackItems = project.attributes.stack || []

    stackItems.forEach(stackItem => {
      if (stackItem.isRemoved) {
        return
      }
      const gearspec = rootGetters[GearspecGetters.FIND_BY]('id', stackItem.gearspec_id)

      if (gearspec) {
        if (typeof result[gearspec.attributes.stack_id] === 'undefined') {
          result[gearspec.attributes.stack_id] = []
        }
        const service = stackItem.service_id ? rootGetters[ServiceGetters.FIND_BY]('id', stackItem.service_id) : null
        /**
         * note, when there is no exact match, service will be null,
         * but we will try to find a good-enough match further down the road;
         * that's why we need to pass over the original serviceId
         */
        result[gearspec.attributes.stack_id].push({
          gearspecId: stackItem.gearspec_id,
          gearspec,
          serviceId: stackItem.service_id,
          service
        })
      }
    })

    /**
     * sort gears by gear role
     */
    Object.keys(result).forEach((stackId) => {
      result[stackId] = result[stackId].sort((a, b) => a.gearspec.attributes.role > b.gearspec.attributes.role ? 1 : (a.gearspec.attributes.role === b.gearspec.attributes.role) ? 0 : -1)
    })

    /**
     * sort stacks by stack id
     */
    return Object.keys(result).sort().reduce((r, key) => {
      // eslint-disable-next-line no-param-reassign
      r[key] = result[key]
      return r
    }, {})
  },

  [Getters.SERVICES_GROUPED_BY_GEARSPEC_ROLE]: (state, getters, rootState, rootGetters) => (project) => {
    var result = {}
    if (project.attributes.stack) {
      project.attributes.stack.forEach((stackMember, idx) => {
        // if (stackMember.isRemoved) {
        //   return
        // }
        const gearspec = rootGetters[GearspecGetters.FIND_BY]('id', stackMember.gearspec_id)
        const serviceId = rootGetters[GearspecGetters.FIND_COMPATIBLE_SERVICE](gearspec, stackMember.service_id)
        const service = serviceId ? rootGetters[ServiceGetters.FIND_BY]('id', serviceId) : null

        if (gearspec && service) {
          // console.log(result, gearspec.attributes.stack_id, gearspec.attributes.role)
          if (typeof result[gearspec.attributes.stack_id] === 'undefined') {
            result[gearspec.attributes.stack_id] = {
              isRemoved: stackMember.isRemoved || false
            }
          }
          result[gearspec.attributes.stack_id][gearspec.attributes.role] = service
        }
      })
    }
    return result
  },

  [Getters.UNUSED_STACKS]: (state, getters, rootState, rootGetters) => (project) => {
    const result = {}

    const projectStack = getters[Getters.SERVICES_GROUPED_BY_GEARSPEC_ROLE](project)

    rootState['stacks'].records.forEach(stack => {
      if (typeof projectStack[stack.id] === 'undefined') {
        result[stack.id] = { stack, isRemoved: false }
      } else if (projectStack[stack.id].isRemoved) {
        // TODO if the services in the removed stack are exactly the same as in the default version of it, show only one option!
        result[stack.id] = { stack, isRemoved: false, isDefault: true }
        result[stack.id + '(removed)'] = { stack, isRemoved: true }
      }
    })

    return result
  }
}
