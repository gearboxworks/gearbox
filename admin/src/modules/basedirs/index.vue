<template>
  <div>
    <header>
      <h2>{{labels.plural}}</h2>
    </header>
    <BasedirsList
      :records="records"
      :labels="labels"
    />
  </div>
</template>

<script>
import moduleConfig from './config'
import BasedirsList from './_components/BasedirsList'

export default {
  name: moduleConfig.moduleName,
  components: {
    BasedirsList
  },
  computed: {
    records () {
      return this.$store.state[moduleConfig.namespace].records
    },
    labels () {
      return moduleConfig.labels
    }
  }
}
</script>

<style scoped>
  header {
    margin: 2rem 1rem 1rem 1rem;
  }
</style>
