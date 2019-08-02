import Vue from 'vue'
import BaseMutations from '../../_base/_store/mutations'
// import { UNSUPPORTED_MUTATION } from '../../_helpers'

import store from '../../../store'
import ProjectMethodNames from './private-types'
import GearspecMethodNames from '../../gearspecs/_store/public-types'
const { MutationTypes: Mutations } = ProjectMethodNames
const { GetterTypes: GearspecGetters } = GearspecMethodNames

const OverrideMutations = {

  [Mutations.SET_STACK] (state, payload) {
    const { project, stack } = payload
    if (!stack) {
      console.warn('Most likely `stack` arg should not be empty!')
    }
    if (project) {
      Vue.set(project.attributes, 'stack', stack)
    }
  },

  [Mutations.UPDATE_HOSTNAME] (state, payload) {
    const { project, hostname } = payload
    if (project && hostname) {
      Vue.set(project.attributes, 'hostname', hostname)
    }
  },

  [Mutations.UPDATE_STATE] (state, payload) {
    const { project, isEnabled } = payload
    if (project) {
      project.attributes.enabled = !!isEnabled
    }
  },

  [Mutations.UPDATE_NOTES] (state, payload) {
    const { project, notes } = payload
    if (project) {
      project.attributes.notes = notes
    }
  },

  [Mutations.ADD_STACK] (state, payload) {
    const { stackId, actualStackId, project, stack } = payload

    if (project && stack && stack.attributes.members.length) {
      if (typeof project.attributes.stack === 'undefined') {
        Vue.set(project.attributes, 'stack', [])
      }
      stack.attributes.members.forEach(m => {
        if (!m.gearspec_id) {
          return true
        }

        const item = project.attributes.stack.find(it => it.gearspec_id === m.gearspec_id)

        if (item && stackId !== actualStackId) {
          // if m.gearspec_id already exists, mark it with isRemoved = false
          Vue.set(item, 'isRemoved', false)
        } else {
          const serviceId = store.getters[GearspecGetters.FIND_COMPATIBLE_SERVICE](
            m.gearspec_id,
            m.default_service
          )
          if (item) {
            Vue.set(item, 'isRemoved', false)
            Vue.set(item, 'service_id', serviceId)
          } else {
            project.attributes.stack.push({
              service_id: serviceId, // a falsy serviceId is OK
              gearspec_id: m.gearspec_id,
              isRemoved: false
            })
          }
        }
      })
    }
  },

  [Mutations.REMOVE_STACK] (state, payload) {
    const { project, stackId } = payload
    const shortStackName = stackId.split('/')[1]

    // if (typeof state.removedStacks[projectId] === 'undefined') {
    //   Vue.set(state.removedStacks, projectId, [])
    // }

    /**
     * We need to remove all elements of project.stack that that have service_id starting with shortStackName, e.g. "wordpress/"
     *
     * For deleting array items in javascript with forEach() and splice())
     * @see https://gist.github.com/chad3814/2924672
     */
    for (let i = project.attributes.stack.length - 1; i >= 0; i--) {
      if (project.attributes.stack[i].gearspec_id.split('/')[1] === shortStackName) {
        Vue.set(project.attributes.stack[i], 'isRemoved', true)
        // state.removedStacks[projectId].push(project.attributes.stack[i])
        // console.log(projectId, stackId, project.attributes.stack[i])
        // Vue.delete(project.attributes.stack, i)
      }
    }
  },

  [Mutations.CHANGE_GEAR] (state, payload) {
    const { project, serviceId, memberIndex } = payload

    if (project && memberIndex >= 0) {
      /**
       * note, serviceId might be an empty string (and that's OK)
       */
      Vue.set(project.attributes.stack[memberIndex], 'service_id', serviceId)
    }
  }
}

export default { ...BaseMutations, ...OverrideMutations }
