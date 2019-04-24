<template>
  <div role="tablist" class="project-stack-list" :id="`${projectBase}stack`">
    <div
      v-for="(services, stackId, stackIndex) in groupProjectServicesByStack(stack)"
      :key="stackId"
      class="project-stack"
    >
      <h2 class="stack-title">{{stackId.replace('gearbox.works/', '')}}</h2>
      <b-button :tabindex="projectIndex*100+stackIndex*10" @click.prevent="removeProjectStack(stackId)" class="js-remove-stack" size="sm" variant="outline-secondary" aria-label="Remove this stack from project" title="Remove this stack from project">&times;</b-button>
      <ul class="service-list">
        <li
            v-for="(service, role, serviceIndex) in services"
            :key="id + service.id"
            class="service-item"
        >
          <project-service :projectId="project.id" :service="service" :projectIndex="projectIndex" :stackIndex="stackIndex" :serviceIndex="serviceIndex"></project-service>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>

import ProjectService from './ProjectService.vue'
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
      ...this.project.attributes
    }
  },
  components: {
    ProjectService
  },
  computed: {
    ...mapGetters(['serviceBy', 'gearspecBy']),
    projectBase () {
      return this.escAttr(this.id) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    groupProjectServicesByStack (projectStack) {
      var result = {}
      projectStack.forEach((stackMember, idx) => {
        const gear = this.gearspecBy('id', stackMember.gearspec_id)
        const service = this.serviceBy('id', stackMember.service_id)
        if (gear && service) {
          // console.log(gear.id, service.id)
          if (typeof result[gear.attributes.stack_id] === 'undefined') {
            result[gear.attributes.stack_id] = {}
          }
          result[gear.attributes.stack_id][this.escAttr(gear.attributes.role)] = service
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
    transition: opacity 400ms;
  }
  .card--project:hover .js-remove-stack {
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
  .add-stack {
/*    color: #fff;
    background-color: #17a2b8;
    border-color: #17a2b8;
*/
  }
/*
  .add-stack:hover {
    color: #fff;
    background-color: #138496;
    border-color: #117a8b;
  }
*/
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
