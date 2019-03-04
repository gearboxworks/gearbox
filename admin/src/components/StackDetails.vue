<template>
  <div class="stack-details">
    <h1>Stack Details</h1>
    <h2>{{ stack.label }}</h2>
    <el-tabs type="border-card">
      <el-tab-pane label="Web">
        <StackMemberWeb :member="stack.members.webserver" />
      </el-tab-pane>
      <el-tab-pane label="Database">
        <StackMemberDB :member="stack.members.dbserver" />
      </el-tab-pane>
      <el-tab-pane label="Cache">
        <StackMemberCache :member="stack.members.cacheserver" />
      </el-tab-pane>
      <el-tab-pane label="Process">
        <StackMemberWorker :member="stack.members.processvm" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>

import StackMemberWeb from './StackMemberWeb.vue'
import StackMemberDB from './StackMemberDB.vue'
import StackMemberCache from './StackMemberCache.vue'
import StackMemberWorker from './StackMemberWorker.vue'

export default {
  name: 'StackDetails',
  computed: {
    stack () {
      return this.$store.state.stacks[this.$route.params.stackName]
    }
  },
  components: {
    StackMemberWeb,
    StackMemberDB,
    StackMemberCache,
    StackMemberWorker
  }
}
</script>

<style scoped>

</style>
