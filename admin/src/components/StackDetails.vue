<template>
  <div class="stack-details">
    <h1>Stack Details</h1>
    <h2>{{ stack.label }}</h2>

    <el-row>
      <el-col :span="4">
        <el-card class="box-card" label="Web">
          <ServiceWeb :service="stack.services.webserver"/>
        </el-card>
      </el-col>

      <el-col :span="4">
        <el-card label="Database">
          <ServiceDB :service="stack.services.dbserver"/>
        </el-card>
      </el-col>

      <el-col :span="4">
        <el-card label="Cache">
          <ServiceCache :service="stack.services.cacheserver"/>
        </el-card>
      </el-col>

      <el-col :span="4">
        <el-card label="Process">
          <ServiceProcessVM :service="stack.services.processvm"/>
        </el-card>
      </el-col>
    </el-row>

  </div>
</template>

<script>

import ServiceWeb from './ServiceWeb.vue'
import ServiceDB from './ServiceDB.vue'
import ServiceCache from './ServiceCache.vue'
import ServiceProcessVM from './ServiceProcessVM.vue'

export default {
  name: 'StackDetails',
  computed: {
    stack () {
      return this.$store.state.stacks[this.$route.params.stackName]
    }
  },
  components: {
    ServiceWeb,
    ServiceDB,
    ServiceCache,
    ServiceProcessVM
  }
}
</script>

<style scoped>

</style>
