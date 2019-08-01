export default {
  install (Vue, options) {
    // eslint-disable-next-line no-param-reassign
    Vue.prototype.$_escapeIDAttr = value => value.replace(/\//g, '-').replace(/\./g, '-')
  }
}
