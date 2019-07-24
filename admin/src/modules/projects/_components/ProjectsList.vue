<template>
  <div
    class="filtered-projects"
  >
    <b-card-group
      v-if="viewMode==='cards'"
      key="projects-view"
      columns
      class="pl-3 pr-3"
    >
      <project-card
        v-for="(project, projectIndex) in projects"
        :key="project.id"
        :project="project"
        :projectIndex="projectIndex"
      />
    </b-card-group>
    <table
      v-else
      key="projects-view"
      class="projects-table"
    >
      <thead>
      <tr>
        <th class="th--state">State</th><th class="th--hostname">Project Name</th><th class="th--location">Location</th><th class="th--stack">Stack</th><th class="th--notes">Notes</th>
      </tr>
      </thead>
      <tbody>
      <project-row
        v-for="(project, projectIndex) in projects"
        :key="project.id"
        :project="project"
        :projectIndex="projectIndex"
      />
      </tbody>
    </table>
  </div>
</template>

<script>
import ProjectCard from './card/ProjectCard'
import ProjectRow from './row/ProjectRow'

export default {
  name: 'ProjectsList',
  components: {
    ProjectCard,
    ProjectRow
  },
  props: {
    projects: {
      type: Array,
      required: true
    },
    viewMode: {
      type: String,
      required: true,
      default: 'cards',
      validator: (value) => {
        return [
          'cards',
          'table'
        ].indexOf(value) !== -1
      }
    }
  }
}
</script>

<style scoped>
  .el-icon-caret-right {
    color: red;
  }
  @media (min-width: 576px) {
    .card-columns {
      -webkit-column-count: 1;
      -moz-column-count: 1;
      column-count: 1;
      -webkit-column-gap: 1.25rem;
      -moz-column-gap: 1.25rem;
      column-gap: 1.25rem;
      orphans: 1;
      widows: 1;
      text-align: center;
    }
    .card-columns .card {
      max-width: 500px;
      /* margin-left: auto; */
      /* margin-right: auto; */
      text-align: left;
    }
  }
  @media (min-width: 800px) {
    .card-columns {
      -webkit-column-count: 2;
      -moz-column-count: 2;
      column-count: 2;
    }
  }
  @media (min-width: 1024px) {
    .card-columns {
      -webkit-column-count: 3;
      -moz-column-count: 3;
      column-count: 3;
    }
  }
  @media (min-width: 1440px) {
    .card-columns {
      -webkit-column-count: 4;
      -moz-column-count: 4;
      column-count: 4;
    }
  }
  @media (min-width: 1960px) {
    .card-columns {
      -webkit-column-count: 5;
      -moz-column-count: 5;
      column-count: 5;
    }
  }
  @media (min-width: 2400px) {
    .card-columns {
      -webkit-column-count: 6;
      -moz-column-count: 6;
      column-count: 6;
    }
  }

  .filtered-projects{
    padding-left: 0;
    padding-right: 1rem;
  }

  .projects-table {
    width: calc(100% - 1rem);
    margin-left: 1rem;
  }

  .th--state {
    width: 50px;
  }

  .th--hostname {
    width: 200px;
  }

  .th--location {
    width: 400px;
  }

  .th--notes {
    width: 300px;
  }
</style>
