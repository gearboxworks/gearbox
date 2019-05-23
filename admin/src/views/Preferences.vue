<template>
  <div>
    <b-form>
      <b-container>
        <h1>Preferences</h1>
        <hr>
        <h2>Project Directories</h2>
        <h4>Known directories:</h4>
          <b-form-row
            v-for="basedir in basedirs.records"
            :key="basedir.id"
          >

            <b-form-input
              :id="`${basedir.id}-basedir`"
              :ref="`${basedir.id}-basedir`"
              :value="basedir.attributes.basedir"
              @keyup="touch(basedir.id)"
              @change="touch(basedir.id)"
              :state="inputState(basedir.id)"
              class="basedir"
              type="text"
              required
              placeholder="Path" />

            <b-button
              type="submit.prevent"
              :variant="touched[basedir.id] ? 'success': 'secondary'"
              :disabled="!touched[basedir.id]"
              @click.prevent="updateBasedir(basedir.id)"
              class="btn--update"
              title="Update directory reference"
            ><font-awesome-icon :icon="['fa', 'check-circle']" /></b-button>

            <b-button
              type="submit.prevent"
              :variant="basedir.id === defaultBasedir.id ? 'secondary':'warning'"
              @click.prevent="deleteBasedir(basedir.id)"
              :disabled="basedir.id === defaultBasedir.id"
              :title="basedir.id === defaultBasedir.id ? 'Cannot delete reference to the default directory' : 'Delete reference to this directory'"
              class="btn--delete"
            >
              <font-awesome-icon :icon="['fa', 'trash-alt']" />
            </b-button>
            <div class="invalid-feedback d-block">{{errors[basedir.id] || '&nbsp;'}}</div>
        </b-form-row>

        <b-form-row>
            <b-form-input
              ref="add-basedir"
              type="text"
              class="basedir"
              placeholder="Input existing directory..."
              @keyup="touch('add')"
              @change="touch('add')"
              :state="inputState('add')"
            />
            <b-button
              type="submit.prevent"
              :variant="touched['add'] ? 'success': 'secondary'"
              @click.prevent="addBasedir"
              :disabled="!touched['add']"
              class="btn--add"
              title="Add new directory reference"
            ><font-awesome-icon :icon="['fa', 'plus-circle']" /></b-button>
          <button class="btn btn--delete is-invisible"></button>
          <div class="invalid-feedback d-block">{{errors['add'] || '&nbsp;'}}</div>
        </b-form-row>

        <hr>

        <h4>Default directory:</h4>
        <b-form-select
          v-if="basedirs.records.length>1"
          id="default-basedir-input"
          value="default"
          @change="onChangeBasedir($event)"
        >
          <option disabled value="">Select directory...</option>
          <option v-for="basedir in basedirs.records" :value="basedir.id" :key="basedir.id">{{basedir.attributes.basedir}}</option>
        </b-form-select>
        <b-form-input
          v-else
          :readonly="true"
          :value="defaultBasedir ? defaultBasedir.attributes.basedir : ''"
        />

      </b-container>
    </b-form>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'

export default {
  name: 'Preferences',
  data () {
    return {
      errors: {},
      touched: {}
    }
  },
  computed: {
    ...mapState(['basedirs']),
    ...mapGetters(['basedirBy']),
    defaultBasedir () {
      return this.basedirBy('id', 'default')
    }
  },
  mounted () {
    this.$store.dispatch('basedirs/loadAll')
  },
  methods: {
    touch (basedirId) {
      this.$set(this.touched, basedirId, true)
      this.$delete(this.errors, basedirId)
    },
    inputState (basedirId) {
      return typeof this.errors[basedirId] === 'undefined'
        ? null
        : (this.errors[basedirId] === 'no error')
    },
    addBasedir () {
      const ctrl = this.$refs['add-basedir']['$el']
      const basedir = ctrl.value

      const recordData = {
        'attributes': {
          basedir
        }
      }
      this.$store.dispatch('basedirs/create', recordData)
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
    },
    updateBasedir (basedirId) {
      const basedir = this.$refs[basedirId + '-basedir'][0]['$el'].value

      if (basedirId && basedir) {
        const recordData = {
          id: basedirId,
          type: 'basedirs',
          attributes: {
            basedir,
            nickname: basedirId
          }
        }
        this.$store.dispatch('basedirs/update', recordData)
          .then((res) => {
            console.log(res)
            this.$set(this.errors, basedirId, 'no error')
            this.$delete(this.touched, basedirId)
          })
          .catch(res => {
            // console.log(res)
            this.$nextTick(function () {
              this.$set(this.errors, basedirId, res.data.errors[0].title || res.statusText)
            })
          })
      }
    },
    deleteBasedir (basedirId) {
      if (basedirId !== 'default') {
        if (basedirId) {
          this.$store.dispatch('basedirs/delete', { id: basedirId }).then((res) => {
            this.$delete(this.errors, basedirId)
            this.$delete(this.touched, basedirId)
          })
        }
      } else {
        // console.log('The default base directory will not be deleted.')
      }
    }
  }
}
</script>

<style scoped>
  .form-row {
    padding: 0.25rem;
  }
  .basedir {
    width: calc(100% - 7rem);
    display: inline-block;
  }
  .is-invisible {
    visibility: hidden;
  }

  .btn--add,
  .btn--update,
  .btn--delete {
    width: 3rem;
    margin-left: 0.5rem;
  }
</style>
