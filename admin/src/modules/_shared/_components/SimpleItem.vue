<template>
  <li>{{record.id}}</li>
</template>

<script>
export default {
  name: 'SimpleItem',
  props: {
    record: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>

</style>
