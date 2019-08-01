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
