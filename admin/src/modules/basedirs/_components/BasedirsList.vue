<template>
  <b-form :class="{'form--preferences': true, 'is-loading': isLoading}">
    <div
      v-if="!isLoading"
      key="basedirs-content"
      class="content-wrap is-ready"
    >
      <basedir-row-edit
        v-for="(basedir) in records"
        :key="basedir.id"
        :basedir="basedir"
        :is-deletable="basedir.id !== 'default'"
      />
    </div>
    <div
      v-else-if="!isLoading && records.length === 0"
      key="basedirs-content"
      class="content-wrap is-empty"
    >
      <font-awesome-icon
        icon="expand"
      />
      &nbsp;<span>{{labels.empty}}</span>
    </div>
    <div
      v-else-if="isLoading"
      key="basedirs-content"
      class="content-wrap is-loading"
    >
      <font-awesome-icon
        icon="circle-notch"
        spin
      />
      &nbsp;<span>{{labels.loading}}</span>
    </div>
    <div
      v-else
      key="basedirs-content"
      class="content-wrap is-failed"
    >
      <font-awesome-icon
        icon="expand"
      />
      &nbsp;<span>{{labels.isFailed}}</span>
    </div>

    <basedir-row-add :tab-offset="records.length*3" />

  </b-form>
</template>

<script>

import BasedirRowEdit from './BasedirRowEdit'
import BasedirRowAdd from './BasedirRowAdd'

export default {
  name: 'BasedirsList',
  components: {
    BasedirRowEdit,
    BasedirRowAdd
  },
  props: {
    records: {
      required: true,
      type: Array
    },

    isLoading: {
      required: false,
      type: Boolean,
      default: false
    },

    isFailed: {
      required: false,
      type: String,
      default: ''
    },

    labels: {
      required: false,
      type: Object,
      default () {
        return {
          title: 'Items',
          singular: 'Item',
          plural: 'Items',
          empty: 'Found no records',
          loading: 'Loading...',
          error: 'Failed to load.'
        }
      }
    }
  },
  data () {
    return {}
  }
}
</script>

<style scoped>
  .form--preferences {
    padding: 1rem;
  }
</style>
