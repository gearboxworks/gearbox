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
        'projects'
      ],
      httpClient: HTTP
    })
  },
  state: {
    stacks: [],
    services: [],
    gearspecs: [],
    projects: [],
    baseDirs: {
      'primary': {
        text: '~/Sites'
      }
    },
    connectionStatus: {
      networkError: null,
      remainingRetries: 5
    }
  },
  getters: {
    stackBy: (state) => (fieldName, fieldValue) => {
      let item = null
      if (fieldName === 'id') {
        item = state.stacks.records.find(p => p.id === fieldValue)
      } else {
        item = state.stacks.records.find(p => p.attributes[fieldName] === fieldValue)
      }
      return item
    },
    serviceBy: (state) => (fieldName, fieldValue) => {
      let item = null
      if (fieldName === 'id') {
        item = state.services.records.find(p => p.id === fieldValue)
      } else {
        item = state.services.records.find(p => p.attributes[fieldName] === fieldValue)
      }
      return item
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
    projectStackMemberIndexBy: (state) => (project, fieldName, fieldValue) => {
      let memberIndex = -1
      project.attributes.stack.find((m, idx) => {
        if (m[fieldName] === fieldValue) {
          memberIndex = idx
          return true
        }
        return false
      })
      return memberIndex
    },
    baseDirsAsOptions: (state) => {
      const options = [{
        'value': 'primary',
        'text': '~/Sites'
      }]
      return options
      // for (const baseDirName in state.baseDirs) {
      //   if (!state.baseDirs.hasOwnProperty(baseDirName)) {
      //     continue
      //   }
      //   options.push({
      //     value: baseDirName,
      //     text: state.baseDirs[baseDirName].text
      //   })
      // }
      // return options
    },
    projectServiceDefaults: (state) => (serviceName, stackName) => {
      const org = serviceName.substring(0, serviceName.indexOf('/'))
      const service = state.gearServices[serviceName]

      /**
       * Resolve default option:
       * - if exact match is found, use it
       * - otherwise, use the last in the list that have the specified name mentioned (hopefully that will be the latest version)
       */
      let firstFound = -1
      let exactFound = -1
      if (service.default) {
        for (var i = service.options.length; i--;) {
          if (service.options[i].indexOf(service.default) !== -1) {
            if (firstFound === -1) {
              firstFound = i
            }
            if (service.options[i] === service.default) {
              exactFound = i
              break
            }
          }
        }
      }
      const defaultOpt = (firstFound !== -1)
        ? service.options[ exactFound !== -1 ? exactFound : firstFound ]
        : ''

      const ver = defaultOpt ? defaultOpt.split(':')[1].split('.') : ''

      const newService = {
        'authority': org,
        'org': org.replace('.', ''),
        'stack': stackName.substring(stackName.indexOf('/') + 1),
        'service_id': defaultOpt ? service.org + '/' + defaultOpt : '',
        'program': defaultOpt ? defaultOpt.substring(0, defaultOpt.indexOf(':')) : '',
        'version': {}
      }

      if (ver.length > 0) {
        newService.version.major = ver[0]
      }
      if (ver.length > 1) {
        newService.version.minor = ver[1]
      }
      if (ver.length > 2) {
        newService.version.patch = ver[2]
      }
      return newService
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
       */
      commit('CHANGE_PROJECT_SERVICE', payload)
    },
    changeProjectState ({ commit }, payload) {
      /**
       * TODO: call the API and commit when it returns
       */
      commit('CHANGE_PROJECT_STATE', payload)
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
        // console.log('Setting stack for project', project.id, project.attributes.stack.length)
        Vue.set(p.attributes, 'stack', project.attributes.stack)
        // p.attributes = project.attributes
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
    SET_GEARSTACK (state, gearstack) {
      const g = this.getters.gearstackBy('id', gearstack.id)
      if (!g) {
        state.gearstacks.records.push(gearstack)
      } else {
        g.attributes = gearstack.attributes
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
    SET_BASEDIRS (state, baseDirs) {
      state.baseDirs = baseDirs
    },
    ADD_BASEDIR (state, baseDir) {
      state.baseDirs[baseDir.value] = baseDir
    },
    ADD_PROJECT_STACK (state, payload) {
      const { projectId, stackId } = payload
      const project = this.getters.projectBy('id', projectId)
      const genericServiceName = stackId.substring(stackId.indexOf('/') + 1)
      if (project) {
        for (const serviceName in this.getters.stackBy('id', stackId)) {
          Vue.set(project.attributes.stack, genericServiceName, this.projectServiceDefaults(stackId, serviceName))
        }
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
      const { projectId, serviceId, gearspecId } = payload

      const project = this.getters.projectBy('id', projectId)

      if (project) {
        const memberIndex = this.getters.projectStackMemberIndexBy(project, 'gearspec_id', gearspecId)
        if (!serviceId) {
          Vue.delete(project.attributes.stack, memberIndex)
        } else {
          Vue.set(project.attributes.stack[memberIndex], 'service_id', serviceId)
        }
      }
    },
    CHANGE_PROJECT_STATE (state, payload) {
      const { projectId, isEnabled } = payload
      const project = this.getters.projectBy('id', projectId)
      if (project) {
        // console.log(project.enabled)
        Vue.set(project.attributes, 'enabled', !!isEnabled)
        // console.log(project.enabled)
      }
    }
  }

})
