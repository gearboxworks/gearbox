<template>
  <div :id="`${projectBase}details`" role="tabpanel">
    <b-form-group
      :id="`${projectBase}location-group`"
      :label-for="`${projectBase}location-input`"
      label=""
      description="Location"
    >

      <project-location :project="project" :projectIndex="projectIndex"></project-location>
    </b-form-group>

    <b-form-group
      id="notesGroup"
      label=""
      label-for="notesInput"
      description="(will be visible only here)"
    >
      <project-notes :project="project" :projectIndex="projectIndex"></project-notes>
    </b-form-group>

    <project-stack-select :project="project" :projectIndex="projectIndex"></project-stack-select>

    <a class="hide-details"
       title="Hide project details"
       @click="$emit('toggle-details')"
    >
      <font-awesome-icon
        :icon="['fa', 'chevron-up']"
      />
      <span>Hide</span>
    </a>

  </div>
</template>

<script>
import ProjectLocation from '../ProjectLocation'
import ProjectNotes from '../ProjectNote'
import ProjectStackSelect from '../ProjectStack'

export default {
  name: 'ProjectCardDetails',
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
  components: {
    ProjectLocation,
    ProjectNotes,
    ProjectStackSelect
  },
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
      selectedService: ''
    }
  },
  computed: {
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
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
    },
    onClosePopoverFor (triggerElementId) {
      this.$root.$emit('bv::hide::popover', triggerElementId)
    }
  }
}
</script>
<style scoped>
  .collapse {
    clear: both;
  }

  .hide-details {
    display: block;
    margin-top: 15px;
    text-align: center;
    margin-bottom: -5px;
    margin-right: auto;
    padding: 1px 6px;
    margin-left: auto;
    color: #1e69b9 !important;
    opacity: 0;
    cursor: pointer;
    transition: opacity 400ms;
  }
  .hide-details span {
    margin-right: 5px;
    margin-left: 5px;
    font-weight: bold;
    font-size: 14px;
  }
  .card--project:hover .hide-details{
    opacity:0.75;
  }
  .card--project:hover .hide-details:hover {
    opacity: 1;
  }
</style>
