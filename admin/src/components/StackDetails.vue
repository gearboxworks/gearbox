<template>
  <div class="stack-details">
    <h1>Stack Details</h1>
    <h2>{{ stack.label }}</h2>

    <el-row>
      <el-col :span="4">
        <el-card class="box-card" label="Web">
          <StackMemberWeb :member="stack.members.webserver"/>
        </el-card>
      </el-col>

      <el-col :span="4">
        <el-card label="Database">
          <StackMemberDB :member="stack.members.dbserver"/>
        </el-card>
      </el-col>

      <el-col :span="4">
        <el-card label="Cache">
          <StackMemberCache :member="stack.members.cacheserver"/>
        </el-card>
      </el-col>

      <el-col :span="4">
        <el-card label="Process">
          <StackMemberWorker :member="stack.members.processvm"/>
        </el-card>
      </el-col>
    </el-row>

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
