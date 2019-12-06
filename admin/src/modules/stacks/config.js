import baseConfig from '../_base/config'

const moduleConfig = {
  moduleName: 'StacksModule',
  namespace: 'stacks',
  apiEndpoint: 'stacks',
  recordType: 'stack',
  labels: {
    singular: 'Stack',
    plural: 'Stacks',
    empty: 'Found no stacks.'
  }
}

moduleConfig.labels = { ...baseConfig.labels, ...moduleConfig.labels }

export default { ...baseConfig, ...moduleConfig }
