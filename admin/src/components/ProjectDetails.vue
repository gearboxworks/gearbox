<template>
  <div :class="{'showing-details': showingDetails, 'not-showing-details': !showingDetails}">
    <b-form>
      <b-form-group
        :id="'hostname-group-'+projectIndex"
        class="hostname-group"
        label=""
        :label-for="'hostname-input-'+projectIndex"
        :description="showingDetails ? 'Hostname' : ''"
      >
        <b-form-input
          :id="'hostname-input-'+projectIndex"
          class="hostname-input"
          type="text"
          v-model="hostname"
          @change="maybeSubmit"
          size="lg"
          v-b-tooltip.hover.bottomright
          title="Expand details"
          required
          @click="showDetails"
          placeholder="" />
      </b-form-group>

      <a target="_blank"
         href="#"
         :title="storedProject.enabled ? 'Stop all services' : 'Run all services'"
         v-b-tooltip.hover
         @click.prevent="onRunStop"
         class="titlebar-icon titlebar-icon--state"
      >
        <font-awesome-icon
          :icon="['fa', storedProject.enabled ? 'stop-circle': 'play-circle']"
        />
      </a>

      <a target="_blank"
         :href="'http://'+hostname+'/'"
         title="Open Frontend"
         v-b-tooltip.hover
         :class="['titlebar-icon', 'titlebar-icon--frontend', {'is-disabled': enabled}]"
      >
        <font-awesome-icon
          :icon="['fa', 'home']"
        />
      </a>

      <a target="_blank"
         :href="'http://'+hostname+'/wp-admin'"
         title="Open Dashboard"
         v-b-tooltip.hover
         :class="['titlebar-icon', 'titlebar-icon--dashboard', {'is-disabled': enabled}]"
      >
        <font-awesome-icon
          :icon="['fa', 'tachometer-alt']"
        />
      </a>

      <!--a target="_blank"
         href="#"
         title="Settings"
         v-b-tooltip.hover
         @click.prevent = ""
         class="titlebar-icon titlebar-icon--details"
         v-b-toggle="project_base + '_advanced'"
      >
        <font-awesome-icon
          :icon="['fa', 'ellipsis-h']"
        />
      </a-->

      <b-collapse :id="project_base + '_advanced'" role="tabpanel" :visible="showingDetails">

        <b-form-group
          :id="project_base + 'location-group'"
          :label-for="project_base+'location-input'"
          label=""
          description="Location"
        >
          <b-form-input
            disabled
            :id="project_base+'location-input'"
            :value="resolvePath(baseDir, path)"
            class="location-input"
          />
          <a target="_blank"
             :id="`${project_base}change-location`"
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
          :target="`${project_base}change-location`"
          :container="project_base + '-advanced'"
          :ref="project_base + 'location-popover'"
          triggers="focus"
          placement="bottom"
        >
          <template slot="title">
            <b-button @click="onClosePopoverFor(project_base + 'location-group')" class="close" aria-label="Close">
              <span class="d-inline-block" aria-hidden="true">&times;</span>
            </b-button>
            Change location
          </template>

          <b-form-group
            id="basedirGroup1"
            label="Base directory"
            label-for="basedirInput"
            description=""
          >
            <b-form-select
              @change="maybeSubmit"
              v-model="baseDir"
              required
              :options="this.$store.getters.baseDirsAsOptions"
            />
          </b-form-group>

          <b-form-group
            id="pathGroup"
            label="Path:"
            label-for="dirNameInput"
            description=""
          >
            <b-form-input
              id="dirNameInput"
              type="text"
              v-model="path"
              @change="maybeSubmit"
              required
              placeholder="" />
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
            @change="maybeSubmit"
            placeholder="Notes..."
            rows="3"
            max-rows="6"
          />
        </b-form-group>
      </b-collapse>
    </b-form>
  </div>
</template>

<script>
// import Vue from 'vue'

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
      showingDetails: false
    }
  },
  computed: {
    project_base () {
      return this.escAttr(this.hostname) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    resolvePath (baseDir, path) {
      return typeof this.$store.state.baseDirs[baseDir] !== 'undefined'
        ? (this.$store.state.baseDirs[baseDir].text + '/' + path)
        : ''
    },
    showDetails () {
      this.showingDetails = true
    },
    maybeSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          'hostname': this.storedProject.hostname,
          'project': this.$data
        }
      ).then(() => {
        // this.$router.push('/project/' + this.hostname)
      })
    },
    onRunStop () {
      console.log(this.enabled)
      this.$store.dispatch(
        'changeProjectState', { 'projectHostname': this.storedProject.hostname, 'isEnabled': !this.enabled }
      )
    },
    onClosePopoverFor (triggerElementId) {
      console.log('onClosePopoverFor', triggerElementId)
      this.$root.$emit('bv::popover::hide', triggerElementId)
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
    cursor: row-resize;
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
