<template>

  <div id="app">

    <b-alert
      :show="!!networkError"
      class="network-error"
      variant="danger"
    >
      <h4>{{$t('message.noConnectionTitle')}}</h4>
      <p>{{$t('message.noConnectionDesc')}} <a href="#" class="retry-link" @click="$router.go()">{{$t('message.noConnectionRetry')}}</a></p>
    </b-alert>

    <the-top-bar />

    <router-view></router-view>

  </div>
</template>

<script>
import TheTopBar from './components/TheTopBar'

import { BasedirActions } from './modules/basedirs/_store/method-names'
import { GearspecActions } from './modules/gearspecs/_store/method-names'
import { StackActions } from './modules/stacks/_store/method-names'
import { ServiceActions } from './modules/services/_store/method-names'
import { ProjectActions } from './modules/projects/_store/method-names'

export default {
  name: 'App',
  components: {
    TheTopBar
  },
  data () {
    return {
      networkError: null
    }
  },
  methods: {
    setError (e) {
      if (e.message === 'Network Error') {
        this.networkError = e
      }
    },
    clearError () {
      this.networkError = null
    }
  },
  mounted () {
    this.$store.dispatch(BasedirActions.LOAD_ALL).catch(e => this.setError(e))
    this.$store.dispatch(StackActions.LOAD_ALL).catch(e => this.setError(e))
    this.$store.dispatch(ServiceActions.LOAD_ALL).catch(e => this.setError(e))
    this.$store.dispatch(GearspecActions.LOAD_ALL).catch(e => this.setError(e))
    this.$store.dispatch(ProjectActions.LOAD_ALL)
      .then(() => {
        return this.$store.dispatch(ProjectActions.LOAD_ALL_DETAILS)
      })
      .then(r => this.clearError)
      .catch(e => this.setError(e))
  }
}
</script>

<style scoped>
  .el-header, .el-footer {
    background-color: #B3C0D1;
    color: #333;
    line-height: 60px;
  }
  .network-error .retry-link {
    color: darkorchid;
    font-weight: bold;
  }
  .network-error {
    margin-bottom: 0;
  }
  .network-error p {
    margin-bottom: 4px;
  }
</style>
<style>
body{
  margin: 0;
}
#app {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Helvetica Neue", Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
  color: #2c3e50;
  margin: 0;
}
</style>
