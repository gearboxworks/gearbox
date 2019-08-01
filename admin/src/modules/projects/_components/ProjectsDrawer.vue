<template>
  <div class="drawer mb-3">
    <div
      class="drawer-handle"
      @click="onToggleDrawer"
      role="tab"
      aria-controls="drawer-contents-id"
    >
      <div class="label small" >
        <span
          tabindex="0"
          @keydown.enter="onToggleDrawer"
        >
          Viewing Options&nbsp;
          <font-awesome-icon
            v-if="expanded"
            key="mode-icon"
            :icon="['fa', 'chevron-up']"
          />
          <font-awesome-icon
            v-else
            key="mode-icon"
            :icon="['fa', 'chevron-down']"
          />
        </span>
      </div>

      <div class="current-filter">
        <b-badge title="Filter by state" :variant="statesVariant">{{labelStates}}</b-badge>
        <b-badge title="Filter by location" :variant="(showBasedirs == 'all') ? 'secondary' : 'warning'" v-if="hasExtraBaseDirs">{{labelBasedirs}}</b-badge>
        <b-badge title="Filter by used stack" :variant="(showStacks == 'all') ? 'secondary' : 'warning'">{{labelStacks}}</b-badge>
        <b-badge title="Filter by used program" :variant="(showPrograms == 'all') ? 'secondary' : 'warning'">{{labelPrograms}}</b-badge>
        <b-badge title="Sorting order">{{labelSorting}}</b-badge>
      </div>
    </div>
    <!--Filter-->
    <div
      v-if="expanded"
      class="drawer-contents clearfix"
      id="drawer-contents-id"
      role="tabpanel"
      :aria-expanded="expanded"
    >
      <div class="left-panel">
        <b-form class="filter-form">

          <b-form-checkbox-group
            class="form-group--states"
            v-model="showStates"
            name="showStates"
            label=""
            label-for="filter-state"
            description="State"
            switches
            stack
          >
            <b-form-checkbox value="running" title="Include projects that are currently RUNNING" @change="onToggleState($event, 'running')">Running</b-form-checkbox>
            <b-form-checkbox value="stopped" title="Include projects that are currently STOPPED" @change="onToggleState($event, 'stopped')">Stopped</b-form-checkbox>
            <b-form-checkbox value="candidates" title="Include projects that are yet to be imported" @change="onToggleState($event, 'candidates')">Candidates</b-form-checkbox>
            <small class="form-text text-muted">Project State</small>
          </b-form-checkbox-group>

          <b-form-group
            class="form-group--location"
            label=""
            label-for="filter-location"
            description="Location"
            v-if="hasExtraBaseDirs"
          >
            <b-select
              id="filter-basedirs"
              variant="secondary"
              v-model="showBasedirs"
              :options="basedirsAsOptions"
              @change="onChangeFilter($event, 'basedir')"
              tabindex="4"
            >
              <template slot="first">
                <option :value="null" disabled>Show projects from...</option>
                <option value="all">All known locations</option>
              </template>
            </b-select>
          </b-form-group>

          <b-form-group
            class="form-group--stacks"
            label=""
            label-for="filter-stacks"
            description="Used Stacks"
          >
            <b-select
              id="filter-stacks"
              variant="secondary"
              v-model="showStacks"
              @change="onChangeFilter($event, 'stacks')"
            >
              <option
                :value="null"
                disabled
              >
                Filter by stacks...
              </option>
              <option
                value="all"
              >
                Any stack
              </option>
              <optgroup label="Specific Stacks">
                <option
                  v-for="item in stacksAsOptions"
                  :key="item.value"
                  :value="item.value"
                >
                  {{item.text.replace('gearbox.works/', '').toUpperCase()}}
                </option>
              </optgroup>
              <!--option value="none">No stacks assigned</option-->
            </b-select>
          </b-form-group>

          <b-form-group
            class="form-group--programs"
            label=""
            label-for="filter-programs"
            description="Used Programs"
          >
            <b-select
              id="filter-programs"
              variant="secondary"
              v-model="showPrograms"
              @change="onChangeFilter($event, 'programs')"
            >
              <option
                :value="null"
                disabled
              >
                Filter by programs...
              </option>
              <option value="all">Any program</option>
              <optgroup label="Specific Program">
                <option
                  v-for="item in programsAsOptions"
                  :key="item.value"
                  :value="item.value"
                >
                  {{item.text.toUpperCase()}}
                </option>
              </optgroup>
              <!--option value="none">No programs assigned</option-->
            </b-select>
          </b-form-group>
        </b-form>
      </div>

      <div class="right-panel">
        <b-form class="sort-form">
          <b-form-group
            id="sort-by-group"
            label=""
            label-for="sort-by"
            description="Sort by"
          >
            <b-select
              id="sort-by"
              variant="secondary"
              v-model="sortByField"
              @change = onChangeSortByField($event)
            >
              <option :value="null" disabled>Sort by...</option>
              <option value="access-date" disabled>Access date</option>
              <option value="creation-date" disabled>Creation date</option>
              <option value="project-title">Project title</option>
            </b-select>
          </b-form-group>
          <b-form-group
            id="sort-order-group"
            label=""
            label-for="sort-order-select"
            description="Order"
          >
            <a target="_blank"
               href="#"
               title="Sort Order"
               class="view-mode view-mode--order"
               @click.prevent="onToggleSortingOrder"
            >
              <font-awesome-icon
                :icon="['fa', sortAscending ? 'sort-alpha-down': 'sort-alpha-up']"
              />
            </a>
          </b-form-group>
        </b-form>

        <b-form class="view-form">
          <b-form-group
            id="view-form"
            label=""
            label-for="view-select"
            description="View Mode"
          >
            <a target="_blank"
               href="#"
               title="Cards View"
               :class="{'view-mode': true, 'view-mode--cards': true, 'is-inactive': (viewMode != 'cards')}"
               @click.prevent="viewMode = 'cards'; $emit('switch-view-mode', $event, 'cards')"
            >
              <font-awesome-icon
                :icon="['fa', 'columns']"
              />
            </a>
            <a target="_blank"
               href="#"
               title="Table View"
               :class="{'view-mode': true, 'view-mode--table': true, 'is-inactive': (viewMode != 'table')}"
               @click.prevent="viewMode = 'table'; $emit('switch-view-mode', $event, 'table')"
            >
              <font-awesome-icon
                :icon="['fa', 'th-list']"
              />
            </a>
          </b-form-group>
        </b-form>
      </div>
    </div>
  </div>
</template>

<script>

import store from '../../../store'
import ProjectMethodTypes from '../_store/public-types'
import StackMethodTypes from '../../stacks/_store/public-types'
import ServiceMethodTypes from '../../services/_store/public-types'
import BasedirMethodTypes from '../../basedirs/_store/public-types'

const { ActionTypes: ProjectActions } = ProjectMethodTypes
const { GetterTypes: BasedirGetters } = BasedirMethodTypes
const { GetterTypes: StackGetters } = StackMethodTypes
const { GetterTypes: ServiceGetters } = ServiceMethodTypes

console.log(store)

export default {
  name: 'ProjectsDrawer',
  data () {
    return {
      expanded: false,

      showStates: ['running', 'stopped', 'candidates'],
      showBasedirs: 'all',
      showStacks: 'all',
      showPrograms: 'all',

      sortByField: 'project-title',
      sortAscending: true,
      viewMode: 'cards'
    }
  },
  computed: {

    // stack () {
    //   return this.$store.getters[StackGetters.FIND_BY]('id', this.gearspec.attributes.stack_id)
    // },
    stacksAsOptions: store.getters[StackGetters.LIST_OPTIONS],
    programsAsOptions: store.getters[ServiceGetters.LIST_PROGRAM_OPTIONS],
    basedirsAsOptions: store.getters[BasedirGetters.LIST_OPTIONS],
    hasExtraBaseDirs: store.getters[BasedirGetters.HAS_EXTRA_BASEDIRS],

    labelStates () {
      const states = this.showStates
      const running = (states.indexOf('running') !== -1) ? 'Running projects' : ''
      const stopped = states.indexOf('stopped') !== -1 ? 'Stopped projects' : ''
      const candidates = states.indexOf('candidates') !== -1

      const projects = (running && stopped)
        ? 'All projects'
        : (running || stopped)
          ? running + stopped
          : ''

      return projects
        ? projects + ((candidates && (running || stopped)) ? '' : ' (no candidates)')
        : (candidates ? 'Project candidates' : '')
    },

    statesVariant () {
      const states = this.showStates
      const running = (states.indexOf('running') !== -1) ? 'Running projects' : ''
      const stopped = states.indexOf('stopped') !== -1 ? 'Stopped projects' : ''
      const candidates = states.indexOf('candidates') !== -1
      return (running && stopped && candidates)
        ? 'secondary'
        : 'warning'
    },

    labelBasedirs () {
      const basedir = (this.showBasedirs !== 'all')
        ? this.$store.getters[BasedirGetters.FIND_BY]('id', this.showBasedirs)
        : null
      return 'From ' + (basedir ? basedir.attributes.basedir : 'all known locations')
    },

    labelStacks () {
      let label
      if (this.showStacks === 'none') {
        label = 'With no stacks assigned'
      } else {
        const stack = (this.showStacks !== 'all')
          ? this.$store.getters[StackGetters.FIND_BY]('id', this.showStacks)
          : null
        label = 'Using ' + (stack ? (stack.attributes.stackname.toUpperCase() + ' stack') : 'any stack')
      }
      return label
    },

    labelPrograms () {
      let label
      if (this.showPrograms === 'none') {
        label = 'With no programs assigned'
      } else {
        const program = (this.showPrograms !== 'all') ? this.showPrograms : null
        label = 'Using ' + (program ? program.toUpperCase() : 'any program')
      }
      return label
    },

    labelSorting () {
      return 'Sorted by ' + this.sortByField.replace('-', ' ') + (this.sortAscending ? '' : ' (reverse)')
    }
  },
  watch: {
    showStates: function (val, oldVal) {
      this.$store.dispatch(ProjectActions.SET_LIST_FILTER, {
        field: 'states',
        values: this.showStates
      })
    }
  },
  methods: {

    onToggleDrawer () {
      this.expanded = !this.expanded
      // if (this.expanded) {
      //   console.log('focus', this.$refs['state-filter'])
      //   this.$refs['state-filter'].$el.focus()
      //   // this.$nextTick(() => {
      //   //   this.$refs['state-filter'].$el.focus()
      //   // })
      // }
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

    onChangeSortByField (value) {
      this.$store.commit(
        ProjectActions.SET_LIST_FILTER_SORT_BY,
        value
      )
    },

    onToggleSortingOrder () {
      this.sortAscending = !this.sortAscending
      this.$store.dispatch(
        ProjectActions.SET_LIST_FILTER_SORT_ASC,
        this.sortAscending
      )
    },

    onChangeFilter (values, field) {
      this.$store.dispatch(
        ProjectActions.SET_LIST_FILTER,
        {
          field,
          values
        }
      )
    }
  }
}
</script>

<style scoped>
.drawer {
  margin-bottom: 2rem;
  display: grid;
  /*grid-template-columns: auto;*/
  /*grid-template-rows: auto;*/
  grid-template-areas:
    "content-area"
    "handle-area";
}

.drawer a {
  color: #1e69b9;
}

.drawer a:hover {
  color: #2379D3;
}

.drawer-contents {
  padding: 1rem;
  background-color: #eee;
  padding-top: 22px;
  padding-bottom: 0;
  border-bottom: 1px solid silver;
  clear: both;
  width: 100%;
  grid-area: content-area;
}

.drawer-handle {
  clear: both;
  width: 100%;
  padding: 0 0 0 1rem;
  grid-area: handle-area;
}

.drawer-handle .current-filter {
  display: inline;
  position: relative;
  top: 6px;
}

.drawer-handle .badge {
  padding: 4px;
  margin-right: 8px;
}

.drawer-handle .label{
  display: inline-block;
  float: right;
  background-color: #1e69b9;
  color: white;
  line-height: 0;
  cursor: pointer;
  padding: 0 16px 0 0;
  vertical-align: text-top;
  height: 30px;
}

.drawer-handle .label:before {
  width: 0;
  content: "";
  height: 0;
  margin-top: 0;
  margin-left: 0;
  display: inline-block;
  border-style: solid;
  border-width: 0 30px 30px 0;
  border-color: transparent transparent white white;
}

.drawer-handle .label span{
  display: inline-block;
  line-height: normal;
  position: relative;
  top: -9px;
  left: 4px;
  margin-left: 4px;
}

.form-group,
.form-group--states {
  display: inline-block;
  margin-right: 2rem;
}

.filter-form {
  float: left;
  display: inline-block;
  margin-top: 10px;
}

.form-group--states{
  vertical-align: top;
  padding-top: 0.5rem;
}

.left-panel {
  float: left;
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
  position: relative;
  top: 6px;
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
</style>
<style>
  .drawer .text-muted {
    opacity: 0;
    transition: all 400ms;
  }

  .drawer:hover .text-muted {
    opacity: 0.7;
  }
</style>
