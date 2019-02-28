<template>
  <div>
    <h1>Project Details</h1>
    <div
      v-if="project"
      class="project-details"
    >
      <h2>{{ project.name }}</h2>
      <dl v-if="project">
        <dt>Hostname:</dt>
        <dd>{{ project.hostname }}</dd>
        <dt>Group:</dt>
        <dd>{{ project.group }}</dd>
        <dt>Enabled:</dt>
        <dd>{{ (true === project.enabled) ? 'true' : (false === project.enabled) ? 'false' : '&ndash;' }}</dd>
      </dl>
    </div>
    <div
      v-else="project"
      class="project-details"
    >
      <h2>{{ this.$route.params.projectName }}</h2>
      <p>This is a dummy project with no actual data!</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProjectDetails',
  computed: {
    project() {
      return this.$store.state.projects.find(p => p.name === this.$route.params.projectName);
    },
  },
};
</script>

<style scoped>

</style>
