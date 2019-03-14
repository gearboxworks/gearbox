<template>
  <b-card-group deck>
    <b-card
      v-for="project in projects"
      :key="project.path"
      :to="{path:'/project/'+project.hostname}"
    >
      <ProjectDetails :storedProject="project"></ProjectDetails>
      <div slot="footer">
        Stack: <strong>{{project.stack[Object.keys(project.stack)[0]].named_stack}}</strong>
        <ul class="service-list">
          <li v-for="service in project.stack" :key="service.stack_role">
            {{service.program}} {{service.version.major+'.'+service.version.minor+'.'+service.version.patch}} <small class="text-muted">({{service.role}})</small>
          </li>
        </ul>
        <!--b-button :to="'/projects/'+project.hostname" variant="primary">Edit</b-button-->
      </div>

    </b-card>
  </b-card-group>
</template>

<script>
import { mapState } from 'vuex'
import ProjectDetails from '../components/ProjectDetails'

export default {
  name: 'ProjectList',
  components: { ProjectDetails },
  computed: {
    ...mapState([
      'projects',
      'baseDirs'
    ])
  },
  mounted () {
    this.$store.dispatch('loadProjectHeaders')
    this.$store.dispatch('loadBaseDirs')
  }
}
</script>
<style scoped>
  .el-icon-caret-right {
    color: green;
  }
</style>
