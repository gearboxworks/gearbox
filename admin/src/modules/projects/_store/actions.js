import api from '../_api'
import { Getters, Actions, Mutations } from './private-types'

const FORCED_DELAY = 2000

const actions = {

  [Actions.LOAD_ALL_HEADERS] ({ commit }, payload) {
    return api.fetchAllHeaders()
      .then((projects) => {
        commit(Mutations.SET_RECORDS, projects)
      }).catch((error) => {
        // eslint-disable-next-line
        console.error(error)
      })
  },

  [Actions.LOAD_DETAILS] ({ commit }, project) {
    return new Promise((resolve, reject) => {
      return api.fetchDetails(project.id)
        .then((projectWithDetails) => {
          if (project.id === projectWithDetails.id) {
            commit(Mutations.SET_STACK, { project, stack: projectWithDetails.attributes.stack })
            resolve(project)
          } else {
            reject(new Error(`Unexpected mismatch in project id: ${project.id} vs ${projectWithDetails.id}!`))
          }
        })
        // .catch(){} // do not do anything here about the rejected calls
    })
  },

  [Actions.LOAD_DETAILS_FOR_ALL] ({ state, dispatch }) {
    /**
     * Note we convert rejection to a resolved error object to get Promise.all to return ALL projects even if some of them cannot be fetched
     * @see https://stackoverflow.com/questions/31424561/wait-until-all-es6-promises-complete-even-rejected-promises#answer-36115549
     */
    const promises = state.records.map((project, idx) => dispatch(Actions.LOAD_DETAILS, project).catch(e => e))
    /**
     * wait for all
     */
    return Promise.all(promises)
  },

  [Actions.UPDATE_DETAILS] ({ commit }, payload) {
    const { projectId, data } = payload

    return new Promise((resolve) => {
      return api.updateDetails(projectId, data)
        .then((project) => {
        /**
         * on success, commit the new project to the projects store before resolving
         */
          if (project) {
            commit(Mutations.UPDATE_DETAILS, payload)
          }
          resolve(project)
        })
    })
  },

  [Actions.ADD_STACK] ({ commit, getters, rootGetters }, payload) {
    console.log('TODO: call API method to add project stack')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const { stackId } = payload
        const actualStackId = stackId.replace('(removed)', '')
        const stack = rootGetters.stackBy('id', actualStackId)

        commit(Mutations.ADD_STACK, {
          ...payload,
          actualStackId,
          stack,
          preselectServiceId: rootGetters.preselectServiceId
        })

        resolve()
      }, FORCED_DELAY)
    })
  },

  [Actions.REMOVE_STACK] ({ commit, getters }, payload) {
    console.log('TODO: call API method to remove project stack')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit(Mutations.REMOVE_STACK, payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  [Actions.UPDATE_NOTES] ({ commit, getters }, payload) {
    console.log('TODO: call API method to update project notes')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit(Mutations.UPDATE_NOTES, payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  [Actions.UPDATE_HOSTNAME]  ({ commit }, payload) {
    console.log('TODO: call API method to update project hostname')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit(Mutations.UPDATE_HOSTNAME, payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  [Actions.UPDATE_STATE] ({ commit, getters }, payload) {
    /**
     * TODO: call the API and commit when it returns
     */
    console.log('TODO: call the API to change project state')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        commit(Mutations.UPDATE_STATE, payload)
        resolve()
      }, FORCED_DELAY)
    })
  },

  [Actions.CHANGE_GEAR] ({ commit, getters }, payload) {
    /**
     * TODO: call the API and commit when it returns
     * TODO: remove delay
     */

    const { project, gearspecId } = payload
    const memberIndex = getters[Getters.PROJECT_STACK_ITEM_INDEX_BY](project, 'gearspec_id', gearspecId)

    return new Promise((resolve, reject) => {
      setTimeout(
        () => {
          commit(Mutations.CHANGE_GEAR, { ...payload, memberIndex })
          resolve()
        }, 100)
    })
  },

  [Actions.SET_LIST_FILTER] ({ commit }, payload) {
    commit(Mutations.SET_LIST_FILTER, payload)
  }
}

export default actions
