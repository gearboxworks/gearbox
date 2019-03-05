<template>
  <StackMemberFields
    :member="member"
    memberType="cache"
  />
</template>

<script>
import StackMemberFields from './StackMemberFields.vue'
export default {
  name: 'StackMemberCache',
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
