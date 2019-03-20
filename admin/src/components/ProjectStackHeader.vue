<template>
  <div class="project-stack-header">
    <strong>{{stackName.replace('gearbox.works/', '')}}</strong> (
    <span v-for="(service, key, index) in stackRoles" :key="key" :title="key.replace(service.stack+'/','')">
              <span v-if="index">,</span> {{service.program}} <small>{{serviceVersion(service.version)}}</small>
    </span>
    )
  </div>
</template>

<script>

export default {
  name: 'ProjectStackHeader',
  props: {
    'stackName': {
      type: String,
      required: true
    },
    'stackRoles': {
      type: Object,
      required: true
    }
  },
  methods: {
    serviceVersion (version) {
      var result = ''
      if (version.major) {
        result += version.major
        if (version.minor) {
          result += '.' + version.minor
          if (version.patch) {
            result += '.' + version.patch
          }
        }
      }
      return result
    }
  }
}
</script>

<style scoped>

</style>
