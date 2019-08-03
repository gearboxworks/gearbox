export function createPublicStoreMethodTypes (PrivateTypes) {
  const { Namespace } = PrivateTypes
  const PublicTypes = {}

  Object.entries(PrivateTypes).forEach(([group, types]) => {
    if (group === 'Namespace') {
      return true
    }

    PublicTypes[group] = {}
    Object.entries(types).forEach(([key, val]) => {
      PublicTypes[group][key] = Namespace + '/' + val
    })
  })

  PublicTypes.Namespace = Namespace

  return PublicTypes
}

export function UNSUPPORTED_GETTER () {
  console.trace('UNSUPPORTED GETTER!')
  throw new Error('UNSUPPORTED GETTER!')
}

export function UNSUPPORTED_ACTION () {
  console.trace('UNSUPPORTED ACTION!')
  throw new Error('UNSUPPORTED ACTION!')
}

export function UNSUPPORTED_MUTATION () {
  console.trace('UNSUPPORTED MUTATION!')
  throw new Error('UNSUPPORTED MUTATION!')
}

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
