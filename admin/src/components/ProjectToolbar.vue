<template>
  <ul class="toolbar-list">

    <li class="toolbar-item">
      <a target="_blank"
         :href="`//${hostname}/`"
         :title="'Open Frontend'+ (isRunning ? '' : ' (not running)')"
         v-b-tooltip.hover
         :class="['toolbar-link', 'toolbar-link--frontend', {'is-disabled': !isRunning}]"
      >
        <font-awesome-icon
          :icon="['fa', 'home']"
        />
      </a>
    </li>

    <li class="toolbar-item">
      <a target="_blank"
         :href="`//${hostname}/wp-admin/`"
         :title="'Open Dashboard'+ (isRunning ? '' : ' (not running)')"
         v-b-tooltip.hover
         :class="['toolbar-link', 'toolbar-link--dashboard', {'is-disabled': !isRunning}]"
      >
        <font-awesome-icon
          :icon="['fa', 'tachometer-alt']"
        />
      </a>
    </li>

    <li class="toolbar-item">
      <a target="_blank"
         href="#"
         :title="isRunning ? 'Stop all services' : 'Run all services'"
         v-b-tooltip.hover
         @click.prevent="onRunStop"
         class="toolbar-link toolbar-link--state"
      >
        <font-awesome-icon
          :icon="['fa', isRunning ? 'stop-circle': 'play-circle']"
        />
      </a>
    </li>

  </ul>
</template>

<script>

export default {
  name: 'ProjectToolbar',
  props: {
    'project': {
      type: Object,
      required: true
    },
    'projectIndex': {
      type: Number,
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
    projectBase () {
      return this.escAttr(this.id) + '-'
    },
    isRunning () {
      return this.project.attributes.enabled
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    onRunStop () {
      console.log('onRunStop')
      this.$store.dispatch(
        'changeProjectState', { 'projectId': this.id, 'isEnabled': !this.isRunning }
      )
    }
  }
}
</script>

<style scoped>
  .toolbar-list {
    list-style: none;
    display: inline-block;
    float: right;
    padding-inline-start: 0;
    margin-bottom: 0;
    margin-top: -5px;
  }
  .toolbar-item {
    display: inline;
    float: left;
  }
  [data-icon="expand"]{
    height: 40px;
    width: 40px;
    margin-bottom: 10px;
    margin-top: 10px;
    color: silver;
  }

  .toolbar-link {
    float: right;
    font-size: 1.25rem;
    cursor: pointer;
    margin-left: 10px;
    color: rgba(42, 85, 130, 0.98);
  }

  .toolbar-link {
    padding-top: 3px;
  }

  .toolbar-link.is-disabled  {
    color: rgba(17, 56, 85, 0.42);
  }

  .toolbar-link--details {
    float: left;
  }

  .toolbar-link--state {
    font-size: 1.65rem;
    padding-top: 0;
  }

  [data-icon="play-circle"] path {
    fill: green;
  }

  [data-icon="stop-circle"] path {
    fill: red;
  }
</style>
