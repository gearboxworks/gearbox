<template>
  <div :class="'stack-member stack-member--'+ member_type">
    <h3>{{ member.label }}</h3>
    <el-radio-group size="small">
      <el-radio
        v-for="example in member.examples"
        :key="example"
        label="example"
        border
        disabled
      >
        {{ example }}
      </el-radio>
    </el-radio-group>
    <dl>
      <dt>Name:</dt>
      <dd>{{ member.name ? member.name : '&ndash;' }}</dd>
      <dt>Short label:</dt>
      <dd>{{ member.short_label ? member.short_label : '&ndash;' }}</dd>
      <dt>Member type:</dt>
      <dd>{{ member.member_type ? member.member_type : '&ndash;' }}</dd>
      <dt>Optional:</dt>
      <dd>{{ (true === member.optional) ? 'true' : (false === member.optional) ? 'false' : '&ndash;' }}</dd>
      <dt>Stack:</dt>
      <dd>{{ member.stack ? member.stack: '&ndash;' }}</dd>
    </dl>
  </div>
</template>

<script>
export default {
  name: 'StackMemberFields',
  props: {
    memberType: {
      type: String,
      required: true
    },
    member: {
      type: Object,
      required: true
    }
  }
}
</script>

<style scoped>

</style>
