<template>
  <b-container>
    <b-card-group deck>
      <project-card
        v-for="(project, projectIndex) in projects"
        :key="project.id"
        :project="project"
        :projectIndex="projectIndex"
      >
      </project-card>
    </b-card-group>
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import ProjectCard from '../components/ProjectCard'

export default {
  name: 'ProjectList',
  components: {
    ProjectCard
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
