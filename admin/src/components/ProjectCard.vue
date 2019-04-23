<template>
  <b-card
    :class="{'card--project': true,'showing-details': showingDetails, 'not-showing-details': !showingDetails}"
    :to="{path:'/project/'+id}"
  >
    <b-form>
      <b-form-group
        :id="`hostname-group-${projectIndex}`"
        class="hostname-group"
        label=""
        :label-for="`hostname-input-${projectIndex}`"
        :description="showingDetails ? 'Hostname' : ''"
      >
        <b-form-input
          :id="`hostname-input-${projectIndex}`"
          class="hostname-input"
          type="text"
          v-model="hostname"
          @change="maybeSubmit"
          size="lg"
          v-b-tooltip.hover.bottomright
          :title="showingDetails ? '' : 'Expand details'"
          required
          @click="showDetails"
          placeholder="" />
      </b-form-group>

      <project-toolbar :project="project" :projectIndex="projectIndex"></project-toolbar>

      <project-details :project="project" :projectIndex="projectIndex"></project-details>

    </b-form>

    <div slot="footer" v-if="project.attributes.stack">
      <project-stack :project="project" :projectIndex="projectIndex"></project-stack>
    </div>

  </b-card>

</template>

<script>

import { mapGetters } from 'vuex'
import ProjectToolbar from './ProjectToolbar'
import ProjectDetails from './ProjectDetails'
import ProjectStack from './ProjectStack'

export default {
  name: 'ProjectCard',
  props: {
    project: {
      type: Object,
      required: true
    },
    projectIndex: {
      type: Number,
      required: true
    }
  },
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
      showingDetails: false
    }
  },
  components: {
    ProjectToolbar,
    ProjectDetails,
    ProjectStack
  },
  computed: {
    ...mapGetters(['groupProjectStacks']),
    projectBase () {
      return this.escAttr(this.hostname) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    showDetails () {
      this.showingDetails = true
    },
    maybeSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          'hostname': this.id,
          'project': {
            id: this.id,
            attributes: this.$data
          }
        }
      ).then(() => {
        // this.$router.push('/project/' + this.hostname)
      })
    }
  }
}
</script>

<style scoped>
  .form-group label {
    font-size: 0.75rem;
    margin-bottom: 0;
  }

  .hostname-group{
    display: inline-block;
    float: left;
    width: calc(100% - 110px);
  }

  .not-showing-details .hostname-group{
    margin-bottom: 0;
  }

  .hostname-input{
    padding-left: 11px;
    font-weight: bold;
    /*padding: 1px 8px 1px 8px;
    margin-left: -10px;*/
  }

  .not-showing-details .hostname-input {
    border: 1px solid transparent;
    cursor: pointer;
  }

  .not-showing-details .hostname-input:hover {
    text-decoration: underline;
  }

  .showing-details .hostname-input{
    cursor: text;
  }

  .titlebar-icon {
    float: right;
    font-size: 1.25rem;
    cursor: pointer;
    margin-left: 10px;
    color: rgba(42, 85, 130, 0.98);
  }

  .titlebar-icon {
    padding-top: 3px;
  }

  .titlebar-icon.is-disabled  {
    color: rgba(17, 56, 85, 0.42);
  }

  .titlebar-icon--details {
    float: left;
  }

  .titlebar-icon--state {
    font-size: 1.65rem;
    padding-top: 0;
  }

  [data-icon="play-circle"] path {
    fill: green;
  }

  [data-icon="stop-circle"] path {
    fill: red;
  }
</style>
