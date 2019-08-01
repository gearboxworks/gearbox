import baseConfig from '../_base/config'

const moduleConfig = {
  moduleName: 'GearspecModule',
  namespace: 'gearspecs',
  apiEndpoint: 'gearspecs',
  recordType: 'gearspec',
  labels: {
    singular: 'Gearspec',
    plural: 'Gearspecs',
    empty: 'Found no gearspecs.'
  }
}

moduleConfig.labels = { ...baseConfig.labels, ...moduleConfig.labels }

export default { ...baseConfig, ...moduleConfig }
