<template>
  <div id="app">
    <el-alert
      v-if="isConnectionProblem"
      title="Connection Problem"
      type="alert"
      :description="'It seems that Gearbox Server is not running. Remaining connection attempts: ' + remainingAttempts"
      show-icon>
    </el-alert>
    <el-alert
      v-if="isUnrecoverableConnectionProblem"
      title="Connection Failed"
      type="error"
      description="Failed to connect to Gearbox Server."
      show-icon>
    </el-alert>

    <el-container
      style="height: 1024px; border: 1px solid #eee"
    >
      <Sidebar />
      <el-container>
        <el-header
          style="text-align: right; font-size: 12px"
        >
          <el-dropdown>
            <i
              class="el-icon-setting"
              style="margin-right: 15px"
            />
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>View</el-dropdown-item>
              <el-dropdown-item>Add</el-dropdown-item>
              <el-dropdown-item>Delete</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <span>Tom</span>
        </el-header>
        <el-main>
          <router-view />
        </el-main>
        <el-footer>
          <p>&copy; Gearbox Works, 2019</p>
        </el-footer>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import Sidebar from './components/Sidebar.vue'

export default {
  name: 'App',
  components: {
    Sidebar
  },
  computed: {
    isConnectionProblem () {
      // console.log('isConnectionProblem', this.$store.state.connectionStatus.networkError, this.$store.state.connectionStatus.remainingRetries)
      return this.$store.state.connectionStatus.networkError && this.$store.state.connectionStatus.remainingRetries > 0
    },
    remainingAttempts () {
      return this.$store.state.connectionStatus.remainingRetries
    },
    isUnrecoverableConnectionProblem () {
      return this.$store.state.connectionStatus.networkError
        ? (this.$store.state.connectionStatus.networkError && this.$store.state.connectionStatus.remainingRetries === 0)
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
