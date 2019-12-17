import NamespacedBaseActions from '../../_base/_store/actions'
import moduleConfig from '../config'
import api from '../_api'
import { FORCED_DELAY } from '../../_helpers'

import { ProjectGetters as Getters, ProjectActions as Actions, ProjectMutations as Mutations } from './method-names'
import { StackGetters } from '../../stacks/_store/method-names'

export default {
  ...NamespacedBaseActions(api, moduleConfig.namespace),

  [Actions.LOAD_ONE] ({ commit }, project) {
    return api.fetchOne(project.id)
      .then((projectWithDetails) => {
        if (project.id === projectWithDetails.id) {
          commit(Mutations.SET_STACK, {
            project,
            stack: projectWithDetails.attributes.stack
          })
          return (project)
        } else {
          throw new Error(`Unexpected mismatch in project id: ${project.id} vs ${projectWithDetails.id}!`)
        }
      })
  },

  [Actions.LOAD_ALL_DETAILS] ({ state, dispatch }) {
    /**
     * Note we convert rejection to a resolved error object to get Promise.all to return ALL projects even if some of them cannot be fetched
     * @see https://stackoverflow.com/questions/31424561/wait-until-all-es6-promises-complete-even-rejected-promises#answer-36115549
     */
    const promises = state.records.map((project, idx) => dispatch(Actions.LOAD_ONE, project).catch(e => e))
    /**
     * wait for all
     */
    return Promise.all(promises)
  },

  [Actions.ADD_STACK] ({ commit, getters, rootGetters }, payload) {
    console.log('TODO: call API method to add project stack')

    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const { stackId } = payload
        const actualStackId = stackId.replace('(removed)', '')
        const stack = rootGetters[StackGetters.FIND_BY]('id', actualStackId)

        commit(Mutations.ADD_STACK, {
          ...payload,
          actualStackId,
          stack
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
  }
}
