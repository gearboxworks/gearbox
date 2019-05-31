<template>
  <div :class="{'project-stack': true, 'is-collapsible': isCollapsible}">
    <h2 class="stack-title" @click="onTitleClicked">{{stackId.replace('gearbox.works/', '')}}</h2>

    <stack-toolbar :project="project" :projectIndex="projectIndex" :stackId="stackId"></stack-toolbar>

    <div class="stack-content">
      <ul class="service-list" v-if="!isCollapsible || !isCollapsed">
        <li
            v-for="(item, itemIndex) in stackItems"
            :key="id + item.gearspec.attributes.role"
            class="service-item"
        >
          <stack-gear :projectId="project.id" :stackItem="item" :projectIndex="projectIndex" :stackIndex="stackIndex" :itemIndex="itemIndex"></stack-gear>
        </li>
      </ul>
      <b-alert
        :key="stackId"
        :show="alertShow"
        :dismissible="alertDismissible"
        :variant="alertVariant"
        @dismissed="alertShow=false"
        fade
      >{{alertContent}}</b-alert>
    </div>
  </div>
</template>

<script>

import StackToolbar from '../stack/StackToolbar.vue'
import StackGear from '../stack/StackGear.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'ProjectStack',
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
    },
    'stackIndex': {
      type: Number,
      required: true
    },
    'stackItems': {
      type: Array,
      required: true
    },
    'isCollapsible': {
      type: Boolean,
      default: false,
      required: false
    }
  },
  data () {
    return {
      id: this.project.id,
      alertShow: false,
      alertContent: 'content',
      alertDismissible: true,
      alertVariant: 'info',
      isCollapsed: true
    }
  },
  components: {
    StackToolbar,
    StackGear
  },
  computed: {
    ...mapGetters(['serviceBy', 'gearspecBy']),
    projectBase () {
      return 'gb-' + this.escAttr(this.id) + '-'
    }
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
}
</script>

<style scoped>
  .project-stack{
    min-width: 310px;
  }

  .stack-title {
    text-transform: uppercase;
    font-size: 1.25rem;
    display: inline-block;
    color: #989898;
    margin-left: 5px;
  }
  .is-collapsible .stack-title {
    cursor: pointer;
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

  .project-stack:last-child .service-list {
    border-bottom: none;
    padding-bottom: 0;
    margin-bottom: 0;
  }

  .service-item {
    display: inline-block;
    list-style: none;
    text-align: center;
    max-width: 110px;
    padding: 5px;
    margin: 5px;
    cursor: pointer;
    border: 1px solid transparent;
    border-radius: 4px;
    transition: all 400ms;
  }
  .service-item:hover {
    border: 1px solid #aaa;
    background-color: #eee;
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
    top: -14px;
    box-shadow: 0px 1px 3px rgba(0,0,0,0.1);
  }
</style>
<style>
  .project-stack .alert-dismissible .close {
    padding: 0.5rem 0.75rem;
    top: 1px;
    right: -4px;
    font-weight: 100;
    font-size: 1.25rem;
  }
</style>
