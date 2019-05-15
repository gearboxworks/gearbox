<template>
  <div class="drawer mb-3 clearfix">
    <!--Filter-->
    <div v-if="expanded" class="drawer-contents clearfix">
      <div class="left-panel">
        <b-form class="filter-form">

          <b-form-checkbox-group
            class="form-group--states"
            v-model="show_states"
            name="show_states"
            label=""
            label-for="filter-state"
            description="State"
            switches
            inline
          >
            <b-form-checkbox value="running" title="Include projects that are currently RUNNING">Running projects</b-form-checkbox>
            <b-form-checkbox value="stopped" title="Include projects that are currently STOPPED">Stopped projects</b-form-checkbox>
            <b-form-checkbox value="candidates" title="Include projects that are yet to be imported">Project candidates</b-form-checkbox>
            <small tabindex="-1" class="form-text text-muted">Project State</small>
          </b-form-checkbox-group>

          <b-form-group
            class="form-group--location"
            label=""
            label-for="filter-location"
            description="Location"
          >
            <b-select id="filter-location" variant="secondary" v-model="show_locations" :options="basedirsAsOptions">
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
            <b-select id="filter-stacks" variant="secondary" v-model="show_stacks">
              <option :value="null" disabled>Filter by stacks...</option>
              <option value="all">Any stack</option>
              <optgroup label="Specific Stacks">
                <option v-for="item in stacksAsOptions" :key="item.value" :value="item.value">{{item.text.toUpperCase()}}</option>
              </optgroup>
              <option value="none">No stacks assigned</option>
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
            <b-select id="sort-by" variant="secondary" v-model="sort_by">
              <option :value="null" disabled>Sort by...</option>
              <option value="access-date">Access date</option>
              <option value="creation-date">Creation date</option>
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
               @click.prevent="sort_ascending = !sort_ascending"
            >
              <font-awesome-icon
                :icon="['fa', sort_ascending ? 'sort-alpha-down': 'sort-alpha-up']"
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
               :class="{'view-mode': true, 'view-mode--cards': true, 'is-inactive': (view_mode != 'cards')}"
               @click.prevent="view_mode = 'cards'"
            >
              <font-awesome-icon
                :icon="['fa', 'columns']"
              />
            </a>
            <a target="_blank"
               href="#"
               title="Table View"
               :class="{'view-mode': true, 'view-mode--table': true, 'is-inactive': (view_mode != 'table')}"
               @click.prevent="view_mode = 'table'"
            >
              <font-awesome-icon
                :icon="['fa', 'th-list']"
              />
            </a>
          </b-form-group>
        </b-form>
      </div>
    </div>
    <div class="drawer-handle" @click="expanded=!expanded">
      <div class="current-filter">
        <b-badge title="Project State">{{states_label}}</b-badge>
        <b-badge title="Project Locations">{{locations_label}}</b-badge>
        <b-badge title="Stacks">{{stacks_label}}</b-badge>
        <b-badge title="Sorting">{{sorting_label}}</b-badge>
      </div>
      <div class="label small"><span>Viewing Options&nbsp;
        <font-awesome-icon
          v-if="expanded"
          :icon="['fa', 'chevron-up']"
        />
        <font-awesome-icon
          v-else
          :icon="['fa', 'chevron-down']"
        />
        </span>
      </div>
    </div>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'

export default {
  name: 'ProjectsDrawer',
  props: {},
  computed: {
    ...mapGetters(['basedirBy', 'stackBy', 'basedirsAsOptions', 'stacksAsOptions']),
    states_label () {
      const states = this.show_states
      const running = (states.indexOf('running') !== -1) ? 'Running projects' : ''
      const stopped = states.indexOf('stopped') !== -1 ? 'Stopped projects' : ''
      const candidates = states.indexOf('candidates') !== -1

      const projects = (running && stopped)
        ? 'All projects'
        : (running || stopped)
          ? running + stopped
          : ''

      return projects
        ? projects + (candidates ? ' and project candidates' : '')
        : (candidates ? 'Project candidates' : '')
    },
    locations_label () {
      const basedir = (this.show_locations !== 'all') ? this.basedirBy('id', this.show_locations) : null
      return 'From ' + (basedir ? basedir.attributes.basedir : 'all known locations')
    },
    stacks_label () {
      let label
      if (this.show_stacks === 'none') {
        label = 'With no stacks assigned'
      } else {
        const stack = (this.show_stacks !== 'all') ? this.stackBy('id', this.show_stacks) : null
        label = 'Using ' + (stack ? (stack.attributes.stackname.toUpperCase() + ' stack') : 'any stack')
      }
      return label
    },
    sorting_label () {
      return 'Sorted by ' + this.sort_by.replace('-', ' ') + (this.sort_ascending ? '' : ' (reverse)')
    }
  },
  data () {
    return {
      expanded: false,
      show_states: ['running', 'stopped', 'candidates'],
      show_locations: 'all',
      show_stacks: 'all',
      sort_by: 'access-date',
      sort_ascending: true,
      view_mode: 'cards'
    }
  },
  methods: {
    // escAttr (value) {
    //   return value.replace(/\//g, '-').replace(/\./g, '-')
    // }
  }
}
</script>

<style scoped>
.drawer {
  margin-bottom: 2rem;
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
}

.drawer-handle {
  padding: 0;
  clear: both;
  height: 30px;
}

.drawer-handle .current-filter {
  display: inline-block;
  margin-top: 0.5rem;
  margin-left: 1rem;
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
