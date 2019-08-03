import NamespacedBaseGetters from '../../_base/_store/getters'
import moduleConfig from '../config'
import { ServiceGetters as Getters } from './method-names'

export default {
  ...NamespacedBaseGetters(moduleConfig.namespace),

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
