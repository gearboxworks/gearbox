<template>
  <b-input-group
    :id="`${projectPrefix}hostname-group`"
    :class="{'input-group--hostname': true, 'is-editing': isEditing, 'is-multimodal': isMultimodal, 'is-modified': isModified, 'is-updating': isUpdating}"
    role="tabpanel"
  >
    <b-form-input
      class="hostname-input"
      type="text"
      v-model="hostname"
      :size="isMultimodal ? 'lg': 'md'"
      required
      @click.prevent="onInputClicked"
      placeholder=""
      :readonly="isUpdating || project.attributes.enabled"
      autocomplete="off"
    />
    <b-input-group-append v-if="isEditing">
      <b-button
        variant="outline-info"
        :title="isModified ? 'Submit the new hostname' : (isMultimodal ? 'Cancel': 'No changes to submit')"
        v-b-tooltip.hover
        class="btn--submit"
        :disabled="(!isMultimodal && !isModified) || isUpdating"
        @click.prevent="onSubmitOrCancel"
      >
        <font-awesome-icon
          v-if="isUpdating"
          icon="circle-notch"
          spin
        />
        <font-awesome-icon
          v-else
          :icon="['fa', isMultimodal ? (isModified ? 'check': 'times') : 'check']"
        />
      </b-button>
    </b-input-group-append>
  </b-input-group>
</template>

<script>

import ProjectMethodTypes from '../../_store/public-types'
const { ActionTypes: ProjectActions } = ProjectMethodTypes

export default {
  name: 'ProjectLocation',
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {
    isMultimodal: {
      type: Boolean,
      required: false,
      default: false
    }
  },
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
      isEditing: !this.isMultimodal,
      isModified: false,
      isUpdating: false
    }
  },
  computed: {},
  watch: {
    hostname: function (val, oldVal) {
      this.isModified = !!val && val !== this.project.id
    }
  },
  methods: {

    onInputClicked () {
      if (!this.isEditing) {
        if (this.project.attributes.enabled) {
          this.$emit('show-alert', 'Hostname cannot be changed while the project is running!')
        } else if (this.isMultimodal) {
          this.isEditing = true
        }
      }
    },

    onSubmitOrCancel () {
      if (this.isModified) {
        // TODO make an API call to change project hostname
        console.log('TODO call the API to change project hostname')
        this.maybeSubmit(this.hostname)
      } else {
        this.isEditing = false
      }
    },

    async maybeSubmit (value) {
      try {
        this.isUpdating = true

        await this.$store.dispatch(ProjectActions.UPDATE_HOSTNAME,
          {
            project: this.project,
            hostname: value
          }
        )

        if (this.isMultimodal) {
          this.isEditing = false
        }
        this.isModified = false
        this.isUpdating = false
      } catch (e) {
        console.error(e.message)
      }
    }
  }
}
</script>
<style scoped>

  .hostname-input{
    padding-left: 11px;
    font-weight: bold;
    border: 1px solid transparent;
    cursor: pointer;
  }

  .is-editing .hostname-input {
    border-color: #ced4da;
    cursor: text;
  }

  .is-multimodal .btn--submit {
    border-color: #ced4da;
    font-size: 22px;
  }

</style>
