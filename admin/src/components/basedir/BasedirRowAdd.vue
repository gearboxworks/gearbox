<template>
  <b-form-row :class="{'form-row--basedir-add': true, 'is-updating': isUpdating}">
    <b-input-group :class="{'input-group--basedir-edit': true}" role="tabpanel">
      <b-form-input
        ref="create-basedir"
        type="text"
        class="basedir"
        placeholder="Input existing directory..."
        @keyup="touch('add')"
        @change="touch('add')"
        :state="inputState('add')"
        :tabindex="tabOffset"
      />
      <b-input-group-append>
        <b-button
          type="submit.prevent"
          :variant="touched['add'] ? 'outline-info': 'outline-secondary'"
          @click.prevent="onCreateBasedir"
          :disabled="!touched['add']"
          class="btn--add"
          title="Add new directory reference"
          :tabindex="tabOffset+1"
        >
          <font-awesome-icon :icon="['fa', 'plus']" />
        </b-button>
      </b-input-group-append>
    </b-input-group>

    <div class="invalid-feedback d-block">{{errors['add'] || '&nbsp;'}}</div>
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
      errors: {},
      touched: {},
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
    }
  },
  methods: {
    ...mapActions({
      'createBasedir': 'basedirs/create'
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
    touch (basedirId) {
      this.$set(this.touched, basedirId, true)
      this.$delete(this.errors, basedirId)
    },
    inputState (basedirId) {
      return typeof this.errors[basedirId] === 'undefined'
        ? null
        : (this.errors[basedirId] === 'no error')
    },
    onCreateBasedir () {
      const ctrl = this.$refs['create-basedir'].$el
      const basedir = ctrl.value

      const recordData = {
        'attributes': {
          basedir
        }
      }
      this.createBasedir(recordData)
        .then((res) => {
          console.log(res)
          this.$set(this.touched, 'add', true)
          this.$delete(this.errors, 'add')
          ctrl.value = ''
        })
        .catch(res => {
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

</style>
