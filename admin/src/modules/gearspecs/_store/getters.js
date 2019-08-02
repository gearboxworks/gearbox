import BaseGetters from '../../_base/_store/getters'
import { UNSUPPORTED_GETTER } from '../../_helpers'

import GearspecMethodTypes from './private-types'
import StackMethodTypes from '../../stacks/_store/public-types'
const { GetterTypes: Getters } = GearspecMethodTypes
const { GetterTypes: StackGetters } = StackMethodTypes

const OverrideGetters = {

  // [Getters.FIND_BY]: (state) => (fieldName, fieldValue) => {
  //   /**
  //    * manipulate arguments
  //    */
  //   const results = BaseGetters[Getters.FIND_BY](state)(fieldName, fieldValue)
  //   /**
  //    * manipulate results
  //    */
  //   return results
  // },

  // [Getters.DEMO_GETTER]: (state) => {
  //   return 'This is the result from DEMO_GETTER.'
  // },

  // [Getters.LIST_FILTERED_BY]: (state) => (state) => (fieldName, allowedValues) => UNSUPPORTED_GETTER(),
  // [Getters.LIST_FILTERED]: (state, getters) => UNSUPPORTED_GETTER(),

  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => () => BaseGetters[Getters.LIST_OPTIONS](state, getters, rootState, rootGetters)('role'),

  [Getters.GEARSPEC_SERVICES]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? getters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const stack = rootGetters[StackGetters.FIND_BY]('id', gearspec.attributes.stack_id)

    let services = []
    stack.attributes.members.find((m, idx) => {
      if (m.gearspec_id === gearspec.id) {
        services = m.services
        return true
      }
      return false
    })
    return services
  },

  [Getters.DEFAULT_GEARSPEC_SERVICE]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? getters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const stack = rootGetters[StackGetters.FIND_BY]('id', gearspec.attributes.stack_id)

    let defaultService = ''

    if (stack) {
      stack.attributes.members.find((m, idx) => {
        if (m.gearspec_id === gearspec.id) {
          defaultService = m.default_service
          return true
        }
        return false
      })
    }
    return defaultService
  },

  /**
   * As an example, for php:7.1.18 it will select php:7.1 or php:7 if exact match is not possible
   */
  [Getters.FIND_COMPATIBLE_SERVICE]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId, requestedServiceId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? getters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const serviceIds = getters[Getters.GEARSPEC_SERVICES](gearspec)
    const defaultServiceId = requestedServiceId || getters[Getters.DEFAULT_GEARSPEC_SERVICE](gearspec)

    /**
     * Resolve default option:
     * - if exact match is found, use it
     * - otherwise, use the last in the list that have the specified name mentioned (hopefully that will be the latest version)
     */

    let firstFound = -1
    let exactFound = -1
    let serviceId = defaultServiceId
    if (serviceId) {
      do {
        for (var i = serviceIds.length; i--;) {
          if (serviceIds[i].indexOf(serviceId) !== -1) {
            if (firstFound === -1) {
              firstFound = i
            }
            if (serviceIds[i] === serviceId) {
              exactFound = i
              break
            }
          }
        }
        if (firstFound === -1) {
          /**
           * drop the part after the last dot
           */
          const parts = serviceId.split('.')
          if (parts.length > 1) {
            delete parts[parts.length - 1]
            serviceId = parts.join('.')
            serviceId = serviceId.substring(0, serviceId.length - 1) // remove the trailing dot
          } else {
            break
          }
        }
      } while (firstFound === -1)
    }
    const selectedService = (firstFound !== -1)
      ? serviceIds[exactFound !== -1 ? exactFound : firstFound]
      : ''

    return selectedService
  },

  [Getters.GEARSPEC_SERVICE_VERSIONS_GROUPED_BY_PROGRAM]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? getters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const services = getters[Getters.GEARSPEC_SERVICES](gearspec)
    const result = {}
    services.forEach(serviceId => {
      /**
       * TODO move this parsing logic to a global helper function
       */
      const programName = serviceId.split(':')[0].replace('gearboxworks/', '')
      if (typeof result[programName] === 'undefined') {
        result[programName] = {}
      }
      const serviceNameAndVersion = serviceId.replace('gearboxworks/', '')
      result[programName][serviceId] = serviceNameAndVersion
    })
    return result
  }

}

export default { ...BaseGetters, ...OverrideGetters }
