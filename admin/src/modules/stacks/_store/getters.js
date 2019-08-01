import BaseGetters from '../../_base/_store/getters'

import StoreMethodTypes from './private-types'
const { GetterTypes: Getters } = StoreMethodTypes

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

  // [Getters.LIST_FILTERED_BY]: (state) => (state) => (fieldName, allowedValues, allowedFields) => UNSUPPORTED_GETTER(),
  // [Getters.LIST_FILTERED]: (state, getters) => UNSUPPORTED_GETTER(),

  [Getters.LIST_OPTIONS]: (state, getters, rootState, rootGetters) => () => BaseGetters[Getters.LIST_OPTIONS](state, getters, rootState, rootGetters)('basedir'),

  /**
   * As an example, for php:7.1.18 it will select php:7.1 or php:7 if exact match is not possible
   */
  [Getters.FIND_COMPATIBLE_SERVICE]: (state, getters, rootState, rootGetters) => (stackIdOrStack, gearspecId, requestedServiceId) => {
    const stack = (typeof stackIdOrStack === 'string')
      ? getters[Getters.FIND_BY]('id', stackIdOrStack)
      : stackIdOrStack

    const defaultServiceId = requestedServiceId || getters[Getters.DEFAULT_SERVICE_FOR_GEARSPEC](stack, gearspecId)
    const serviceIds = getters[Getters.GEARSPEC_SERVICES](stack, gearspecId)

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

  [Getters.GEARSPEC_SERVICES]: (state, getters, rootState, rootGetters) => (stackIdOrStack, gearspecId) => {
    const stack = (typeof stackIdOrStack === 'string')
      ? getters[Getters.FIND_BY]('id', stackIdOrStack)
      : stackIdOrStack

    let services = []
    stack.attributes.members.find((m, idx) => {
      if (m.gearspec_id === gearspecId) {
        services = m.services
        return true
      }
      return false
    })
    return services
  },

  [Getters.DEFAULT_SERVICE_FOR_GEARSPEC]: (state, getters, rootState, rootGetters) => (stackIdOrStack, gearspecId) => {
    const stack = (typeof stackIdOrStack === 'string')
      ? getters[Getters.FIND_BY]('id', stackIdOrStack)
      : stackIdOrStack

    let defaultService = ''

    if (stack) {
      stack.attributes.members.find((m, idx) => {
        if (m.gearspec_id === gearspecId) {
          defaultService = m.default_service
          return true
        }
        return false
      })
    }
    return defaultService
  }
}

export default { ...BaseGetters, ...OverrideGetters }
