import axios from 'axios'
import { attach as raxAttach, getConfig as raxConfig } from 'retry-axios'

const HTTP = axios.create({
  baseURL: `http://127.0.0.1:9999/`
})

HTTP.defaults.raxConfig = {

  // If you are using a non static instance of Axios you need
  // to pass that instance here (const ax = axios.create())
  instance: HTTP,

  // Retry 3 times on requests that return a response (500, etc) before giving up.  Defaults to 3.
  retry: 5,

  // Retry twice on errors that don't return a response (ENOTFOUND, ETIMEDOUT, etc).
  noResponseRetries: 5,

  // Milliseconds to delay at first.  Defaults to 100.
  retryDelay: 3000,

  // HTTP methods to automatically retry.  Defaults to:
  // ['GET', 'HEAD', 'OPTIONS', 'DELETE', 'PUT']
  httpMethodsToRetry: ['GET', 'HEAD', 'OPTIONS', 'DELETE', 'PUT'],

  // The response status codes to retry.  Supports a double
  // array with a list of ranges.  Defaults to:
  // [[100, 199], [429, 429], [500, 599]]
  statusCodesToRetry: [[100, 199], [429, 429], [500, 599]],
  // Override the decision making process on if you should retry
  shouldRetry: (err) => {
    const config = raxConfig(err)

    if (!config || config.retry === 0) {
      return false
    }

    // Check if this error has no response (ETIMEDOUT, ENOTFOUND, etc)
    if (!err.response && ((config.currentRetryAttempt || 0) >= config.noResponseRetries)) {
      return false
    }
    // Only retry with configured HttpMethods.
    if (!err.config.method || Object.values(config.httpMethodsToRetry).indexOf(err.config.method.toUpperCase()) < 0) {
      return false
    }

    // If this wasn't in the list of status codes where we want
    // to automatically retry, return.
    if (err.response && err.response.status) {
      let isInRange = false
      for (const [min, max] of config.statusCodesToRetry) {
        const status = err.response.status
        if (status >= min && status <= max) {
          isInRange = true
          break
        }
      }
      if (!isInRange) {
        return false
      }
    }

    // If we are out of retry attempts, return
    config.currentRetryAttempt = config.currentRetryAttempt || 0
    if (config.currentRetryAttempt >= config.retry) {
      return false
    }

    return true
  }
}

raxAttach(HTTP)

export default HTTP
