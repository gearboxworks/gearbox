<template>
  <div
    :class="{'clearfix': true, 'input-group--note': true, 'is-collapsed': isCollapsed, 'is-expanded': !isCollapsed, 'is-modified': isModified, 'is-updating': isUpdating, 'is-empty': !notes, 'is-editing': isEditing}"
    :id="`$(projectPrefix)panel-notes`"
    role="tabpanel"
  >

    <b-form-textarea
      v-if="!isCollapsed && isEditing"
      :ref="`${projectPrefix}note`"
      v-model="notes"
      placeholder="Add note..."
      class="notes-input"
      rows="7"
      :readonly="isUpdating"
      autocomplete="off"
      autofocus
      @keyup.esc="isCollapsed = true"
      role="tab"
      aria-controls="$(projectPrefix)panel-notes`"
    />

    <div
      v-if="!isCollapsed && !isEditing && notes"
      class="notes-wrapper"
    >
      <VueMarkdown>{{notes}}</VueMarkdown>
    </div>

    <b-button
      :id="`${projectPrefix}button`"
      :variant="isCollapsed ? (notes ? 'outline-warning': 'outline-info') : (isEditing ? 'outline-info': 'outline-warning')"
      :class="{'btn--submit': true, 'btn--add': isCollapsed}"
      @click.prevent="onSwitchMode"
      :disabled="isUpdating"
    >
      <font-awesome-icon
        :icon="['fa', isCollapsed ? 'sticky-note': 'times']"
      />
      <span v-if="!isUpdating && !notes">{{isCollapsed ? '+' : ''}}</span>
    </b-button>

    <b-tooltip
      triggers="hover"
      :target="`${projectPrefix}button`"
      placement="top"
    >
      {{isCollapsed ? ( notes ? 'View notes' : 'Add notes' ) : ((isEditing && isModified) ? 'Cancel changes!' : 'Hide notes')}}
    </b-tooltip>

    <b-button
      :id="`${projectPrefix}notes-delete`"
      v-show="!isCollapsed && isEditing && (project.attributes.notes || deletedNotes || isRestoring)"
      variant="outline-warning"
      :class="{'btn--delete': true}"
      @click.prevent="onDelete"
      :disabled="isUpdating"
    >
      <font-awesome-icon
        v-if="isDeleting || isRestoring"
        key="trash-icon"
        icon="circle-notch"
        spin
      />
      <font-awesome-icon
        v-else
        key="trash-icon"
        :icon="['fa', deletedNotes ? 'trash-restore-alt': 'trash-alt' ]"
      />
    </b-button>

    <b-tooltip
      triggers="hover"
      :target="`${projectPrefix}notes-delete`"
      placement="top"
    >
      {{isDeleting ? 'Deleting...' : (isRestoring ? 'Restoring...' : (deletedNotes ? 'Restore notes' : 'Delete notes'))}}
    </b-tooltip>

    <b-button
      :id="`${projectPrefix}notes-edit`"
      v-show="!isCollapsed"
      :variant="isEditing ? 'outline-info': 'outline-warning'"
      :class="{'btn--edit': true}"
      @click.prevent="onEdit"
      :disabled="isUpdating"
    >
      <font-awesome-icon
        v-if="isUpdating && notes && !isRestoring && !isDeleting"
        key="status-icon"
        icon="circle-notch"
        spin
      />
      <font-awesome-icon
        v-else
        key="status-icon"
        :icon="['fa', isEditing ? 'check': 'pencil-alt' ]"
      />
    </b-button>

    <b-tooltip
      :target="`${projectPrefix}notes-edit`"
      triggers="hover"
      placement="top"
    >
      {{isEditing ? ( isModified ? (isUpdating ? 'Updating...' : 'Save changes') : 'Make some changes first' ) : 'Edit notes'}}
    </b-tooltip>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import VueMarkdown from 'vue-markdown'

export default {
  name: 'ProjectNote',
  components: {
    VueMarkdown
  },
  inject: ['project', 'projectPrefix'],
  props: {},
  data () {
    return {
      id: this.project.id,
      notes: this.project.attributes.notes,
      isCollapsed: true,
      isModified: false,
      isUpdating: false,
      isEditing: false,
      isRestoring: false,
      isDeleting: false,
      deletedNotes: ''
    }
  },
  computed: {},
  watch: {
    notes: function (val, oldVal) {
      this.isModified = val !== this.project.attributes.notes
    }
  },
  methods: {
    ...mapActions({
      updateProjectNotes: 'projects/updateNotes'
    }),
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    onSwitchMode () {
      if (this.isCollapsed) {
        this.isCollapsed = false
        this.isEditing = !this.notes
      } else {
        if (this.isModified) {
          this.notes = this.project.attributes.notes
        } else {
          this.isCollapsed = true
        }
      }
    },
    onEdit () {
      if (this.isEditing && !this.isModified) {
        return
      }

      if (!this.isEditing) {
        // console.log(`${this.projectPrefix}note-input`, this.$refs[`${this.projectPrefix}note-input`])
        // this.$refs[`${this.projectPrefix}note-input`].focus()
        this.isEditing = true
      } else if (this.isModified) {
        this.maybeSubmit()
      }
    },

    onDelete () {
      if (this.isEditing) {
        if (this.deletedNotes) {
          this.notes = this.deletedNotes
          this.deletedNotes = ''
          this.isRestoring = true
        } else {
          this.deletedNotes = this.notes
          this.notes = ''
          this.isDeleting = true
        }
        this.maybeSubmit()
      }
    },

    maybeSubmit () {
      this.isUpdating = true
      this.updateProjectNotes(
        {
          project: this.project,
          notes: this.notes
        }
      ).then(() => {
        if (!this.isDeleting && !this.isRestoring) {
          this.isCollapsed = true
        }
        this.isModified = false
        this.isUpdating = false
        this.isDeleting = false
        this.isRestoring = false
      })
    }
  }
}
</script>
<style scoped>

  .input-group--note:not(.is-collapsed) {
    position: relative;
    width: 100%;
    z-index: 3;
    margin: -45px 0 7px 0;
  }

  .notes-wrapper {
    display: block;
    min-height: 182px;
    max-height: 182px;
    overflow-y: auto;
    background-color: #f7f1c3;
    padding: 10px 40px 10px 12px;
    border-radius: 5px;
    border: 1px solid #ffc107;
  }

  .notes-wrapper, textarea {
    float: left;
    width: calc(100% - 50px);
  }

  .btn--submit,
  .btn--edit {
    margin-top: 0;
    margin-left: 8px;
    width: 38px;
  }

  .is-collapsed .btn--submit {
    margin-top: 10px;
  }

  .is-expanded .btn--submit {
    margin-bottom: 10px;
  }

  .is-expanded .btn--submit,
  .is-expanded .btn--edit {
    float: right;
  }

  .is-editing .btn--edit {
    transition: opacity 400ms ease-in;
    opacity: 1;
    position: absolute;
    bottom: 0;
    right: 0;
  }

  .is-editing .btn--edit svg{
    margin-left: -2px;
  }

  .is-expanded:not(.is-editing) .btn--edit {
    opacity: 1;
  }

  .is-editing:not(.is-modified) .btn--edit {
    opacity: 0.5;
  }

  .is-collapsed {
    height: 35px;
  }

  .btn-outline-info {
    border-color: #ced4da;
  }

  .btn--add {
    position: relative;
    top: -10px;
  }

  .btn--delete {
    position: absolute;
    transition: opacity 400ms ease-in;
    top: 50%;
    right: -1px;
    transform: translateY(-50%);
  }

  .btn--add svg {
    position: relative;
    left: -2px;
    top: 2px;
  }

  .btn--add span {
    position: absolute;
    right: 6px;
    font-size: 17px;
    top: -2px;
  }

  .is-collapsed .btn-outline-warning,
  .is-collapsed .btn-outline-info {
    border-color: transparent;
    border-top-left-radius: 0.25rem;
    border-bottom-left-radius: 0.25rem;
  }
</style>
