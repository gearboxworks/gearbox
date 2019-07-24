<template>
  <b-form :class="{'form--preferences': true, 'is-loading': isLoading}">
    <h1>Preferences</h1>
    <h2>Project Directories</h2>

    <div
      v-if="!isLoading"
      key="basedirs-content"
      class="basedirs-wrap"
    >
      <basedir-row-edit
        v-for="(basedir,index) in basedirs"
        :key="basedir.id"
        :basedir="basedir"
        :tab-offset="index"
        :is-deletable="basedir.id !== 'default'"
      />
    </div>
    <div
        v-else
        key="basedirs-content"
        class="basedirs-loading"
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
  components: {
    BasedirRowEdit,
    BasedirRowAdd
  },
  data () {
    return {
      isLoading: true,
      errors: {},
      touched: {}
    }
  },
  computed: {
    ...mapGetters({
      'basedirs': 'basedirs/all'
    })
  },
  mounted () {
    this.doLoadBasedirs()
      .then(() => {
        this.isLoading = false
      })
  },
  methods: {
    ...mapActions({
      'doLoadBasedirs': 'basedirs/loadAll'
    })
  }
}
</script>

<style scoped>
  .form--preferences {
    padding: 1rem;
  }
</style>
