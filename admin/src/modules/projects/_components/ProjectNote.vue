<template>
  <div
    :class="{'clearfix': true, 'input-group--note': true, 'is-collapsed': isCollapsed, 'is-expanded': !isCollapsed, 'is-modified': isModified, 'is-updating': isUpdating, 'is-empty': !notes, 'is-editing': isEditing}"
    role="tabpanel"
  >

    <b-form-textarea
      :ref="`${projectBase}note`"
      v-model="notes"
      placeholder="Add note..."
      class="notes-input"
      rows="7"
      v-if="!isCollapsed && isEditing"
      :readonly="isUpdating"
      autocomplete="off"
      autofocus
      @keyup.esc="isCollapsed = true"
    />

    <div class="notes-wrapper" v-if="!isCollapsed && !isEditing && notes && !isModified">
      <VueMarkdown>{{notes}}</VueMarkdown>
    </div>

    <b-button
      :id="`${projectBase}button`"
      :variant="isCollapsed ? (notes ? 'outline-warning': 'outline-info') : (isEditing ? 'outline-info': 'outline-warning')"
      :class="{'btn--submit': true, 'btn--add': isCollapsed}"
      @click.prevent="onButtonClicked"
      :disabled="isUpdating"
    >
      <font-awesome-icon
        :icon="['fa', isCollapsed ? 'sticky-note': 'times']"
      />
      <span v-if="!isUpdating && !notes">{{isCollapsed ? '+' : ''}}</span>
    </b-button>

    <b-tooltip triggers="hover" :target="`${projectBase}button`" placement="top">
      {{isCollapsed ? ( notes ? 'View notes' : 'Add notes' ) : 'Cancel changes!'}}
    </b-tooltip>

    <b-button
      :id="`${projectBase}button-edit`"
      v-show="!isCollapsed"
      :variant="isEditing ? 'outline-info': 'outline-warning'"
      :class="{'btn--edit': true}"
      @click.prevent="onButtonEditClicked"
      :disabled="isUpdating"
    >
      <font-awesome-icon
        v-if="isUpdating"
        icon="circle-notch"
        spin
      />
      <font-awesome-icon
        v-else
        :icon="['fa', isEditing ? 'check': 'pencil-alt' ]"
      />
    </b-button>

    <b-tooltip triggers="hover" :target="`${projectBase}button-edit`" placement="top">
      {{isEditing ? ( isModified ? (isUpdating ? 'Updating...' : 'Save changes') : 'Make some changes first' ) : 'Edit notes'}}
    </b-tooltip>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import VueMarkdown from 'vue-markdown'

export default {
  name: 'ProjectNote',
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
    VueMarkdown
  },
  data () {
    return {
      id: this.project.id,
      notes: this.project.attributes.notes,
      isCollapsed: true,
      isModified: false,
      isUpdating: false,
      isEditing: false
    }
  },
  computed: {
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    }
  },
  methods: {
    ...mapActions({ updateProjectNotes: 'projects/updateNotes' }),
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    onButtonClicked () {
      if (this.isCollapsed) {
        this.isCollapsed = false
        this.isEditing = !this.notes
      } else {
        this.isCollapsed = true
      }
    },
    onButtonEditClicked () {
      if (this.isEditing && !this.isModified) {
        return
      }

      if (!this.isEditing) {
        // console.log(`${this.projectBase}note-input`, this.$refs[`${this.projectBase}note-input`])
        // this.$refs[`${this.projectBase}note-input`].focus()
        this.isEditing = true
      } else if (this.isModified) {
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
        this.isCollapsed = true
        this.isModified = false
        this.isUpdating = false
      })
    }
  },
  watch: {
    notes: function (val, oldVal) {
      this.isModified = !!val
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
