<template>
  <b-form :class="{'form--preferences': true, 'is-loading': isLoading}">
    <h1>Preferences</h1>
    <h2>Project Directories</h2>

    <div class="basedirs-wrap" v-if="!isLoading">
      <basedir-row-edit
        v-for="(basedir,index) in basedirs"
        :key="basedir.id"
        :basedir="basedir"
        :tab-offset="index"
        :is-deletable="basedir.id !== 'default'"
      />
    </div>
    <div class="basedirs-loading"
         v-else
    >
      <font-awesome-icon
        icon="circle-notch"
        spin
      />
      &nbsp;<span>Loading directories...</span>
    </div>

    <basedir-row-add :tab-offset="basedirs.length*3" />
  </b-form>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import BasedirRowEdit from '../components/basedir/BasedirRowEdit'
import BasedirRowAdd from '../components/basedir/BasedirRowAdd'

export default {
  name: 'Preferences',
  data () {
    return {
      isLoading: true,
      errors: {},
      touched: {}
    }
  },
  computed: {
    ...mapGetters({ 'basedirs': 'basedirs/all' })
  },
  methods: {
    ...mapActions({ 'doLoadBasedirs': 'basedirs/loadAll' })
  },
  components: {
    BasedirRowEdit,
    BasedirRowAdd
  },
  mounted () {
    this.doLoadBasedirs()
      .then(() => {
        this.isLoading = false
      })
  }
}
</script>

<style scoped>
  .form--preferences {
    padding: 1rem;
  }
</style>
