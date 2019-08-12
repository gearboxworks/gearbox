<template>
  <b-popover
    :target="gearControlId"
    :ref="`${gearControlId}-popover`"
    triggers="focus"
    placement="bottom"
  >
    <template slot="title">
      <b-button @click="$emit('close-popover')" class="close" :aria-label="$t('button.close')">
        <span class="d-inline-block" aria-hidden="true">&times;</span>
      </b-button>
      {{$t('services.change')}}
    </template>

    <b-form-group>
      <label :for="`${gearControlId}-input`">{{role}}:</label>

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
        <option v-if="!defaultService" value="">{{$t('services.doNotRun')}}</option>
        <option disabled :value="null">{{$t('services.select')}}</option>
        <optgroup
          v-for="(serviceVersions, program) in serviceVersions"
          :label="program"
          :key="program"
        >
          <option
            v-for="(serviceVersion, serviceId) in serviceVersions"
            :value="serviceId"
            :key="serviceId"
            :disabled="isProjectEnabled"
          >
            {{serviceVersion}}
          </option>
        </optgroup>
      </b-form-select>
      <b-alert :show="isProjectEnabled">{{$t('services.readonlyWhileRunning')}}</b-alert>
    </b-form-group>
  </b-popover>
</template>

<script>
export default {
  name: 'GearBadgePopover',
  props: {
    isProjectEnabled: {
      type: Boolean,
      require: true
    },
    gearControlId: {
      type: String,
      require: true
    },
    role: {
      type: String,
      require: true
    },
    serviceVersions: {
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
  },
  methods: {
    onChangeProjectGear (selectedServiceId) {
      this.$emit('change-gear', selectedServiceId)
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
