import api from './methods'

function RepositoryFactory (moduleConfig) {
  const endpoint = moduleConfig.apiEndpoint

  return {
    recordType: moduleConfig.recordType,
    apiEndpoint: endpoint,

    fetchAll: () => {
      return api.fetchAll(endpoint)
    },
    fetchOne: (id) => {
      return api.fetchOne(endpoint, id)
    },
    create: (recordData) => {
      return api.create(endpoint, recordData)
    },
    update: (record, recordData) => {
      return api.update(endpoint, record, recordData)
    },
    delete: (id) => {
      return api.delete(endpoint, id)
    }
  }
}

export default RepositoryFactory
