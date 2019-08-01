import BaseMutations from '../../_base/_store/mutations'
// import { UNSUPPORTED_MUTATION } from '../../_helpers'
// import StoreMethodNames from './private-types'
// const { MutationTypes: Mutations } = StoreMethodNames

const OverrideMutations = {
  // [Mutations.DEMO_MUTATION] (state, payload) {
  //   console.error('Called DEMO_MUTATION mutation:', payload)
  // },
  //
  // [Mutations.SET_ALL]: (state, payload) => UNSUPPORTED_MUTATION(),
  // [Mutations.SET_ONE]: (state, payload) => UNSUPPORTED_MUTATION(),
  // [Mutations.UPDATE]: (state, payload) => UNSUPPORTED_MUTATION(),
  // [Mutations.DELETE]: (state, payload) => UNSUPPORTED_MUTATION()
}

export default { ...BaseMutations, ...OverrideMutations }
