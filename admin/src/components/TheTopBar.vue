<template functional>
  <b-navbar toggleable="lg" type="dark" variant="secondary">
    <b-navbar-brand to="/projects">Gearbox</b-navbar-brand>

    <b-navbar-toggle target="nav_collapse" />

    <b-collapse is-nav id="nav_collapse">
      <b-navbar-nav>
        <b-nav-item to="/projects">Projects</b-nav-item>
      </b-navbar-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item right to="/preferences">
          <font-awesome-icon
            title="Preferences"
            :icon="['fa', 'cog']"
          /> Basedirs
        </b-nav-item>
        <b-nav-item right to="/stacks">
          <font-awesome-icon
            title="Stacks"
            :icon="['fa', 'layer-group']"
          /> Stacks
        </b-nav-item>
        <b-nav-item right to="/services">
          <font-awesome-icon
            title="Services"
            :icon="['fa', 'layer-group']"
          /> Services
        </b-nav-item>
        <b-nav-item right to="/gearspecs">
          <font-awesome-icon
            title="Gearspecs"
            :icon="['fa', 'layer-group']"
          /> Gearspecs
        </b-nav-item>

      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>
