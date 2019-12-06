<template>
  <div class="projects-container">
    <projects-drawer
      visible="false"
      @switch-view-mode="switchViewMode"
    />

    <projects-list
      key="projects-content"
      :records="records"
      :labels="labels"
      :viewMode="viewMode"
    />

  </div>
</template>

<script>

// import storeConfig from './_store'
import moduleConfig from './config'

import ProjectsDrawer from './_components/ProjectsDrawer'
import ProjectsList from './_components/ProjectsList'

import { ProjectGetters } from './_store/method-names'

export default {
  name: moduleConfig.moduleName,
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
    records () {
      return this.$store.getters[ProjectGetters.LIST_FILTERED]()
    },
    labels () {
      return moduleConfig.labels
    }
  },
  methods: {
    switchViewMode (viewMode) {
      this.viewMode = viewMode
    }
  }
}
</script>

<style scoped>
</style>
