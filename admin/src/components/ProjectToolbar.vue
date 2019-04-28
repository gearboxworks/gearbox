<template>
  <ul class="toolbar-list">

    <li class="toolbar-item">
      <a target="_blank"
         href="#"
         :title="isRunning ? 'Stop all services' : 'Run all services'"
         v-b-tooltip.hover
         @click.prevent="onRunStop"
         class="toolbar-link toolbar-link--state"
      >
        <font-awesome-icon
          :icon="['fa', isRunning ? 'stop': 'play']"
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
      return 'gb-' + this.escAttr(this.id) + '-'
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

  .toolbar-link {
    float: right;
    font-size: 1.25rem;
    cursor: pointer;
    margin-left: 13px;
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
    margin-top: -3px;
    padding-top: 0;
  }

  .toolbar-link--state:hover {
    color: rgba(28, 76, 166, 0.91);
  }
</style>
