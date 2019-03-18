import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import VueAxios from 'vue-axios'
import { getConfig as raxConfig } from 'retry-axios'
import HTTP from './http-common'

Vue.use(Vuex)
Vue.use(VueAxios, axios)

export default new Vuex.Store({
  /**
   * In strict mode, whenever Vuex state is mutated outside of mutation handlers, an error will be thrown.
   */
  strict: true,
  state: {
    projects: [],
    stacks: [],
    stack_members: [],
    gearStacks: {},
    gearRoles: {},
    gearServices: {},
    connectionStatus: {
      networkError: null,
      remainingRetries: 5
    },
    baseDirs: []
  },
  getters: {
    projectByName: (state) => (projectName) => {
      return state.projects.find(p => p.name === projectName)
    },
    projectBy: (state) => (fieldName, fieldValue) => {
      return state.projects.find(p => p[fieldName] === fieldValue)
    },
    stackRoles: (state) => (stack) => {
      var result = {}
      var key
      for (key in state.gearRoles) {
        if (state.gearRoles.hasOwnProperty(key) && (key.indexOf(stack) !== -1)) {
          result[key] = state.gearRoles[key]
        }
      }
      return result
    },
    stackServices: (state) => (stack) => {
      var result = {}
      var key
      for (key in state.gearServices) {
        if (state.gearServices.hasOwnProperty(key) && (key.indexOf(stack) !== -1)) {
          result[key] = state.gearServices[key]
        }
      }
      console.log(result)
      return result
    },
    baseDirsAsOptions: (state) => {
      const options = []
      for (const baseDirName in state.baseDirs) {
        if (!state.baseDirs.hasOwnProperty(baseDirName)) {
          continue
        }
        options.push({
          value: baseDirName,
          text: state.baseDirs[baseDirName].text
        })
      }
      return options
    }
  },
  actions: {
    loadBaseDirs ({ commit }) {
      try {
        HTTP.get(
          'basedirs',
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
          .then(r => r ? r.data.data : null)
          .then((baseDirs) => {
            if (baseDirs) {
              const bd = {}
              for (const dirName in baseDirs) {
                if (!baseDirs.hasOwnProperty(dirName)) {
                  continue
                }

                bd[dirName] = {
                  'value': dirName,
                  'text': baseDirs[dirName]
                }
              }

              commit('SET_BASEDIRS', bd)
            // this.dispatch('loadProjectDetails')
            }
          })
      } catch (e) {
        console.log(e)
      }
    },
    loadProjectHeaders ({ commit }) {
      try {
        HTTP.get(
          'projects/with-details',
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
          .then(r => r ? r.data.data : null)
          .then((projects) => {
            if (projects) {
              // console.log(projects)

              const p = []

              for (const hostname in projects) {
                if (!projects.hasOwnProperty(hostname)) {
                  continue
                }
                let project = projects[hostname]
                let data = project.data

                p.push(
                  {
                    baseDir: data.basedir,
                    path: data.path,
                    hostname: data.hostname,
                    fullPath: data.project_dir,
                    enabled: data.enabled,
                    notes: data.notes,
                    aliases: data.aliases,
                    stack: data.stack
                  }
                )
              }

              commit('SET_PROJECTS', p)
              // this.dispatch('loadProjectDetails')
            }
          })
      } catch (e) {
        console.log(e)
      }
    },
    loadProjectDetails ({ commit }) {
      this.state.projects.forEach((project, index) => {
        try {
          HTTP.get(
            'projects/' + project.hostname,
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
            .then(r => r ? r.data.data : null)
            .then((p) => {
              // console.log(projects)

              const project = {
                path: p.hostname,
                enabled: p.enabled
              }
              commit('SET_PROJECT_DETAILS', project)
            })
        } catch (e) {
          console.log(e)
        }
      })
    },
    loadStacks ({ commit }) {
      HTTP.get(
        'stacks',
        { crossDomain: true }
      ).catch((error) => {
        // handle error
        // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        console.log('rejected', error)
      })
        .then(r => r ? r.data.data : null)
        .then((stacks) => {
          if (stacks) {
            const s = []
            for (const stackName in stacks) {
              if (!stacks.hasOwnProperty(stackName)) {
                continue
              }
              let stack = stacks[stackName]
              s.push(
                {
                  name: stackName,
                  label: stack.label,
                  examples: stack.examples,
                  stack: stack.stack,
                  optional: stack.optional,
                  shortLabel: stack.short_label,
                  memberType: stack.member_type
                }
              )
            }

            commit('SET_STACKS', stacks)
          }
        })
    },
    loadGears ({ commit }) {
      axios
        .get(
          'https://raw.githubusercontent.com/gearboxworks/gearbox/master/assets/gears.json',
          { crossDomain: true }
        )
        .catch((error) => {
          // handle error
          // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
        })
        .then(r => r.data)
        .then((data) => {
          commit('SET_GEAR_STACKS', data.stacks)
          commit('SET_GEAR_ROLES', data.roles)
          commit('SET_GEAR_SERVICES', data.role_services)
        })
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
    }
  },
  mutations: {
    /**
     * Names of mutation functions should be all-caps -- that's "idiomatic Vue"
     */
    SET_PROJECTS (state, projects) {
      state.projects = projects
    },
    SET_PROJECT_DETAILS (state, project) {
      const p = this.getters.projectBy('path', project.path)
      p.enabled = project.enabled
    },
    SET_STACKS (state, stacks) {
      state.stacks = stacks
    },
    SET_GEAR_STACKS (state, stacks) {
      state.gearStacks = stacks
    },
    SET_GEAR_ROLES (state, roles) {
      state.gearRoles = roles
    },
    SET_GEAR_SERVICES (state, services) {
      state.gearServices = services
    },
    UPDATE_PROJECT (state, args) {
      const { hostname, project } = args
      console.log(args)
      const p = this.getters.projectBy('hostname', hostname)
      p.hostname = project.hostname
      p.notes = project.notes
      p.baseDir = project.baseDir
      p.path = project.path
      p.fullPath = project.fullPath
      p.enabled = project.enabled
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
    }
  }

})
