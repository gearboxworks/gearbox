import Vue from 'vue'
import { library } from '@fortawesome/fontawesome-svg-core'
import {
  faPlay,
  faStop,
  faCheck,
  faTachometerAlt,
  faHome,
  faCog,
  faExpand,
  faChevronUp,
  faChevronRight,
  faChevronDown,
  faChevronLeft,
  faEllipsisV,
  faEnvelope,
  faDatabase,
  faFolder,
  faFolderOpen,
  faTrashAlt,
  faCheckCircle,
  faPlus,
  faColumns,
  faThList,
  faSortAlphaUp,
  faSortAlphaDown,
  faStickyNote,
  faLayerGroup,
  faClone,
  faTimes,
  faReply,
  faCircleNotch
} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlay)
library.add(faStop)
library.add(faCheck)
library.add(faTachometerAlt)
library.add(faHome)
library.add(faCog)
library.add(faExpand)
library.add(faChevronUp)
library.add(faChevronRight)
library.add(faChevronDown)
library.add(faChevronLeft)
library.add(faEllipsisV)
library.add(faEnvelope)
library.add(faDatabase)
library.add(faFolder)
library.add(faFolderOpen)
library.add(faTrashAlt)
library.add(faCheckCircle)
library.add(faPlus)
library.add(faColumns)
library.add(faThList)
library.add(faSortAlphaUp)
library.add(faSortAlphaDown)
library.add(faStickyNote)
library.add(faLayerGroup)
library.add(faClone)
library.add(faTimes)
library.add(faReply)
library.add(faCircleNotch)

Vue.component('font-awesome-icon', FontAwesomeIcon)
