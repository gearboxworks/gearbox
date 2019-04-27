<template>
  <div role="tablist" class="project-stack-list" :id="`${projectBase}stack`">
    <div
      v-for="(stackItems, stackId, stackIndex) in groupedStackItems(projectStackItems)"
      :key="stackId"
      class="project-stack"
    >
      <h2 class="stack-title">{{stackId.replace('gearbox.works/', '')}}</h2>
      <b-button :tabindex="projectIndex*100+stackIndex*10" @click.prevent="removeProjectStack(stackId)" class="js-remove-stack" size="sm" variant="outline-secondary" aria-label="Remove these services" title="Remove these services">&times;</b-button>
      <ul class="service-list">
        <li
            v-for="(item, itemIndex) in stackItems"
            :key="id + item.gear.attributes.role"
            class="service-item"
        >
          <project-gear :projectId="project.id" :stackItem="item" :projectIndex="projectIndex" :stackIndex="stackIndex" :itemIndex="itemIndex"></project-gear>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>

import ProjectGear from './ProjectGear.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'ProjectStack',
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
    ProjectGear
  },
  computed: {
    ...mapGetters(['serviceBy', 'gearBy']),
    projectBase () {
      return this.escAttr(this.id) + '-'
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
    },
    removeProjectStack (stackId) {
      this.$store.dispatch('removeProjectStack', { 'projectId': this.id, stackId })
    }
  }
}
</script>

<style scoped>
  .js-remove-stack {
    float: right;
    opacity: 0;
    transition: all 400ms;
    padding: 2px 8px;
    margin-right: 1px;
  }
  .project-stack:hover .js-remove-stack {
    opacity: 1;
  }
  .project-stack{
    margin-bottom: 8px;
  }
  .project-stack:not(:first-child){
    border-top: 1px solid #e0e0e0;
    padding-top: 8px;
  }
  .stack-title {
    text-transform: uppercase;
    font-size: 1.25rem;
    display: inline-block;
    color: #989898;
  }
  .service-list{
    margin-bottom: 0;
    margin-left: auto;
    margin-right: auto;
    list-style: none;
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
