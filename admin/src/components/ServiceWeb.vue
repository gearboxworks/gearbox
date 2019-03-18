<template>
  <StackMemberFields
    :member="member"
    memberType="web"
  />
</template>

<script>
import StackMemberFields from './StackMemberFields.vue'
export default {
  name: 'StackMemberWeb',
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
