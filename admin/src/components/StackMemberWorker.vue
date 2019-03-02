<template>
  <StackMemberFields :member="member" member_type="worker" />
</template>

<script>
import StackMemberFields from './StackMemberFields.vue';
export default {
  name: 'StackMemberWorker',
  components: {
    StackMemberFields
  },
  props: {
    member: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style scoped>

</style>
