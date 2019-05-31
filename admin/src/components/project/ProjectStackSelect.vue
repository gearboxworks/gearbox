<template>
  <b-form-select
    class="add-stack"
    v-model="selectedService"
    :disabled="!hasStacksNotInProject"
    :required="true"
    @change="addProjectStack"
  >
    <option value="" disabled>{{hasStacksNotInProject ? 'Add stack...' : 'All stacks already added'}}</option>
    <option
      v-for="(stack,stackId) in stacksNotInProject"
      :key="stackId"
      :value="stackId"
    >{{stack.attributes.stackname}}</option>
  </b-form-select>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'ProjectStackSelect',
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
      selectedService: ''
    }
  },
  computed: {
    ...mapGetters({ serviceBy: 'serviceBy', gearspecBy: 'gearspecBy', allGearspecs: 'gearspecs/all', allStacks: 'stacks/all', hasExtraBasedirs: 'hasExtraBasedirs' }),
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    },
    stacksNotInProject () {
      const result = {}

      const projectStack = this.project.attributes.stack
        ? this.groupProjectServicesByStack(this.project.attributes.stack)
        : {}

      for (const idx in this.allStacks) {
        const stack = this.allStacks[idx]
        if (typeof projectStack[stack.id] === 'undefined') {
          result[stack.id] = stack
        }
      }
      return result
    },
    hasStacksNotInProject () {
      return Object.entries(this.stacksNotInProject).length > 0
    },
    servicesInProject () {
      const result = {}
      for (let idx = 0; idx > this.stack.length; idx++) {
        const s = this.serviceBy('id', this.stack[idx].service_id)
        if (s) {
          result[this.stack[idx].service_id] = s
        }
      }
      return result
    },
    gearsInProject () {
      const result = {}
      for (let idx = 0; idx > this.stack.length; idx++) {
        const g = this.gearspecBy('id', this.stack[idx].gearspec_id)
        if (g) {
          result[this.stack[idx].gearspec_id] = g
        }
      }
      return result
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    groupProjectServicesByStack (projectStack) {
      var result = {}
      projectStack.forEach((stackMember, idx) => {
        const gearspec = this.gearspecBy('id', stackMember.gearspec_id)
        const service = this.serviceBy('id', stackMember.service_id)
        if (gearspec && service) {
          // console.log(result, gearspec.attributes.stack_id, gearspec.attributes.role)
          if (typeof result[gearspec.attributes.stack_id] === 'undefined') {
            result[gearspec.attributes.stack_id] = {}
          }
          result[gearspec.attributes.stack_id][gearspec.attributes.role] = service
        }
      })
      // console.log('groupProjectStacks', result)
      return result
    },
    addProjectStack (stackId) {
      this.selectedService = ''
      this.$store.dispatch('addProjectStack', { 'projectId': this.id, stackId })
    }
  }
}
</script>
<style scoped>
</style>
