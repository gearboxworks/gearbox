<template>
  <el-menu
    :default-openeds="['1', '1-1']"
    :router="true"
  >
    <el-submenu index="1">
      <template slot="title">
        <i class="el-icon-menu" />
        Projects
      </template>
      <el-menu-item
        v-for="project in projects"
        :index="project.name"
        :key="project.name"
        :route="{path:'/project/'+project.name}"
      >
        {{ project.name }}
        <i class="el-icon-caret-right" />
      </el-menu-item>
    </el-submenu>
    <el-submenu index="1-2">
      <template slot="title">
        <i class="el-icon-menu" />
        Stopped Projects
      </template>
      <el-menu-item
        index="1-2-1"
        :route="{ path: '/project/stopped1' }"
      >
        Dummy Stopped 1
      </el-menu-item>
    </el-submenu>
    <el-submenu index="1-3">
      <template slot="title">
        <i class="el-icon-menu" />
        Candidate Projects
      </template>
      <el-menu-item
        index="1-3-1"
        :route="{ path: '/project/candidate1' }"
      >
        Dummy Candidate 1
      </el-menu-item>
    </el-submenu>
  </el-menu>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'ProjectList',
  computed: {
    ...mapState([
      'projects'
    ])
  },
  mounted () {
    this.$store.dispatch('loadProjects')
  }
}
</script>
<style scoped>
  .el-icon-caret-right {
    color: green;
  }
</style>
