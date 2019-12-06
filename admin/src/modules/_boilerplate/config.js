import baseConfig from '../_base/config'

const moduleConfig = {
  moduleName: 'ItemsModule',
  namespace: 'items',
  apiEndpoint: 'items',
  recordType: 'item',
  labels: {
    singular: 'Item',
    plural: 'Items',
    empty: 'Found no items.'
  }
}

moduleConfig.labels = { ...baseConfig.labels, ...moduleConfig.labels }

export default { ...baseConfig, ...moduleConfig }
