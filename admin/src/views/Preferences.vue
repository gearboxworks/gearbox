<template>
  <div>
    <h1>Preferences</h1>
    <b-form>
      <h2>Base Directories</h2>
      <b-container fluid>
        <b-form-row class="my-1">
          <b-col sm="4">
            <label for="input-none">Name:</label>
          </b-col>
          <b-col sm="7">
            <label for="input-none">Path:</label>
          </b-col>
          <b-col sm="1">
            <label for="input-none">Action:</label>
          </b-col>
        </b-form-row>

        <b-form-row
          v-for="(baseDir, value) in $store.state.baseDirs"
          :key="value"
          class="my-1">
          <b-col sm="4">
            <b-form-input
              :id="value+'BaseDirName'"
              type="text"
              v-model="baseDir.value"
              required
              placeholder="Name" />
          </b-col>

          <b-col sm="7">
            <b-form-input
              :id="value+'BaseDirPath'"
              type="text"
              v-model="baseDir.text"
              required
              placeholder="Path" />
          </b-col>
        </b-form-row>

        <b-form-row class="my-1">
          <b-col sm="4">
            <b-form-input
              v-model="baseDirName"
              type="text"
              required
              placeholder="" />
          </b-col>

          <b-col sm="7">
            <b-form-input
              v-model="baseDirPath"
              type="text"
              required
              placeholder="" />
          </b-col>
          <b-col sm="1">
            <b-button
              type="submit.prevent"
              variant="success"
              @click="addNewBaseDir"
            >
              Add
            </b-button>
          </b-col>
        </b-form-row>
      </b-container>
    </b-form>
  </div>
</template>

<script>
// import { mapState } from 'vuex'

export default {
  name: 'Preferences',
  data () {
    return {
      baseDirName: '',
      baseDirPath: ''
    }
  },
  computed: {
    // ...mapState(['baseDirs'])
  },
  mounted () {
    this.$store.dispatch('loadBaseDirs')
  },
  methods: {
    addNewBaseDir () {
      this.$store.dispatch('addBaseDir', { name: this.baseDirName, path: this.baseDirPath })
      this.baseDirName = ''
      this.baseDirPath = ''
      // console.log(this.baseDirs)
      // this.$refs['baseDirName'].$el.value = ''
      // this.$refs['baseDirPath'].$el.value = ''
    }
  }
}
</script>

<style scoped>

</style>
