<template>
  <div role="tablist" class="project-stack-list" :id="`${projectBase}stack`">
    <div
      v-for="(stackItems, stackId, stackIndex) in groupedStackItems(projectStackItems)"
      :key="stackId"
      class="project-stack"
    >
      <h2 class="stack-title">{{stackId.replace('gearbox.works/', '')}}</h2>

      <stack-toolbar :project="project" :projectIndex="projectIndex" :stackId="stackId"></stack-toolbar>

      <ul class="service-list">
        <li
            v-for="(item, itemIndex) in stackItems"
            :key="id + item.gear.attributes.role"
            class="service-item"
        >
          <stack-gear :projectId="project.id" :stackItem="item" :projectIndex="projectIndex" :stackIndex="stackIndex" :itemIndex="itemIndex"></stack-gear>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>

import StackToolbar from './StackToolbar.vue'
import StackGear from './StackGear.vue'
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
    }
  },
  data () {
    return {
      id: this.project.id,
      projectStackItems: this.project.attributes.stack
    }
  },
  components: {
    StackToolbar,
    StackGear
  },
  computed: {
    ...mapGetters(['serviceBy', 'gearBy']),
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
        const gear = this.gearBy('id', stackItem.gearspec_id)
        if (gear) {
          if (typeof result[gear.attributes.stack_id] === 'undefined') {
            result[gear.attributes.stack_id] = []
          }
          const service = stackItem.service_id ? this.serviceBy('id', stackItem.service_id) : null
          result[gear.attributes.stack_id].push({ gear, service })
        }
      })
      // console.log('groupProjectStacks', result)
      return result
    }
  }
}
</script>

<style scoped>
  .project-stack-list {
    margin-top: 10px;
  }
  .project-stack:not(:first-child){
    border-top: 1px solid #e0e0e0;
    padding-top: 20px;
  }
  .stack-title {
    text-transform: uppercase;
    font-size: 1.25rem;
    display: inline-block;
    color: #989898;
    margin-left: 5px;
  }
  .service-list{
    margin-bottom: 15px;
    margin-top: 0px;
    list-style: none;
    padding-inline-start: 0;
    clear: both;
  }
  .service-item {
    display: inline-block;
    list-style: none;
    text-align: center;
    max-width: 110px;
    padding: 5px;
    margin: 5px;
    cursor: pointer;
    border: 1px solid transparent;
    border-radius: 4px;
    transition: all 400ms;
  }
  .service-item:hover {
    border: 1px solid #aaa;
    background-color: #eee;
  }

</style>
