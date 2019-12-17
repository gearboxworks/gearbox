import baseConfig from '../_base/config'

const moduleConfig = {
  moduleName: 'BasedirsModule',
  namespace: 'basedirs',
  apiEndpoint: 'basedirs',
  recordType: 'basedir',
  labels: {
    singular: 'Basedir',
    plural: 'Basedirs',
    empty: 'Found no basedirs.'
  }
}

moduleConfig.labels = { ...baseConfig.labels, ...moduleConfig.labels }

export default { ...baseConfig, ...moduleConfig }
