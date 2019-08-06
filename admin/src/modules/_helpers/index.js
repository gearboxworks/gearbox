/**
 * TODO: verify there are no artificial delays in production code
 */
export const FORCED_DELAY = 2000

export function namespaceKeys (methods, namespace) {
  // return methods
  const result = {}
  for (let [key, value] of Object.entries(methods)) {
    result[namespace + '/' + key] = value
  }
  return result
}

// export function denamespaceKeys (methods) {
//   // return methods
//   const result = {}
//   for (let [key, value] of Object.entries(methods)) {
//     result[key.split('/')[1]] = value
//   }
//   return result
// }

export function namespaceValues (constants, namespace) {
  // return constants
  const result = {}
  for (let [key, value] of Object.entries(constants)) {
    result[key] = namespace + '/' + value
  }
  return result
}

export function escapeIDAttr (value) {
  return value.replace(/\//g, '-').replace(/\./g, '-')
}

export function programFromServiceId (serviceId) {
  return serviceId ? serviceId.split('/')[1].split(':')[0] : ''
}

export function programAndVersionFromServiceId (serviceId) {
  return serviceId ? serviceId.split('/')[1] : ''
}

export function versionFromServiceId (serviceId) {
  return serviceId ? serviceId.split(':')[1] : ''
}

export function stackNameFromStackId (stackId) {
  return stackId.split('/')[1]
}

export function stackNameFromGearspecId (gearspecId) {
  return gearspecId.split('/')[1]
}

export function setCookie (name, value) {
  const d = new Date()
  d.setTime(d.getTime() + (3 * 24 * 60 * 60 * 1000)) // 3 days
  const expires = 'expires=' + d.toUTCString()

  let cleaned = value
  if (typeof value === 'object' || Array.isArray(value)) {
    cleaned = JSON.stringify({ ...value })
  }

  document.cookie = name + '=' + cleaned + ';' + expires + ';path=/'
}

export function getCookie (name) {
  let value = '; ' + document.cookie
  const parts = value.split('; ' + name + '=')
  if (parts.length === 2) {
    value = parts.pop().split(';').shift()
  }
  if (value) {
    try {
      return JSON.parse(value)
    } catch (e) {
      return value
    }
  }
  return null
}
