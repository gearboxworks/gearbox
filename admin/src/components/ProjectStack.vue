<template>
  <div role="tablist" class="project-stack">
    <b-card
      v-for="(stackName, index) in this.$store.state.gearStacks"
      :key="index"
      no-body
      class="mb-1"
    >
      <b-card-header header-tag="header" class="p-1" role="tab">
        <b-button block href="#" v-b-toggle="project_base + '_accordion_' + index" variant="info">{{stackName}}</b-button>
      </b-card-header>
      <b-collapse :id="project_base + '_accordion_' + index" visible :accordion="project_base + '_accordion'" role="tabpanel">
        <b-card-body>
          <b-form>
            <b-form-group
              v-for="(service, serviceRole) in stackServices(stackName)"
              :key="project_base + escAttr(serviceRole)"
              :label="stackRoles(stackName)[serviceRole].short_label"
              :label-for="project_base + escAttr(serviceRole)+'_input'"
              :description="stackRoles(stackName)[serviceRole].label"
              label-cols-sm="4"
              label-cols-lg="3"
            >
              <b-form-select
                type="text"
                :id="project_base + escAttr(serviceRole)+'_input'"
                :options="mapOptions(service.options)"
                required
                placeholder="" />
            </b-form-group>
          </b-form>
        </b-card-body>
      </b-collapse>
    </b-card>
  </div>
</template>

<script>

// import ServiceWeb from './ServiceWeb.vue'
// import ServiceDB from './ServiceDB.vue'
// import ServiceCache from './ServiceCache.vue'
// import ServiceProcessVM from './ServiceProcessVM.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'ProjectStack',
  props: {
    'projectHostname': {
      type: String
    },
    'projectStack': {
      type: Object
    }
  },
  computed: {
    ...mapGetters(['stackRoles', 'stackServices']),
    project_base () {
      return this.escAttr(this.projectHostname)
    },
    project () {
      return this.$store.getters.projectBy('hostname', this.projectHostname)
    },
    stack () {
      return this.project.stack
    },
    serviceName () {
      return this.$route.params.service
    },
    service () {
      const serviceName = this.serviceName
      return serviceName ? this.stack[ serviceName ] : ''
    },
    stacks () {
      return this.$store.state.gearStacks
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '_').replace(/\./g, '_')
    },
    mapOptions (options) {
      const result = []
      for (const value in options) {
        result.push({
          value,
          text: options[value]
        })
      }
      return result
    }
  },
  components: {
  }
}
</script>

<style scoped>

</style>
