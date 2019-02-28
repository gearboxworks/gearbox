<template>
  <div class="stack-list">
    <h3>Stacks</h3>
    <ul>
      <li
        v-for="( stack, index ) in stacks"
        :key="index"
      >
        {{ stack.label }}
      </li>
    </ul>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'StackList',
  computed: {
    ...mapState([
      'stacks',
    ]),
  },
  mounted() {
    this.$store.dispatch('loadStacks');
  },
};
</script>
