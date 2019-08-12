import Vue from 'vue'
import VueClipboard from 'vue-clipboard2'

// import VueI18n from 'vue-i18n'
import i18n from './i18n'

import App from './App.vue'
import router from './router'
import store from './store'

import './plugins/bootstrap'
import './plugins/fontawesome'
import './plugins/treeselect'

Vue.use(VueClipboard)

Vue.config.productionTip = false

// Vue.use(VueSelect, {
//   theme: 'bootstrap' // or 'material-design'
// })

new Vue({
  i18n,
  router,
  store,
  render: h => h(App)
}).$mount('#app')

// Vue.use(VueI18n)
