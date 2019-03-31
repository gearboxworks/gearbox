<template>
  <div class="project-stack-header">
    <button @click.prevent="removeProjectStack" class="js-remove-stack" aria-label="Remove this stack from project" title="Remove this stack from project">&times;</button>
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
    'projectHostname': {
      type: String,
      required: true
    },
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
    },
    removeProjectStack () {
      this.$store.dispatch('removeProjectStack', { 'projectHostname': this.projectHostname, 'stackName': this.stackName })
    }
  }
}
</script>

<style scoped>
  .js-remove-stack {
    float: left;
    margin-top: -3px;
    margin-left: -9px;
  }
</style>
