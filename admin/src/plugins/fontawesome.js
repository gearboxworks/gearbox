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
  faChevronDown,
  faEllipsisV,
  faEnvelope,
  faDatabase,
  faFolder,
  faFolderOpen,
  faTrashAlt,
  faCheckCircle,
  faPlusCircle,
  faColumns,
  faThList,
  faSortAlphaUp,
  faSortAlphaDown,
  faStickyNote,
  faLayerGroup,
  faClone,
  faReply,
  faCircleNotch,
  faSpinner
} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlay)
library.add(faStop)
library.add(faTachometerAlt)
library.add(faHome)
library.add(faCog)
library.add(faExpand)
library.add(faChevronUp)
library.add(faChevronDown)
library.add(faEllipsisV)
library.add(faEnvelope)
library.add(faDatabase)
library.add(faFolder)
library.add(faFolderOpen)
library.add(faTrashAlt)
library.add(faCheckCircle)
library.add(faPlusCircle)
library.add(faColumns)
library.add(faThList)
library.add(faSortAlphaUp)
library.add(faSortAlphaDown)
library.add(faStickyNote)
library.add(faLayerGroup)
library.add(faClone)
library.add(faReply)
library.add(faCircleNotch)
library.add(faSpinner)

Vue.component('font-awesome-icon', FontAwesomeIcon)
