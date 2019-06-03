<template>
  <b-input-group
    :id="`${projectBase}note`"
    :class="{'input-group--note': true, 'is-collapsed': isCollapsed, 'is-modified': isModified, 'is-updating': isUpdating}"
    role="tabpanel"
  >
    <b-form-input
      :id="`${projectBase}note-input`"
      v-model="notes"
      placeholder="Add note..."
      class="notes-input"
      v-if="!isCollapsed"
      :readonly="isUpdating"
    />
    <b-input-group-append>
      <b-button
        variant="outline-info"
        :title="isCollapsed ? 'Add a note' : (isModified ? 'Submit the new note': 'Please enter some text first or Click to cancel')"
        v-b-tooltip.hover
        :id="`${projectBase}submit-note`"
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
        <span v-if="!isUpdating">{{isCollapsed ? '+' : ''}}</span>
      </b-button>
    </b-input-group-append>
  </b-input-group>
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
  .is-collapsed {
    height: 35px;
  }
  .btn-outline-info {
    border-color: #ced4da;
  }
  .btn--add {
    position:relative;
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

  .is-collapsed .btn-outline-info {
    border-color: transparent;
    border-top-left-radius: 0.25rem;
    border-bottom-left-radius: 0.25rem;
  }
</style>
