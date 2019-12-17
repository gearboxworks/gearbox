<template>
  <div>
    <header>
      <h2>{{labels.title}}</h2>
    </header>
    <SimpleList :labels="labels" :records="records"/>
  </div>
</template>

<script>
import moduleConfig from './config'
import SimpleList from '../_shared/_components/SimpleList'

export default {
  name: moduleConfig.moduleName,
  components: {
    SimpleList
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
