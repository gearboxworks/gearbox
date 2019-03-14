<template>
  <div>
    <b-form>

      <b-form-group
        id="hostnameGroup"
        label="Hostname:"
        label-for="hostnameInput"
        description=""
        label-cols-sm="4"
        label-cols-lg="3"
      >
        <b-form-input
          id="hostnameInput"
          type="text"
          v-model="hostname"
          @change="maybeSubmit"
          size="lg"
          required
          placeholder="" />
      </b-form-group>

      <b-form-group
        id="basedirGroup1"
        label="Base Directory:"
        label-for="basedirInput"
        description=""
        label-cols-sm="4"
        label-cols-lg="3"
      >
        <b-form-select
          @change="maybeSubmit"
          v-model="baseDir"
          required
          :options="this.$store.getters.baseDirsAsOptions"
          class="mt-3" />
      </b-form-group>

      <b-form-group
        id="pathGroup"
        label="Subdirectory:"
        label-for="dirNameInput"
        :description="'Full path: ' + this.$store.state.baseDirs[baseDir].text + '/' + path"
        label-cols-sm="4"
        label-cols-lg="3"
      >
        <b-form-input
          id="dirNameInput"
          type="text"
          v-model="path"
          @change="maybeSubmit"
          required
          placeholder="" />
      </b-form-group>

      <b-form-group
        id="enabledGroup"
        label="Status:"
        label-for="enabledInput"
        description=""
        label-cols-sm="4"
        label-cols-lg="3"
      >
        <b-form-radio v-model="enabled" @change="maybeSubmit" value="true" name="enabledInput">Enabled</b-form-radio>
        <b-form-radio v-model="enabled" @change="maybeSubmit" value="false" name="enabledInput">Disabled</b-form-radio>
      </b-form-group>

      <b-form-group
        id="notesGroup"
        label="Notes:"
        label-for="notesInput"
        description=""
        label-cols-sm="4"
        label-cols-lg="3"
      >
        <b-form-textarea
          id="textarea"
          v-model="notes"
          @change="maybeSubmit"
          placeholder="Enter something..."
          rows="3"
          max-rows="6"
        />
      </b-form-group>
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
    }
  },
  data () {
    return {
      ...this.storedProject
    }
  },
  methods: {
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
    }
  }
}
</script>

<style scoped>

</style>
