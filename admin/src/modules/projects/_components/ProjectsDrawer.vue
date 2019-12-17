<template>
  <div class="drawer mb-3 clearfix">
    <h2 class="filter-heading">{{$t('projects.filterHeading')}}</h2>
    <div class="left-panel">
      <b-form class="filter-form">

        <b-form-checkbox-group
          v-model="showStates"
          class="form-group--states"
          label=""
          label-for="filter-state"
          description=""
          switches
          stack
        >
          <b-form-checkbox value="running" :title="$t('projects.filterRunningTitle')" v-b-tooltip.hover @change="onToggleState($event, 'running')">{{$t('projects.filterRunning')}}</b-form-checkbox>
          <b-form-checkbox value="stopped" :title="$t('projects.filterStoppedTitle')" v-b-tooltip.hover @change="onToggleState($event, 'stopped')">{{$t('projects.filterStopped')}}</b-form-checkbox>
          <b-form-checkbox value="candidates" :title="$t('projects.filterCandidatesTitle')" v-b-tooltip.hover @change="onToggleState($event, 'candidates')">{{$t('projects.filterCandidates')}}</b-form-checkbox>
        </b-form-checkbox-group>

        <b-form-group
          v-if="basedirsAsOptions.length > 2"
          class="form-group--location"
          label=""
          label-for="filter-location"
          description=""
          :title="$t('projects.filterByLocation')"
          v-b-tooltip.hover
        >
          <treeselect
            v-model="showBasedirs"
            instanceId="filter-location"
            class="vue-treeselect--location"
            :options="basedirsAsOptions"
            :normalizer="optionNormalizer"
            :placeholder="$t('projects.filterAllLocations')"
            :clearable="false"
          />
        </b-form-group>

        <b-form-group
          class="form-group--stack"
          name="showStacks"
          label=""
          label-for="filter-stack"
          description=""
          :title="$t('projects.filterByStack')"
          v-b-tooltip.hover
        >
          <treeselect
            v-model="showStacks"
            instanceId="filter-stack"
            class="vue-treeselect--stack"
            :options="stacksAsOptions"
            :normalizer="optionNormalizer"
            :placeholder="$t('projects.filterAllStacks')"
            :clearable="false"
          />
        </b-form-group>

        <b-form-group
          class="form-group--program"
          name="showProgram"
          label=""
          label-for="filter-program"
          description=""
          :title="$t('projects.filterByProgram')"
          v-b-tooltip.hover.top
        >
          <treeselect
            v-model="showPrograms"
            instanceId="filter-program"
            class="vue-treeselect--program"
            :options="programsAsOptions"
            :normalizer="optionNormalizer"
            :placeholder="$t('projects.filterAllPrograms')"
            :clearable="false"
          />
        </b-form-group>
      </b-form>
    </div>

    <div class="right-panel">
      <b-form class="sort-form">
        <b-form-group
          id="form-group--sort-by"
          label=""
          label-for="sort-by"
          description=""
          :title="$t('projects.sortBy')"
          v-b-tooltip.hover
        >
          <treeselect
            class="vue-treeselect--sort-by"
            instance-id="sort-by"
            v-model="sortByField"
            :options="[{id: 'access_date', label: $t('projects.sortByAccess')}, {id: 'creation_date', label: $t('projects.sortByCreation')}, {id: 'id', label: $t('projects.sortById')}]"
            :clearable = "false"
          />
        </b-form-group>

        <b-form-group
          class="form-group--sort-order"
          label=""
          label-for="sort-order-select"
          description=""
          :title="sortAscending ? $t('projects.sortAscending') : $t('projects.sortDescending')"
          v-b-tooltip.hover
        >
          <a target="_blank"
             href="#"
             :title="$t('projects.sortOrder')"
             class="view-mode view-mode--order"
             @click.prevent="onSetSortingOrder"
          >
            <font-awesome-icon
              :icon="['fa', sortAscending ? 'sort-alpha-down': 'sort-alpha-up']"
            />
          </a>
        </b-form-group>
      </b-form>

      <b-form class="view-form">
        <b-form-group
          class="form-group--view-form"
        >
          <a target="_blank"
             href="#"
             :title="(viewMode === 'cards') ? this.$t('projects.viewingCards') : this.$t('projects.viewAsCards')"
             v-b-tooltip.hover
             :class="{'view-mode': true, 'view-mode--cards': true, 'is-inactive': (viewMode !== 'cards')}"
             @click.prevent="onViewModeChange('cards')"
          >
            <font-awesome-icon
              :icon="['fa', 'columns']"
            />
          </a>
          <a target="_blank"
             href="#"
             :title="(viewMode === 'table') ? this.$t('projects.viewingTable') : this.$t('projects.viewAsTable')"
             v-b-tooltip.hover
             :class="{'view-mode': true, 'view-mode--table': true, 'is-inactive': (viewMode !== 'table')}"
             @click.prevent="onViewModeChange('table')"
          >
            <font-awesome-icon
              :icon="['fa', 'th-list']"
            />
          </a>
        </b-form-group>
      </b-form>
    </div>
  </div>
</template>

<script>

import { ProjectActions } from '../_store/method-names'
import { BasedirGetters } from '../../basedirs/_store/method-names'
import { StackGetters } from '../../stacks/_store/method-names'
import { ServiceGetters } from '../../services/_store/method-names'

import { getCookie, setCookie } from '../../_helpers'

const COOKIE = 'Gearbox-view-settings'
const dataDefaults = {
  showStates: ['running', 'stopped', 'candidates'],
  showBasedirs: 'all',
  showStacks: 'all',
  showPrograms: 'all',
  sortByField: 'id',
  sortAscending: true,
  viewMode: 'cards'
}

let _loadingCookie = false

export default {
  name: 'ProjectsDrawer',
  data () {
    return dataDefaults
  },
  computed: {

    stacksAsOptions () {
      const stacksOptions = this.$store.getters[StackGetters.LIST_OPTIONS]('stackname')
      stacksOptions.unshift({ value: 'all', text: this.$t('projects.filterAllStacks') })
      return stacksOptions
    },

    programsAsOptions () {
      const programOptions = this.$store.getters[ServiceGetters.LIST_PROGRAM_OPTIONS]()
      programOptions.unshift({ value: 'all', text: this.$t('projects.filterAllPrograms') })
      return programOptions
    },

    basedirsAsOptions () {
      const basedirOptions = this.$store.getters[BasedirGetters.LIST_OPTIONS]('basedir')
      basedirOptions.unshift({ value: 'all', text: this.$t('projects.filterAllLocations') })
      return basedirOptions
    },

    hasExtraBaseDirs () {
      return this.$store.getters[BasedirGetters.HAS_EXTRA_BASEDIRS]()
    }

  },
  watch: {
    showStates: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER, {
        field: 'states',
        values: this.showStates
      }).then(() => this.updateCookie())
    },
    showBasedirs: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER, {
        field: 'basedir',
        values: this.showBasedirs
      }).then(() => this.updateCookie())
    },
    showStacks: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER, {
        field: 'stacks',
        values: this.showStacks
      }).then(() => this.updateCookie())
    },
    showPrograms: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER, {
        field: 'programs',
        values: this.showPrograms
      }).then(() => this.updateCookie())
    },
    sortByField: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER_SORT_BY, val).then(() => this.updateCookie())
    },
    sortAscending: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER_SORT_ASC, val).then(() => this.updateCookie())
    }
  },
  methods: {

    optionNormalizer (option) {
      return {
        id: option.value,
        label: option.text
        /* children: option.subOptions, */
      }
    },

    onToggleState (value, attribute) {
      const states = this.showStates
      const running = states.indexOf('running') !== -1
      const stopped = states.indexOf('stopped') !== -1
      const candidates = states.indexOf('candidates') !== -1

      /**
       * All unselected would be and invalid state, therefore
       * make sure either candidates or running/stopped is selected
       */
      if ((attribute === 'candidates') && !running && !stopped && candidates) {
        this.showStates = ['running', 'stopped']
      } else if ((attribute === 'running') && running && !stopped && !candidates) {
        this.showStates = ['candidates']
      } else if ((attribute === 'stopped') && !running && stopped && !candidates) {
        this.showStates = ['candidates']
      }
    },

    onSetSortingOrder () {
      this.sortAscending = !this.sortAscending
    },

    onViewModeChange (mode) {
      this.viewMode = mode
      this.$emit('switch-view-mode', mode)
    },

    updateCookie () {
      // console.log('Update cookie', { ...this.$data })
      if (!_loadingCookie) {
        setCookie(COOKIE, this.$data)
      }
    },

    loadCookie () {
      const values = getCookie(COOKIE)
      const filterFields = ['showStates', 'showBasedirs', 'showPrograms', 'showStacks']
      _loadingCookie = true
      if (values && Object.entries(values).length) {
        Object.entries(values).forEach(([key, val]) => {
          if (typeof this.$data[key] !== 'undefined') {
            if (filterFields.includes(key)) {
              this.$data[key] = val
            } else if (key === 'sortByField') {
              this.$data[key] = val
            } else if (key === 'sortAscending') {
              this.$data[key] = val
            } else if (key === 'viewMode') {
              this.onViewModeChange(val)
            } else {
              console.log('A value from cookie will be ignored:', key, val)
            }
          }
        })
      }
      _loadingCookie = false
    }
  },
  created () {
    this.loadCookie()
  }

}
</script>

<style scoped>
.drawer {
  display: block;
  width: 100%;
  clear: both;
  margin-bottom: 0;
  padding: 10px 1rem 0 1rem;
  background-color: #f9f9f9;
  border-bottom: 1px solid #e8e8e8;
  position: relative;
}
.filter-heading {
  position: absolute;
  top: 14px;
  left: 1rem;
  /* font-weight: bold; */
  font-size: 20px;
  color: #7e7e7e;
}

.drawer a {
  color: #1e69b9;
}

.drawer a:hover {
  color: #2379D3;
}

.drawer >>> .text-muted {
  opacity: 0;
  transition: all 400ms;
}

.drawer:hover >>> .text-muted {
  opacity: 0.7;
}

>>> .form-text {
  margin-top: 0;
}

.drawer >>> .form-group {
  margin-bottom: 9px;
}

.filter-form {
  float: left;
  display: inline-block;
  margin-top: 0;
}

.form-group {
  display: inline-block;
  vertical-align: top;
  margin-right: 2rem;
}

.form-group--states {
  display: inline-block;
  margin-right: 13px;
  margin-top: 5px;
}

.left-panel {
  float: left;
  margin-left: 99px;
}
.right-panel {
  float: right;
  min-width: 339px;
}

.view-form,
.sort-form {
  display: inline-block;
}

.view-mode {
  font-size: 200%;
  line-height: 1;
  display: inline-block;
}

.view-mode--cards{
  margin-right: 0.5rem;
}

.view-mode--table {
  font-size: 195%;
}

.view-mode.is-inactive {
  color: #6c757d !important;
}

.form-group--location {
  max-width: 25rem;
}

.short {
  display: none;
}

@media (max-width: 1200px) {
}

@media (max-width: 1515px) {
  .right-panel{
    float: left;
  }
}

.vue-treeselect--state,
.vue-treeselect--location,
.vue-treeselect--stack,
.vue-treeselect--program {
 /* display: inline-block;*/
}

.vue-treeselect--state,
.vue-treeselect--location {
  max-width: 18.5rem;
}

.vue-treeselect--stack,
.vue-treeselect--program {
  max-width: 10rem;
}

.vue-treeselect--sort-by {
  max-width: 10rem;
}

</style>
