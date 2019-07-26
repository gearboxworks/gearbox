<template>
  <b-form-row :class="{'form-row--basedir-add': true, 'is-updating': isUpdating}">
    <b-input-group
      :class="{'input-group--basedir-edit': true}"
      role="tabpanel"
    >
      <b-form-input
        v-model="currentValue"
        type="text"
        class="basedir"
        placeholder="Input directory..."
        @keyup="touch('add')"
        @change="touch('add')"
        :state="inputState('add')"
        tabindex="0"
      />
      <b-input-group-append>
        <b-button
          type="submit.prevent"
          :variant="touched['add'] ? 'outline-info': 'outline-secondary'"
          @click.prevent="onAddBasedir"
          :disabled="!touched['add']"
          class="btn--add"
          :title="currentValue ? 'Add directory': 'First, input some directory!'"
          tabindex="0"
        >
          <font-awesome-icon :icon="['fa', 'plus']" />
        </b-button>
      </b-input-group-append>
    </b-input-group>

    <div :class="{confirmation: true, visible: notfound[currentValue]}">
      This dir does not exist. Would you like to create it?
      <a class="yes" @click.stop="createDir(currentValue)" title="Create directory">Yes</a>
      |
      <a class="no" @click.stop="notfound[currentValue]=0" title="Try a different dir">No</a>
    </div>

    <div class="invalid-feedback d-block">
      {{errors['add'] || '&nbsp;'}}
    </div>
  </b-form-row>

</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'BasedirRowEdit',
  props: {
    tabOffset: {
      type: Number,
      required: true
    }
  },
  data () {
    return {
      isUpdating: false,
      currentValue: '',
      errors: {},
      touched: {},
      notfound: {},
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'warning'
    }
  },
  computed: {
    ...mapGetters([
      'basedirBy'
    ])
  },
  methods: {
    ...mapActions({
      doCreateBasedir: 'basedirs/create',
      getDirectory: 'getDirectory'
    }),

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

    touch (basedirId) {
      if (this.currentValue) {
        this.$set(this.touched, basedirId, true)
        this.$delete(this.errors, basedirId)
      } else {
        this.$delete(this.touched, basedirId)
        this.$delete(this.errors, basedirId)
      }
    },

    inputState (basedirId) {
      return typeof this.errors[basedirId] === 'undefined'
        ? null
        : (this.errors[basedirId] === 'no error')
    },

    onAddBasedir () {
      const basedir = this.currentValue
      if (!basedir) {
        return
      }

      this.getDirectory({ 'dir': basedir })
        .then(r => r ? r.data : null)
        .then(response => {
          this.createDir(basedir)
        })
        .catch(e => {
        /**
         * TODO deal with a code which indicates that the dir is invalid! Maybe 409?
         */
          if (e.response.status === 404) {
            this.$set(this.notfound, basedir, 1)
          } else {
            this.$delete(this.notfound, basedir)
          }
        })
    },

    createDir (basedir) {
      const recordData = {
        'attributes': {
          basedir
        }
      }

      this.doCreateBasedir(recordData)
        .then((res) => {
          // console.log(res, this)
          this.$set(this.touched, 'add', true)
          this.$delete(this.errors, 'add')
          this.$delete(this.notfound, this.currentValue)
          this.currentValue = ''
        })
        .catch(res => {
          // console.log(res, this)
          this.$set(this.errors, 'add', res.data.errors[0].title || res.statusText)
          this.$delete(this.touched, 'add')
        })
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
  .is-invisible {
    visibility: hidden;
  }

  .btn--add {
    width: 3rem;
  }

  .confirmation{
    margin-top: 7px;
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
