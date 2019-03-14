import Vue from 'vue'
import Router from 'vue-router'
import Welcome from './views/Welcome.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'welcome',
      component: Welcome
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/About.vue')
    },
    {
      path: '/preferences',
      name: 'preferences',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "preferences" */ './views/Preferences.vue')
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import(/* webpackChunkName: "project" */ './views/Projects.vue')
      // children: [
      //   {
      //     path: ':hostname',
      //     component: () => import(/* webpackChunkName: "project" */ './components/ProjectDetails.vue')
      //   }
      // ]
    },
    {
      path: '/stack/:stackName',
      name: 'stack',
      component: () => import(/* webpackChunkName: "stack" */ './components/StackDetails.vue')
    },
    {
      path: '/gear/:gearName',
      name: 'gear',
      component: () => import(/* webpackChunkName: "gear" */ './components/GearDetails.vue')
    }

  ]
})
