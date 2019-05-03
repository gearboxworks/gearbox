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

            <b-form-input
              :id="`${basedir.id}-basedir`"
              :ref="`${basedir.id}-basedir`"
              :value="basedir.attributes.basedir"
              class="basedir"
              type="text"
              required
              placeholder="Path" />

            <b-button
              type="submit.prevent"
              variant="success"
              @click.prevent="updateBasedir(basedir.id)"
              class="btn--update"
              title="Update directory reference"
            ><font-awesome-icon :icon="['fa', 'check-circle']" /></b-button>
        </b-form-row>

        <b-form-row>
            <button class="btn btn--delete is-invisible"></button>
            <b-form-input
              ref="new-basedir"
              type="text"
              class="basedir"
              placeholder="Input existing directory..."
            />
            <b-button
              type="submit.prevent"
              variant="success"
              @click.prevent="addBasedir"
              class="btn--add"
              title="Add new directory reference"
            ><font-awesome-icon :icon="['fa', 'plus-circle']" /></b-button>
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
      baseDirName: '',
      baseDirPath: ''
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
    addBasedir () {
      const recordData = {
        'attributes': {
          basedir: this.$refs['new-basedir']['$el'].value
        }
      }
      this.$store.dispatch('basedirs/create', recordData).then(() => {
        const basedir = this.$store.getters['basedirs/lastCreated']
        console.log(basedir)
      })
    },
    updateBasedir (basedirId) {
      const basedir = this.$refs[basedirId + '-basedir'][0]['$el'].value

      if (basedirId && basedir) {
        const recordData = {
          id: basedirId,
          attributes: {
            basedir
          }
        }
        // const record = this.basedirBy(basedirId)
        // record.attributes.basedir = basedir
        console.log(basedirId, recordData)
        this.$store.dispatch('basedirs/update', recordData).then((res) => {
          console.log(res)
        })
      }
    },
    deleteBasedir (basedirId) {
      console.log('deleteBasedir', basedirId)
      if (basedirId !== 'default') {
        if (basedirId) {
          this.$store.dispatch('basedirs/delete', { id: basedirId }).then((res) => {
            console.log(res)
          })
        }
      } else {
        console.log('The default base directory will not be deleted.')
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
    margin-right: 0.5rem;
    margin-left: 0.5rem;
  }
  .is-invisible {
    visibility: hidden;
  }
  .btn--delete {
    width: 3rem;
  }
  .btn--add,
  .btn--update {
    width: 3rem;
  }
</style>
