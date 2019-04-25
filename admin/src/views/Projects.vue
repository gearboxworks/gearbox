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
    this.$store.dispatch('basedirs/loadAll').then(() => {
      // const stacks = this.$store.getters['stacks/all']
    })

    this.$store.dispatch('stacks/loadAll').then(() => {
      // const stacks = this.$store.getters['stacks/all']
    })

    this.$store.dispatch('services/loadAll').then(() => {
      // const services = this.$store.getters['services/all']
    })

    this.$store.dispatch('gearspecs/loadAll').then(() => {
      // const gearspecs = this.$store.getters['gearspecs/all']
    })

    this.$store.dispatch('projects/loadAll').then(() => {
      // const projects = this.$store.getters['projects/all']
    }).then(() => {
      this.$store.dispatch('loadProjectDetails')
    })
  }
}
</script>
<style scoped>
  .el-icon-caret-right {
    color: red;
  }
</style>
