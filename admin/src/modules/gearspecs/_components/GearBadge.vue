<template>
  <div
    :id="gearControlId"
    tabindex="0"
    class="project-gear"
  >

  <service-icon
    :gear-control-id = "gearControlId"
    :service-id="this.stackItem.serviceId"
    :service = "service"
    :role = "role"
    :changingStatus = "changingStatus"
    @image-loaded = "onImageLoaded"
  />

   <gear-badge-popover
     :is-project-enabled="this.project.attributes.enabled"
     :key = "gearControlId"
     :role = "role"
     :gear-control-id = "gearControlId"
     :service-versions = "serviceVersionsGroupedByProgram"
     :default-service = "defaultService"
     :compatible-service-id = "compatibleServiceId"
     :version-mismatch-message = "versionMismatchMessage"
     @change-gear="onChangeProjectGear"
     @close-popover="closePopover"
   />
  </div>
</template>

<script>
import ServiceIcon from '../../services/_components/ServiceIcon'
import GearBadgePopover from './GearBadgePopover'

import { GearspecGetters } from '../_store/method-names'
import { StackGetters } from '../../stacks/_store/method-names'
import { ServiceGetters } from '../../services/_store/method-names'
import { ProjectActions } from '../../projects/_store/method-names'

import { escapeIDAttr, versionFromServiceId, programFromServiceId } from '../../_helpers'

export default {
  name: 'GearBadge',
  components: {
    ServiceIcon,
    GearBadgePopover
  },
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {
    'stackItem': {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      isLoaded: false,
      isSwitching: true,
      isSwitchingSame: false,
      isSwitchingSameAgain: false
    }
  },
  computed: {
    gearControlId () {
      return this.projectPrefix + escapeIDAttr((this.stack ? this.stack.attributes.stackname + '-' : '') + this.role)
    },

    changingStatus () {
      return {
        isLoaded: this.isLoaded,
        isSwitching: this.isSwitching,
        isSwitchingSame: this.isSwitchingSame,
        isSwitchingSameAgain: this.isSwitchingSameAgain
      }
    },

    gearspec () {
      if (!this.stackItem.gearspec) {
        throw new Error('Gearspec object is expected to be resolved by now!')
      }
      return this.stackItem.gearspec
    },

    role () {
      return this.gearspec.attributes.role
    },

    serviceVersionsGroupedByProgram () {
      return this.$store.getters[GearspecGetters.SERVICE_VERSIONS_GROUPED_BY_PROGRAM](this.gearspec)
    },

    versionMismatchMessage () {
      if ((this.defaultService || !!this.stackItem.serviceId) && !this.stackItem.service) {
        const requestedVersion = versionFromServiceId(this.stackItem.serviceId)
        const compatibleVersion = versionFromServiceId(this.compatibleServiceId)
        return `Could not find the requested version (v.${requestedVersion}), will use the closest match (v.${compatibleVersion}) instead.`
      }
      return ''
    },

    stack () {
      return this.$store.getters[StackGetters.FIND_BY](
        'id',
        this.gearspec.attributes.stack_id
      )
    },

    defaultService () {
      return this.$store.getters[GearspecGetters.DEFAULT_SERVICE](
        this.stackItem.gearspec
      )
    },

    compatibleServiceId () {
      return this.$store.getters[GearspecGetters.FIND_COMPATIBLE_SERVICE](
        this.stackItem.gearspec,
        this.stackItem.serviceId
      )
    },

    service () {
      let service = null
      if (this.stackItem.service) {
        service = this.stackItem.service
      } else if (this.stackItem.serviceId) {
        const compatibleServiceId = this.compatibleServiceId
        if (compatibleServiceId) {
          service = this.$store.getters[ServiceGetters.FIND_BY](
            'id',
            compatibleServiceId
          )
        }
      }
      return service
    }
  },
  methods: {
    onImageLoaded () {
      this.isSwitching = false
      this.isLoaded = true
    },

    async onChangeProjectGear (selectedServiceId) {
      const previousServiceId = this.service ? this.service.id : ''
      const oldProgram = programFromServiceId(previousServiceId)
      const newProgram = programFromServiceId(selectedServiceId)

      if (oldProgram !== newProgram) {
        this.isLoaded = false
        this.isSwitching = true
        this.isSwitchingSame = false
        this.isSwitchingSameAgain = false
      } else {
        if (previousServiceId !== selectedServiceId) {
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
        /**
         * TODO emit event and let Project module do the actual changing
         */
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
  .project-gear {
    text-align: center;
    max-width: 110px;
    padding: 5px;
    margin: 5px;
    cursor: pointer;
    border: 1px solid transparent;
    border-radius: 4px;
    transition: all 400ms;
  }
  .project-gear:hover,
  .project-gear:focus {
    border: 1px solid #aaa;
    background-color: #eee;
  }
  .project-gear{}
</style>
