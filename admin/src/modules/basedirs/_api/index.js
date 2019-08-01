import BaseMethodsFactory from '../../_base/_api'
import moduleConfig from '../config'
import HTTP from '../../../http-common'

const ExtraMethods = {
  checkDirectory (dir) {
    return HTTP.head('directories/' + encodeURI(dir)).then((r) => {
      console.log('Result:', r)
      return (r)
    }).catch(e => {
      if (e.response && e.response.status) {
        return (e.response.status)
      } else {
        throw e
      }
    })
  },
  // @see basedirs/CREATE
  // createDirectory (dir) {
  //   return HTTP.post('basedirs/new' + encodeURI(dir))
  // },
  openDirectory (dir) {
    console.log('TODO: implement API method to open directory in system file explorer', dir)
    return HTTP.post('directories/' + encodeURI(dir), { data: { action: 'open' } })
  }
}

// const merged = { ...BaseMethodsFactory(moduleConfig), ...ExtraMethods }
// console.dir(merged)

export default { ...BaseMethodsFactory(moduleConfig), ...ExtraMethods }
