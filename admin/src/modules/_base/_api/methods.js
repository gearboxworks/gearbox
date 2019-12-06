import HTTP from '../../_helpers/http-common'

export default {

  fetchAll: (endpoint) => HTTP.get(endpoint)
    .then(
      response => {
        if (response.data && response.data.data) {
          return (response.data.data)
        } else {
          throw new Error('Cannot fetch all records')
        }
      }
    ),

  fetchOne: (endpoint, id) => HTTP.get(endpoint + '/' + id)
    .then(response => {
      if (response.data && response.data.data) {
        return (response.data.data)
      } else {
        throw new Error('Could not fetch details for ' + id)
      }
    }),

  create: (endpoint, recordData) => {
    // console.log('create', endpoint, recordData)
    return HTTP({
      method: 'post',
      url: endpoint + '/new',
      data: recordData
    }).then(results => {
      if (results.data && results.data.data) {
        return results.data.data
      } else {
        throw new Error(`Could not create ${endpoint} record`)
      }
    })
  },

  update: (endpoint, record, recordData) => {
    // console.log('update', endpoint, record, recordData)
    return HTTP({
      method: 'patch',
      url: endpoint + '/' + record.id,
      data: recordData
    }).then(results => {
      if (results.data && results.data.data) {
        return results.data.data
      } else {
        throw new Error('Could not update details for record ' + record.id)
      }
    })
  },

  delete: (endpoint, id) => HTTP({
    method: 'delete',
    url: endpoint + '/' + id
  }).then(results => {
    if (results.data && results.data.data) {
      return results.data.data
    } else {
      throw new Error('Could not delete record ' + id)
    }
  })
}
