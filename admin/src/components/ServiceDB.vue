<template>
  <StackMemberFields
    :member="member"
    memberType="database"
  />
</template>

<script>
import StackMemberFields from './StackMemberFields.vue'
export default {
  name: 'StackMemberDB',
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
