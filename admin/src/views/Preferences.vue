<template>
  <div>
    <h1>Preferences</h1>
    <b-form v-if="baseDirs">
      <h2>Base Dirs</h2>

      <b-row class="my-1">
        <b-col sm="4" md="3" lg="2">
          <label for="input-none">Name:</label>
        </b-col>
        <b-col sm="8" md="3" lg="4">
          <label for="input-none">Path:</label>
        </b-col>
      </b-row>

      <b-row
        v-for="(baseDir, index) in baseDirs"
        :key="baseDir.name"
        class="my-1">
        <b-col sm="4" md="3" lg="2">
          <b-form-input
            :id="'basedirInput'+index"
            type="text"
            v-model="baseDir.name"
            required
            placeholder="Name" />
        </b-col>

        <b-col sm="8" md="3" lg="4">
          <b-form-input
            :id="'basedirInputPath'+index"
            type="text"
            v-model="baseDir.path"
            required
            placeholder="Path" />
        </b-col>
      </b-row>
    </b-form>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Preferences',
  computed: {
    ...mapState([
      'baseDirs'
    ])
  },
  mounted () {
    this.$store.dispatch('loadBaseDirs')
  }
}
</script>

<style scoped>

</style>
