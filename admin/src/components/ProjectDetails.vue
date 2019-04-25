<template>
  <b-collapse :id="`${projectBase}advanced`" role="tabpanel" :visible="true">
    <b-form-group
      :id="`${projectBase}location-group`"
      :label-for="`${projectBase}location-input`"
      label=""
      description="Location"
    >
      <b-form-input
        disabled
        :id="`${projectBase}location-input`"
        :value="resolveDir(basedir, path)"
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
      :container="`${projectBase}advanced`"
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
          :value="this.$store.state.baseDirs[basedir] ? this.$store.state.baseDirs[basedir].text : ''"
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
      <option value="" disabled>{{hasStacksNotInProject ? 'Add Stack...' : 'No more stacks to add'}}</option>
      <option
        v-for="(stack,stackId) in stacksNotInProject"
        :key="stackId"
        :value="stackId"
      >{{stack.attributes.stackname}}</option>
    </b-form-select>

  </b-collapse>
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
    ...mapGetters({ serviceBy: 'serviceBy', gearspecBy: 'gearspecBy', allGearspecs: 'gearspecs/all', allStacks: 'stacks/all' }),
    projectBase () {
      return this.escAttr(this.id) + '-'
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
        const gear = this.gearspecBy('id', stackMember.gearspec_id)
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
    resolveDir (basedir, projectPath) {
      const dir = this.$store.state.baseDirs[basedir]
      const slash = (typeof dir !== 'undefined')
        ? ((dir.text.indexOf('/') !== -1) ? '/' : '\\')
        : '/'
      // console.log(this.$store.state.baseDirs[basedir])
      return (typeof dir !== 'undefined')
        ? (dir.text + slash + projectPath)
        : ''
    },
    showDetails () {
      this.showingDetails = true
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
    addProjectStack (stackName) {
      this.selectedService = ''
      this.$store.dispatch('addProjectStack', { 'projectId': this.id, stackName })
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
</style>
