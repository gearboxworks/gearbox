import PrivateTypes from './private-types'

export const prefix = 'Project'
export const namespace = 'projects'

const PublicTypes = {}
Object.entries(PrivateTypes).forEach(([group, types]) => {
  PublicTypes[group] = {}
  Object.entries(types).forEach(([key, val]) => {
    PublicTypes[group][key] = namespace + '/' + val
  })
})

export const {
  Getters: ProjectGetters,
  Actions: ProjectActions,
  Mutations: ProjectMutations
} = PublicTypes

export default PublicTypes
