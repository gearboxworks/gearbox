import { Getters } from './private-types'

const getters = {
  [Getters.PROJECT_BY]: (state) => (fieldName, fieldValue) => {
    return (fieldName === 'id')
      ? state.records.find(p => p.id === fieldValue)
      : state.records.find(p => p.attributes[fieldName] === fieldValue)
  },

  [Getters.FILTER_PROJECTS_BY]: (state) => (fieldName, allowedValues) => {
    const attrs = ['basedir', 'enabled', 'filepath', 'hostname', 'path', 'project_dir']
    let valuesArray = Array.isArray(allowedValues) ? allowedValues : [allowedValues]
    // 'notes' and 'stack' are not included on purpose because simple comparison does not work on them
    let projects = []

    if (fieldName === 'id') {
      projects = state.records.filter(p => valuesArray.indexOf(p.id) !== -1)
    } else if (attrs.indexOf(fieldName) !== -1) {
      projects = state.records.filter(p => valuesArray.indexOf(p.attributes[fieldName]) !== -1)
    } else if (fieldName === 'stacks') {
      projects = state.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.gearspec_id.indexOf(val) > -1)))
    } else if (fieldName === 'programs') {
      projects = state.records.filter(p => p.attributes.stack.some(s => valuesArray.some(val => s.service_id.split('/')[1].split(':')[0] === val)))
    }

    return projects
  },

  [Getters.FILTERED_PROJECTS]: (state, getters) => {
    let projects = state.records
    const sortAscending = !!state.sortOrder
    // const sortBy = state.sortBy
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
            projects = projects.filter(p => getters[Getters.FILTER_PROJECTS_BY]('enabled', true).includes(p))
          }
          if (values.indexOf('stopped') > -1) {
            projects = projects.filter(p => getters[Getters.FILTER_PROJECTS_BY]('enabled', false).includes(p))
          }
          // TODO merge candidates into projects array
          // if (values.indexOf('candidates') > -1) {
          //   projects = projects.filter(p => getters.filterProjectsBy('candidate', true).includes(p))
          // }
        }
        continue
      }
      projects = (projects.filter(p => getters[Getters.FILTER_PROJECTS_BY](field, values).includes(p)))
    }
    return projects.concat().sort((a, b) => a.id > b.id ? (sortAscending ? 1 : -1) : (a.id === b.id) ? 0 : (sortAscending ? -1 : 1))
  },

  [Getters.PROJECT_STACK_ITEM_INDEX_BY]: (state) => (project, fieldName, fieldValue) => {
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
  }
}

export default getters
