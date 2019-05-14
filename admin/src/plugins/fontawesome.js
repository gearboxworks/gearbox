import Vue from 'vue'
import { library } from '@fortawesome/fontawesome-svg-core'
import {
  faPlay,
  faStop,
  faTachometerAlt,
  faHome,
  faCog,
  faExpand,
  faChevronUp,
  faEllipsisV,
  faEnvelope,
  faDatabase,
  faFolder,
  faTrashAlt,
  faCheckCircle,
  faPlusCircle,
  faColumns,
  faThList,
  faSortAlphaUp,
  faSortAlphaDown
} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlay)
library.add(faStop)
library.add(faTachometerAlt)
library.add(faHome)
library.add(faCog)
library.add(faExpand)
library.add(faChevronUp)
library.add(faEllipsisV)
library.add(faEnvelope)
library.add(faDatabase)
library.add(faFolder)
library.add(faTrashAlt)
library.add(faCheckCircle)
library.add(faPlusCircle)
library.add(faColumns)
library.add(faThList)
library.add(faSortAlphaUp)
library.add(faSortAlphaDown)


Vue.component('font-awesome-icon', FontAwesomeIcon)
