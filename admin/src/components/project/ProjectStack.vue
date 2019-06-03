<template>
  <b-input-group
    :id="`${projectBase}stack`"
    :class="{'input-group--stack': true, 'is-collapsed': isCollapsed, 'is-modified': isModified, 'is-updating': isUpdating}"
    role="tabpanel"
  >
    <b-form-select
      class="select-stack"
      v-model="selectedStack"
      :disabled="!hasStacksNotInProject || isUpdating"
      :required="true"
      @change="isModified=true"
      v-show="!isCollapsed"
      :ref="`${projectBase}-select`"
      autofocus
    >
      <option value="" disabled>{{hasStacksNotInProject ? 'Add stack...' : 'All stacks already added'}}</option>
      <option
        v-for="(stack,stackId) in stacksNotInProject"
        :key="stackId"
        :value="stackId"
      >{{stack.attributes.stackname}}</option>
    </b-form-select>
    <b-input-group-append>
      <b-button
        variant="outline-info"
        :title="isUpdating ? 'Updating...' : (isCollapsed ? 'Add a stack' : (isModified ? 'Add the selected stack': 'Please select some stack first or Click to cancel'))"
        v-b-tooltip.hover
        :disabled="isUpdating"
        :class="{'btn--submit': true, 'btn--add': isCollapsed}"
        @click.prevent="onButtonClicked"
      >
        <font-awesome-icon
          v-if="isUpdating"
          icon="circle-notch"
          spin
        />
        <font-awesome-icon
          v-else
          :icon="['fa', (isCollapsed ? 'layer-group' : (isModified ? 'check' : 'times'))]"
        />
        <span>{{(isCollapsed && !isUpdating) ? '+' : ''}}</span>
      </b-button>
    </b-input-group-append>
  </b-input-group>

</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'ProjectStack',
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
      selectedStack: '',
      isCollapsed: true,
      isModified: false,
      isUpdating: false
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
    ...mapActions(['addProjectStack']),
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
    maybeAddProjectStack (stackId) {
      if (!stackId) {
        return
      }
      this.isUpdating = true
      this.addProjectStack({ 'projectId': this.id, stackId }).then(() => {
        this.isUpdating = false
        this.isCollapsed = true
        this.selectedStack = ''
        this.isModified = false
      })
    },
    onButtonClicked () {
      if (this.isCollapsed) {
        this.isCollapsed = false
        this.$nextTick(() => {
          this.$refs[`${this.projectBase}-select`].$el.focus()
        })
      } else {
        if (this.isModified) {
          this.maybeAddProjectStack(this.selectedStack)
        } else {
          this.isCollapsed = true
        }
      }
    }
  }
}
</script>
<style scoped>
  .btn--add {
    position:relative;
  }

  .btn--add svg {
    position: relative;
    left: -2px;
    top: 2px;
  }
  .btn--add span {
    position: absolute;
    right: 6px;
    font-size: 17px;
    top: 0px;
  }

  .btn-outline-info {
    border-color: #ced4da;
  }

  .is-collapsed .btn-outline-info {
    border-color: transparent;
    border-top-left-radius: 0.25rem;
    border-bottom-left-radius: 0.25rem;
  }

</style>
