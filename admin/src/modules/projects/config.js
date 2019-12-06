import baseConfig from '../_base/config'

const moduleConfig = {
  moduleName: 'ProjectsModule',
  namespace: 'projects',
  apiEndpoint: 'projects',
  recordType: 'project',
  labels: {
    singular: 'Project',
    plural: 'Projects',
    empty: 'Found no projects.'
  }
}

moduleConfig.labels = { ...baseConfig.labels, ...moduleConfig.labels }

export default { ...baseConfig, ...moduleConfig }
