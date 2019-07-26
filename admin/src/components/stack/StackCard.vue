<template>
  <div
    :class="{'stack-card': true, 'is-collapsible': isCollapsible, 'is-collapsed': isCollapsed}"
    role="tab"
    :aria-expanded="!isCollapsed"
  >
    <h2
      class="stack-title"
      @click="onTitleClicked"
      @keydown.enter="onTitleClicked"
      tabindex="0"
      :title="isCollapsed ? 'Show services' : 'Hide services'"
      v-b-tooltip.hover
    >
      <font-awesome-icon
        :icon="['fa', 'layer-group']"
        style="color:silver;"
      />&nbsp;
      {{stackId.replace('gearbox.works/', '')}}
    </h2>

    <stack-toolbar
      v-if="!isCollapsed"
      :stackId="stackId"
      @show-alert="showAlert"
    />

    <div class="stack-content">
      <ul
          v-if="!isCollapsible || !isCollapsed"
          class="service-list"
      >
        <li
            v-for="item in stackItems"
            :key="id + item.gearspec.attributes.role"
            class="service-item"
        >
          <stack-gear :stackItem="item" />
        </li>
      </ul>
      <b-alert
        :key="stackId"
        :show="alertShow"
        :dismissible="alertDismissible"
        :variant="alertVariant"
        @dismissed="alertShow=false"
        fade
      >
        {{alertContent}}
      </b-alert>
    </div>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'
import StackToolbar from './StackToolbar.vue'
import StackGear from './StackGear.vue'

export default {
  name: 'StackCard',
  components: {
    StackToolbar,
    StackGear
  },
  inject: ['project', 'projectPrefix'],
  props: {
    'stackId': {
      type: String,
      required: true
    },
    'stackItems': {
      type: Array,
      required: true
    },
    'isCollapsible': {
      type: Boolean,
      default: true,
      required: false
    },
    'startCollapsed': {
      type: Boolean,
      required: false,
      default: true
    }
  },
  data () {
    return {
      id: this.project.id,
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'info',
      isCollapsed: this.startCollapsed
    }
  },
  computed: {
    ...mapGetters([
      'serviceBy',
      'gearspecBy'
    ])
  },
  methods: {
    escAttr (value) {
      return value.replace(/\//g, '-').replace(/\./g, '-')
    },

    onTitleClicked () {
      this.isCollapsed = !this.isCollapsed
      // this.showAlert('Clicked ' + this.stackId.replace('gearbox.works/', ''))
    },

    showAlert (alert) {
      if (typeof alert === 'string') {
        this.alertContent = alert
      } else {
        this.alertVariant = alert.variant || this.alertVariant
        this.alertDismissible = alert.dismissible || this.alertDismissible
        this.alertContent = alert.content || this.alertContent
      }
      this.alertShow = true
    }
  }
  // ,
  // watch: {
  //   startCollapsed: function (val, oldVal) {
  //     if (!val) {
  //       this.isCollapsed = false
  //     }
  //   }
  // }
}
</script>

<style scoped>
  .stack-card:not(.is-collapsed){
    min-width: 310px;
  }

  .stack-card.is-collapsible{
    margin-right: 0.5rem;
    margin-bottom: 0.75rem;
    border: 1px solid #dedede;
    border-radius: 5px;
    padding: 10px;
    display: inline-block;
    align-self: flex-start;
    background-color :#eaeaea;
  }

  .stack-card.is-collapsible.is-collapsed {
    padding-bottom: 2px;
    color: #17a2b8;
  }

  .stack-card.is-collapsible.is-collapsed:hover {
    background-color: #17a2b8;
    color: white;
  }

  .stack-title {
    text-transform: uppercase;
    font-size: 1rem;
    display: inline-block;
    color: #989898;
    margin-left: 5px;
  }

  .stack-card.is-collapsible .stack-title {
    cursor: pointer;
    margin-bottom: 0;
    line-height: 0;
    /* margin-top: 0; */
    position: relative;
    top: -2px;
    left: -3px;
  }

  .stack-card.is-collapsible.is-collapsed:hover .stack-title{
    color: white;
  }

  .stack-card.is-collapsible.is-collapsed .stack-title {
    color: #17a2b8;
  }

  .stack-card.is-collapsible:not(.is-collapsed) .stack-title:hover {
    color: #17a2b8;
  }

  .stack-content{
  }

  .service-list{
    margin-top: 0px;

    padding-bottom: 10px;
    border-bottom: 1px solid #e0e0e0;
    margin-bottom: 10px;

    list-style: none;
    padding-inline-start: 0;
    clear: both;
  }

  .stack-card.is-collapsible .service-list,
  .stack-card:last-child .service-list {
    border-bottom: none;
    padding-bottom: 0;
    margin-bottom: 0;
  }

  .service-item {
    display: inline-block;
    list-style: none;
  }

  .alert {
    margin-left: 0;
    width: calc(100%);
    margin-top: 0;
    margin-bottom: 0;
    border-radius: 0;
    font-size: 12px;
    padding-top: 10px;
    padding-left: 10px;
    padding-bottom: 10px;
    position: relative;
    top: 0;
    box-shadow: 0px 1px 3px rgba(0,0,0,0.1);
  }
</style>
<style>
  .stack-card .alert-dismissible .close {
    padding: 0.5rem 0.75rem;
    top: 1px;
    right: -4px;
    font-weight: 100;
    font-size: 1.25rem;
  }
</style>
