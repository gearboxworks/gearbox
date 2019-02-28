<template>
  <div class="project-list">
    <h3>Projects</h3>
    <ul>
      <li
        v-for="( project, index ) in projects"
        :key="index"
      >
        {{ project.name }}
      </li>
    </ul>
  </div>
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
