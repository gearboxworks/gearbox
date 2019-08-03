import NamespacedBaseGetters from '../../_base/_store/getters'
import moduleConfig from '../config'

import { GearspecGetters as Getters } from './method-names'
import { StackGetters } from '../../stacks/_store/method-names'

export default {
  ...NamespacedBaseGetters(moduleConfig.namespace),

  [Getters.STACK]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? rootGetters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    return rootGetters[StackGetters.FIND_BY]('id', gearspec.attributes.stack_id)
  },

  [Getters.SERVICES]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? rootGetters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const stack = rootGetters[Getters.STACK](gearspec)

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

  [Getters.DEFAULT_SERVICE]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? rootGetters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const stack = rootGetters[Getters.STACK](gearspec)

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
    const serviceIds = rootGetters[Getters.SERVICES](gearspecOrGearspecId)
    const defaultServiceId = requestedServiceId || rootGetters[Getters.DEFAULT_SERVICE](gearspecOrGearspecId)

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

  [Getters.SERVICE_VERSIONS_GROUPED_BY_PROGRAM]: (state, getters, rootState, rootGetters) => (gearspecOrGearspecId) => {
    const gearspec = (typeof gearspecOrGearspecId === 'string')
      ? rootGetters[Getters.FIND_BY]('id', gearspecOrGearspecId)
      : gearspecOrGearspecId

    const services = rootGetters[Getters.SERVICES](gearspec)
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
