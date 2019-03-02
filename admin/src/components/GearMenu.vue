<template>
  <el-menu>
    <el-submenu index="3">
      <template slot="title">
        <i class="el-icon-setting" />
        Gears
      </template>
      <el-menu-item
        v-for="(gear, index) in gears"
        :index="index"
        :key="index"
      >
        {{ gear.name }}
      </el-menu-item>
    </el-submenu>
  </el-menu>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'GearList',
  computed: {
    ...mapState([
      'gears',
    ]),
  },
  mounted() {
    this.$store.dispatch('loadGears');
  },
};
</script>
