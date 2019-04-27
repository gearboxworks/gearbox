import Vue from 'vue'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faPlayCircle, faStopCircle, faTachometerAlt, faHome, faCog, faExpand, faChevronUp, faEllipsisH } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlayCircle)
library.add(faStopCircle)
library.add(faTachometerAlt)
library.add(faHome)
library.add(faCog)
library.add(faExpand)
library.add(faChevronUp)
library.add(faEllipsisH)

Vue.component('font-awesome-icon', FontAwesomeIcon)
