import Vue from 'vue'
import { Mutations } from './private-types'

const mutations = {

  [Mutations.SET_RECORDS] (state, projects) {
    if (!projects || projects.length === 0) {
      console.warn('Most likely `projects` arg should not be empty!')
    }
    Vue.set(state, 'records', projects)
  },

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
    const { stackId, actualStackId, project, stack, preselectServiceId } = payload

    if (project && stack && stack.attributes.members.length) {
      if (typeof project.attributes.stack === 'undefined') {
        Vue.set(project.attributes, 'stack', [])
      }
      stack.attributes.members.forEach((el, idx) => {
        if (el.gearspec_id) {
          const item = project.attributes.stack.find(
            it => it.gearspec_id === el.gearspec_id)
          if (item && stackId !== actualStackId) {
            // if el.gearspec_id already exists, mark it with isRemoved = false
            Vue.set(item, 'isRemoved', false)
          } else {
            // reactive!
            const serviceId = preselectServiceId(el.services, el.default_service)
            if (item) {
              Vue.set(item, 'isRemoved', false)
              Vue.set(item, 'service_id', serviceId)
            } else {
              project.attributes.stack.push({
                service_id: serviceId, // it's ok if serviceId is empty
                gearspec_id: el.gearspec_id,
                isRemoved: false
              })
            }
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
      if (project.attributes.stack[i].gearspec_id.split('/')[1] ===
        shortStackName) {
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
  },

  [Mutations.SET_LIST_FILTER] (state, payload) {
    const { field, values } = payload
    Vue.set(state.showProjectsHaving, field, values)
  },

  [Mutations.SET_LIST_FILTER_SORT_BY] (state, sortBy) {
    state.sortBy = sortBy
  },

  [Mutations.SET_LIST_FILTER_SORT_ORDER] (state, isAscending) {
    state.sortOrder = isAscending
  }
}

export default mutations
