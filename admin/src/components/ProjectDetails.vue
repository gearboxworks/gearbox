<template>
  <div>
    <h1>Project Details</h1>
    <el-form
      v-if="project"
      label-width="120px"
    >
      <h2>{{ project.name }}</h2>

      <el-form-item label="Project Name">
        <i class="el-icon-info" />
        <el-input
          placeholder="Please input"
          v-model="name"
        />
      </el-form-item>

      <el-form-item label="Hostname">
        <i class="el-icon-info" />
        <el-input
          placeholder="Please input"
          v-model="hostname"
        />
      </el-form-item>

      <el-form-item label="Root Dir">
        <i class="el-icon-info" />
        <el-input
          placeholder="Please input"
          v-model="group"
        />
      </el-form-item>

      <el-form-item label="Enabled">
        <i class="el-icon-info" />
        <el-switch
          v-model="enabled"
          active-color="#13ce66"
          inactive-color="#ff4949"
        />
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          @click="onSubmit"
        >
          Submit
        </el-button>
        <el-button disabled>
          Reset
        </el-button>
      </el-form-item>
    </el-form>

    <div
      v-else
      class="project-details"
    >
      <h2>{{ this.$route.params.projectName }}</h2>
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
      name: '',
      hostname: '',
      group: 0,
      enabled: null,
      groups: [
        {
          value: '0',
          label: 'Group 0'
        },
        {
          value: '1',
          label: 'Group 1'
        },
        {
          value: '2',
          label: 'Group 2'
        }
      ]
    }
  },
  watch: {
    '$route.params.projectName': {
      handler: function (projectName) {
        const p = this.projectByName(projectName)
        this.name = p.name
        this.hostname = p.hostname
        this.group = p.group
        this.enabled = p.enabled
      },
      deep: true,
      immediate: true
    }
  },
  computed: {
    ...mapGetters([
      'projectByName'
    ]),
    project () {
      return this.projectByName(this.$route.params.projectName)
    }
  },
  methods: {
    updateProjectHostname (value) {
      this.hostname = value
    },
    onSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          'projectName': this.project.name,
          'project': {
            'name': this.name,
            'hostname': this.hostname,
            'group': this.group,
            'enabled': this.enabled
          }
        }
      ).then(() => {
        this.$router.push('/project/' + this.name)
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
