<template>
  <div class="projects-container">
    <projects-drawer visible="false" @switch-view-mode="switchViewMode"></projects-drawer>
    <div v-if="projects.length" class="filtered-projects">
      <b-card-group columns class="pl-3 pr-3" v-if="viewMode==='cards'">
        <project-card
          v-for="(project, projectIndex) in projects"
          :key="project.id"
          :project="project"
          :projectIndex="projectIndex"
        >
        </project-card>
      </b-card-group>
      <table class="projects-table" v-else>
        <thead>
          <tr>
            <th class="th--state">State</th><th class="th--hostname">Project Name</th><th class="th--location">Location</th><th class="th--stack">Stack</th><th class="th--notes">Notes</th>
          </tr>
        </thead>
        <tbody>
          <project-row
            v-for="(project, projectIndex) in projects"
            :key="project.id"
            :project="project"
            :projectIndex="projectIndex"
          >
          </project-row>
        </tbody>
      </table>
    </div>
    <div v-else class="filtered-projects is-empty">
      <h5>No projects match the current criteria.</h5>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import ProjectsDrawer from '../components/ProjectsDrawer'
import ProjectCard from '../components/project/card/ProjectCard'
import ProjectRow from '../components/project/row/ProjectRow'

export default {
  name: 'ProjectList',
  data () {
    return {
      viewMode: 'cards'
    }
  },
  components: {
    ProjectsDrawer,
    ProjectCard,
    ProjectRow
  },
  computed: {
    ...mapGetters({
      'projects': 'filteredProjects'
    })
  },
  methods: {
    switchViewMode ($ev, viewMode) {
      this.viewMode = viewMode
    }
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
  @media (min-width: 576px) {
    .card-columns {
      -webkit-column-count: 1;
      -moz-column-count: 1;
      column-count: 1;
      -webkit-column-gap: 1.25rem;
      -moz-column-gap: 1.25rem;
      column-gap: 1.25rem;
      orphans: 1;
      widows: 1;
      text-align: center;
    }
    .card-columns .card {
      max-width: 500px;
      /* margin-left: auto; */
      /* margin-right: auto; */
      text-align: left;
    }
  }
  @media (min-width: 800px) {
    .card-columns {
      -webkit-column-count: 2;
      -moz-column-count: 2;
      column-count: 2;
    }
  }
  @media (min-width: 1024px) {
    .card-columns {
      -webkit-column-count: 3;
      -moz-column-count: 3;
      column-count: 3;
    }
  }
  @media (min-width: 1440px) {
    .card-columns {
      -webkit-column-count: 4;
      -moz-column-count: 4;
      column-count: 4;
    }
  }
  @media (min-width: 1960px) {
    .card-columns {
      -webkit-column-count: 5;
      -moz-column-count: 5;
      column-count: 5;
    }
  }
  @media (min-width: 2400px) {
    .card-columns {
      -webkit-column-count: 6;
      -moz-column-count: 6;
      column-count: 6;
    }
  }

  .filtered-projects{
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .projects-table {
    width: 100%;
  }

  .th--state {
    width: 50px;
  }

  .th--hostname {
    width: 300px;
  }

  .th--location {
    width: 400px;
  }

  .th--notes {
    width: 300px;
  }

  .is-empty{
    margin-left: 1rem;
  }

</style>
