<template>
  <b-card-group deck>
    <b-card
      v-for="project in projects"
      :title="project.path"
      :sub-title="project.hostname"
      :key="project.path"
      :to="{path:'/project/'+project.hostname}"
    >
      <b-card-text>
        <dl>
          <dt>Enabled:</dt><dd>{{project.enabled}}</dd>
          <dt>Directory:</dt><dd>{{project.fullPath}}</dd>
        </dl>
      </b-card-text>
      <div slot="footer">
        Stack: <strong>{{project.stack[Object.keys(project.stack)[0]].named_stack}}</strong>
        <ul class="service-list">
          <li v-for="service in project.stack" :key="service.stack_role">
            {{service.program}} {{service.version.major+'.'+service.version.minor+'.'+service.version.patch}} <small class="text-muted">({{service.role}})</small>
          </li>
        </ul>
        <b-button :to="'/projects/'+project.hostname" variant="primary">Edit</b-button>
      </div>
    </b-card>
  </b-card-group>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'ProjectList',
  computed: {
    ...mapState([
      'projects'
    ])
  },
  mounted () {
    this.$store.dispatch('loadProjectHeaders')
  }
}
</script>
<style scoped>
  .el-icon-caret-right {
    color: green;
  }
</style>
