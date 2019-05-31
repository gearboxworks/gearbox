<template>
  <div role="tablist" class="project-stack-list" :id="`${projectBase}stack`">
    <project-stack
      v-for="(stackItems, stackId, stackIndex) in groupedStackItems(projectStackItems)"
      :key="stackId"
      :stackId="stackId"
      :stackIndex="stackIndex"
      :stackItems="stackItems"
      :project="project"
      :projectIndex="projectIndex"
      :is-collapsible="isCollapsible"
    >
    </project-stack>
  </div>
</template>

<script>

import ProjectStack from './ProjectStack.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'ProjectStackList',
  props: {
    'project': {
      type: Object,
      required: true
    },
    'projectIndex': {
      type: Number,
      required: true
    },
    'isCollapsible': {
      type: Boolean,
      required: false,
      default: false
    }
  },
  data () {
    return {
      id: this.project.id,
      projectStackItems: this.project.attributes.stack
    }
  },
  components: {
    ProjectStack
  },
  computed: {
    ...mapGetters(['serviceBy', 'gearspecBy']),
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    groupedStackItems (stackItems) {
      var result = {}
      stackItems.forEach(stackItem => {
        const gearspec = this.gearspecBy('id', stackItem.gearspec_id)
        if (gearspec) {
          if (typeof result[gearspec.attributes.stack_id] === 'undefined') {
            result[gearspec.attributes.stack_id] = []
          }
          const service = stackItem.service_id ? this.serviceBy('id', stackItem.service_id) : null
          /**
           * grouping project's services by stack
           */
          result[gearspec.attributes.stack_id].push({ gearspec, service })
        }
      })
      return result
    }
  }
}
</script>

<style scoped>
  .project-stack-list {
    margin-top: 10px;
  }
</style>
