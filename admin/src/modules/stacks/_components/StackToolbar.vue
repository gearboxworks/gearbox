<template>
  <div class="toolbar-list">
    <div
      class="toolbar-item toolbar-item--more"
      :key="projectPrefix + stackId + 'more'"
    >
      <font-awesome-icon
        v-if="isDeleting"
        key="action-menu"
        spin
        :icon="['fa', 'circle-notch']"
        class="spinner"
      />
      <b-dropdown
        class="dropdown--more"
        v-else
        key="action-menu"
        variant="link"
        no-caret
        title="More stack actions"
        v-b-tooltip.hover
      >
        <template slot="button-content">
          <span tabindex="0">
            <font-awesome-icon :icon="['fa', 'ellipsis-v']" />
            <span class="sr-only">More actions</span>
          </span>
        </template>
        <b-dropdown-item @click.prevent="onRemoveProjectStack">Remove stack</b-dropdown-item>
      </b-dropdown>
    </div>

    <transition
      v-if="isWordPress"
      name="icons"
      tag="ul"
    >
      <li
        v-if="isRunning"
        :class="['toolbar-item', 'toolbar-item--frontend']"
        :key="projectPrefix + stackId + 'frontend'"
      >
        <a target="_blank"
           :href="`//${hostname}/`"
           :title="'Open Frontend'+ (isRunning ? '' : ' (not running)')"
           v-b-tooltip.hover
           class="toolbar-link"
        >
          <font-awesome-icon
            :icon="['fa', 'home']"
          />
        </a>
      </li>
    </transition>

    <transition
      v-if="isWordPress"
      name="icons"
      tag="ul"
    >
      <li
        v-if="isRunning"
        :class="['toolbar-item', 'toolbar-item--dashboard']"
        :key="projectPrefix + stackId + 'dashboard'"
      >
        <a target="_blank"
           :href="`//${hostname}/wp-admin/`"
           :title="'Open Dashboard' + (isRunning ? '' : ' (not running)')"
           v-b-tooltip.hover
           class="toolbar-link"
        >
          <font-awesome-icon
            :icon="['fa', 'tachometer-alt']"
          />
        </a>
      </li>
    </transition>

    <transition
      v-if="isWordPress"
      name="icons"
      tag="ul"
    >
      <li
        v-if="isRunning"
        :class="['toolbar-item', 'toolbar-item--adminer']"
        :key="projectPrefix + stackId + 'adminer'"
      >
        <a target="_blank"
           :href="`//${hostname}/wp-admin/`"
           :title="'Open Adminer' + (isRunning ? '' : ' (not running)')"
           v-b-tooltip.hover
           class="toolbar-link"
        >
          <font-awesome-icon
            :icon="['fa', 'database']"
          />
        </a>
      </li>
    </transition>

    <transition
      v-if="isWordPress"
      name="icons"
      tag="ul"
    >
      <li
        v-if="isRunning"
        :class="['toolbar-item', 'toolbar-item--frontend']"
        :key="projectPrefix + stackId + 'mailhog'"
      >
        <a target="_blank"
           :href="`//${hostname}:4003`"
           :title="'Open Mailhog' + (isRunning ? '' : ' (not running)')"
           v-b-tooltip.hover
           class="toolbar-link"
        >
          <font-awesome-icon
            :icon="['fa', 'envelope']"
          />
        </a>
      </li>
    </transition>

  </div>
</template>

<script>
import StoreMethodTypes from '../../projects/_store/public-types'
const { ActionTypes: ProjectActions } = StoreMethodTypes

export default {
  name: 'StackToolbar',
  inject: [
    'project',
    'projectPrefix'
  ],
  props: {
    'stackId': {
      type: String,
      required: true
    }
  },
  data () {
    return {
      isWordPress: this.stackId.indexOf('/wordpress') !== -1,
      isDeleting: false
    }
  },
  computed: {
    hostname () {
      return this.project.attributes.hostname
    },

    isRunning () {
      return this.project.attributes.enabled
    }
  },
  methods: {

    async onRemoveProjectStack (stackId) {
      if (this.project.attributes.enabled) {
        this.$emit('show-alert', 'Cannot remove stack while the project is running!')
      } else {
        this.isDeleting = true
        try {
          await this.$store.dispatch(
            ProjectActions.REMOVE_STACK,
            {
              project: this.project,
              stackId: this.stackId
            }
          )
        } catch (e) {
          console.error(e.message)
        }
      }
    }
  }
}
</script>

<style scoped>
  .spinner {
    margin-right: 4px;
    margin-top: 6px;
  }

  .dropdown--more >>> .btn {
    padding: 0;
  }

  .dropdown--more >>> .btn > span {
    padding: 5px 8px 5px 5px;
  }
  .toolbar-list {
    list-style: none;
    display: inline-block;
    float: right;
    padding-inline-start: 0;
    margin-bottom: 0;
    margin-top: -5px;
    margin-right: -5px;
  }
  .toolbar-list ul{
    float: right;
    margin-bottom: 0;
    margin-top: 2px;
    padding-inline-start: 0;
  }
  .toolbar-item {
    display: inline-block;
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

  /*.toolbar-item:not(.is-disabled) {*/
  /*  opacity: 0;*/
  /*  animation: slide-in 300ms ease-in forwards;*/
  /*  animation-iteration-count: 1;*/
  /*  !*animation-delay: unset;*!*/
  /*}*/
  /*.toolbar-item.is-disabled  {*/
  /*  color: rgba(17, 56, 85, 0.42);*/
  /*  visibility: hidden;*/
  /*  opacity: 0;*/
  /*}*/

  .toolbar-item--more {
    margin-left: 2px;
    margin-right: 2px;
    margin-top: -3px;
    float: right;
  }

  .toolbar-link {
    color: rgba(42, 85, 130, 0.98);
  }
  .toolbar-link:hover {
    color:rgba(28, 76, 166, 0.91);
  }

  .btn-link [data-icon="ellipsis-v"] path {
    fill: rgba(42, 85, 130, 0.98);
  }
  .btn-link:hover [data-icon="ellipsis-v"] path {
    fill: rgba(28, 76, 166, 0.91);
  }

  /**
   * Animation
   */
  /* appearing */
  .icons-enter-active {
    opacity: 0;
    animation: icons-enter 300ms ease-in forwards;
    animation-iteration-count: 1;
  }

  @keyframes icons-enter {
    0% {transform: translateY(-10px); opacity: 0;}
    50% {transform: translateY(5px); opacity: 0.75;}
    100% { transform: translateY(0px); opacity: 1;}
  }

  /* disappearing */
  .icons-leave-active {
    animation: icons-leave 300ms ease-out forwards;
    animation-iteration-count: 1;
  }

  @keyframes icons-leave {
    0% {transform: translateY(0px); opacity: 1;}
    50% {transform: translateY(5px); opacity: 0.75;}
    100% { transform: translateY(-10px); opacity: 0;}
  }

  /* delays */
  .icons-enter-active:nth-child(5) {
    animation-delay: 25ms;
  }
  .icons-enter-active:nth-child(4) {
    animation-delay: 50ms;
  }
  .icons-enter-active:nth-child(3) {
    animation-delay: 75ms;
  }
  .icons-enter-active:nth-child(2) {
    animation-delay: 100ms;
  }

  .icons-leave-active:nth-child(2) {
    animation-delay: 25ms;
  }
  .icons-leave-active:nth-child(3) {
    animation-delay: 50ms;
  }
  .icons-leave-active:nth-child(4) {
    animation-delay: 75ms;
  }
  .icons-leave-active:nth-child(5) {
    animation-delay: 100ms;
  }
</style>
