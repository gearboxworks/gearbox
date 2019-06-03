<template>
  <div class="project-note-list">
    <b-button
      variant="outline-warning"
      title="Dummy note 1"
      v-b-tooltip.hover
      href="#"
      class="btn--view-note"
      v-if="Math.random() > 0.3"
    >
      <font-awesome-icon
        :icon="['fa', 'sticky-note']"
      />
    </b-button>
    <b-button
      variant="outline-warning"
      title="Dummy note 2"
      v-b-tooltip.hover
      href="#"
      class="btn--view-note"
      v-if="Math.random() > 0.5"
    >
      <font-awesome-icon
        :icon="['fa', 'sticky-note']"
      />
    </b-button>
    <b-button
      variant="outline-warning"
      title="Dummy note 3"
      v-b-tooltip.hover
      href="#"
      class="btn--view-note"
      v-if="Math.random() > 0.7"
    >
      <font-awesome-icon
        :icon="['fa', 'sticky-note']"
      />
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'ProjectNoteList',
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
      isCollapsed: true
    }
  },
  computed: {
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    maybeSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          id: this.id,
          attributes: this.$data
        }
      ).then(() => {
        // this.$router.push('/project/' + this.hostname)
      })
    }
  }
}
</script>
<style scoped>
  .btn-outline-warning {
    border-color: transparent;
    color: #e4bd77;
  }

  .btn-outline-warning:hover {
    color: #212529;
    background-color: #e4bd77;
  }
</style>
