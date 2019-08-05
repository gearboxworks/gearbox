<template>
  <div>
    <img
      v-if="service"
      :src="require('../_assets/'+service.attributes.program+'.svg')"
      :class="{'service-program': true, 'is-loaded': changingStatus.isLoaded, 'is-switching': changingStatus.isSwitching, 'is-switching-same': changingStatus.isSwitchingSame, 'is-switching-same-again': changingStatus.isSwitchingSameAgain}"
      @load="$emit('image-loaded')"
    />
    <font-awesome-icon
      class="service-program is-unassigned"
      v-else
      :icon="['fa', 'expand']"
    />

    <h6 class="gear-role">
      {{role}}
    </h6>

    <b-tooltip
      triggers="hover"
      :target="gearControlId"
      :key="gearControlId + '-' + (service ? service.id : 'unselected')"
      :title="programTooltip"
    />
  </div>
</template>

<script>
import { versionFromServiceId, programFromServiceId } from '../../_helpers'

export default {
  name: 'ServiceIcon',
  props: {
    /**
     * Note, serviceId might be different from service.id!
     */
    'serviceId': {
      type: String,
      required: true
    },
    'service': {
      type: Object,
      required: true
    },
    'gearControlId': {
      type: String,
      required: true
    },
    'role': {
      type: String,
      required: true
    },
    'changingStatus': {
      type: Object,
      required: true
    }
  },

  computed: {
    programTooltip () {
      const serviceId = this.serviceId
      const attributes = (serviceId && this.service) ? this.service.attributes : null

      let program = attributes ? attributes.program : ''
      let version = attributes ? attributes.version : ''

      if (serviceId && (!attributes || (this.service && this.service.id !== serviceId))) {
        program = programFromServiceId(serviceId)
        version = versionFromServiceId(serviceId)
      }

      return (program && version)
        ? (program + ' ' + version)
        : 'Service not selected'
    }
  }
}
</script>

<style scoped>
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
