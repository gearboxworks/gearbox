<template>
  <div
    :id="gearControlId"
    :tabindex="projectIndex*100+stackIndex*10+itemIndex+1"
    class="project-gear"
  >

    <img
      v-if="service"
      :src="require('../assets/'+service.attributes.program+'.svg')"
      :class="{'service-program': true, 'is-loaded': isLoaded, 'is-switching': isSwitching, 'is-switching-same': isSwitchingSame, 'is-switching-same-again': isSwitchingSameAgain }"
      @load="onImageLoaded"
    />
    <font-awesome-icon
      v-else
      :icon="['fa', 'expand']"
    />

    <h6 class="gear-role">{{gear.attributes.role}}</h6>

    <b-tooltip
      triggers="hover"
      :target="gearControlId"
      :key="gearControlId+'-'+(service?service.id:'unselected')"
      :title="programTooltip"
    />

    <b-popover
      :target="gearControlId"
      :container="`${projectBase}stack`"
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

      <div>
        <label :for="`${gearControlId}-input`">{{gear.attributes.role}}:</label>
        <b-form-select
          :id="`${gearControlId}-input`"
          :value="service ? service.id : ''"
          :tabindex="projectIndex*100+stackIndex*10+itemIndex+9"
          @change="onChangeService($event)"
        >
          <option value="" v-if="!defaultService">Do not run this service</option>
          <option disabled value="">Select service...</option>
          <optgroup v-for="(services, groupLabel) in groupedGearServices" :label="groupLabel" :key="groupLabel">
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
  name: 'StackGear',
  props: {
    'projectId': {
      type: String,
      required: true
    },
    'stackItem': {
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
    'itemIndex': {
      type: Number,
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
    ...mapGetters(['gearBy', 'stackBy', 'stackDefaultServiceByGear', 'stackServicesByGear']),
    projectBase () {
      return 'gb-' + this.escAttr(this.projectId) + '-'
    },
    gear () {
      return this.stackItem.gear
    },
    service () {
      return this.stackItem.service
    },
    stack () {
      return this.stackBy('id', this.gear.attributes.stack_id)
    },
    gearControlId () {
      return this.projectBase + (this.stack ? this.stack.attributes.stackname + '-' : '') + this.gear.attributes.role
    },
    defaultService () {
      return this.stackDefaultServiceByGear(this.stack, this.gear.id)
    },
    // preselectGearService () {
    //   const defaultService = this.defaultService()
    //
    //   /**
    //    * Resolve default option:
    //    * - if exact match is found, use it
    //    * - otherwise, use the last in the list that have the specified name mentioned (hopefully that will be the latest version)
    //    */
    //   let firstFound = -1
    //   let exactFound = -1
    //   if (defaultService) {
    //     for (var i = gearServices.length; i--;) {
    //       if (gearServices[i].indexOf(defaultService) !== -1) {
    //         if (firstFound === -1) {
    //           firstFound = i
    //         }
    //         if (gearServices[i] === defaultService) {
    //           exactFound = i
    //           break
    //         }
    //       }
    //     }
    //   }
    //   const serviceId = (firstFound !== -1)
    //     ? gearServices[ exactFound !== -1 ? exactFound : firstFound ]
    //     : ''
    //
    //   return serviceId
    // },
    groupedGearServices () {
      const services = this.stackServicesByGear(this.stack, this.gear.id)
      // console.log(services)
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
    programTooltip () {
      const attributes = this.service ? this.service.attributes : null
      return attributes
        ? (attributes.program + ' ' + attributes.version)
        : 'Service not selected'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    onChangeService (selectedServiceId) {
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
      this.$store.dispatch('changeProjectService', { 'projectId': this.projectId, gearId: this.gear.id, serviceId: selectedServiceId })
      this.closePopover()
    },
    closePopover () {
      this.$root.$emit('bv::hide::popover', this.gearControlId)
    },
    onImageLoaded (a) {
      this.isSwitching = false
      this.isLoaded = true
    }
  }
}
</script>

<style scoped>
  .project-gear{
    outline: none;
  }
  .gear-role{
    margin-top:5px;
    margin-bottom: 0;
    clear: both;
  }
  .service-program {
    height: 64px;
    width: 64px;
  }
  .service-program.is-loaded {
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

</style>
