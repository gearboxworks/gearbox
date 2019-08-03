<template>
  <b-form-row
    :class="{'form-row--basedir-edit': true, 'is-updating': isUpdating, 'is-modified': isModified,'is-deleting': isDeleting}"
  >

    <b-input-group
      :class="{'input-group--basedir-edit': true}"
      role="tabpanel"
    >
      <b-form-input
        v-model="currentValue"
        @keyup="touch"
        @change="touch"
        :state="isValid"
        class="basedir"
        type="text"
        required
        placeholder="Path"
        tabindex="0"
      />
      <b-input-group-append
        v-if="isModified"
      >
        <b-button
          type="submit.prevent"
          :variant="isModified ? 'outline-info': 'outline-secondary'"
          :disabled="!isModified"
          @click.prevent="onUpdateBasedir"
          class="btn--update"
          title="Update directory reference"
          :tabindex="0"
        >
          <font-awesome-icon
            v-if="isUpdating"
            key="status-icon"
            spin
            :icon="['fa', 'circle-notch']"
          />
          <font-awesome-icon
            v-else
            key="status-icon"
            :icon="['fa', 'check']"
          />
        </b-button>
      </b-input-group-append>
    </b-input-group>

    <b-button-group class="button-group--extras">
      <b-button
        variant="outline-info"
        title="Copy to clipboard"
        v-b-tooltip.hover
        class="btn--copy-dir"
        v-clipboard:copy="currentValue"
        v-clipboard:success="onCopyToClipboard"
      >
        <font-awesome-icon
          :icon="['fa', 'clone']"
        />
      </b-button>
      <b-button
        variant="outline-info"
        title="Open in file manager"
        v-b-tooltip.hover
        class="btn--open-dir"
        @click.prevent="onOpenDirectory"
      >
        <font-awesome-icon
          :icon="['fa', 'folder-open']"
        />
      </b-button>
    </b-button-group>

    <b-button
      type="submit.prevent"
      :variant="basedir.id === isDeletable ? 'outline-secondary':'outline-warning'"
      @click.prevent="onDeleteBasedir"
      :disabled="!isDeletable"
      :title="isDeletable ? 'Delete this directory': 'Cannot delete the default directory' "
      class="btn--delete"
      :tabindex="0"
    >
      <font-awesome-icon
        v-if="isDeleting"
        key="status-icon"
        spin
        :icon="['fa', 'circle-notch']"
      />
      <font-awesome-icon
        v-else
        key="status-icon"
        :icon="['fa', 'trash-alt']"
      />
    </b-button>

    <div v-if="!errors" :class="{confirmation: true, visible: notfound[currentValue]}">
      This dir does not exist. Would you like to create it?
      <a class="yes" @click.stop="updateBasedir(currentValue)" title="Create directory">Yes</a>
      |
      <a class="no" @click.stop="notfound[currentValue]=0" title="Try a different dir">No</a>
    </div>

    <div
      v-if="errors"
      class="invalid-feedback d-block"
    >
      {{errors}}
    </div>
  </b-form-row>

</template>

<script>
import { BasedirActions } from '../_store/method-names'

export default {
  name: 'BasedirRowEdit',
  props: {
    basedir: {
      type: Object,
      required: true
    },
    isDeletable: {
      type: Boolean,
      required: true,
      default: true
    }
  },
  data () {
    return {
      currentValue: this.basedir.attributes.basedir,
      errors: '',
      notfound: {},
      isModified: false,
      isUpdating: false,
      isDeleting: false,
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'warning'
    }
  },
  computed: {
    isValid () {
      return (this.errors === '') ? null : false
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

    touch () {
      const basedir = this.currentValue
      // console.log(this.currentValue, this.basedir.attributes.basedir)
      this.isModified = basedir && (basedir !== this.basedir.attributes.basedir)
      this.errors = ''
    },

    async onUpdateBasedir () {
      const basedir = this.currentValue

      if (!this.basedir.id || !basedir) {
        return
      }

      try {
        const statusCode = await this.$store.dispatch(BasedirActions.CHECK_DIRECTORY, basedir)
        if (statusCode > 200) {
          this.$set(this.notfound, basedir, 1)
        } else {
          this.$delete(this.notfound, basedir)
          this.updateBasedir(basedir)
        }
      } catch (e) {
        console.log(e)
        /**
         * TODO deal with a code which indicates that the dir is invalid! Maybe 409?
         */
      }
    },

    async updateBasedir (basedir) {
      if (!this.isModified) {
        return
      }
      const recordData = {
        id: this.basedir.id,
        type: 'basedirs',
        attributes: {
          basedir,
          nickname: this.basedir.id
        }
      }
      try {
        this.isUpdating = true

        await this.$store.dispatch(BasedirActions.UPDATE, { record: this.basedir, recordData })

        // this.errors = ''
        this.isModified = false
        this.isUpdating = false
        this.$delete(this.notfound, this.currentValue)
        // this.currentValue = ''
        // console.log('success', res)
      } catch (res) {
        this.isUpdating = false
        // console.log(res)
        // this.errors = res.data.errors[0].title || res.statusText
        console.error('fail', res)
      }
    },

    async onDeleteBasedir () {
      try {
        this.isDeleting = true

        await this.$store.dispatch(BasedirActions.DELETE, this.basedir.id)

        this.errors = ''
        this.isModified = false
        this.isDeleting = false
      } catch (res) {
        this.errors = res
        this.isDeleting = false
        // this.isModified = false
      }
    },

    onCopyToClipboard (e) {
      /**
       * TODO: Show something to the user to indicate that it was copied successfully
       */
      console.log('Copied to clipboard:', e.text)
    },

    onOpenDirectory () {
      this.$store.dispatch(BasedirActions.OPEN_DIRECTORY, this.currentValue)
    }
  }
}
</script>

<style scoped>
  .form-row {
    padding: 0.25rem;
    max-width: 35rem;
  }
  .input-group {
    width: calc(100% - 10rem);
  }

  .btn-group .btn,
  .btn--update,
  .btn--delete {
    width: 3rem;
  }

  .btn-group,
  .btn--delete {
    margin-left: 0.5rem;
  }

  .confirmation{
    margin-top: 7px;
    margin-bottom: 7px;
    font-size: 93%;
    transition: opacity 400ms ease-in;
    opacity: 0;
    color: var(--gray)
  }
  .confirmation.visible{
    opacity: 1;
  }
  .no, .yes {
    cursor: pointer;
  }

  .yes {
    color: var(--success) !important;
  }

  .no {
    color: var(--warning) !important;
  }
</style>
