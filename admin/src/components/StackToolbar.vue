<template>
  <ul class="toolbar-list">

    <li class="toolbar-item" v-if="isWordPress" :key="hostname+'frontend'">
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

    <li class="toolbar-item" v-if="isWordPress" :key="hostname+'dashboard'">
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

    <li class="toolbar-item" v-if="isWordPress" :key="hostname+'adminer'">
      <a target="_blank"
         :href="`//${hostname}/wp-admin/`"
         :title="'Open Adminer'+ (isRunning ? '' : ' (not running)')"
         v-b-tooltip.hover
         :class="['toolbar-link', 'toolbar-link--database', {'is-disabled': !isRunning}]"
      >
        <font-awesome-icon
          :icon="['fa', 'database']"
        />
      </a>
    </li>

    <li class="toolbar-item" v-if="isWordPress" :key="hostname+'mailhog'">
      <a target="_blank"
         :href="`//${hostname}:4003`"
         :title="'Open Mailhog'+ (isRunning ? '' : ' (not running)')"
         v-b-tooltip.hover
         :class="['toolbar-link', 'toolbar-link--mail', {'is-disabled': !isRunning}]"
      >
        <font-awesome-icon
          :icon="['fa', 'envelope']"
        />
      </a>
    </li>

    <li class="toolbar-item" :key="hostname+'other'">
      <b-dropdown
        variant="link"
        no-caret
        title="More stack actions"
        v-b-tooltip.hover
      >
        <template slot="button-content"><font-awesome-icon
          :icon="['fa', 'ellipsis-v']"
        />
          <span class="sr-only">More actions</span>
        </template>
        <b-dropdown-item href="#" @click.prevent="removeProjectStack">Remove stack</b-dropdown-item>
      </b-dropdown>
    </li>
  </ul>
</template>

<script>

export default {
  name: 'StackToolbar',
  props: {
    'project': {
      type: Object,
      required: true
    },
    'projectIndex': {
      type: Number,
      required: true
    },
    'stackId': {
      type: String,
      required: true
    }
  },
  data () {
    return {
      isWordPress: this.stackId.indexOf('/wordpress') !== -1
    }
  },
  computed: {
    projectBase () {
      return 'gb-' + this.escAttr(this.project.id) + '-'
    },
    hostname () {
      return this.project.attributes.hostname
    },
    isRunning () {
      return this.project.attributes.enabled
    }
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },
    removeProjectStack (stackId) {
      this.$store.dispatch('removeProjectStack', { 'projectId': this.project.id, 'stackId': this.stackId })
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
    margin-right: -5px;
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
  .toolbar-link:hover {
    color: rgba(28, 76, 166, 0.91);
  }
  .toolbar-link {
    padding-top: 3px;
    /*opacity: 1;*/
    /*transition: opacity 400ms;*/
  }
  .toolbar-link:not(.is-disabled) {
    opacity: 0;
    animation: slide-in 300ms ease-in forwards;
    animation-iteration-count: 1;
    /*animation-delay: unset;*/
  }
  .toolbar-link.is-disabled  {
    color: rgba(17, 56, 85, 0.42);
    visibility: hidden;
    opacity: 0;
  }

  .toolbar-item:nth-child(4) .toolbar-link:not(.is-disabled) {
    animation-delay: 25ms;
  }
  .toolbar-item:nth-child(3) .toolbar-link:not(.is-disabled) {
    animation-delay: 50ms;
  }
  .toolbar-item:nth-child(2) .toolbar-link:not(.is-disabled) {
    animation-delay: 75ms;
  }
  .toolbar-item:nth-child(1) .toolbar-link:not(.is-disabled) {
    animation-delay: 100ms;
  }

  @keyframes slide-in {
    0% {transform: translateY(-10px); opacity: 0;}
    50% {transform: translateY(3px); opacity: 0.5;}
    100% { transform: translateY(0px); opacity: 1;}
  }

  .btn-link [data-icon="ellipsis-v"] path {
    fill: rgba(42, 85, 130, 0.98);
  }
  .btn-link:hover [data-icon="ellipsis-v"] path {
    fill: rgba(28, 76, 166, 0.91);
  }
</style>
