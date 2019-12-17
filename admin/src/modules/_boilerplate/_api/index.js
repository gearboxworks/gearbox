import BaseMethodsFactory from '../../_base/_api'
import config from '../config'

const ExtraMethods = {
  demoMethod (arg1, arg2) {
    return new Promise((resolve, reject) => {
      console.log('Called extra method', arg1, arg2)
      return 'Result from demoMethod'
    })
  }
}

export default { ...BaseMethodsFactory(config.recordType, config), ...ExtraMethods }
