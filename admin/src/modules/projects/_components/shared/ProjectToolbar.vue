<template>
  <ul class="toolbar-list">

    <li class="toolbar-item">
      <a target="_blank"
         href="#"
         :title="isSwitching ? 'Switching state...': (isRunning ? 'Stop all services' : 'Run all services')"
         v-b-tooltip.hover
         @click.prevent="$emit('run-stop-project')"
         class="toolbar-link toolbar-link--state"
      >
        <font-awesome-icon
          v-if="isSwitching"
          key="status-icon"
          :icon="['fa', 'circle-notch']"
          spin
        />
        <font-awesome-icon
          v-else
          key="status-icon"
          :icon="['fa', isRunning ? 'stop': 'play']"
        />
      </a>
    </li>

  </ul>
</template>

<script>

export default {
  name: 'ProjectToolbar',
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {
    'isUpdating': {
      type: Boolean,
      required: true
    }
  },
  data () {
    return {
      id: this.project.id,
      hostname: this.project.attributes.hostname
    }
  },
  computed: {
    isRunning () {
      return this.project.attributes.enabled
    },

    isSwitching () {
      return this.isUpdating
    }
  },
  methods: {}
}
</script>

<style scoped>
  .toolbar-list {
    list-style: none;
    display: inline-block;
    float: right;
    padding-inline-start: 0;
    margin-bottom: 0;
    margin-top: 0;
  }
  .toolbar-item {
    display: inline;
    float: left;
  }

  .toolbar-link {
    float: right;
    font-size: 1.25rem;
    cursor: pointer;
    padding-left: 6px;
    padding-right: 6px;
    color: rgba(42, 85, 130, 0.98);
  }

  .toolbar-link {
    padding-top: 3px;
  }

  .toolbar-link.is-disabled  {
    color: rgba(17, 56, 85, 0.42);
  }

  .toolbar-link--state {
    font-size: 29px;
    margin-top: 1px;
    padding-top: 0;
  }

  .toolbar-link--state:hover {
    color: rgba(28, 76, 166, 0.91);
  }
</style>
