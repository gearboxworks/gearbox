<template>
  <div
    :id="gearControlId"
    :tabindex="projectIndex*100+stackIndex*10+itemIndex+1"
    class="project-gear"
  >

    <img
      v-if="service"
      :src="require('../../assets/'+service.attributes.program+'.svg')"
      :class="{'service-program': true, 'is-loaded': isLoaded, 'is-switching': isSwitching, 'is-switching-same': isSwitchingSame, 'is-switching-same-again': isSwitchingSameAgain }"
      @load="onImageLoaded"
    />
    <font-awesome-icon
      v-else
      :icon="['fa', 'expand']"
    />

    <h6 class="gear-role">{{gearspec.attributes.role}}</h6>

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

      <b-form-group>
        <label :for="`${gearControlId}-input`">{{gearspec.attributes.role}}:</label>
        <b-form-select
          :ref="`${gearControlId}-select`"
          :value="preselectClosestGearServiceId"
          :tabindex="projectIndex*100+stackIndex*10+itemIndex+9"
          @change="onChangeService($event)"
        >
          <option value="" v-if="!defaultService">Do not run this service</option>
          <option disabled :value="null">Select service...</option>
          <optgroup v-for="(services, groupLabel) in servicesGroupedByRole" :label="groupLabel" :key="groupLabel">
            <option v-for="serviceId in services" :value="serviceId" :key="serviceId" :disabled="project.attributes.enabled">{{serviceId.replace('gearboxworks/','')}}</option>
          </optgroup>
        </b-form-select>
        <b-alert :show="!stackItem.service" variant="warning">Note, the currently selected version of the service is different from what is in project specification!</b-alert>
        <b-alert :show="project.attributes.enabled">Note, you cannot change this service while the project is running!</b-alert>
      </b-form-group>
    </b-popover>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'

export default {
  name: 'StackGear',
  props: {
    'project': {
      type: Object,
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
    ...mapGetters(['serviceBy', 'gearspecBy', 'stackBy', 'stackDefaultServiceByRole', 'stackServicesByRole', 'preselectServiceId']),
    projectBase () {
      return 'gb-' + this.escAttr(this.project.id) + '-'
    },
    gearspec () {
      return this.stackItem.gearspec
    },
    service () {
      let service = null
      if (this.stackItem.service) {
        service = this.stackItem.service
      } else if (this.stackItem.serviceId) {
        const closestServiceId = this.preselectClosestGearServiceId
        if (closestServiceId) {
          service = this.serviceBy('id', closestServiceId)
        }
      }
      return service
    },
    stack () {
      return this.stackBy('id', this.gearspec.attributes.stack_id)
    },
    gearControlId () {
      return this.projectBase + (this.stack ? this.stack.attributes.stackname + '-' : '') + this.gearspec.attributes.role
    },
    defaultService () {
      return this.stackDefaultServiceByRole(this.stack, this.stackItem.gearspecId)
    },
    preselectClosestGearServiceId () {
      /**
       * As an example, for php:7.1.18 it will select php:7.1 or php:7 if exact match is not possible
       */
      return this.preselectServiceId(
        this.stackServicesByRole(this.stack, this.stackItem.gearspecId),
        this.defaultService,
        this.stackItem.serviceId
      )
    },
    servicesGroupedByRole () {
      const services = this.stackServicesByRole(this.stack, this.gearspec.id)
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
      this.$store.dispatch('changeProjectService', { 'projectId': this.project.id, gearspecId: this.gearspec.id, serviceId: selectedServiceId })
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

  .alert {
    margin-top: 1rem;
    margin-bottom: 0;
    padding: 0.5rem;
  }

</style>
