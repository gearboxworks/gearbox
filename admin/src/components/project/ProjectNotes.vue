<template>
  <b-form-textarea
    v-model="notes"
    placeholder="Notes..."
    :rows="isCollapsed ? 1: 3"
    @click="isCollapsed = false"
    max-rows="6"
    @change="maybeSubmit"
    :class="{'project-notes':true, 'is-collapsed': isCollapsed}"
  />
</template>

<script>
export default {
  name: 'ProjectNotes',
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
  .is-collapsed {
    height: 35px;
  }
</style>
