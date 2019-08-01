import BaseMethodsFactory from '../../_base/_api'
import moduleConfig from '../config'

const ExtraMethods = {
  demoMethod (arg1, arg2) {
    return new Promise((resolve, reject) => {
      console.log('Called extra method', arg1, arg2)
      return 'Result from demoMethod'
    })
  }
}

export default { ...BaseMethodsFactory(moduleConfig), ...ExtraMethods }
