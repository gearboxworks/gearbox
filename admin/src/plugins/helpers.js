// import Vue from 'vue'

function escapeIDAttr (value) {
  value.replace(/\//g, '-').replace(/\./g, '-')
}

export default {
  install (Vue, options) {
    // eslint-disable-next-line no-param-reassign
    Vue.prototype.$escapeIDAttr = escapeIDAttr
  }
}
