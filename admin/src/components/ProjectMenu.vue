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
      <el-submenu index="1-1">
        <template slot="title">
          Gearbox Projects
        </template>
        <el-menu-item
          v-for="project in projects"
          :index="project.name"
          :key="project.name"
          :route="{path:'/project/'+project.name}"
        >
          {{ project.name }}
        </el-menu-item>
      </el-submenu>
      <el-submenu index="1-2">
        <template slot="title">
          Candidate Projects
        </template>
        <el-menu-item
          index="1-2-1"
          :route="{ path: '/project/candidate1' }"
        >
          Candidate Project 1
        </el-menu-item>
      </el-submenu>
    </el-submenu>
  </el-menu>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'ProjectList',
  computed: {
    ...mapState([
      'projects',
    ]),
  },
  mounted() {
    this.$store.dispatch('loadProjects');
  },
};
</script>
