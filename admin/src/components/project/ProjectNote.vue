<template>
  <div
    :class="{'clearfix': true, 'input-group--note': true, 'is-collapsed': isCollapsed, 'is-modified': isModified, 'is-updating': isUpdating, 'is-empty': !!notes}"
    role="tabpanel"
  >
    <b-form-textarea
      :ref="`${projectBase}note-input`"
      v-model="notes"
      placeholder="Add note..."
      class="notes-input"
      rows="6"
      v-if="!isCollapsed"
      :readonly="isUpdating"
      autocomplete="off"
      autofocus
      @keyup.esc="isCollapsed = true"
    />
    <b-button
      :variant="isCollapsed ? (notes ? 'outline-warning': 'outline-info') : 'outline-info'"
      :title="isCollapsed ? ( notes ? notes: 'Add a note' ) : (isModified ? 'Submit the note': 'Please enter some text first or Click to cancel')"
      v-b-tooltip.hover
      :class="{'btn--submit': true, 'btn--add': isCollapsed}"
      @click.prevent="onButtonClicked"
      :disabled="isUpdating"
    >
      <font-awesome-icon
        v-if="isUpdating"
        icon="circle-notch"
        spin
      />
      <font-awesome-icon
        v-else
        :icon="['fa', isCollapsed ? 'sticky-note': (isModified ? 'check': 'times')]"
      />
      <span v-if="!isUpdating && !notes">{{isCollapsed ? '+' : ''}}</span>
    </b-button>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
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
  data () {
    return {
      id: this.project.id,
      ...this.project.attributes,
      isCollapsed: true,
      isModified: false,
      isUpdating: false
    }
  },
  computed: {
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    }
  },
  methods: {
    ...mapActions(['addProjectNote']),
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    onButtonClicked () {
      if (this.isCollapsed) {
        this.isCollapsed = false
        // console.log(`${this.projectBase}note-input`, this.$refs[`${this.projectBase}note-input`])
        // this.$refs[`${this.projectBase}note-input`].focus()
      } else {
        if (this.isModified) {
          this.maybeSubmit()
        } else {
          this.isCollapsed = true
        }
      }
    },
    maybeSubmit () {
      this.isUpdating = true
      /**
       * TODO: deal with timestamp
       */
      this.addProjectNote(
        {
          projectId: this.id,
          text: this.notes
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
    /* top: -45px; */
    left: 0px;
    margin: -45px 0 7px 0;
  }

  .is-collapsed {
    height: 35px;
  }
  .btn-outline-info {
    border-color: #ced4da;
  }

  .btn--submit {
    float: right;
    margin-top: 10px;
    margin-bottom: 0;
  }

  .btn--add {
    position:relative;
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
