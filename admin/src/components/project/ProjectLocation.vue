<template>
  <b-input-group :id="`${projectBase}location`"  :class="{'input-group--location': true, 'is-collapsed': isCollapsed, 'is-multimodal': isMultimodal}" role="tabpanel">
    <b-form-input
      disabled
      :id="`${projectBase}location-input`"
      :value="resolveDir(currentBasedir, path)"
      class="location-input"
      v-if="!isCollapsed"
    />
    <b-input-group-append>
      <b-button
        variant="outline-info"
        title="Copy to clipboard"
        v-b-tooltip.hover
        href="#"
        class="btn--copy-dir"
        v-if="!isCollapsed"
        @click="onCopyToClipboard"
      >
        <font-awesome-icon
          :icon="['fa', 'clone']"
        />
      </b-button>
      <b-button
        variant="outline-info"
        :title="isCollapsed ? 'View project location' : 'Open in file manager'"
        v-b-tooltip.hover
        href="#"
        class="btn--open-dir"
        @click="onButtonClicked"
      >
        <font-awesome-icon
          :icon="['fa', isCollapsed ? 'folder': 'folder-open']"
        />
      </b-button>
    </b-input-group-append>
  </b-input-group>
</template>

<script>
import { mapGetters } from 'vuex'

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
      default: true
    }
  },
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
      isCollapsed: this.isMultimodal
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
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    resolveDir (dir, path) {
      return dir + ((dir.indexOf('/') !== -1) ? '/' : '\\') + path
    },
    onButtonClicked () {
      if (this.isMultimodal && this.isCollapsed) {
        this.isCollapsed = false
      } else {
        // TODO call API method to open directory in file manager
        console.log('TODO call API method to open directory in file manager')
        if (this.isMultimodal) {
          this.$nextTick(() => {
            this.isCollapsed = true
          })
        }
      }
    },
    onCopyToClipboard () {
      // @TODO implement copy to clipboard
      // @see https://github.com/Inndy/vue-clipboard2
      console.log('TODO implement copy to clipboard')
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

  .is-collapsed .btn-outline-info {
    border-color: transparent;
    border-top-left-radius: 0.25rem;
    border-bottom-left-radius: 0.25rem;
  }

</style>
