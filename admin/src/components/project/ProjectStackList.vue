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
      /**
       * returns project's services grouped by stack (indexed by stack_id)
       */
      var result = {}
      stackItems.forEach(stackItem => {
        const gearspec = this.gearspecBy('id', stackItem.gearspec_id)
        if (gearspec) {
          if (typeof result[gearspec.attributes.stack_id] === 'undefined') {
            result[gearspec.attributes.stack_id] = []
          }
          const service = stackItem.service_id ? this.serviceBy('id', stackItem.service_id) : null
          /**
           * note, when there is no exact match, service will be null,
           * but we will try to find a good-enough match further down the road;
           * that's why we need to pass over the original serviceId
           */
          result[gearspec.attributes.stack_id].push({
            gearspecId: stackItem.gearspec_id,
            gearspec,
            serviceId: stackItem.service_id,
            service
          })
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
