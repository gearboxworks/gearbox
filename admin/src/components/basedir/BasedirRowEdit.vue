<template>
  <b-form-row
    :class="{'form-row--basedir-edit': true, 'is-updating': isUpdating, 'is-modified': isModified,'is-deleting': isDeleting}"
  >

    <b-input-group :class="{'input-group--basedir-edit': true}" role="tabpanel">

      <b-form-input
        :ref="`${basedir.id}-basedir`"
        :value="basedir.attributes.basedir"
        @keyup="touch"
        @change="touch"
        :state="isValid"
        class="basedir"
        type="text"
        required
        placeholder="Path"
        tabindex="tabOffset"
      />
      <b-input-group-append v-if="isModified">
        <b-button
          type="submit.prevent"
          :variant="isModified ? 'outline-info': 'outline-secondary'"
          :disabled="!isModified"
          @click.prevent="onUpdateBasedir"
          class="btn--update"
          title="Update directory reference"
          :tabindex="tabOffset+1"
        >
          <font-awesome-icon v-if="isUpdating" spin :icon="['fa', 'circle-notch']" />
          <font-awesome-icon v-else :icon="['fa', 'check']" />
        </b-button>
      </b-input-group-append>
    </b-input-group>

    <b-button-group class="button-group--extras">
      <b-button
        variant="outline-info"
        title="Copy to clipboard"
        v-b-tooltip.hover
        class="btn--copy-dir"
        @click.prevent="onCopyToClipboard"
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
      :title="isDeletable ? 'Delete reference to this directory': 'Cannot delete the last reference to projects directory' "
      class="btn--delete"
      :tabindex="tabOffset+2"
    >
      <font-awesome-icon v-if="isDeleting" spin :icon="['fa', 'circle-notch']" />
      <font-awesome-icon v-else :icon="['fa', 'trash-alt']" />
    </b-button>
    <div class="invalid-feedback d-block">{{errors || '&nbsp;'}}</div>
  </b-form-row>

</template>

<script>
import { mapGetters, mapActions } from 'vuex'

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
    },
    tabOffset: {
      type: Number,
      required: true
    }
  },
  data () {
    return {
      id: this.basedir.id,
      errors: '',
      isModified: false,
      isUpdating: false,
      isDeleting: false,
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'warning'
    }
  },
  components: {
  },
  computed: {
    ...mapGetters(['basedirBy']),
    ctrlBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    },
    isValid () {
      return (this.errors === '') ? null : false
    }
  },
  methods: {
    ...mapActions({
      'doUpdateBasedir': 'basedirs/update',
      'doDeleteBasedir': 'basedirs/delete'
    }),
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
    touch () {
      const basedir = this.$refs[this.id + '-basedir'].$el.value
      // console.log(basedir, this.basedir)
      this.isModified = basedir && (basedir !== this.basedir.attributes.basedir)
      this.errors = ''
    },
    onUpdateBasedir () {
      const basedir = this.$refs[this.id + '-basedir'].$el.value

      if (this.id && basedir) {
        const recordData = {
          id: this.id,
          type: 'basedirs',
          attributes: {
            basedir,
            nickname: this.id
          }
        }
        if (this.isModified) {
          this.isUpdating = true
          this.doUpdateBasedir(recordData).then((res) => {
            this.errors = ''
            this.isModified = false
            this.isUpdating = false
          }).catch(res => {
            this.isUpdating = false
            // console.log(res)
            this.$nextTick(function () {
              this.errors = res.data.errors[0].title || res.statusText
            })
          })
        }
      }
    },
    onDeleteBasedir () {
      this.isDeleting = true
      this.doDeleteBasedir({ id: this.id })
        .then((res) => {
          this.errors = ''
          this.isModified = false
          this.isDeleting = false
        })
        .catch((res) => {
          this.errors = res
          this.isDeleting = false
          // this.isModified = false
        })
    },
    onCopyToClipboard () {
      // @TODO implement copy to clipboard
      // @see https://github.com/Inndy/vue-clipboard2
      console.log('TODO: implement copy to clipboard')
    },
    onOpenDirectory () {
      // TODO call API method to open directory in file manager
      console.log('TODO: call API method to open directory in file manager')
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
</style>
