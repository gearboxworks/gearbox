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
        v-for="(stackItems, stackId) in projectStacks"
        :key="stackId"
        :stackId="stackId"
        :stackItems="stackItems"
        :is-expanded="isExpanded(stackId)"
        @expand-collapse="onExpandCollapseStack"
      />
    </div>
  </div>
</template>

<script>

import StackCard from '../../../stacks/_components/StackCard.vue'
import { ProjectGetters } from '../../_store/method-names'

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
    return {}
  },
  computed: {

    projectStacks () {
      const stacks = this.$store.getters[ProjectGetters.GEARS_GROUPED_BY_STACK](this.project)
      return stacks
    },

    countProjectStacks () {
      return Object.entries(this.projectStacks).length
    },

    isLoading () {
      return typeof this.project.attributes.stack === 'undefined'
    }
  },
  methods: {
    isExpanded (stackId) {
      const expanded = this.expandedStackIds
      return (expanded[stackId] && expanded[stackId] > 0) ||
        ((!expanded[stackId] || expanded[stackId] === 0) && this.countProjectStacks === 1)
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
