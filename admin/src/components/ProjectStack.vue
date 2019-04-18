<template>
  <div role="tablist" class="project-stack-list" :id="`${projectBase}stack`">
    <div
      v-for="(projectServices, stackName, stackIndex) in groupProjectStacks(projectStack)"
      :key="stackName"
      class="project-stack"
    >
      <h2 class="stack-title">{{stackName.replace('gearbox.works/', '')}}</h2>
      <b-button :tabindex="projectIndex*100+stackIndex*10" @click.prevent="removeProjectStack(stackName)" class="js-remove-stack" size="sm" variant="outline-secondary" aria-label="Remove this stack from project" title="Remove this stack from project">&times;</b-button>
      <!--project-stack-header :projectHostname = "projectHostname" :stackName = "stackName" :stackRoles = "projectServices"></project-stack-header-->
      <ul class="service-list">
        <li
            v-for="(service, serviceRole, serviceIndex) in stackServices(stackName)"
            :key="projectBase + serviceRole"
            :id="projectBase + serviceRole"
            class="service-item"
            :tabindex="projectIndex*100+stackIndex*10+serviceIndex+1"

        >
          <img v-if="projectServices[serviceRole] && projectServices[serviceRole].program" :src="require('../assets/'+projectServices[serviceRole].program+'.svg')" class="service-program" />
          <font-awesome-icon
            v-else
            :icon="['fa', 'expand']"
          />

          <h6 class="service-role">{{stackRoles(stackName)[serviceRole].label}}</h6>
          <b-tooltip
            triggers="hover"
            :target="projectBase + serviceRole"
            :title="programTooltip(projectServices[serviceRole])"
            :key="projectServices[serviceRole] ? projectServices[serviceRole].service_id: null"
          />
          <b-popover
            :target="projectBase + serviceRole"
            :container="`${projectHostname}stack`"
            :ref="projectBase + serviceRole + '-popover'"
            triggers="focus"
            placement="bottom"
          >
            <template slot="title">
              <b-button @click="onClosePopoverFor(projectBase + serviceRole)" class="close" aria-label="Close">
                <span class="d-inline-block" aria-hidden="true">&times;</span>
              </b-button>
              Change service
            </template>

            <div>
              <label :for="projectBase + serviceRole+'_input'">
              {{stackRoles(stackName)[serviceRole] ? stackRoles(stackName)[serviceRole].program:''}}:
              </label>
              <b-form-select
                :id="projectBase + serviceRole+'_input'"
                :value="projectServices[serviceRole] ? projectServices[serviceRole].service_id: null"
                :tabindex="projectIndex*100+stackIndex*10+serviceIndex+9"
                @change="changeProjectService(serviceRole,$event)"
              >
                <option value="" v-if="!$store.state.gearServices[serviceRole].default">Do not run this service</option>
                <option disabled value="">Select service...</option>
                <optgroup v-for="(options, groupLabel) in optionGroups(service.options)" :label="groupLabel" :key="groupLabel">
                  <option v-for="serviceVer in options" :value="service.org + '/' + serviceVer" :key="serviceVer">{{serviceVer}}</option>
                </optgroup>
              </b-form-select>
            </div>
          </b-popover>
        </li>
      </ul>
    </div>
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
    },
    'projectIndex': {
      type: Number,
      required: true
    }
  },
  data () {
    return {

    }
  },
  components: {
    // ProjectStackHeader
  },
  computed: {
    ...mapGetters(['groupProjectStacks', 'stackRoles', 'stackServices']),
    projectBase () {
      return this.escAttr(this.projectHostname) + '-'
    },
    project () {
      return this.$store.getters.projectBy('hostname', this.projectHostname)
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
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
    programTooltip (service) {
      return service ? (service.program ? service.service_id.replace('gearboxworks/', '') : 'Not selected') : 'Not selected'
    },
    changeProjectService (serviceName, serviceId) {
      this.$store.dispatch('changeProjectService', { 'projectHostname': this.projectHostname, serviceName, serviceId })
      this.onClosePopoverFor(this.projectBase + this.escAttr(serviceName))
    },
    onClosePopoverFor (triggerElementId) {
      this.$root.$emit('bv::hide::popover', triggerElementId)
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
    opacity: 0;
    transition: opacity 400ms;
  }
  .card--project:hover .js-remove-stack {
    opacity: 1;
  }
  .project-stack{
    margin-bottom: 8px;
  }
  .project-stack:not(:first-child){
    border-top: 1px solid #e0e0e0;
    padding-top: 8px;
  }
  .stack-title {
    text-transform: uppercase;
    font-size: 1.25rem;
    display: inline-block;
    color: #989898;
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
    cursor: pointer;
    border: 1px solid transparent;
    border-radius: 4px;
    transition: all 400ms;
  }
  .service-item:hover {
    border: 1px solid #aaa;
    background-color: #eee;
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
  [data-icon="expand"]{
    height: 40px;
    width: 40px;
    margin-bottom: 10px;
    margin-top: 10px;
    color: silver;
  }
</style>
