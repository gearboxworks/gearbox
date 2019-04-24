<template>
  <div
    :id="serviceControlId"
    :tabindex="projectIndex*100+stackIndex*10+serviceIndex+1"
    class="project-service"
  >

    <img v-if="program" :src="require('../assets/'+program+'.svg')" class="service-program" />
    <font-awesome-icon
      v-else
      :icon="['fa', 'expand']"
    />

    <h6 class="service-role">{{role}}</h6>

    <b-tooltip
      triggers="hover"
      :target="serviceControlId"
      :title="programTooltip(program, version)"
      :key="id"
    />

    <b-popover
      :target="serviceControlId"
      :container="`${projectControlId}stack`"
      :ref="`${serviceControlId}-popover`"
      triggers="focus"
      placement="bottom"
    >
      <template slot="title">
        <b-button @click="onClosePopoverFor(serviceControlId)" class="close" aria-label="Close">
          <span class="d-inline-block" aria-hidden="true">&times;</span>
        </b-button>
        Change service
      </template>

      <div>
        <label :for="`${serviceControlId}-input`">{{role}}:</label>
        <b-form-select
          :id="`${serviceControlId}-input`"
          :value="id"
          :tabindex="projectIndex*100+stackIndex*10+serviceIndex+9"
          @change="changeProjectService(gearspec_id,$event,role)"
        >
          <option value="" v-if="!stackRoleDefaultService(role)">Do not run this service</option>
          <option disabled value="">Select service...</option>
          <optgroup v-for="(services, groupLabel) in serviceOptionGroups(stackRoleServices(role))" :label="groupLabel" :key="groupLabel">
            <option v-for="serviceId in services" :value="serviceId" :key="serviceId">{{serviceId.replace('gearboxworks/','')}}</option>
          </optgroup>
        </b-form-select>
      </div>
    </b-popover>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'

export default {
  name: 'ProjectService',
  props: {
    'projectId': {
      type: String,
      required: true
    },
    'service': {
      type: Object,
      required: true
    },
    'projectIndex': {
      type: Number,
      required: true
    },
    'stackIndex': {
      type: Number,
      required: true
    },
    'serviceIndex': {
      type: Number,
      required: true
    }
  },
  data () {
    return {
      id: this.service.id,
      ...this.service.attributes
    }
  },
  computed: {
    ...mapGetters(['gearspecBy', 'stackBy']),
    projectControlId () {
      return this.escAttr(this.projectId) + '-'
    },
    serviceControlId () {
      return this.projectControlId + (this.stack ? this.stack.attributes.stackname + '-' : '') + this.role
    },
    gear () {
      return this.gearspecBy('id', this.gearspec_id)
    },
    role () {
      const gear = this.gear
      return gear ? this.escAttr(gear.attributes.role) : ''
    },
    stack () {
      const gear = this.gear
      return gear ? this.stackBy('id', gear.attributes.stack_id) : null
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    stackRoleObject (role) {
      const gear = this.gear
      const stack = gear ? this.stackBy('id', gear.attributes.stack_id) : null
      const members = stack ? stack.attributes.members : []
      return members.find(m => m.role === role)
    },
    stackRoleServices (role) {
      const stackRole = this.stackRoleObject(role)
      return stackRole ? stackRole.services : []
    },
    stackRoleDefaultService (role) {
      const stackRole = this.stackRoleObject(role)
      return stackRole
        ? ((typeof stackRole.default_service !== 'undefined')
          ? stackRole.default_service
          : '')
        : ''
    },
    stackRoleSmartDefaultService (role) {
      const stackRole = this.stackRoleObject(role)
      const defaultService = this.stackRoleDefaultService(role)

      /**
       * Resolve default option:
       * - if exact match is found, use it
       * - otherwise, use the last in the list that have the specified name mentioned (hopefully that will be the latest version)
       */
      let firstFound = -1
      let exactFound = -1
      if (defaultService) {
        for (var i = stackRole.services.length; i--;) {
          if (stackRole.services[i].indexOf(defaultService) !== -1) {
            if (firstFound === -1) {
              firstFound = i
            }
            if (stackRole.services[i] === defaultService) {
              exactFound = i
              break
            }
          }
        }
      }
      const smartDefault = (firstFound !== -1)
        ? stackRole.services[ exactFound !== -1 ? exactFound : firstFound ]
        : ''

      return smartDefault
    },
    serviceOptionGroups (services) {
      const result = {}
      for (const index in services) {
        const base = services[index].split(':')[0].replace('gearboxworks/', '')
        if (typeof result[base] === 'undefined') {
          result[base] = {}
        }
        result[base][index] = services[index]
      }
      return result
    },
    removeProjectStack () {
      this.$store.dispatch('removeProjectStack', { 'projectId': this.projectId, 'stackName': this.stackName })
    },
    programTooltip (program, version) {
      return (program && version)
        ? (program + ' ' + version)
        : 'Service not selected'
    },
    changeProjectService (gearspecId, serviceId, role) {
      this.$store.dispatch('changeProjectService', { 'projectId': this.projectId, gearspecId, serviceId })
      this.onClosePopoverFor(this.serviceControlId)
    },
    onClosePopoverFor (triggerElementId) {
      this.$root.$emit('bv::hide::popover', triggerElementId)
    }
  }
}
</script>

<style scoped>
  .project-service{
    outline: none;
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
