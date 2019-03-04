<template>
  <el-menu
    :router="true"
  >
    <el-submenu index="2">
      <template slot="title">
        <i class="el-icon-tickets" />
        Stacks
      </template>
      <el-menu-item
        v-for="stack in stacks"
        :index="'2-'+stack.name"
        :key="stack.name"
        :route="{path:'/stack/'+stack.name}"
      >
        {{ stack.label }}
      </el-menu-item>
    </el-submenu>
  </el-menu>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'StackList',
  computed: {
    ...mapState([
      'stacks'
    ])
  },
  mounted () {
    this.$store.dispatch('loadStacks')
  }
}
</script>
