<template>
  <b-collapse :id="`${projectBase}advanced`" role="tabpanel" visible="true">
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
      :disabled="!hasUnusedStacks"
      :required="true"
      @change="addProjectStack"
    >
      <option value="" disabled>{{hasUnusedStacks ? 'Add Stack...' : 'No more stacks to add'}}</option>
      <option
        v-for="(stack,stackName) in stacksNotUnusedInProject"
        :key="stackName"
        :value="stackName"
      >{{stackName.replace('gearbox.works/', '')}}</option>
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
    ...mapGetters(['groupProjectStacks']),
    projectBase () {
      return this.escAttr(this.id) + '-'
    },
    hasUnusedStacks () {
      return Object.entries(this.stacksNotUnusedInProject).length > 0
    },
    stacksNotUnusedInProject () {
      const result = {}
      const projectStacks = this.groupProjectStacks(this.stack)
      for (const index in this.$store.state.gearStacks) {
        const stackName = this.$store.state.gearStacks[index]
        if (typeof projectStacks[stackName] === 'undefined') {
          result[stackName] = this.$store.state.gearStacks[stackName]
        }
      }
      return result
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
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
          'hostname': this.id,
          'project': {
            id: this.id,
            attributes: this.$data
          }
        }
      ).then(() => {
        // this.$router.push('/project/' + this.hostname)
      })
    },
    onRunStop () {
      this.$store.dispatch(
        'changeProjectState', { 'projectHostname': this.id, 'isEnabled': !this.enabled }
      )
    },
    onClosePopoverFor (triggerElementId) {
      this.$root.$emit('bv::hide::popover', triggerElementId)
    },
    addProjectStack (stackName) {
      this.selectedService = ''
      this.$store.dispatch('addProjectStack', { 'projectHostname': this.id, stackName })
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
