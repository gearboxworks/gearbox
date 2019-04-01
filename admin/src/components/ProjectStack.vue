<template>
  <div role="tablist" class="project-stack-list" :id="project_base + '-stack'">
    <div
      v-for="(projectServices, stackName) in groupProjectStacks(projectStack)"
      :key="stackName"
      class="project-stack"
    >
      <strong>{{stackName.replace('gearbox.works/', '')}}</strong>
      <b-button @click.prevent="removeProjectStack(stackName)" class="js-remove-stack" size="sm" variant="outline-secondary" aria-label="Remove this stack from project" title="Remove this stack from project">&times;</b-button>
      <!--project-stack-header :projectHostname = "projectHostname" :stackName = "stackName" :stackRoles = "projectServices"></project-stack-header-->
      <ul class="service-list">
        <li
            v-for="(service, serviceRole) in stackServices(stackName)"
            :key="project_base + escAttr(serviceRole)"
            :id="project_base + escAttr(serviceRole)"
            class="service-item"
        >
          <img :src="require('../assets/'+projectServices[serviceRole].program+'.svg')" class="service-program" />
          <h6 class="service-role">{{stackRoles(stackName)[serviceRole].label}}</h6>
          <b-popover
            :target="project_base + escAttr(serviceRole)"
            :container="projectHostname+'-stack'"
            :ref="project_base + escAttr(serviceRole)+'_popover'"
            triggers="click blur"
            placement="bottom"
          >
            <template slot="title">
              <b-button @click="onClosePopoverFor(project_base + escAttr(serviceRole))" class="close" aria-label="Close">
                <span class="d-inline-block" aria-hidden="true">&times;</span>
              </b-button>
              {{stackRoles(stackName)[serviceRole].name}}
            </template>

            <div>
              <!--{{'Service: '+projectServices[serviceRole].service_id}}
              <div><small>Click to pick a different one...</small></div>-->
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
            </div>
          </b-popover>
        </li>
      </ul>
    </div>

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

// import ProjectStackHeader from './ProjectStackHeader.vue'
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
    // ProjectStackHeader
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
    },
    onClosePopoverFor (triggerElementId) {
      console.log('onClosePopoverFor', triggerElementId)
      this.$root.$emit('bv::popover::hide', triggerElementId)
    },
    removeProjectStack (stackName) {
      this.$store.dispatch('removeProjectStack', { 'projectHostname': this.projectHostname, stackName })
    }
  }
}
</script>

<style scoped>
  .js-remove-stack {
    float: right;
  }
  .project-stack{
    margin-bottom: 8px;
  }
  .project-stack:not(:first-child){
    border-top: 1px solid silver;
    padding-top: 8px;
  }
  .add-stack {
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
  .service-list{
    margin-bottom: 0;
  }
  .service-item{
    display: inline-block;
    list-style: none;
    text-align: center;
    max-width: 110px;
    padding: 5px;
    margin: 5px;
  }
  .service-role{
    margin-top:5px;
    margin-bottom: 0;
    clear: both;
  }
  .service-program {
    height: 64px;
    width: 64px;
  }
</style>
