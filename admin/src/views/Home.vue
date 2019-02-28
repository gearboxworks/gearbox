<template>
  <div class="home">
    <h1>{{ title }}</h1>
    <ProjectList />
    <StackList />
  </div>
</template>

<script>
import ProjectList from '../components/ProjectList.vue';
import StackList from '../components/StackList.vue';

export default {
  props: {
    title: {
      type: String,
      default: 'Default Home Title',
    },
  },
  components: {
    ProjectList,
    StackList,
  },
};
</script>
