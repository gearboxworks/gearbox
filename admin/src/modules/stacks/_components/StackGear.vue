<template>
  <div
    :id="gearControlId"
    tabindex="0"
    class="project-gear"
  >

    <img
      v-if="service"
      :src="require('../_assets/'+service.attributes.program+'.svg')"
      :class="{'service-program': true, 'is-loaded': isLoaded, 'is-switching': isSwitching, 'is-switching-same': isSwitchingSame, 'is-switching-same-again': isSwitchingSameAgain}"
      @load="onImageLoaded"
    />
    <font-awesome-icon
      class="service-program is-unassigned"
      v-else
      :icon="['fa', 'expand']"
    />

    <h6 class="gear-role">
      {{gearspec.attributes.role}}
    </h6>

    <b-tooltip
      triggers="hover"
      :target="gearControlId"
      :key="gearControlId + '-' + (service ? service.id : 'unselected')"
      :title="programTooltip"
    />

   <stack-gear-popover
     :key = "gearControlId"
     :role = "gearspec.attributes.role"
     :gearControlId = "gearControlId"
     :serviceVersions = "serviceVersionsGroupedByProgram"
     :defaultService = "defaultService"
     :compatibleServiceId = "compatibleServiceId"
     :versionMismatchMessage = "versionMismatchMessage"
     @change-gear="onChangeProjectGear"
     @close-popover="closePopover"
   />
  </div>
</template>

<script>
import StackGearPopover from './StackGearPopover'

import StackMethodTypes from '../_store/public-types'
import GearspecMethodNames from '../../gearspecs/_store/public-types'
import ServiceMethodTypes from '../../services/_store/public-types'
import ProjectMethodTypes from '../../projects/_store/public-types'
const { GetterTypes: StackGetters } = StackMethodTypes
const { GetterTypes: ServiceGetters } = ServiceMethodTypes
const { GetterTypes: GearspecGetters } = GearspecMethodNames
const { ActionTypes: ProjectActions } = ProjectMethodTypes

export default {
  name: 'StackGear',
  components: {
    StackGearPopover
    // CoolSelect
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
      return this.projectPrefix + this.$_escapeIDAttr((this.stack ? this.stack.attributes.stackname + '-' : '') + this.gearspec.attributes.role)
    },

    gearspec () {
      if (!this.stackItem.gearspec) {
        throw new Error('Gearspec object is expected to be resolved by now!')
      }
      return this.stackItem.gearspec
    },

    serviceVersionsGroupedByProgram () {
      return this.$store.getters[GearspecGetters.GEARSPEC_SERVICE_VERSIONS_GROUPED_BY_PROGRAM](this.gearspec)
    },

    versionMismatchMessage () {
      let message = ''
      console.log()
      if ((this.defaultService || !!this.stackItem.serviceId) && !this.stackItem.service) {
        const requested = this.stackItem.serviceId.split(':')[1]
        const compatible = this.compatibleServiceId.split(':')[1]
        message = `Could not find the requested version (v.${requested}), will use the closest match (v.${compatible}) instead.`
      }
      return message
    },

    stack () {
      return this.$store.getters[StackGetters.FIND_BY](
        'id',
        this.gearspec.attributes.stack_id
      )
    },

    defaultService () {
      return this.$store.getters[GearspecGetters.DEFAULT_GEARSPEC_SERVICE](
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
    },

    programTooltip () {
      const serviceId = this.stackItem.serviceId
      const attributes = (serviceId && this.service) ? this.service.attributes : null

      let program = attributes ? attributes.program : ''
      let version = attributes ? attributes.version : ''

      if (serviceId && (!attributes || (this.service && this.service.id !== serviceId))) {
        /**
         * TODO move this logic to a global helper function
         */
        program = serviceId.split('/')[1].split(':')[0]
        version = serviceId.split('/')[1].split(':')[1]
      }

      return (program && version)
        ? (program + ' ' + version)
        : 'Service not selected'
    }

  },
  methods: {
    onImageLoaded () {
      this.isSwitching = false
      this.isLoaded = true
    },

    async onChangeProjectGear (selectedServiceId) {
      const previousId = this.service ? this.service.id : ''
      /**
       * TODO move this logic to a global helper function
       */
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
  .gear-role{
    margin-top:5px;
    margin-bottom: 0;
    clear: both;
  }
  .service-program {
    height: 64px;
    width: 64px;
  }
  .service-program.is-unassigned {
    width: 32px;
    height: 50px;
    color: var(--gray);
  }
/*  .service-program.is-loaded {
    animation-duration: 0.5s;
    animation-timing-function: cubic-bezier(0.075, 0.82, 0.165, 1);
    animation-delay: 0s;
    animation-iteration-count: 1;
    animation-direction: normal;
    animation-fill-mode: none;
    animation-play-state: running;
    animation-name: full-zoom;
  }
  @keyframes full-zoom {
    from {
      transform:scale(0)
    }
    to {
      transform: scale(1);
    }
  }
  .service-program.is-switching {
    animation-duration: 0.5s;
    animation-timing-function: cubic-bezier(0.075, 0.82, 0.165, 1);
    animation-delay: 0s;
    animation-iteration-count: 1;
    animation-direction: normal;
    animation-fill-mode: forwards;
    animation-play-state: running;
    animation-name: full-zoom-out;
  }
  @keyframes full-zoom-out {
    from {
      transform:scale(1)
    }
    to {
      transform: scale(0);
    }
  }
  .service-program.is-switching-same,
  .service-program.is-switching-same-again {
    animation-duration: 0.5s;
    animation-timing-function: cubic-bezier(0.075, 0.82, 0.165, 1);
    animation-delay: 0s;
    animation-iteration-count: 1;
    animation-direction: alternate;
    animation-fill-mode: none;
    animation-play-state: running;
  }
  .service-program.is-switching-same {
    animation-name: zoom-in-out;
  }
  @keyframes zoom-in-out {
    0% {
      transform:scale(1)
    }
    50% {
      transform:scale(1.1)
    }
    100% {
      transform: scale(0.75);
    }
  }
  .service-program.is-switching-same-again {
    animation-name: zoom-out-in;
  }
  @keyframes zoom-out-in {
    0% {
      transform:scale(1)
    }
    50% {
      transform:scale(1.1)
    }
    100% {
      transform: scale(0.75);
    }
  }
*/
</style>
