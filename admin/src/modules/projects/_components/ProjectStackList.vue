<template>
  <div
    :class="{'project-stack-list': true, 'start-collapsed': startCollapsed, 'is-loading': isLoading}"
    :id="`${projectBase}stack`" role="tablist"
  >
    <font-awesome-icon v-if="isLoading" icon="circle-notch" spin title="Loading project details..."/>
    <div v-else class="project-stack-list-wrap">
      <stack-card
        v-for="(stackItems, stackId, stackIndex) in groupedStackItems"
        :key="stackId"
        :stackId="stackId"
        :stackIndex="stackIndex"
        :stackItems="stackItems"
        :project="project"
        :projectIndex="projectIndex"
        :start-collapsed="startCollapsed || (!startCollapsed && Object.entries(groupedStackItems).length > 1)"
      >
      </stack-card>
    </div>
  </div>
</template>

<script>

import StackCard from '../../../components/stack/StackCard.vue'
// import StackCardSelect from './StackCardSelect'

import { mapGetters } from 'vuex'

export default {
  name: 'StackCardList',
  props: {
    'project': {
      type: Object,
      required: true
    },
    'projectIndex': {
      type: Number,
      required: true
    },
    'startCollapsed': {
      type: Boolean,
      required: false,
      default: false
    }
  },
  data () {
    return {
      id: this.project.id
    }
  },
  components: {
    StackCard
    // StackCardSelect
  },
  computed: {
    ...mapGetters(['serviceBy', 'gearspecBy']),
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    },
    isLoading () {
      return typeof this.project.attributes.stack === 'undefined'
    },
    groupedStackItems () {
      /**
       * returns project's services grouped by stack (indexed by stack_id)
       */
      var result = {}
      const stackItems = this.project.attributes.stack || []
      stackItems.forEach(stackItem => {
        if (stackItem.isRemoved) {
          return
        }
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
      // console.log(result)

      /**
       * sort gears by gear role
       */
      Object.keys(result).forEach((stackId) => {
        result[stackId] = result[stackId].sort((a, b) => a.gearspec.attributes.role > b.gearspec.attributes.role ? 1 : (a.gearspec.attributes.role === b.gearspec.attributes.role) ? 0 : -1)
      })

      // console.log('groupedStackItems', result)
      /**
       * sort stacks by stack id
       */
      return Object.keys(result).sort().reduce((r, key) => {
        // eslint-disable-next-line no-param-reassign
        r[key] = result[key]
        return r
      }, {})
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    }
  }
}
</script>

<style scoped>
  .project-stack-list.is-loading {
    color: #17a2b8;
    margin-left: 10px;
    /*display: inline-flex;*/
  }
</style>
