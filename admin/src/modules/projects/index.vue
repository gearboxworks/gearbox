<template>
  <div class="projects-container">
    <projects-drawer
      visible="false"
      @switch-view-mode="switchViewMode"
    />

    <projects-list
      v-if="projects.length"
      key="projects-content"
      :projects="projects"
      :viewMode="viewMode"
    />
    <div
      v-else
      key="projects-content"
      class="is-empty"
    >
      <h5>No projects match the current criteria.</h5>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import store from './_store'
// import { Getters } from './_store/private-types'
// import ProjectTypes from './_store/types'
import { ProjectGetters, ProjectActions } from './_store/public-types'
import ProjectsDrawer from './_components/ProjectsDrawer'
import ProjectsList from './_components/ProjectsList'

export default {
  name: 'ProjectsModule',
  components: {
    ProjectsDrawer,
    ProjectsList
  },
  data () {
    return {
      viewMode: 'cards'
    }
  },
  computed: {
    ...mapGetters({
      'projects': ProjectGetters.FILTERED_PROJECTS
    })
  },
  created () {
    const STORE_KEY = 'projects'
    // eslint-disable-next-line no-underscore-dangle
    if (!(STORE_KEY in this.$store._modules.root._children)) {
      this.$store.registerModule(STORE_KEY, store)
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

    this.$store.dispatch(ProjectActions.LOAD_ALL_HEADERS).then(() => {
      // console.log('after loadAllHeaders', this.$store.state.projects.records)
      this.$store.dispatch(ProjectActions.LOAD_DETAILS_FOR_ALL).then((result) => {
        // console.log('all details loaded', result)
      })
    })
  },
  methods: {
    switchViewMode ($ev, viewMode) {
      this.viewMode = viewMode
    }
  }
}
</script>

<style scoped>
.is-empty{
  margin-left: 1rem;
  padding-left: 0;
  padding-right: 1rem;
}
</style>
