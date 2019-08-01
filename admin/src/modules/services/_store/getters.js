import BaseGetters from '../../_base/_store/getters'
import { UNSUPPORTED_GETTER } from '../../_helpers'

import StoreMethodTypes from './private-types'
const { GetterTypes: Getters } = StoreMethodTypes

const OverrideGetters = {

  [Getters.FIND_BY]: (state) => (fieldName, fieldValue) => {
    /**
     * manipulate arguments
     */
    const results = BaseGetters[Getters.FIND_BY](state)(fieldName, fieldValue)
    /**
     * manipulate results
     */
    return results
  },

  [Getters.LIST_PROGRAM_OPTIONS]: (state) => () => {
    const programs = []
    const options = []
    state.records.forEach((el, idx) => {
      const program = el.attributes.program
      if (programs.indexOf(program) === -1) {
        programs.push(program)
        options.push({
          value: program,
          text: program
        })
      }
    })
    return options.sort((a, b) => a.value > b.value ? 1 : (a.value === b.value) ? 0 : -1)
  }
}

export default { ...BaseGetters, ...OverrideGetters }
