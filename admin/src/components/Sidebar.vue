<template>
  <el-aside
    width="200px"
    style="background-color: rgb(238, 241, 246)"
  >
    <h3>
      <router-link
        to="/"
      >
        Gearbox Admin
      </router-link>
    </h3>

    <ProjectMenu />
    <StackMenu />

    <el-menu>
      <el-submenu index="4">
        <template slot="title">
          <i class="el-icon-more-outline" />
          <router-link
            to="/preferences"
          >
            Preferences
          </router-link>
        </template>
      </el-submenu>
    </el-menu>
  </el-aside>
</template>

<script>
import ProjectMenu from '../views/Projects.vue'
import StackMenu from './StackMenu.vue'

export default {
  name: 'Sidebar',
  components: {
    ProjectMenu,
    StackMenu
  }
}
</script>

<style scoped>
.el-aside h3 {
  text-align: center;
}
.el-aside h3 a {
  text-decoration: none;
}

</style>
