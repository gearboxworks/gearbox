<template>
  <el-aside
    width="200px"
    style="background-color: rgb(238, 241, 246)"
  >
    <h3>Gearbox Admin</h3>

    <ProjectMenu />
    <StackMenu />
    <GearMenu />

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
import ProjectMenu from './ProjectMenu.vue';
import StackMenu from './StackMenu.vue';
import GearMenu from './GearMenu.vue';

export default {
  name: 'Sidebar',
  components: {
    ProjectMenu,
    StackMenu,
    GearMenu,
  },
};
</script>

<style scoped>

</style>
