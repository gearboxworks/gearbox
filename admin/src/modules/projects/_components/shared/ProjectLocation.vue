<template>
  <b-input-group
    :id="`${projectPrefix}location`"
    :class="{'input-group--location': true, 'is-collapsed': isCollapsed, 'is-multimodal': isMultimodal}"
    role="tabpanel"
  >
    <b-form-input
      readonly
      :ref="`${projectPrefix}location`"
      :value="location"
      class="location-input"
      v-show="!isCollapsed"
      autocomplete="off"
      @keyup.esc="isCollapsed = true"
    />
    <b-input-group-append>
      <b-button
        variant="outline-info"
        :title="$t('projects.fieldLocationCopy')"
        v-b-tooltip.hover
        href="#"
        class="btn--copy-dir"
        v-if="!isCollapsed"
        @click.stop="onCopyToClipboard"
      >
        <font-awesome-icon
          :icon="['fa', 'clone']"
        />
      </b-button>
      <b-button
        variant="outline-info"
        :title="isCollapsed ? $t('projects.fieldLocationView') : $t('projects.fieldLocationOpen')"
        v-b-tooltip.hover
        href="#"
        class="btn--open-dir"
        @click="onViewLocation"
      >
        <font-awesome-icon
          :icon="['fa', isCollapsed ? 'folder': 'folder-open']"
        />
      </b-button>
    </b-input-group-append>
  </b-input-group>
</template>

<script>
import { BasedirGetters } from '../../../basedirs/_store/method-names'

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
      default: true
    }
  },
  data () {
    return {
      isCollapsed: this.isMultimodal
    }
  },
  computed: {
    location () {
      const { basedir: basedirId, path } = this.project.attributes
      const basedir = this.$store.getters[BasedirGetters.FIND_BY]('id', basedirId)
      const dir = basedir ? basedir.attributes.basedir : ''
      return dir + ((dir.indexOf('/') !== -1) ? '/' : '\\') + path
    }
  },
  methods: {

    onViewLocation () {
      if (this.isMultimodal && this.isCollapsed) {
        this.isCollapsed = false
        this.$nextTick(() => {
          const el = this.$refs[`${this.projectPrefix}location`].$el
          el.focus()
          el.setSelectionRange(0, 9999)
        })
      } else {
        // TODO call API method to open directory in file manager
        console.log('TODO: call API method to open directory in file manager')
        if (this.isMultimodal) {
          this.$nextTick(() => {
            this.isCollapsed = true
          })
        }
      }
    },

    onCopyToClipboard () {
      /**
       * @see https://github.com/Inndy/vue-clipboard2
       */
      this.$copyText(this.$refs[`${this.projectPrefix}location`].$el.value).then((e) => console.log('Copied:', e.text))

      if (this.isMultimodal) {
        this.$nextTick(() => {
          this.isCollapsed = true
        })
      }
    }
  }
}
</script>
<style scoped>

  .btn-outline-info {
    border-color: #ced4da;
  }

  .is-collapsed .btn-outline-info {
    border-color: transparent;
    border-top-left-radius: 0.25rem;
    border-bottom-left-radius: 0.25rem;
  }

</style>
