<template>

  <div id="app">
    <b-alert :show="isConnectionProblem" variant="warning">
      <h4>Connection Problem</h4>
      <p>It seems that Gearbox Server is not running. Remaining connection attempts: {{remainingRetries}}</p>
    </b-alert>
    <b-alert show v-if="isUnrecoverableConnectionProblem" variant="danger">
      <h4>Connection Problem</h4>
      <p>Failed to connect to Gearbox Server.</p>
    </b-alert>
    <the-top-bar />
    <router-view></router-view>
  </div>
</template>

<script>
import TheTopBar from './components/TheTopBar'

export default {
  name: 'App',
  components: {
    TheTopBar
  },
  computed: {
    isConnectionProblem () {
      // console.log('isConnectionProblem', this.$store.state.connectionStatus.networkError, this.$store.state.connectionStatus.remainingRetries)
      return this.$store.state.connectionStatus.networkError && this.$store.state.connectionStatus.remainingRetries > 0
    },
    remainingRetries () {
      return this.$store.state.connectionStatus.remainingRetries
    },
    isUnrecoverableConnectionProblem () {
      return this.$store.state.connectionStatus.networkError
        ? (this.$store.state.connectionStatus.remainingRetries === 0)
        : ''
    }
  }
}
</script>

<style>
body{
  margin: 0;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: left;
  color: #2c3e50;
  margin: 0;
}
.el-header, .el-footer {
  background-color: #B3C0D1;
  color: #333;
  line-height: 60px;
}
.el-aside {
  color: #333;
}
</style>
