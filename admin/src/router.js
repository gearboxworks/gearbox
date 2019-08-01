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
      // component: () => import(/* webpackChunkName: "preferences" */ './views/Preferences.vue')
      component: () => import(/* webpackChunkName: "basedirs" */ './modules/basedirs/index.vue')
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import(/* webpackChunkName: "projects" */ './modules/projects/index.vue')
      // children: [
      //   {
      //     path: ':hostname/stack',
      //     component: () => import(/* webpackChunkName: "projectstack" */ './modules/projects/_components/shared/ProjectStackList.vue')
      //   }
      // ]
    },
    {
      path: '/gearspecs',
      name: 'gearspecs',
      component: () => import(/* webpackChunkName: "gearspecs" */ './modules/gearspecs/index.vue')
    },
    {
      path: '/stacks',
      name: 'stacks',
      component: () => import(/* webpackChunkName: "stacks" */ './modules/stacks/index.vue')
    },
    {
      path: '/services',
      name: 'services',
      component: () => import(/* webpackChunkName: "services" */ './modules/services/index.vue')
    }
  ]
})
