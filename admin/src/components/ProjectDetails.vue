<template>
  <div :id="`${projectBase}details`" role="tabpanel">
    <b-form-group
      :id="`${projectBase}location-group`"
      :label-for="`${projectBase}location-input`"
      label=""
      description="Location"
    >
      <b-form-input
        disabled
        :id="`${projectBase}location-input`"
        :value="resolveDir(currentBasedir, path)"
        class="location-input"
      />
      <a target="_blank"
         :id="`${projectBase}change-location`"
         href="#"
         title="Change location"
         :class="['cog-icon']"
         @click.prevent=""
      >
        <font-awesome-icon
          :icon="['fa', 'cog']"
        />
      </a>
    </b-form-group>

    <b-popover
      :target="`${projectBase}change-location`"
      :container="`${projectBase}details`"
      :ref="`${projectBase}location-popover`"
      triggers="focus"
      placement="bottom"
    >
      <template slot="title">
        <b-button @click="onClosePopoverFor(`${projectBase}change-location`)" class="close" aria-label="Close">
          <span class="d-inline-block" aria-hidden="true">&times;</span>
        </b-button>
        Change location
      </template>

      <b-form-group
        label="Base directory:"
        label-for="basedirInput"
        description="Go to Preferences page to add more directories."
      >
        <b-form-select
          @change="maybeSubmit"
          v-model="basedir"
          required
          v-if="Object.entries(this.$store.getters.baseDirsAsOptions).length>1"
          :options="this.$store.getters.baseDirsAsOptions"
        />
        <b-form-input
          @change="maybeSubmit"
          required
          disabled
          v-else
          :value="currentBasedir"
        />

      </b-form-group>

      <b-form-group
        label="Path:"
        label-for="dirNameInput"
        description=""
      >
        <b-form-input
          type="text"
          v-model="path"
          required
          placeholder=""
          @input = "path = sanitizePath(path)"
          @change="maybeSubmit"
        />
      </b-form-group>

    </b-popover>

    <b-form-group
      id="notesGroup"
      label=""
      label-for="notesInput"
      description="(will be visible only here)"
    >
      <b-form-textarea
        id="textarea"
        v-model="notes"
        placeholder="Notes..."
        rows="3"
        max-rows="6"
        @change="maybeSubmit"
      />
    </b-form-group>

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

    <a class="hide-details"
       title="Hide project details"
       @click="$emit('toggle-details')"
    >
      <font-awesome-icon
        :icon="['fa', 'chevron-up']"
      />
      <span>Hide</span>
    </a>

  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import filenamify from 'filenamify'

export default {
  name: 'ProjectDetails',
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
    ...mapGetters({ basedirBy: 'basedirBy', serviceBy: 'serviceBy', gearBy: 'gearBy', allGearspecs: 'gearspecs/all', allStacks: 'stacks/all' }),
    projectBase () {
      return this.escAttr(this.id) + '-'
    },
    currentBasedir () {
      const basedir = this.basedirBy('id', this.basedir)
      return basedir ? basedir.attributes.basedir : ''
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
        const g = this.gearBy('id', this.stack[idx].gearspec_id)
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
        const gear = this.gearBy('id', stackMember.gearspec_id)
        const service = this.serviceBy('id', stackMember.service_id)
        if (gear && service) {
          // console.log(result, gear.attributes.stack_id, gear.attributes.role)
          if (typeof result[gear.attributes.stack_id] === 'undefined') {
            result[gear.attributes.stack_id] = {}
          }
          result[gear.attributes.stack_id][gear.attributes.role] = service
        }
      })
      // console.log('groupProjectStacks', result)
      return result
    },
    resolveDir (dir, path) {
      return dir + ((dir.indexOf('/') !== -1) ? '/' : '\\') + path
    },
    maybeSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          id: this.id,
          attributes: this.$data
        }
      ).then(() => {
        // this.$router.push('/project/' + this.hostname)
      })
    },
    onClosePopoverFor (triggerElementId) {
      this.$root.$emit('bv::hide::popover', triggerElementId)
    },
    addProjectStack (stackId) {
      this.selectedService = ''
      this.$store.dispatch('addProjectStack', { 'projectId': this.id, stackId })
    },
    sanitizePath (path) {
      const sanitized = filenamify(path).trim()
      return sanitized || 'project'
    }
  }
}
</script>
<style scoped>
  .collapse {
    clear: both;
  }
  .location-input{
    width: calc(100% - 38px);
    display: inline-block;
  }
  .cog-icon {
    padding: 6px 10px;
    border-top-right-radius: 3px;
    border-bottom-right-radius: 3px;
    background-color: #e9ecef;
    border: 1px solid #ced4da;
    position: relative;
    /* top: 7px; */
    left: -3px;
    display: inline-block;
  }
  .hide-details {
    display: block;
    margin-top: 15px;
    text-align: center;
    margin-bottom: -5px;
    margin-right: auto;
    padding: 1px 6px;
    margin-left: auto;
    color: #1e69b9 !important;
    opacity: 0;
    cursor: pointer;
    transition: opacity 400ms;
  }
  .hide-details span {
    margin-right: 5px;
    margin-left: 5px;
    font-weight: bold;
    font-size: 14px;
  }
  .card--project:hover .hide-details{
    opacity:0.75;
  }
  .card--project:hover .hide-details:hover {
    opacity: 1;
  }

</style>
