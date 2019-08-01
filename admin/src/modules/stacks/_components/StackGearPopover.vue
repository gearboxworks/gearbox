<template>
  <b-popover
    :target="gearControlId"
    :container="gearControlId"
    :ref="`${gearControlId}-popover`"
    triggers="focus"
    placement="bottom"
  >
    <template slot="title">
      <b-button @click="closePopover" class="close" aria-label="Close">
        <span class="d-inline-block" aria-hidden="true">&times;</span>
      </b-button>
      Change service
    </template>

    <b-form-group>
      <label :for="`${gearControlId}-input`">{{gearspec.attributes.role}}:</label>

      <!--        <cool-select-->
      <!--          value="Option2"-->
      <!--          :items="[{'text':'Option1'}, {'text':'Option2'}]"-->
      <!--        />-->

      <b-form-select
        :ref="`${gearControlId}-select`"
        :value="compatibleServiceId"
        tabindex="0"
        @change="onChangeProjectGear($event)"
      >
        <option v-if="!defaultService" value="">Do not run this service</option>
        <option disabled :value="null">Select service...</option>
        <optgroup
          v-for="(services, groupLabel) in servicesGroupedByRole"
          :label="groupLabel"
          :key="groupLabel"
        >
          <option
            v-for="serviceId in services"
            :value="serviceId"
            :key="serviceId"
            :disabled="project.attributes.enabled"
          >
            {{serviceId.replace('gearboxworks/','')}}
          </option>
        </optgroup>
      </b-form-select>
      <b-alert :show="(defaultService || !!stackItem.serviceId) && !stackItem.service" variant="warning">Could not find the requested version (v.{{stackItem.serviceId.split(':')[1]}}), will use the closest match (v.{{compatibleServiceId.split(':')[1]}}) instead.</b-alert>
      <b-alert :show="project.attributes.enabled">Note, you cannot change this service while the project is running!</b-alert>
    </b-form-group>
  </b-popover>
</template>

<script>
import StackMethodTypes from '../../stacks/_store/public-types'
import ProjectMethodTypes from '../../projects/_store/public-types'
const { GetterTypes: StackGetters } = StackMethodTypes
const { ActionTypes: ProjectActions } = ProjectMethodTypes

export default {
  name: 'StackGearPopover',
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {
    gearControlId: {
      type: String,
      require: true
    },
    stackItem: {
      type: Object,
      require: true
    },
    gearspec: {
      type: Object,
      require: true
    },
    service: {
      type: Object,
      require: true
    },
    stack: {
      type: Object,
      require: true
    },
    defaultService: {
      type: String,
      require: true
    },
    compatibleServiceId: {
      type: String,
      require: true
    }
  },
  computed: {

    servicesGroupedByRole () {
      const services = this.$store.getters[StackGetters.GEARSPEC_SERVICES](
        this.stack,
        this.gearspec.id
      )
      // console.log(services)
      const result = {}
      if (Array.isArray(services)) {
        for (const index in services) {
          const base = services[index].split(':')[0].replace('gearboxworks/', '')
          if (typeof result[base] === 'undefined') {
            result[base] = {}
          }
          result[base][index] = services[index]
        }
      }
      return result
    },

    programTooltip () {
      const serviceId = this.stackItem.serviceId
      const attributes = (serviceId && this.service) ? this.service.attributes : null

      let program = attributes ? attributes.program : ''
      let version = attributes ? attributes.version : ''

      if (serviceId && (!attributes || (this.service && this.service.id !== serviceId))) {
        program = serviceId.split('/')[1].split(':')[0]
        version = serviceId.split('/')[1].split(':')[1]
      }

      return (program && version)
        ? (program + ' ' + version)
        : 'Service not selected'
    }
  },
  methods: {

    async onChangeProjectGear (selectedServiceId) {
      const previousId = this.service ? this.service.id : ''
      const program1 = previousId ? previousId.split('/')[1].split(':')[0] : ''
      const program2 = selectedServiceId ? selectedServiceId.split('/')[1].split(':')[0] : ''

      if (program1 !== program2) {
        this.isLoaded = false
        this.isSwitching = true
        this.isSwitchingSame = false
        this.isSwitchingSameAgain = false
      } else {
        if (previousId !== selectedServiceId) {
          if (!this.isSwitchingSame && !this.isSwitchingSameAgain) {
            this.isSwitchingSame = true
            this.isSwitchingSameAgain = false
          } else {
            this.isSwitchingSame = !this.isSwitchingSame
            this.isSwitchingSameAgain = !this.isSwitchingSameAgain
          }
        }
      }
      try {
        await this.$store.dispatch(
          ProjectActions.CHANGE_GEAR,
          {
            project: this.project,
            gearspecId: this.gearspec.id,
            serviceId: selectedServiceId
          }
        )
      } catch (e) {
        console.error(e.message)
      }
      this.closePopover()
    },

    closePopover () {
      this.$root.$emit('bv::hide::popover', this.gearControlId)
    }
  }
}
</script>

<style scoped>
  .alert {
    margin-top: 1rem;
    margin-bottom: 0;
    padding: 0.5rem;
  }
</style>
