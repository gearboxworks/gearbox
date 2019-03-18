<template>
  <b-container>
    <b-card-group deck>
      <b-card
        v-for="project in projects"
        :key="project.path"
        :to="{path:'/project/'+project.hostname}"
      >
        <project-details :storedProject="project"></project-details>
        <div slot="footer">
          Stack: <strong>{{project.stack[Object.keys(project.stack)[0]].stack}}</strong> (
          <span v-for="(service, index) in project.stack" :key="index" :title="index.replace(service.stack+'/','')">
            {{service.program}} <small>{{service.version.major+'.'+service.version.minor+'.'+service.version.patch}}</small>
          </span>
          )
          <project-stack :projectHostname="project.hostname" :projectStack="project.stack"></project-stack>
          <!--b-button title="Configure service stack" :to="'/projects/'+project.hostname+'/stack'" variant="primary">Customize</b-button-->
        </div>

      </b-card>
    </b-card-group>
  </b-container>
</template>

<script>
import { mapState } from 'vuex'
import ProjectDetails from '../components/ProjectDetails'
import ProjectStack from '../components/ProjectStack'

export default {
  name: 'ProjectList',
  components: { ProjectDetails, ProjectStack },
  computed: {
    ...mapState([
      'projects',
      'baseDirs'
    ])
  },
  mounted () {
    this.$store.dispatch('loadProjectHeaders')
    this.$store.dispatch('loadBaseDirs')
    this.$store.dispatch('loadGears')
  }
}
</script>
<style scoped>
  .el-icon-caret-right {
    color: green;
  }
</style>
