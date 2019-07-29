<template>
  <b-input-group
    :id="`${projectPrefix}stack`"
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
      :ref="`${projectPrefix}-select`"
      autofocus
    >
      <option value="" disabled>
        {{hasStacksNotInProject ? 'Add stack...' : 'All stacks already added'}}
      </option>
      <option
        v-for="(item,stackId) in stacksNotInProject"
        :key="stackId"
        :value="stackId"
      >
        {{item.stack.attributes.stackname + (item.isRemoved ? ' (removed)': '') + (item.isDefault? ' (default)': '')}}
      </option>
    </b-form-select>
    <b-input-group-append>
      <b-button
        variant="outline-info"
        :title="isUpdating ? 'Updating...' : (isCollapsed ? 'Add a stack' : (isModified ? 'Add the selected stack': 'Please select some stack first or Click to cancel'))"
        v-b-tooltip.hover
        :disabled="isUpdating"
        :class="{'btn--submit': true, 'btn--add': isCollapsed}"
        @click.prevent="onAddProjectStack"
      >
        <font-awesome-icon
          v-if="isUpdating"
          key="status-icon"
          icon="circle-notch"
          spin
        />
        <font-awesome-icon
          v-else
          key="status-icon"
          :icon="['fa', (isCollapsed ? 'layer-group' : (isModified ? 'check' : 'times'))]"
        />
        <span>{{(isCollapsed && !isUpdating) ? '+' : ''}}</span>
      </b-button>
    </b-input-group-append>
  </b-input-group>

</template>

<script>
import { mapGetters } from 'vuex'
import { ProjectActions } from '../../_store/public-types'

export default {
  name: 'ProjectStackAdd',
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {},
  data () {
    return {
      id: this.project.id,
      // ...this.project.attributes,
      selectedStack: '',
      isCollapsed: true,
      isModified: false,
      isUpdating: false
    }
  },
  computed: {
    ...mapGetters({
      serviceBy: 'serviceBy',
      gearspecBy: 'gearspecBy',
      stackBy: 'stackBy',
      allGearspecs: 'gearspecs/all',
      allStacks: 'stacks/all',
      hasExtraBasedirs: 'hasExtraBasedirs',
      stackDefaultServiceByRole: 'stackDefaultServiceByRole',
      stackServicesByRole: 'stackServicesByRole',
      preselectServiceId: 'preselectServiceId'
    }),

    projectGearsGroupedByStack () {
      var result = {}
      if (this.project.attributes.stack) {
        this.project.attributes.stack.forEach((stackMember, idx) => {
          // if (stackMember.isRemoved) {
          //   return
          // }
          const gearspec = this.gearspecBy('id', stackMember.gearspec_id)

          const stack = this.stackBy('id', gearspec.attributes.stack_id)
          const serviceId = this.preselectClosestGearServiceId(stack, gearspec.id, stackMember.service_id)
          const service = serviceId ? this.serviceBy('id', serviceId) : null

          if (gearspec && service) {
            // console.log(result, gearspec.attributes.stack_id, gearspec.attributes.role)
            if (typeof result[gearspec.attributes.stack_id] === 'undefined') {
              result[gearspec.attributes.stack_id] = {
                isRemoved: stackMember.isRemoved || false
              }
            }
            result[gearspec.attributes.stack_id][gearspec.attributes.role] = service
          }
        })
      }
      return result
    },

    stacksNotInProject () {
      const result = {}

      const projectStack = this.projectGearsGroupedByStack

      for (const idx in this.allStacks) {
        const stack = this.allStacks[idx]
        if (typeof projectStack[stack.id] === 'undefined') {
          result[stack.id] = { stack, isRemoved: false }
        } else if (projectStack[stack.id].isRemoved) {
          // TODO if the services in the removed stack are exactly the same as in the default version of it, show only one option!
          result[stack.id] = { stack, isRemoved: false, isDefault: true }
          result[stack.id + '(removed)'] = { stack, isRemoved: true }
        }
      }

      return result
    },

    hasStacksNotInProject () {
      return Object.entries(this.stacksNotInProject).length > 0
    }
    //
    // servicesInProject () {
    //   const result = {}
    //   for (let idx = 0; idx > this.stack.length; idx++) {
    //     if (this.stack[idx].isRemoved) {
    //       continue
    //     }
    //     const s = this.serviceBy('id', this.stack[idx].service_id)
    //     if (s) {
    //       result[this.stack[idx].service_id] = s
    //     }
    //   }
    //   return result
    // },
    // gearsInProject () {
    //   const result = {}
    //   for (let idx = 0; idx > this.stack.length; idx++) {
    //     const g = this.gearspecBy('id', this.stack[idx].gearspec_id)
    //     if (g) {
    //       result[this.stack[idx].gearspec_id] = g
    //     }
    //   }
    //   return result
    // }
  },
  methods: {

    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },

    preselectClosestGearServiceId (stack, gearspecId, requestedServiceId) {
      const defaultService = this.stackDefaultServiceByRole(stack, gearspecId)
      /**
       * As an example, for php:7.1.18 it will select php:7.1 or php:7 if exact match is not possible
       */
      return this.preselectServiceId(
        this.stackServicesByRole(stack, gearspecId),
        defaultService,
        requestedServiceId
      )
    },

    async maybeAddProjectStack (stackId) {
      if (!stackId) {
        return
      }

      try {
        this.isUpdating = true
        await this.$store.dispatch(
          ProjectActions.ADD_STACK,
          {
            project: this.project, stackId
          }
        )

        this.isUpdating = false
        this.isCollapsed = true
        this.selectedStack = ''
        this.isModified = false

        this.$emit('maybe-hide-alert', 'Please add some stacks first!')
        this.$emit('added-stack', stackId)
      } catch (e) {
        console.error(e.message)
      }
    },

    onAddProjectStack () {
      if (this.isCollapsed) {
        this.isCollapsed = false
        this.$nextTick(() => {
          this.$refs[`${this.projectPrefix}-select`].$el.focus()
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
