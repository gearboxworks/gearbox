<template>
  <div role="tablist" class="project-stack">
    <b-card
      v-for="(projectServices, stackName, index) in groupProjectStacks(projectStack)"
      :key="stackName"
      no-body
      class="mb-1"
    >
      <b-card-header header-tag="header" class="p-1" role="tab">
        <b-button block href="#" v-b-toggle="project_base + '_accordion_' + index" variant="info">
          <project-stack-header :stackName = "stackName" :stackRoles = "projectServices"></project-stack-header>
        </b-button>
      </b-card-header>
      <b-collapse :id="project_base + '_accordion_' + index" :accordion="project_base + '_accordion'" role="tabpanel">
        <b-card-body>
          <b-form>
            <b-form-group
              v-for="(service, serviceRole) in stackServices(stackName)"
              :key="project_base + escAttr(serviceRole)"
              :label="stackRoles(stackName)[serviceRole].short_label"
              :label-for="project_base + escAttr(serviceRole)+'_input'"
              :description="stackRoles(stackName)[serviceRole].label"
              label-cols-sm="4"
              label-cols-lg="3"
            >
              <b-form-select
                :id="project_base + escAttr(serviceRole)+'_input'"
              >
                <optgroup v-for="(options, groupLabel) in optionGroups(service.options)" :label="groupLabel" :key="groupLabel">
                  <option v-for="(serviceVer) in options" :value="serviceVer" :key="serviceVer" :selected="stackIncludesService(projectServices, service.org+'/'+serviceVer)">{{serviceVer}}</option>
                </optgroup>
              </b-form-select>
            </b-form-group>
          </b-form>
        </b-card-body>
      </b-collapse>
    </b-card>
    <br>
    <b-form-select v-model='stackToAdd' @change="addProjectStack">
      <option :value="null">Add Stack...</option>
      <option
        v-for="(stack,stackName) in unusedProjectStacks(projectStack)"
        :key="stackName"
      >{{stackName.replace('gearbox.works/', '')}}</option>
    </b-form-select>

  </div>
</template>

<script>

import ProjectStackHeader from './ProjectStackHeader.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'ProjectStack',
  props: {
    'projectHostname': {
      type: String,
      required: true
    },
    'projectStack': {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      stackToAdd: null
    }
  },
  components: {
    ProjectStackHeader
  },
  computed: {
    ...mapGetters(['groupProjectStacks', 'stackRoles', 'stackServices']),
    project_base () {
      return this.escAttr(this.projectHostname)
    },
    project () {
      return this.$store.getters.projectBy('hostname', this.projectHostname)
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '_').replace(/\./g, '_')
    },
    stackIncludesService (services, serviceId) {
      let result = false
      for (const serviceName in services) {
        if (serviceId === services[serviceName].service_id) {
          result = true
          break
        }
      }
      if (result) {
        // console.log('found', serviceId)
      }
      return result
    },
    unusedProjectStacks (projectStack) {
      const result = {}
      const projectStacks = this.groupProjectStacks(projectStack)
      for (const index in this.$store.state.gearStacks) {
        const stackName = this.$store.state.gearStacks[index]
        // console.log(stackName, projectStack)
        if (typeof projectStacks[stackName] === 'undefined') {
          result[stackName] = this.$store.state.gearStacks[stackName]
        }
      }
      return result
    },
    mapOptions (options) {
      const result = []
      for (const value in options) {
        result.push({
          value,
          text: options[value]
        })
      }
      return result
    },
    optionGroups (options) {
      const result = {}
      for (const index in options) {
        const base = options[index].split(':')[0]
        if (typeof result[base] === 'undefined') {
          result[base] = {}
        }
        result[base][index] = options[index]
      }
      return result
    },
    addProjectStack (stackName) {
      console.log('Selected', this.stackToAdd, stackName)
      this.$store.dispatch('addProjectStack', stackName)
      this.stackToAdd = null
    }
  }
}
</script>

<style scoped>

</style>
