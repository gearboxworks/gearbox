<template>
  <div>
    <h1>Preferences</h1>
    <BasedirsList/>
  </div>
</template>

<script>

import BasedirsList from '../modules/basedirs/_components/BasedirsList'

export default {
  name: 'Preferences',
  components: {
    BasedirsList
  },
  data () {
    return {}
  },
  computed: {},
  mounted () {
  }
}</script>

<style scoped>
</style>
