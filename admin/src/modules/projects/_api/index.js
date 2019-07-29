import HTTP from '../../../http-common'

const api = {

  fetchAllHeaders: () => HTTP.get(
    'projects',
    {
      crossDomain: true
      // raxConfig: {
      //   // You can detect when a retry is happening, and figure out how many
      //   // retry attempts have been made
      //   onRetryAttempt: (err) => {
      //     const cfg = raxConfig(err)
      //     commit('SET_NETWORK_ERROR', err.message)
      //     commit('SET_REMAINING_RETRIES', cfg.retry - cfg.currentRetryAttempt)
      //   }
      // }
    }
  ).then(response => {
    if (response.data && response.data.data) {
      return (response.data.data)
    } else {
      throw new Error('Cannot fetch all project headers')
    }
  }),

  fetchDetails: (projectId) => HTTP.get(
    `projects/${projectId}`,
    {
      crossDomain: true
    }
  ).then(response => {
    if (response.data && response.data.data) {
      return (response.data.data)
    } else {
      throw new Error('Could not fetch project details for ' + projectId)
    }
  }),

  updateDetails: (projectId, data) => HTTP({
    method: 'post',
    url: 'project/' + projectId,
    data: data
  }).then(results => {
    if (results.data && results.data.data) {
      return results.data.data
    } else {
      throw new Error('Could not update details for project ' + projectId)
    }
  })
}

export default api
