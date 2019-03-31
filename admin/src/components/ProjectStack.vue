groupProjectStacks<template>
  <div role="tablist" class="project-stack">
    <b-card
      v-for="(projectServices, stackName, index) in groupProjectStacks(projectStack)"
      :key="stackName"
      no-body
      class="mb-1"
    >
      <b-card-header header-tag="header" class="p-1" role="tab">
        <b-button block href="#" v-b-toggle="project_base + '_accordion_' + index" variant="info">
          <project-stack-header :projectHostname = "projectHostname" :stackName = "stackName" :stackRoles = "projectServices"></project-stack-header>
        </b-button>
      </b-card-header>
      <b-collapse :id="project_base + '_accordion_' + index" :accordion="project_base + '_accordion'" role="tabpanel">
        <b-card-body>
          <b-form>
            <b-form-group
              v-for="(service, serviceRole) in stackServices(stackName)"
              :key="project_base + escAttr(serviceRole)"
              :label="stackRoles(stackName)[serviceRole].label"
              :label-for="project_base + escAttr(serviceRole)+'_input'"
              :description="stackRoles(stackName)[serviceRole].name"
              label-cols-sm="4"
              label-cols-lg="3"
            >
              <b-form-select
                :id="project_base + escAttr(serviceRole)+'_input'"
                :value="projectServices[serviceRole] ? projectServices[serviceRole].service_id: null"
                @change="changeProjectService(serviceRole,$event)"
              >
                <option disabled value="">Please select one...</option>
                <optgroup v-for="(options, groupLabel) in optionGroups(service.options)" :label="groupLabel" :key="groupLabel">
                  <option v-for="serviceVer in options" :value="service.org + '/' + serviceVer" :key="serviceVer">{{serviceVer}}</option>
                </optgroup>
              </b-form-select>
            </b-form-group>
          </b-form>
        </b-card-body>
      </b-collapse>
    </b-card>
    <b-form-select v-model='stackToAdd' @change="addProjectStack" v-if="hasUnusedStacks" class="add-stack">
      <option :value="null" disabled>Add Stack...</option>
      <option
        v-for="(stack,stackName) in stacksNotUnusedInProject"
        :key="stackName"
        :value="stackName"
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
    },
    hasUnusedStacks () {
      return Object.entries(this.stacksNotUnusedInProject).length > 0
    },
    stacksNotUnusedInProject () {
      const result = {}
      const projectStacks = this.groupProjectStacks(this.projectStack)
      for (const index in this.$store.state.gearStacks) {
        const stackName = this.$store.state.gearStacks[index]
        if (typeof projectStacks[stackName] === 'undefined') {
          result[stackName] = this.$store.state.gearStacks[stackName]
        }
      }
      return result
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
        console.log('found', serviceId)
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
      this.$store.dispatch('addProjectStack', { 'projectHostname': this.projectHostname, stackName })
      this.stackToAdd = null
    },
    changeProjectService (serviceName, serviceId) {
      this.$store.dispatch('changeProjectService', { 'projectHostname': this.projectHostname, serviceName, serviceId })
    }
  }
}
</script>

<style scoped>
  .add-stack {
    margin-left: 15px;
    width: calc(100% - 30px);
/*    color: #fff;
    background-color: #17a2b8;
    border-color: #17a2b8;
*/
  }
/*
  .add-stack:hover {
    color: #fff;
    background-color: #138496;
    border-color: #117a8b;
  }
*/
</style>
