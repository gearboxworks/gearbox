<template>
  <tr class="row--project">
    <td class="td--state">
      <project-toolbar :project="project" :projectIndex="projectIndex" @run-stop="onRunStop"></project-toolbar>
    </td>

    <td class="td--hostname">
      <project-hostname :project="project" :projectIndex="projectIndex" :is-multimodal="false"></project-hostname>
    </td>

    <td class="td--location">
      <project-location :project="project" :projectIndex="projectIndex" :is-multimodal="false"></project-location>
    </td>

    <td class="td--stack">
      <project-stack-list :project="project" :projectIndex="projectIndex" :cards-mode="true" v-if="project.attributes.stack.length"></project-stack-list>
      <project-stack-select :project="project" :projectIndex="projectIndex"></project-stack-select>
    </td>

    <td class="td--notes">
      <project-note-list :project="project" :projectIndex="projectIndex"></project-note-list>
      <project-note :project="project" :projectIndex="projectIndex"></project-note>
    </td>

  </tr>

</template>

<script>

import ProjectHostname from '../ProjectHostname'
import ProjectToolbar from '../ProjectToolbar'
import ProjectStackList from '../ProjectStackList'
import ProjectLocation from '../ProjectLocation'
import ProjectNote from '../ProjectNote'
import ProjectNoteList from '../ProjectNoteList'
import ProjectStackSelect from '../ProjectStack'

export default {
  name: 'ProjectRow',
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
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'warning'
    }
  },
  components: {
    ProjectHostname,
    ProjectToolbar,
    ProjectLocation,
    ProjectNote,
    ProjectStackSelect,
    ProjectStackList,
    ProjectNoteList
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
    onRunStop () {
      if (this.project.attributes.stack && this.project.attributes.stack.length > 0) {
        this.$store.dispatch(
          'changeProjectState', { 'projectId': this.id, 'isEnabled': !this.isRunning }
        )
      } else {
        this.showAlert('Please add some stacks first!')
      }
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

</style>
<style>
  .row--project td {
    padding: 10px 15px 10px 0;
  }

  .td--state {
    max-width: 50px;
  }

  .td--hostname {
    max-width: 300px;
  }

  .td--location {
    max-width: 400px;
  }

  .td--notes {
    max-width: 300px;
  }
</style>
