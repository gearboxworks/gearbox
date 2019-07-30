import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import './plugins/bootstrap'
import './plugins/fontawesome'

import HelpersPlugin from './plugins/helpers'

// import VueSelect from 'vue-cool-select'
import VueClipboard from 'vue-clipboard2'

Vue.use(VueClipboard)
Vue.use(HelpersPlugin)

Vue.config.productionTip = false

// Vue.use(VueSelect, {
//   theme: 'bootstrap' // or 'material-design'
// })

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
