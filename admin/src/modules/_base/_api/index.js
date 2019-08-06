import api from './methods'

function RepositoryFactory (moduleConfig) {
  const endpoint = moduleConfig.apiEndpoint

  return {
    recordType: moduleConfig.recordType,
    apiEndpoint: endpoint,

    fetchAll: () => api.fetchAll(endpoint),
    fetchOne: (id) => api.fetchOne(endpoint, id),
    create: (recordData) => api.create(endpoint, recordData),
    update: (record, recordData) => api.update(endpoint, record, recordData),
    delete: (id) => api.delete(endpoint, id)
  }
}

export default RepositoryFactory
