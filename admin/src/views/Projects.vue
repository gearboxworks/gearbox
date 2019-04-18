<template>
  <b-container>
    <b-card-group deck>
      <b-card
        v-for="(project, projectIndex) in projects"
        :key="project.attributes.path"
        :to="{path:'/project/'+project.id}"
        class="card--project"
      >
        <project-details :storedProject="project" :projectIndex="projectIndex"></project-details>
        <div slot="footer" v-if="project.stack">
          <project-stack :projectHostname="project.id" :projectStack="project.stack" :projectIndex="projectIndex"></project-stack>
        </div>
      </b-card>
    </b-card-group>
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import ProjectDetails from '../components/ProjectDetails'
import ProjectStack from '../components/ProjectStack'

export default {
  name: 'ProjectList',
  components: {
    ProjectDetails,
    ProjectStack
  },
  computed: {
    ...mapGetters({
      'projects': 'projects/all'
    })
  },
  mounted () {
    this.$store.dispatch('projects/loadAll').then(() => {
      const projects = this.$store.getters['projects/all']
      console.log(projects[0].id)
    })
    this.$store.dispatch('stacks/loadAll').then(() => {
      const stacks = this.$store.getters['stacks/all']
      console.log(stacks[0].id)
    })
    // this.$store.dispatch('loadProjectHeaders')
    // this.$store.dispatch('loadBaseDirs')
    // this.$store.dispatch('loadGears')
  }
}
</script>
<style scoped>
  .el-icon-caret-right {
    color: red;
  }
</style>
