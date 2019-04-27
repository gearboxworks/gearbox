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
          required
          @click="showingDetails = true"
          placeholder="" />
      </b-form-group>

      <project-toolbar :project="project" :projectIndex="projectIndex"></project-toolbar>

      <project-details :project="project" :projectIndex="projectIndex" v-if="showingDetails" @toggle-details="toggleDetails"></project-details>

      <a class="show-details"
         title="Show project details"
         @click="showingDetails = true"
      >
        <font-awesome-icon
          :icon="['fa', 'ellipsis-h']"
        />
      </a>

    </b-form>

    <div slot="footer" v-if="project.attributes.stack && project.attributes.stack.length">
      <project-stack :project="project" :projectIndex="projectIndex"></project-stack>
    </div>

  </b-card>

</template>

<script>

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
    projectBase () {
      return this.escAttr(this.id) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    toggleDetails () {
      this.showingDetails = !this.showingDetails
    },
    maybeSubmit (ev) {
      this.$store.dispatch(
        'updateProject',
        {
          id: this.id,
          attributes: this.$data
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
    margin-top: -10px;
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
    width: auto;
  }

  .not-showing-details .hostname-input:hover {
    text-decoration: underline;
  }

  .showing-details .hostname-input{
    cursor: text;
    width: 100%;
  }

  .show-details {
    display: block;
    position:relative;
    top: -5px;
    text-align: left;
    margin-bottom: -8px;
    line-height: 0;
    margin-left: 7px;
    padding: 1px 6px;
    color: #1e69b9 !important;
    opacity: 0;
    cursor: pointer;
    transition: opacity 400ms;
    clear: both;
  }
  .show-details span {
    margin-left: 5px;
    margin-right: 5px;
    font-weight: bold;
    font-size: 14px;
  }

  .showing-details .show-details {
    display: none;
  }

  .card--project:hover .show-details{
    opacity:0.75;
  }
  .card--project:hover .show-details:hover {
    opacity: 1;
  }

</style>
