<template>
  <b-container>
    <b-card-group deck>
      <b-card
        v-for="(project, projectIndex) in projects"
        :key="project.path"
        :to="{path:'/project/'+project.hostname}"
        class="card--project"
      >
        <project-details :storedProject="project" :projectIndex="projectIndex"></project-details>
        <div slot="footer" v-if="Object.entries(project.stack).length">
          <project-stack :projectHostname="project.hostname" :projectStack="project.stack" :projectIndex="projectIndex"></project-stack>
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
  components: {
    ProjectDetails,
    ProjectStack
  },
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
