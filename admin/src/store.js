import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import VueAxios from 'vue-axios'
// import { getConfig as raxConfig } from 'retry-axios'
import HTTP from './http-common'
// import projectsStore from './modules/projects/_store'
import { mapResourceModules } from '@reststate/vuex'

Vue.use(Vuex)
Vue.use(VueAxios, axios)

export default new Vuex.Store({
  /**
   * In strict mode, whenever Vuex state is mutated outside of mutation handlers, an error will be thrown.
   */
  strict: true,
  modules: {
    ...mapResourceModules({
      names: [
        'stacks',
        'services',
        'gearspecs',
        'basedirs'
      ],
      httpClient: HTTP
    })
    // projects: projectsStore
  },
  state: {
    stacks: [],
    removedStacks: {},
    services: [],
    gearspecs: [],
    // projects: [],
    basedirs: [],
    connectionStatus: {
      networkError: null,
      remainingRetries: 5
    }
  },
  getters: {
    basedirBy: (state) => (fieldName, fieldValue) => {
      return (fieldName === 'id')
        ? state.basedirs.records.find(p => p.id === fieldValue)
        : state.basedirs.records.find(p => p.attributes[fieldName] === fieldValue)
    },
    stackBy: (state) => (fieldName, fieldValue) => {
      return (fieldName === 'id')
        ? state.stacks.records.find(p => p.id === fieldValue)
        : state.stacks.records.find(p => p.attributes[fieldName] === fieldValue)
    },
    serviceBy: (state) => (fieldName, fieldValue) => {
      return (fieldName === 'id')
        ? state.services.records.find(p => p.id === fieldValue)
        : state.services.records.find(p => p.attributes[fieldName] === fieldValue)
    },
    gearspecBy: (state) => (fieldName, fieldValue) => {
      return (fieldName === 'id')
        ? state.gearspecs.records.find(p => p.id === fieldValue)
        : state.gearspecs.records.find(p => p.attributes[fieldName] === fieldValue)
    },
    stackDefaultServiceByRole: (state) => (stack, gearspecId) => {
      let defaultService = ''
      stack.attributes.members.find((m, idx) => {
        if (m.gearspec_id === gearspecId) {
          defaultService = m.default_service
          return true
        }
        return false
      })
      return defaultService
    },
    stackServicesByRole: (state) => (stack, gearspecId) => {
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
    stacksAsOptions: (state) => {
      const options = []
      state.stacks.records.forEach((el, idx) => {
        options.push({
          value: el.id,
          text: el.attributes.stackname
        })
      })
      return options.sort((a, b) => a.value > b.value ? 1 : (a.value === b.value) ? 0 : -1)
    },
    servicesAsOptions: (state) => {
      const options = []
      state.services.records.forEach((el, idx) => {
        options.push({
          value: el.id,
          text: el.id
        })
      })
      return options
    },
    programsAsOptions: (state) => {
      const programs = []
      const options = []
      state.services.records.forEach((el, idx) => {
        const program = el.attributes.program
        if (programs.indexOf(program) === -1) {
          programs.push(program)
          options.push({
            value: program,
            text: program
          })
        }
      })
      return options.sort((a, b) => a.value > b.value ? 1 : (a.value === b.value) ? 0 : -1)
    },
    basedirsAsOptions: (state) => {
      const options = []
      state.basedirs.records.forEach((el, idx) => {
        options.push({
          value: el.id,
          text: el.attributes.basedir
        })
      })
      return options.sort((a, b) => a.value > b.value ? 1 : (a.value === b.value) ? 0 : -1)
    },
    hasExtraBasedirs: (state) => {
      return state.basedirs.records.length > 1
    },
    preselectServiceId: (state) => (serviceIds, defaultServiceId, providedServiceId) => {
      /**
       * Resolve default option:
       * - if exact match is found, use it
       * - otherwise, use the last in the list that have the specified name mentioned (hopefully that will be the latest version)
       */
      let firstFound = -1
      let exactFound = -1
      let serviceId = providedServiceId || defaultServiceId
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
        ? serviceIds[ exactFound !== -1 ? exactFound : firstFound ]
        : ''

      return selectedService
    }
  },
  actions: {

    addBaseDir ({ commit }, payload) {
      const { name, path } = payload
      commit('ADD_BASEDIR', {
        value: name,
        text: path
      })
      HTTP({
        method: 'post',
        url: 'basedirs/new',
        data: payload
      }).then(r => r.data).then((baseDir) => {
        // commit('ADD_BASEDIR', baseDir)
      }).catch((error) => {
        console.log('rejected', error)
        // resolve();
      })
    },
    getDirectory ({ commit }, payload) {
      return HTTP.head(
        'directories/' + encodeURI(payload.dir)
      )
    }
  },
  mutations: {
    SET_STACK (state, stack) {
      const s = this.getters.stackBy('id', stack.id)
      if (!s) {
        state.stacks.records.push(stack)
      } else {
        s.attributes = stack.attributes
      }
    },
    SET_SERVICE (state, service) {
      const s = this.getters.serviceBy('id', service.id)
      if (!s) {
        state.services.records.push(service)
      } else {
        s.attributes = service.attributes
      }
    },
    SET_GEARSPEC (state, gearspec) {
      const g = this.getters.gearspecBy('id', gearspec.id)
      if (!g) {
        state.gearspecs.records.push(gearspec)
      } else {
        g.attributes = gearspec.attributes
      }
    },
    SET_NETWORK_ERROR (state, message) {
      state.connectionStatus.networkError = message
    },
    CLEAR_NETWORK_ERROR (state) {
      state.connectionStatus.networkError = ''
    },
    SET_REMAINING_RETRIES (state, remainingRetries) {
      state.connectionStatus.remainingRetries = remainingRetries
    },
    ADD_BASEDIR (state, baseDir) {
      state.baseDirs[baseDir.value] = baseDir
    }
  }

})
