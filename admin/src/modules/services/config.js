import baseConfig from '../_base/config'

const moduleConfig = {
  moduleName: 'ServicesModule',
  namespace: 'services',
  apiEndpoint: 'services',
  recordType: 'service',
  labels: {
    singular: 'Service',
    plural: 'Services',
    empty: 'Found no services.'
  }
}

moduleConfig.labels = { ...baseConfig.labels, ...moduleConfig.labels }

export default { ...baseConfig, ...moduleConfig }
