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
    gears: [],
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
          .then((basedirs) => {
            if (basedirs) {
            // console.log(projects)

              const bd = []

              for (const dir in basedirs) {
                if (!basedirs.hasOwnProperty(dir)) {
                  continue
                }

                bd.push(
                  {
                    name: dir,
                    path: basedirs[dir]
                  }
                )
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
      // axios
      //   .get(
      //     'http://127.0.0.1:9999/gears',
      //     { crossDomain: true },
      //   )
      //   .catch((error) => {
      //     // handle error
      //     // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
      //   })
      //   .then(r => r.data.data)
      //   .then((gears) => {
      //     commit('SET_GEARS', gears);
      //   });
    },
    updateProject ({ commit }, payload) {
      const { projectName, project } = payload

      commit('UPDATE_PROJECT', { projectName, project })

      HTTP({
        method: 'post',
        url: 'project/' + projectName,
        data: project
      }).then(r => r.data).then((project) => {
        // move commit here
        // resolve()
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
    SET_GEARS (state, gears) {
      state.gears = gears
    },
    UPDATE_PROJECT (state, args) {
      const { hostname, project } = args
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
    SET_BASEDIRS (state, basedirs) {
      state.baseDirs = basedirs
    }
  }

})
