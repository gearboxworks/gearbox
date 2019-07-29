<template>
  <b-card
    class="card--project"
  >
    <div class="clearfix">
      <project-hostname
        :is-multimodal="true"
        @show-alert="showAlert"
      />

      <project-toolbar
        @run-stop-project="onRunStop"
        :is-updating="isUpdating"
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
        :expanded-stack-ids="expandedStackIds"
        @expand-collapse-stack="onExpandCollapseStack"
      />

      <project-stack-add
        @maybe-hide-alert="maybeHideAlert"
        @added-stack="onAddedStack"
      />

      <project-location />

      <project-note />

    </div>
  </b-card>

</template>

<script>
import Vue from 'vue'
import { ProjectActions, ProjectMutations } from '../_store/public-types'

import ProjectToolbar from './shared/ProjectToolbar'
import ProjectHostname from './shared/ProjectHostname'
import ProjectLocation from './shared/ProjectLocation'
import ProjectNote from './shared/ProjectNote'
import ProjectStackAdd from './shared/ProjectStackAdd'
import ProjectStackList from './shared/ProjectStackList'

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
      isUpdating: false,
      expandedStackIds: {}
    }
  },

  provide () {
    return {
      project: this.project,
      projectPrefix: 'gb-' + this.escAttr(this.id) + '-'
    }
  },
  computed: {
    isRunning () {
      return this.project.attributes.enabled
    },
    theTypes () {
      return ProjectActions
    }
  },
  methods: {

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

    async onRunStop () {
      if (this.project.attributes.stack && this.project.attributes.stack.length > 0) {
        this.isUpdating = true
        try {
          await this.$store.dispatch(
            ProjectActions.UPDATE_STATE,
            {
              project: this.project,
              isEnabled: !this.isRunning
            }
          )
          this.isUpdating = false
          this.hideAlert()
        } catch (e) {
          console.error(e.message)
        }
      } else {
        this.showAlert('Please add some stacks first!')
      }
    },

    onAddedStack (stackId) {
      const cleanStackId = stackId.replace('(removed)', '')
      const newObject = { [cleanStackId]: 1, [cleanStackId + '(removed)']: 1 }

      for (const key in this.expandedStackIds) {
        // const val = this.expandedStackIds[key]
        if ((key !== stackId && key !== cleanStackId)) {
          newObject[key] = -1
        }
      }
      Vue.set(this, 'expandedStackIds', newObject)
    },

    onExpandCollapseStack (stackId, isExpanded) {
      Vue.set(this.expandedStackIds, stackId, isExpanded ? -1 : 1)
    }

  },
  mounted () {
    this.$store.subscribe((mutation, state) => {
      switch (mutation.type) {
        case ProjectMutations.REMOVE_PROJECT_STACK:
          const { stackId, project } = mutation.payload
          if (project.id === this.project.id) {
            Vue.set(this.expandedStackIds, stackId, 0)
            Vue.set(this.expandedStackIds, stackId + '(removed)', 0)
          }
          break
      }
    })
  }
}
</script>

<style scoped>
  .card--project {
    margin-bottom: 1.5rem;
    transition: box-shadow 400ms;
  }

  .card--project:active,
  .card--project:focus,
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
  .card-footer >>> .stack-card.is-collapsible.is-expanded {
    width: 100%;
    margin-right: 0;
  }
</style>
