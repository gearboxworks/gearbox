import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import './plugins/bootstrap'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faPlayCircle, faStopCircle, faTachometerAlt, faHome, faCog, faExpand } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlayCircle)
library.add(faStopCircle)
library.add(faTachometerAlt)
library.add(faHome)
library.add(faCog)
library.add(faExpand)

Vue.component('font-awesome-icon', FontAwesomeIcon)

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
