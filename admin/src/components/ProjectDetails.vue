<template>
  <div :class="{'showing-details': showingDetails, 'not-showing-details': !showingDetails}">
    <b-form>
      <b-form-group
        :id="`hostname-group-${projectIndex}`"
        class="hostname-group"
        label=""
        :label-for="`hostname-input-${projectIndex}`"
        :description="showingDetails ? 'Hostname' : ''"
      >
        <b-form-input
          :id="`hostname-input-${projectIndex}`"
          class="hostname-input"
          type="text"
          v-model="attributes.hostname"
          @change="maybeSubmit"
          size="lg"
          v-b-tooltip.hover.bottomright
          :title="showingDetails ? '' : 'Expand details'"
          required
          @click="showDetails"
          placeholder="" />
      </b-form-group>

      <a target="_blank"
         href="#"
         :title="storedProject.attributes.enabled ? 'Stop all services' : 'Run all services'"
         v-b-tooltip.hover
         @click.prevent="onRunStop"
         class="titlebar-icon titlebar-icon--state"
      >
        <font-awesome-icon
          :icon="['fa', storedProject.attributes.enabled ? 'stop-circle': 'play-circle']"
        />
      </a>

      <a target="_blank"
         :href="`http://${attributes.hostname}/`"
         title="Open Frontend"
         v-b-tooltip.hover
         :class="['titlebar-icon', 'titlebar-icon--frontend', {'is-disabled': attributes.enabled}]"
      >
        <font-awesome-icon
          :icon="['fa', 'home']"
        />
      </a>

      <a target="_blank"
         :href="`http://${attributes.hostname}/wp-admin/`"
         title="Open Dashboard"
         v-b-tooltip.hover
         :class="['titlebar-icon', 'titlebar-icon--dashboard', {'is-disabled': attributes.enabled}]"
      >
        <font-awesome-icon
          :icon="['fa', 'tachometer-alt']"
        />
      </a>

      <b-collapse :id="`${projectBase}advanced`" role="tabpanel" :visible="showingDetails">
        <b-form-group
          :id="`${projectBase}location-group`"
          :label-for="`${projectBase}location-input`"
          label=""
          description="Location"
        >
          <b-form-input
            disabled
            :id="`${projectBase}location-input`"
            :value="resolveDir(attributes.baseDir, attributes.path)"
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
              v-model="attributes.baseDir"
              required
              v-if="Object.entries(this.$store.getters.baseDirsAsOptions).length>1"
              :options="this.$store.getters.baseDirsAsOptions"
            />
            <b-form-input
              @change="maybeSubmit"
              required
              disabled
              v-else
              :value="this.$store.state.baseDirs[attributes.baseDir] ? this.$store.state.baseDirs[attributes.baseDir].text : ''"
            />

          </b-form-group>

          <b-form-group
            label="Path:"
            label-for="dirNameInput"
            description=""
          >
            <b-form-input
              type="text"
              v-model="attributes.path"
              required
              placeholder=""
              @input = "attributes.path = sanitizePath(attributes.path)"
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
            v-model="attributes.notes"
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
    </b-form>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'
import filenamify from 'filenamify'

export default {
  name: 'ProjectDetails',
  props: {
    storedProject: {
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
      ...this.storedProject,
      showingDetails: false,
      selectedService: ''
    }
  },
  computed: {
    ...mapGetters(['groupProjectStacks']),
    projectBase () {
      return this.escAttr(this.attributes.hostname) + '-'
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
    resolveDir (baseDir, projectPath) {
      const dir = this.$store.state.baseDirs[baseDir]
      const slash = (typeof dir !== 'undefined')
        ? ((dir.text.indexOf('/') !== -1) ? '/' : '\\')
        : '/'
      // console.log(this.$store.state.baseDirs[baseDir])
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
          'hostname': this.storedProject.id,
          'project': this.$data
        }
      ).then(() => {
        // this.$router.push('/project/' + this.hostname)
      })
    },
    onRunStop () {
      this.$store.dispatch(
        'changeProjectState', { 'projectHostname': this.storedProject.id, 'isEnabled': !this.attributes.enabled }
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
  .form-group label {
    font-size: 0.75rem;
    margin-bottom: 0;
  }

  .hostname-group{
    display: inline-block;
    float: left;
    width: calc(100% - 110px);
  }

  .not-showing-details .hostname-group{
    margin-bottom: 0;
  }

  .hostname-input{
    padding-left: 11px;
    font-weight: bold;
    /*padding: 1px 8px 1px 8px;
    margin-left: -10px;*/
  }

  .not-showing-details .hostname-input {
    border: 1px solid transparent;
    cursor: pointer;
  }

  .not-showing-details .hostname-input:hover {
    text-decoration: underline;
  }

  .showing-details .hostname-input{
    cursor: text;
  }

  .titlebar-icon {
    float: right;
    font-size: 1.25rem;
    cursor: pointer;
    margin-left: 10px;
    color: rgba(42, 85, 130, 0.98);
  }

  .titlebar-icon {
    padding-top: 3px;
  }

  .titlebar-icon.is-disabled  {
    color: rgba(17, 56, 85, 0.42);
  }

  .titlebar-icon--details {
    float: left;
  }

  .titlebar-icon--state {
    font-size: 1.65rem;
    padding-top: 0;
  }

  [data-icon="play-circle"] path {
    fill: green;
  }
  [data-icon="stop-circle"] path {
    fill: red;
  }
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
