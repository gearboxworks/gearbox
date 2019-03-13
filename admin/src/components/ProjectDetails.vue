<template>
  <div>
    <h1>Project Details</h1>
    <b-form v-if="project">
      <h2>{{ project.hostname }}</h2>

      <b-form-group
        id="basedirGroup1"
        label="Base Dir:"
        label-for="basedirInput"
        description="Base dir"
      >
        <b-form-input
          id="basedirInput"
          type="text"
          v-model="baseDir"
          required
          placeholder="Enter base dir" />
      </b-form-group>

      <b-form-group
        id="pathGroup"
        label="Dir Name:"
        label-for="dirNameInput"
        description=""
      >
        <b-form-input
          id="dirNameInput"
          type="text"
          v-model="path"
          required
          placeholder="" />
      </b-form-group>

      <b-form-group
        id="fullPathGroup"
        label="Full path:"
        label-for="fullPathInput"
        description=""
      >
        <b-form-input
          id="fullPathInput"
          type="text"
          v-model="fullPath"
          required
          placeholder="" />
      </b-form-group>

      <b-form-group
        id="hostnameGroup"
        label="Hostname:"
        label-for="hostnameInput"
        description=""
      >
        <b-form-input
          id="hostnameInput"
          type="text"
          v-model="hostname"
          required
          placeholder="" />
      </b-form-group>

      <b-form-group
        id="notesGroup"
        label="Notes:"
        label-for="notesInput"
        description=""
      >
        <b-form-textarea
          id="textarea"
          v-model="notes"
          placeholder="Enter something..."
          rows="3"
          max-rows="6"
        />
      </b-form-group>

      <b-form-group
        id="enabledGroup"
        label="Status:"
        label-for="enabledInput"
        description=""
      >
        <b-form-radio value="true" v-model="enabled" name="enabledInput">Enabled</b-form-radio>
        <b-form-radio value="false" v-model="enabled" name="enabledInput">Disabled</b-form-radio>
      </b-form-group>
    </b-form>

    <div
      v-else
      class="project-details"
    >
      <h2>{{ this.$route.params.hostname }}</h2>
      <p>This is a dummy project with no actual data!</p>
    </div>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'

export default {
  name: 'ProjectDetails',
  data () {
    return {
      hostname: '',
      notes: '',
      baseDir: '',
      path: '',
      fullPath: '',
      enabled: null,
      stack: {}
    }
  },
  watch: {
    '$route.params.hostname': {
      handler: function (hostname) {
        const p = this.projectBy('hostname', hostname)
        if (p) {
          // console.log(projectName, p.baseDir)
          this.hostname = p.hostname
          this.notes = p.notes
          this.baseDir = p.baseDir
          this.path = p.path
          this.fullPath = p.fullPath
          this.enabled = p.enabled
          this.stack = p.stack
        }
      },
      deep: true,
      immediate: true
    }
  },
  computed: {
    ...mapGetters([
      'projectBy',
      'projectByName'
    ]),
    project () {
      return this.projectBy('hostname', this.$route.params.hostname)
    }
  },
  methods: {
    onSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          'hostname': this.project.hostname,
          'project': {
            'hostname': this.hostname,
            'notes': this.notes,
            'baseDir': this.baseDir,
            'path': this.path,
            'enabled': this.enabled,
            'fullPath': this.fullPath
          }
        }
      ).then(() => {
        this.$router.push('/project/' + this.hostname)
      })
    }
  }
}
</script>

<style scoped>
super {
  color: #ffffff;
  background-color: #1a81ef;
  border-radius: 50%;
  height: 1rem;
  width: 12px;
  padding: 0 0 0 4px;
  line-height: 16px;
  display: inline-block;
  margin-right: 5px;
}
</style>
