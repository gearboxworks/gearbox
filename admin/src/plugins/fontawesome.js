import Vue from 'vue'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faPlay, faStop, faTachometerAlt, faHome, faCog, faExpand, faChevronUp, faEllipsisH, faEllipsisV, faEnvelope, faDatabase } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faPlay)
library.add(faStop)
library.add(faTachometerAlt)
library.add(faHome)
library.add(faCog)
library.add(faExpand)
library.add(faChevronUp)
library.add(faEllipsisH)
library.add(faEllipsisV)
library.add(faEnvelope)
library.add(faDatabase)

Vue.component('font-awesome-icon', FontAwesomeIcon)
