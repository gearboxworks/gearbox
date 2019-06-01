import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { getConfig as raxConfig } from 'retry-axios'
import HTTP from './http-common'
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
        'projects',
        'basedirs'
      ],
      httpClient: HTTP
    })
  },
  state: {
    stacks: [],
    services: [],
    gearspecs: [],
    projects: [],
    basedirs: [],
    connectionStatus: {
      networkError: null,
      remainingRetries: 5
    },
    showProjectsHaving: {
      'states': ['running', 'stopped', 'candidates'],
      'basedir': 'all',
      'stacks': 'all',
      'programs': 'all'
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
    projectBy: (state) => (fieldName, fieldValue) => {
      return (fieldName === 'id')
        ? state.projects.records.find(p => p.id === fieldValue)
        : state.projects.records.find(p => p.attributes[fieldName] === fieldValue)
    },
    filterProjectsBy: (state) => (fieldName, allowedValues) => {
      const attrs = ['basedir', 'enabled', 'filepath', 'hostname', 'path', 'project_dir']
      let valuesArray = Array.isArray(allowedValues) ? allowedValues : [allowedValues]
      // 'notes' and 'stack' are not included on purpose because simple comparison does not work on them
      let projects = []

      if (fieldName === 'id') {
        projects = state.projects.records.filter(p => valuesArray.indexOf(p.id) !== -1)
      } else if (attrs.indexOf(fieldName) !== -1) {
        projects = state.projects.records.filter(p => valuesArray.indexOf(p.attributes[fieldName]) !== -1)
      } else if (fieldName === 'stacks') {
        projects = state.projects.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.gearspec_id.indexOf(val) > -1)))
      } else if (fieldName === 'programs') {
        projects = state.projects.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.service_id.split('/')[1].split(':')[0] === val)))
      }

      return projects
    },
    filteredProjects: (state, getters) => {
      let projects = state.projects.records
      for (const field in state.showProjectsHaving) {
        const values = state.showProjectsHaving[field]
        if (values === 'all') {
          continue
        }
        if (field === 'states') {
          if (values.length === 3) {
            continue
          } else {
            if (values.indexOf('running') > -1) {
              projects = projects.filter(p => getters.filterProjectsBy('enabled', true).includes(p))
            }
            if (values.indexOf('stopped') > -1) {
              projects = projects.filter(p => getters.filterProjectsBy('enabled', false).includes(p))
            }
            // TODO merge candidates into projects array
            // if (values.indexOf('candidates') > -1) {
            //   projects = projects.filter(p => getters.filterProjectsBy('candidate', true).includes(p))
            // }
          }
          continue
        }
        projects = projects.filter(p => getters.filterProjectsBy(field, values).includes(p))
      }
      return projects
    },
    projectStackItemIndexBy: (state) => (project, fieldName, fieldValue) => {
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
      return options
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
      return options
    },
    basedirsAsOptions: (state) => {
      const options = []
      state.basedirs.records.forEach((el, idx) => {
        options.push({
          value: el.id,
          text: el.attributes.basedir
        })
      })
      return options
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
    loadProjectDetails ({ commit }) {
      for (const idx in this.state.projects.records) {
        const project = this.state.projects.records[idx]
        try {
          HTTP.get(
            'projects/' + project.id,
            {
              crossDomain: true,
              raxConfig: {
                // You can detect when a retry is happening, and figure out how many
                // retry attempts have been made
                onRetryAttempt: (err) => {
                  const cfg = raxConfig(err)
                  commit('SET_NETWORK_ERROR', err.message)
                  commit('SET_REMAINING_RETRIES', cfg.retry - cfg.currentRetryAttempt)
                }
              }
            }
          ).catch((error, config) => {
            // handle error
            // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
            console.log('rejected', error)
            // if (error.message === 'Network Error') {
            //   commit('SET_NETWORK_ERROR', error.message)
            // }
          })
            .then(r => r ? r.data : null)
            .then(response => {
              const project = response.data
              commit('SET_PROJECT', project)
              if (response.included.length) {
                for (const idx in response.included) {
                  const item = response.included[idx]
                  if (item.type === 'service') {
                    commit('SET_SERVICE', item)
                  }
                  if (item.type === 'stack') {
                    commit('SET_STACK', item)
                  }
                }
              }
            })
        } catch (e) {
          console.log(e)
        }
      }
    },
    updateProject ({ commit }, payload) {
      const { hostname, project } = payload

      commit('UPDATE_PROJECT', { hostname, project })

      HTTP({
        method: 'post',
        url: 'project/' + hostname,
        data: project
      }).then(r => r.data).then((project) => {
        // move commit here
        // resolve()
      }).catch((error) => {
        console.log('rejected', error)
        // resolve();
      })
    },
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
    addProjectStack ({ commit }, payload) {
      /**
       * TODO: call the API and commit when it returns
       */
      commit('ADD_PROJECT_STACK', payload)
    },
    removeProjectStack ({ commit }, payload) {
      /**
       * TODO: call the API and commit when it returns
       */
      commit('REMOVE_PROJECT_STACK', payload)
    },
    changeProjectService ({ commit }, payload) {
      /**
       * TODO: call the API and commit when it returns
       * TODO: remove delay
       */
      setTimeout(() => commit('CHANGE_PROJECT_SERVICE', payload), 1000)
      // commit('CHANGE_PROJECT_SERVICE', payload)
    },
    changeProjectState ({ commit }, payload) {
      /**
       * TODO: call the API and commit when it returns
       */
      commit('CHANGE_PROJECT_STATE', payload)
    },
    setProjectsFilter ({ commit }, payload) {
      commit('SET_PROJECTS_FILTER', payload)
    }
  },
  mutations: {
    /**
     * Names of mutation functions should be all-caps -- that's "idiomatic Vue"
     */
    SET_PROJECT (state, project) {
      const p = this.getters.projectBy('id', project.id)
      if (!p) {
        state.projects.records.push(project)
      } else {
        Vue.set(p.attributes, 'stack', project.attributes.stack)
      }
    },
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
    },
    ADD_PROJECT_STACK (state, payload) {
      const { projectId, stackId } = payload
      const project = this.getters.projectBy('id', projectId)
      const stack = this.getters.stackBy('id', stackId)
      if (project && stack && stack.attributes.members.length) {
        if (typeof project.attributes.stack === 'undefined') {
          Vue.set(project.attributes, 'stack', [])
        }
        stack.attributes.members.forEach((el, idx) => {
          const serviceId = this.getters.preselectServiceId(el.services, el.default_service)
          if (el.gearspec_id) {
            // reactive!
            project.attributes.stack.push({
              service_id: serviceId, // it's ok if it is empty
              gearspec_id: el.gearspec_id
            })
          }
        })
      }
    },
    REMOVE_PROJECT_STACK (state, payload) {
      const { projectId, stackId } = payload
      const project = this.getters.projectBy('id', projectId)
      if (project) {
        /**
         * We need to remove all elements of project.stack that that have service_id starting with shortStackName, e.g. "wordpress/"
         *
         * For deleting array items in javascript with forEach() and splice())
         * @see https://gist.github.com/chad3814/2924672
         */
        const shortStackName = stackId.split('/')[1]
        for (let i = project.attributes.stack.length - 1; i >= 0; i--) {
          if (project.attributes.stack[i].gearspec_id.split('/')[1] === shortStackName) {
            Vue.delete(project.attributes.stack, i)
          }
        }
      }
    },
    CHANGE_PROJECT_SERVICE (state, payload) {
      /**
       * Payload is of this form:
       * {projectId: "project1", gearspecId: "gearbox.works/wordpress/webserver", serviceId: "gearboxworks/apache:2.4"}
       */
      const { projectId, gearspecId, serviceId } = payload
      const project = this.getters.projectBy('id', projectId)

      if (project) {
        const memberIndex = this.getters.projectStackItemIndexBy(project, 'gearspec_id', gearspecId)
        /**
         * note, serviceId might be an empty string (and that's OK)
         */
        Vue.set(project.attributes.stack[memberIndex], 'service_id', serviceId)
      }
    },
    CHANGE_PROJECT_STATE (state, payload) {
      const { projectId, isEnabled } = payload
      const project = this.getters.projectBy('id', projectId)
      if (project) {
        project.attributes.enabled = !!isEnabled
      }
    },
    SET_PROJECTS_FILTER (state, payload) {
      const { field, values } = payload
      Vue.set(state.showProjectsHaving, field, values)
    }
  }

})
