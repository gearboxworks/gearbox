<template functional>
  <b-navbar toggleable="lg" type="dark" variant="secondary">
    <b-navbar-brand to="/projects">Gearbox</b-navbar-brand>

    <b-navbar-toggle target="nav_collapse" />

    <b-collapse is-nav id="nav_collapse">
      <b-navbar-nav>
        <b-nav-item to="/projects">{{parent.$t('label.projects')}}</b-nav-item>
      </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item right to="/preferences">
          <font-awesome-icon
            :title="parent.$t('label.preferences')"
            :icon="['fa', 'cog']"
          /> {{parent.$t('label.basedirs')}}
        </b-nav-item>
        <b-nav-item right to="/stacks">
          <font-awesome-icon
            :title="parent.$t('label.stacks')"
            :icon="['fa', 'layer-group']"
          /> {{parent.$t('label.stacks')}}
        </b-nav-item>
        <b-nav-item right to="/services">
          <font-awesome-icon
            :title="parent.$t('label.services')"
            :icon="['fa', 'layer-group']"
          /> {{parent.$t('label.services')}}
        </b-nav-item>
        <b-nav-item right to="/gearspecs">
          <font-awesome-icon
            :title="parent.$t('label.gearspecs')"
            :icon="['fa', 'layer-group']"
          /> {{parent.$t('label.gearspecs')}}
        </b-nav-item>

      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>
