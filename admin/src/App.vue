<template>

  <div id="app">

    <b-alert
      :show="isConnectionProblem"
      variant="warning"
    >
      <h4>Connection Problem</h4>
      <p>It seems that Gearbox Server is not running. Remaining connection attempts: {{remainingRetries}}</p>
    </b-alert>
    <b-alert
      v-if="isUnrecoverableConnectionProblem"
      show
      variant="danger"
    >
      <h4>Connection Problem</h4>
      <p>Failed to connect to Gearbox Server.</p>
    </b-alert>

    <the-top-bar />

    <router-view></router-view>

  </div>
</template>

<script>
import TheTopBar from './components/TheTopBar'

import BasedirMethodTypes from './modules/basedirs/_store/public-types'
import GearspecMethodTypes from './modules/gearspecs/_store/public-types'
import StackMethodTypes from './modules/stacks/_store/public-types'
import ServiceMethodTypes from './modules/services/_store/public-types'
import ProjectMethodTypes from './modules/projects/_store/public-types'

const { ActionTypes: BasedirActions } = BasedirMethodTypes
const { ActionTypes: GearspecActions } = GearspecMethodTypes
const { ActionTypes: StackActions } = StackMethodTypes
const { ActionTypes: ServiceActions } = ServiceMethodTypes
const { ActionTypes: ProjectActions } = ProjectMethodTypes

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
  },

  mounted () {
    // console.log(
    //   BasedirActions.LOAD_ALL,
    //   StackActions.LOAD_ALL,
    //   ServiceActions.LOAD_ALL,
    //   GearspecActions.LOAD_ALL,
    //   ProjectActions.LOAD_ALL,
    //   ProjectActions.LOAD_ALL_DETAILS
    // )

    this.$store.dispatch(BasedirActions.LOAD_ALL)
    this.$store.dispatch(StackActions.LOAD_ALL)
    this.$store.dispatch(ServiceActions.LOAD_ALL)
    this.$store.dispatch(GearspecActions.LOAD_ALL)

    this.$store.dispatch(ProjectActions.LOAD_ALL).then(() => {
      return this.$store.dispatch(ProjectActions.LOAD_ALL_DETAILS)
    })
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
