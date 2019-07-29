<template>
  <div
    :class="{'project-stack-list': true, 'start-expanded': startExpanded, 'is-loading': isLoading}"
    :id="`${projectPrefix}stack`"
    role="tablist"
  >
    <font-awesome-icon
      v-if="isLoading"
      icon="circle-notch"
      spin
      title="Loading project details..."
    />
    <div
      v-else
      class="project-stack-list-wrap"
    >
      <stack-card
        v-for="(stackItems, stackId) in groupedStackItems"
        :key="stackId"
        :stackId="stackId"
        :stackItems="stackItems"
        :is-expanded="(expandedStackIds[stackId] && expandedStackIds[stackId] > 0) || ((!expandedStackIds[stackId] || expandedStackIds[stackId] === 0) && Object.entries(groupedStackItems).length === 1)"
        @expand-collapse="onExpandCollapseStack"
      />
    </div>
  </div>
</template>

<script>
// || (!stackToExpand && ( startExpanded || (startExpanded && Object.entries(groupedStackItems).length > 1)))
import { mapGetters } from 'vuex'
import StackCard from '../../../../components/stack/StackCard.vue'

export default {
  name: 'StackCardList',
  components: {
    StackCard
  },
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {
    startExpanded: {
      type: Boolean,
      required: false,
      default: false
    },
    expandedStackIds: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      id: this.project.id,
      singularCollapsedStackId: ''
    }
  },
  computed: {
    ...mapGetters([
      'serviceBy',
      'gearspecBy'
    ]),

    isLoading () {
      return typeof this.project.attributes.stack === 'undefined'
    },

    groupedStackItems () {
      /**
       * returns project's services grouped by stack (i.e. indexed by stack_id)
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
    },

    onExpandCollapseStack (stackId, isExpanded) {
      this.$emit('expand-collapse-stack', stackId, isExpanded)
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
