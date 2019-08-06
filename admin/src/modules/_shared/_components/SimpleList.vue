<template>
  <section>
    <ul
      v-if="!isLoading && records.length > 0"
      key="content-wrap"
      class="content-wrap is-ready"
    >
      <simple-item
        v-for="record in records"
        :record="record"
        :key="record.id"
      />
    </ul>
    <div
      v-else-if="isLoading"
      key="content-wrap"
      class="content-wrap .is-loading"
    >
      <font-awesome-icon
        icon="circle-notch"
        spin
      />
     &nbsp;<span>{{labels.loading}}</span>
    </div>
    <div
      v-else-if="isFailed"
      key="content-wrap"
      class="content-wrap .is-failed"
    >
      <font-awesome-icon
        icon="expand"
      />
      &nbsp;<span>{{isFailed}}</span>
    </div>
    <div
      v-else
      key="content-wrap"
      class="content-wrap .is-empty"
    >
      <font-awesome-icon
        icon="expand"
      />
      &nbsp;<span>{{labels.empty}}</span>
    </div>
  </section>
</template>

<script>

import SimpleItem from './SimpleItem'

export default {
  name: 'SimpleList',
  components: {
    SimpleItem
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
  .content-wrap {
    margin-left: 1rem;
    font-size: 150%;
  }
</style>
