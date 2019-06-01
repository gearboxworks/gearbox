<template>
  <b-input-group :id="`${projectBase}location`" role="tabpanel">
    <b-form-input
      disabled
      :id="`${projectBase}location-input`"
      :value="resolveDir(currentBasedir, path)"
      class="location-input"
    />
    <b-input-group-append>
      <b-button
        variant="outline-info"
        title="Open project directory"
        v-b-tooltip.hover
        :id="`${projectBase}open-location`"
        href="#"
        class="folder-icon"
      >
        <font-awesome-icon
          :icon="['fa', 'folder']"
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
    }
  },
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
      selectedService: ''
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
    }
  }
}
</script>
<style scoped>
  .location-input {
    /*text-align: right;*/
  }
  .btn-outline-info {
    border-color: #ced4da;
  }
</style>
