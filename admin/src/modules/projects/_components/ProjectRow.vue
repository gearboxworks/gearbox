<template>
  <tr class="row--project" role="row">
    <td class="td--state">
      <project-toolbar
        @run-stop-project="onRunStop"
        :is-updating="isUpdating"
      />
    </td>

    <td class="td--hostname">
      <project-hostname
        :is-multimodal="false"
        @show-alert="showAlert"
      />
    </td>

    <td class="td--location">
      <project-location :is-multimodal="false" />
    </td>

    <td class="td--stack">
      <project-stack-list
        :start-expanded="false"
        :expanded-stack-ids="expandedStackIds"
        @expand-collapse-stack="onExpandCollapseStack"
      />
      <project-stack-add
        @added-stack="onAddedStack"
      />
    </td>

    <td class="td--notes">
      <project-note />
    </td>

  </tr>

</template>

<script>
import Vue from 'vue'
import { ProjectActions } from '../_store/public-types'

import ProjectHostname from './shared/ProjectHostname'
import ProjectToolbar from './shared/ProjectToolbar'
import ProjectLocation from './shared/ProjectLocation'
import ProjectNote from './shared/ProjectNote'
import ProjectStackAdd from './shared/ProjectStackAdd'
import ProjectStackList from './shared/ProjectStackList'

export default {
  name: 'ProjectRow',
  components: {
    ProjectHostname,
    ProjectToolbar,
    ProjectLocation,
    ProjectNote,
    ProjectStackAdd,
    ProjectStackList
  },
  props: {
    project: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
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
      projectPrefix: 'gb-' + this.$escapeIDAttr(this.id) + '-'
    }
  },
  computed: {
    isRunning () {
      return this.project.attributes.enabled
    }
  },
  methods: {

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

    async maybeSubmit (ev) {
      await this.$store.dispatch(ProjectActions.UPDATE_DETAILS, { projectId: this.id, attributes: this.$data })
    },

    async onRunStop () {
      if (this.project.attributes.stack && this.project.attributes.stack.length > 0) {
        this.isUpdating = true
        try {
          await this.$store.dispatch(ProjectActions.UPDATE_STATE, { 'project': this.project, 'isEnabled': !this.isRunning })
          this.isUpdating = false
        } catch (e) {
          console.error(e.message)
        }
      } else {
        this.showAlert('Please add some stacks first!')
      }
    },

    onExpandCollapseStack (stackId, isExpanded) {
      Vue.set(this.expandedStackIds, stackId, isExpanded ? -1 : 1)
    }
  }
}
</script>

<style scoped>
  .row--project {
    border-top: 1px solid #f3f3f3;
    margin-bottom: 1.5rem;
    vertical-align: top;
  }

  .alert {
    margin-left: -1.25rem;
    width: calc(100% + 2.5rem);
    margin-top: 0.5rem;
    margin-bottom: -14px;
    border-top-left-radius: 0;
    border-top-right-radius: 0;
    font-size: 12px;
    padding-top: 1rem;
    padding-left: 25px;
    padding-bottom: 1rem;
  }

  .alert-dismissible .close {
    padding: 0.5rem 0.75rem;
    right: 0px;
    top: 3px;
  }

  .hostname-input {
    max-width: 300px;
  }

  .input-group--stack {
    position: relative;
    top: -4px;
  }

  .input-group--stack,
  .input-group--note.is-collapsed {
    display: inline-flex;
    width: auto;
  }

  .project-stack-list,
  .project-note-list {
    display: inline-block;
  }

  >>> .toolbar-link--state {
    font-size: 18px;
    margin-top: 6px;
  }

  >>> .input-group .btn-outline-info {
    border-color: #ced4da;
  }

  >>> .input-group.is-collapsed .btn-outline-info {
    border-color: transparent;
  }

  .input-group--note:not(.is-collapsed) {
    position: absolute;
    width: 100%;
    top: 11px;
    left: 0;
  }

  >>> .input-group--note:not(.is-collapsed) .input-group-append {
    background-color: rgb(246,246,246)
  }

  >>> .stack-card.is-collapsible {
    background-color: #fafafa;
  }
</style>
<style>
  .row--project td {
    padding: 10px 15px 10px 0;
  }

  .td--state {
    max-width: 50px;
  }

  .td--hostname {
    max-width: 200px;
  }

  .td--location {
    max-width: 400px;
  }

  .td--notes {
    max-width: 300px;
    position: relative;
  }
</style>
