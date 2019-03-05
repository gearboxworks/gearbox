<template>
  <StackMemberFields
    :member="member"
    memberType="worker"
  />
</template>

<script>
import StackMemberFields from './StackMemberFields.vue'
export default {
  name: 'StackMemberWorker',
  components: {
    StackMemberFields
  },
  props: {
    member: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>

</style>
