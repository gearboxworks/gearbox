<template>
  <b-input-group
    :id="`${projectBase}hostname-group`"
    :class="{'input-group--hostname': true, 'is-editing': isEditing, 'is-multimodal': isMultimodal, 'is-modified': isModified, 'is-updating': isUpdating}"
    role="tabpanel"
  >
    <b-form-input
      :id="`hostname-input-${projectIndex}`"
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
        @click.prevent="onButtonClicked"
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
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'ProjectLocation',
  props: {
    project: {
      type: Object,
      required: true
    },
    projectIndex: {
      type: Number,
      required: true
    },
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
  computed: {
    ...mapGetters({ basedirBy: 'basedirBy' }),
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    },
    currentBasedir () {
      const basedir = this.basedirBy('id', this.basedir)
      return basedir ? basedir.attributes.basedir : ''
    }
  },
  methods: {
    ...mapActions(['updateProjectHostname']),
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    resolveDir (dir, path) {
      return dir + ((dir.indexOf('/') !== -1) ? '/' : '\\') + path
    },
    onInputClicked () {
      if (!this.isEditing) {
        if (this.project.attributes.enabled) {
          this.$emit('show-alert', 'Hostname cannot be changed while the project is running!')
        } else if (this.isMultimodal) {
          this.isEditing = true
        }
      }
    },
    onButtonClicked () {
      if (this.isModified) {
        // TODO make an API call to change project hostname
        console.log('TODO call the API to change project hostname')
        this.maybeSubmit(this.hostname)
      } else {
        this.isEditing = false
      }
    },
    maybeSubmit (value) {
      this.isUpdating = true
      /**
       * TODO: deal with timestamp
       * TODO: check validity
       */
      this.updateProjectHostname(
        {
          projectId: this.id,
          hostname: value
        }
      ).then(() => {
        if (this.isMultimodal) {
          this.isEditing = false
        }
        this.isModified = false
        this.isUpdating = false
      })
    }
  },
  watch: {
    hostname: function (val, oldVal) {
      this.isModified = !!val && val !== this.project.id
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
