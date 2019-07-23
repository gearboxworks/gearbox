import api from '../_api'

const FORCED_DELAY = 2000

const actions = {

  loadAllHeaders ({ commit }, payload) {
    return api.fetchAllHeaders()
      .then((projects) => {
        commit('SET_RECORDS', projects)
      }).catch((error) => {
        // eslint-disable-next-line
        console.error(error)
      })
  },

  loadDetails ({ commit }, project) {
    return new Promise((resolve, reject) => {
      return api.fetchDetails(project.id)
        .then((projectWithDetails) => {
          if (project.id === projectWithDetails.id) {
            commit('SET_PROJECT_STACK', { project, stack: projectWithDetails.attributes.stack })
            resolve(project)
          } else {
            reject(new Error(`Unexpected mismatch in project id: ${project.id} vs ${projectWithDetails.id}!`))
          }
        })
        // .catch(){} // do not do anything here about the rejected calls
    })
  },

  loadDetailsForAll ({ state, dispatch }) {
    /**
     * Note we convert rejection to a resolved error object to get Promise.all to return ALL projects even if some of them cannot be fetched
     * @see https://stackoverflow.com/questions/31424561/wait-until-all-es6-promises-complete-even-rejected-promises#answer-36115549
     */
    const promises = state.records.map((project, idx) => dispatch('loadDetails', project).catch(e => e))
    /**
     * wait for all
     */
    return Promise.all(promises)
  },

  updateDetails ({ commit }, payload) {
    const { projectId, data } = payload

    return new Promise((resolve) => {
      return api.updateDetails(projectId, data)
        .then((project) => {
        /**
         * on success, commit the new project to the projects store before resolving
         */
          if (project) {
            commit('UPDATE_PROJECT_DETAILS', payload)
          }
          resolve(project)
        })
    })
  },

  addStack ({ commit, getters, rootGetters }, payload) {
    console.log('TODO: call API method to add project stack')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const { stackId } = payload
        const actualStackId = stackId.replace('(removed)', '')
        const stack = rootGetters.stackBy('id', actualStackId)

        commit('ADD_PROJECT_STACK', { ...payload, actualStackId, stack, preselectServiceId: rootGetters.preselectServiceId })
        resolve()
      }, FORCED_DELAY)
    })
  },

  removeStack ({ commit, getters }, payload) {
    console.log('TODO: call API method to remove project stack')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit('REMOVE_PROJECT_STACK', payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  updateNotes ({ commit, getters }, payload) {
    console.log('TODO: call API method to update project notes')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit('UPDATE_PROJECT_NOTES', payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  updateHostname ({ commit }, payload) {
    console.log('TODO: call API method to update project hostname')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit('UPDATE_PROJECT_HOSTNAME', payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  updateState ({ commit, getters }, payload) {
    /**
     * TODO: call the API and commit when it returns
     */
    console.log('TODO: call the API to change project state')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit('UPDATE_PROJECT_STATE', payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  changeService ({ commit, getters }, payload) {
    /**
     * TODO: call the API and commit when it returns
     * TODO: remove delay
     */

    const { project, gearspecId } = payload
    const memberIndex = getters.projectStackItemIndexBy(project, 'gearspec_id', gearspecId)

    setTimeout(() => commit('CHANGE_PROJECT_SERVICE', { ...payload, memberIndex }), FORCED_DELAY)
    // commit('CHANGE_PROJECT_SERVICE', payload)
  },

  setProjectsFilter ({ commit }, payload) {
    commit('SET_PROJECTS_FILTER', payload)
  }
}

export default actions
