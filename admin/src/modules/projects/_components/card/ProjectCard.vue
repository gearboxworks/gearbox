<template>
  <b-card
    class="card--project"
  >
    <div class="clearfix">
      <project-hostname
        :project="project"
        :projectIndex="projectIndex"
        :is-multimodal="true"
        @show-alert="showAlert"
      />

      <project-toolbar
        :project="project"
        :projectIndex="projectIndex"
        @run-stop="onRunStop"
        :isUpdating="isUpdating"
      />
    </div>

    <b-alert
      :show="alertShow"
      :dismissible="alertDismissible"
      :variant="alertVariant"
      @dismissed="alertShow=false"
      fade
    >
      {{alertContent}}
    </b-alert>

    <div class="clearfix" slot="footer">

      <project-stack-list
        :project="project"
        :projectIndex="projectIndex"
      />

      <project-stack-add
        :project="project"
        :projectIndex="projectIndex"
        @maybe-hide-alert="maybeHideAlert"
      />

      <project-location
        :project="project"
        :projectIndex="projectIndex"
      />

      <project-note
        :project="project"
        :projectIndex="projectIndex"
      />

    </div>

  </b-card>

</template>

<script>
import { createNamespacedHelpers } from 'vuex'
import ProjectToolbar from '../ProjectToolbar'
import ProjectHostname from '../ProjectHostname'
import ProjectLocation from '../ProjectLocation'
import ProjectNote from '../ProjectNote'
import ProjectStackAdd from '../ProjectStackAdd'
import ProjectStackList from '../ProjectStackList'

const { mapActions } = createNamespacedHelpers('projects')

export default {
  name: 'ProjectCard',
  components: {
    ProjectToolbar,
    ProjectHostname,
    ProjectLocation,
    ProjectStackList,
    ProjectStackAdd,
    ProjectNote
  },
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
      showingDetails: false,
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'warning',
      isUpdating: false
    }
  },
  computed: {
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    },
    isRunning () {
      return this.project.attributes.enabled
    }
  },
  methods: {
    ...mapActions({ updateProjectState: 'updateState' }),
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    showAlert (alert) {
      if (typeof alert === 'string') {
        this.alertContent = alert
      } else {
        this.alertVariant = alert.variant || this.alertVariant
        this.alertDismissible = alert.dismissible || this.alertDismissible
        this.alertContent = alert.content || this.alertContent
      }
      this.alertShow = true
    },
    hideAlert () {
      this.alertContent = ''
      this.alertShow = false
    },
    maybeHideAlert (alert) {
      if (this.alertContent === alert) {
        this.hideAlert()
      }
    },
    onRunStop () {
      if (this.project.attributes.stack && this.project.attributes.stack.length > 0) {
        this.isUpdating = true
        this.updateProjectState({ project: this.project, isEnabled: !this.isRunning })
          .then((status) => {
            this.isUpdating = false
            this.hideAlert()
          })
      } else {
        this.showAlert('Please add some stacks first!')
      }
    }
  }
}
</script>

<style scoped>
  .card--project {
    margin-bottom: 1.5rem;
    transition: box-shadow 400ms;
  }

  .card--project:hover {
    box-shadow: 1px 1px 4px rgba(0, 0, 0, 0.1);
  }

  .card-body {
    padding: 0.75rem;
  }

  .card-body .alert {
    margin-left: -0.75rem;
    width: calc(100% + 1.5rem);
    margin-top: 0.5rem;
    margin-bottom: -14px;
    border-top-left-radius: 0;
    border-top-right-radius: 0;
    font-size: 12px;
    padding-top: 1rem;
    padding-left: 25px;
    padding-bottom: 1rem;
  }

  .card-body .alert-dismissible .close {
    padding: 0.5rem 0.75rem;
    right: 0px;
    top: 3px;
  }

  .form-group label {
    font-size: 0.75rem;
    margin-bottom: 0;
  }

  .card-footer {
    background-color: rgb(246,246,246);
    padding: 0.75rem 0.75rem 0.25rem 0.75rem;
  }
  .btn-outline-warning,
  .btn-outline-info {
    border-color: transparent;
  }

  .btn--submit-hostname {
    border-color:#ced4da;
  }
  .btn-outline-warning {
    color: #e4bd77;
  }

  .btn-outline-warning:hover {
    color: #212529;
    background-color: #e4bd77;
  }
  .btn--add {
    position:relative;
  }
  .btn--add span {
    position: absolute;
    right: 6px;
    font-size: 17px;
    top: -2px;
  }

  .btn--add svg {
    position: relative;
    left: -2px;
    top: 2px;
  }

  .project-note-list {
    float:right;
  }

  .input-group--stack,
  .input-group--location {
    margin-bottom: 0.5rem;
  }

  .input-group--stack:not(.is-collapsed),
  .input-group--location:not(.is-collapsed) {
    position: absolute;
    width: calc(100% - 1.5rem);
    z-index: 3;
    bottom: 0.25rem;
    left: 0.75rem;
    background-color: rgb(246,246,246);
  }

  .input-group--location:not(.is-collapsed){
    background-color: white;
  }

  .input-group.is-collapsed {
    display: inline-flex;
    width: auto;
    margin-right: 0.5rem;
  }

  .card--project:hover .input-group--hostname:not(.is-editing) .hostname-input {
    text-decoration: underline;
  }

  .input-group--hostname{
    margin-right: 0;
    float: left;
    /**
     * leaving some room for the Run/Stop icon
     */
    width: calc(100% - 50px);
  }

  .input-group--note.is-collapsed{
    margin-right: 0;
    float:right;
  }

</style>

<style>
.card-footer .stack-card.is-collapsible:not(.is-collapsed) {
  width: 100%;
  margin-right: 0;
}

</style>
