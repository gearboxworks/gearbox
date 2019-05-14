<template>
  <div class="drawer">
    <!--Filter-->
    <div v-if="expanded" class="drawer-contents clearfix">
      <div class="left-panel">
        <b-form class="filter-form">
          <b-form-group
            class="form-group--location"
            label=""
            label-for="filter-location"
            description="Location"
          >
            <b-select id="filter-location" variant="secondary">
              <option disabled>Filter by location...</option>
              <option :value="null">Projects from all known locations</option>
              <option>C:\Users\user\Gearbox Sites\</option>
            </b-select>
          </b-form-group>

          <b-form-group
            class="form-group--state"
            label=""
            label-for="filter-state"
            description="State"
          >
            <b-select id="filter-state" variant="secondary">
              <option disabled>Filter by State ...</option>
              <option :value="null">Running & Stopped</option>
              <option>Running only</option>
              <option>Stopped only</option>
            </b-select>

          </b-form-group>

          <b-form-group
            class="form-group--stacks"
            label=""
            label-for="filter-stacks"
            description="Stacks"
          >
            <b-select id="filter-stacks" variant="secondary">
              <option disabled>Filter by Used Stacks...</option>
              <option :value="null">Projects with any Stack</option>
              <option :value="null">Projects with no Stacks</option>
              <option>WordPress</option>
              <option>ElasticSearch</option>
              <option>LAMP</option>
              <option>LEMP</option>
            </b-select>
          </b-form-group>

          <b-form-group
            class="form-group--candidates"
            label=""
            label-for="filter-candidates"
            description=" "
            title="Check to see all projects that can be imported"
          >
            <b-form-checkbox
              id="filter-candidates"
              name="filter-candidates"
              value="show-candidates"
              unchecked-value="hide-candidates"
            >
              Project Candidates
            </b-form-checkbox>
          </b-form-group>
        </b-form>
      </div>

      <div class="right-panel">
        <b-form class="sort-form">
          <b-form-group
            id="sort-by-group"
            label=""
            label-for="sort-select"
            description="Sort by"
          >
            <b-select id="sort-by" variant="secondary">
              <option disabled>Sort by...</option>
              <option :value="null">Access date</option>
              <option>Creation date</option>
              <option>Project title</option>
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
               title="Table View"
               class="view-mode view-mode--table"
               @click.prevent="toggleSortOrder"
            >
              <font-awesome-icon
                :icon="['fa', 'sort-alpha-up']"
              />
            </a>
          </b-form-group>
        </b-form>

        <b-form class="view-form">
          <b-form-group
            id="view-form"
            label=""
            label-for="view-select"
            description="View"
          >
            <a target="_blank"
               href="#"
               title="Cards View"
               class="view-mode view-mode--cards"
               @click.prevent="switchView('cards')"
            >
              <font-awesome-icon
                :icon="['fa', 'columns']"
              />
            </a>
            <a target="_blank"
               href="#"
               title="Table View"
               class="view-mode view-mode--table is-inactive"
               @click.prevent="switchView('table')"
            >
              <font-awesome-icon
                :icon="['fa', 'th-list']"
              />
            </a>
          </b-form-group>
        </b-form>
      </div>
    </div>
    <div class="drawer-handle" @click="expanded=!expanded"><span class="small">Viewing Options {{expanded?'^':'v'}}</span></div>
  </div>
</template>

<script>

export default {
  name: 'ProjectsFilter',
  props: {},
  data () {
    return {
      expanded: false
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
  border-bottom: 1px solid silver;
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
}

.drawer-handle {
  display: inline-block;
  float: right;
  clear: both;
  background-color: #1e69b9;
  padding: 0;
  height: 30px;
  color: white;
  line-height: 0;
  cursor: pointer;
}

.drawer-handle:before {
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

.drawer-handle span {
  padding: 0 16px 0 12px;
  vertical-align: text-top;
}

.form-group{
  display: inline-block;
  margin-right: 2rem;
}

.filter-form {
  float: left;
  display: inline-block;
  margin-top: 10px;
}

.form-group--candidates{
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
