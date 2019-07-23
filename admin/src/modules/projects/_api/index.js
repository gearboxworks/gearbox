import HTTP from '../../../http-common'

const api = {

  fetchAllHeaders: () => new Promise((resolve, reject) => {
    const call = HTTP.get(
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
    )

    call.then(response => {
      if (response.data && response.data.data) {
        resolve(response.data.data)
      } else {
        reject(new Error('Cannot fetch all project headers'))
      }
    })

    call.catch((error, config) => {
      // alert('Please make sure Gearbox API is running at \nhttp://127.0.0.1:9999/');
      reject(error)
      // if (error.message === 'Network Error') {
      //   commit('SET_NETWORK_ERROR', error.message)
      // }
    })

    return call
  }),

  fetchDetails: (projectId) => new Promise((resolve, reject) => {
    const call = HTTP.get(
      `projects/${projectId}`,
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
    )

    call.then(response => {
      if (response.data && response.data.data) {
        resolve(response.data.data)
      } else {
        reject(new Error('Could not fetch project details for ' + projectId))
      }
    })

    call.catch((error, config) => {
      // if (error.message === 'Network Error') {
      //   commit('SET_NETWORK_ERROR', error.message)
      // }
      reject(error)
    })

    return call
  }),

  updateDetails: (projectId, data) => new Promise((resolve, reject) => {
    const call = HTTP({
      method: 'post',
      url: 'project/' + projectId,
      data: data
    })

    call.then(results => {
      if (results.data && results.data.data) {
        resolve(results.data.data)
      } else {
        reject(new Error('Could not update details for project ' + projectId))
      }
    })

    call.catch((error) => {
      reject(error)
    })

    return call
  })
}

export default api
