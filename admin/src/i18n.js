import Vue from 'vue'
import VueI18n from 'vue-i18n'

Vue.use(VueI18n)

const messages = {
  'en': {
    label: {
      projects: 'Projects',
      basedirs: 'Basedirs',
      services: 'Services',
      stacks: 'Stacks',
      gearspecs: 'Gearspecs',
      preferences: 'Preferences'
    },
    message: {
      welcome: 'Welcome',
      noConnectionTitle: 'Connection Problem',
      noConnectionDesc: 'It seems that Gearbox Server is not running.',
      noConnectionRetry: 'Retry?'
    },
    projects: {
      noResults: 'No projects match the current criteria.',
      filterHeading: 'Projects:',
      filterRunning: 'Running',
      filterRunningTitle: 'Include projects that are currently RUNNING',
      filterStopped: 'Stopped',
      filterStoppedTitle: 'Rodyti projektus kurie šiuo metu yra SUSTABDYTI',
      filterCandidates: 'Candidates',
      filterCandidatesTitle: 'Include projects that are yet to be imported',
      filterByLocation: 'Filter projects by location',
      filterAllLocations: 'All locations',
      filterByStack: 'Filter by used stack',
      filterAllStacks: 'All stacks',
      filterByProgram: 'Filter by used program',
      filterAllPrograms: 'All programs',
      sortBy: 'Sort by',
      sortOrder: 'Sort order',
      sortAscending: 'Sort in ascending order',
      sortDescending: 'Sort in descending order',
      sortByCreation: 'Creation Date',
      sortByAccess: 'Access Date',
      sortById: 'Project Title',
      viewingCards: 'Cards view',
      viewingTable: 'Table view',
      viewAsCards: 'Switch to cards view',
      viewAsTable: 'Switch to table view',
      fieldHostname: 'Project Name',
      fieldHostnameSubmit: 'Submit the new hostname',
      fieldHostnameUnmodified: 'No changes to submit',
      fieldHostnameReadonly: 'Hostname cannot be changed while the project is running!',
      fieldLocation: 'Location',
      fieldLocationView: 'View project location',
      fieldLocationCopy: 'Copy to clipboard',
      fieldLocationOpen: 'Open in file manager',
      fieldNotes: 'Notes',
      fieldNotesView: 'View notes',
      fieldNotesAdd: 'Add notes',
      fieldNotesEdit: 'Edit notes',
      fieldNotesCancel: 'Cancel changes!',
      fieldNotesHide: 'Hide notes',
      fieldNotesRestore: 'Restore notes',
      fieldNotesDelete: 'Delete notes',
      fieldStack: 'Services',
      fieldStackAdd: 'Add stack...',
      fieldStackAddSome: 'Please add some stacks first!',
      fieldStackAddOne: 'Add a stack',
      fieldStackAddSelected: 'Add the selected stack',
      fieldStackUnmodified: 'Please select some stack first or Click to cancel',
      fieldStackAllAdded: 'All stacks already added',
      fieldState: 'Project State',
      fieldStateSwitching: 'Switching state...',
      fieldStateStop: 'Stop all services',
      fieldStateRun: 'Run all services'
    },
    gearspecs: {
      versionMismatch: 'Could not find the requested version (v.{0}), will use the closest match (v.{1}) instead.'
    },
    services: {
      doNotRun: 'Do not run this service',
      select: 'Select service...',
      readonlyWhileRunning: 'Note, you cannot change this service while the project is running!',
      change: 'Change service',
      open: 'Open {0}{1}',
      notRunning: '(not running)',
      frontend: 'Frontend',
      dashboard: 'Dashboard'
    },
    stacks: {
      hide: 'Hide services',
      show: 'Show services',
      moreActions: 'More actions',
      moreStackActions: 'More stack actions',
      remove: 'Remove stack',
      readonlyWhileRunning: 'Cannot remove stack while the project is running!'
    },
    button: {
      cancel: 'Cancel',
      addNotes: 'Add Notes',
      addStack: 'Add Stack',
      viewLocation: 'View project location',
      saveChanges: 'Save changes',
      makeChanges: 'Make some changes first',
      close: 'Close'
    },
    process: {
      loading: 'Loading...',
      updating: 'Updating...',
      deleting: 'Deleting...',
      restoring: 'Restoring...'
    }
  },
  'lt': {
    label: {
      projects: 'Projektai',
      basedirs: 'Direktorijos',
      services: 'Servisai',
      stacks: 'Stekai',
      gearspecs: 'Gearspekai',
      preferences: 'Nustatymai'
    },
    message: {
      welcome: 'Sveiki',
      noConnectionTitle: 'Nepavyko prisijungti',
      noConnectionDesc: 'Panašu, kad Gearbox serveris nepaleistas.',
      noConnectionRetry: 'Tikrinti dar kartą?'

    },
    projects: {
      noResults: 'Nerasta projektų, kurie atitiktų pasirinktus kriterijus.',
      filterHeading: 'Projektai:',
      filterRunning: 'Paleisti',
      filterRunningTitle: 'Rodyti projektus, kurie šiuo metu yra PALEISTI',
      filterStopped: 'Sustadyti',
      filterStoppedTitle: 'Rodyti projektus, kurie šiuo metu yra SUSTABDYTI',
      filterCandidates: 'Neimportuoti',
      filterCandidatesTitle: 'Rodyti projektus, kurie dar neimportuoti',
      filterByLocation: 'Filtruoti pagal vietą diske',
      filterAllLocations: 'Visos vietos',
      filterByStack: 'Filtruoti pagal naudojamą steką',
      filterAllStacks: 'Visi stekai',
      filterByProgram: 'Filtruoti pagal naudojamas programas',
      filterAllPrograms: 'Visos programos',
      sortBy: 'Rikiuoti pagal',
      sortOrder: 'Rikiavimo tvarka',
      sortAscending: 'Rikiuoti didėjimo tvarka',
      sortDescending: 'Rikiuoti mažėjimo tvarka',
      sortByCreation: 'Sukūrimo laiką',
      sortByAccess: 'Paleidimo laiką',
      sortById: 'Projekto pavadinimą',
      viewingCards: 'Projektai rodomi kortelėse',
      viewingTable: 'Projektai rodomi lentelėje',
      viewAsCards: 'Rodyti kaip korteles',
      viewAsTable: 'Rodyti kaip lentelę',
      fieldHostname: 'Pavadinimas',
      fieldHostnameSubmit: 'Išsaugoti naująjį pavadinimą',
      fieldHostnameUnmodified: 'Niekas nepakeista',
      fieldHostnameReadonly: 'Pavadinimas negali būti keičiamas kol projektas yra paleistas!',
      fieldLocation: 'Vieta diske',
      fieldLocationView: 'Žiūrėti projekto vietą',
      fieldLocationCopy: 'Kopijuoti į atmintinę',
      fieldLocationOpen: 'Atidaryti failų tvarkyklėje',
      fieldNotes: 'Užrašai',
      fieldNotesView: 'Peržiūrėti užrašus',
      fieldNotesAdd: 'Pridėti užrašus',
      fieldNotesEdit: 'Redaguoti užrašus',
      fieldNotesCancel: 'Atšaukti pakeitimus!',
      fieldNotesHide: 'Slėpti užrašus',
      fieldNotesRestore: 'Atstatyti užrašus',
      fieldNotesDelete: 'Pašalinti užrašus',
      fieldStack: 'Servisai',
      fieldStackAdd: 'Pridėti servisus...',
      fieldStackAddSome: 'Pirmiau prašom pridėti kažkokius servisus!',
      fieldStackAddOne: 'Pridėti servisus',
      fieldStackAddSelected: 'Pridėti pasirinktus servisus',
      fieldStackUnmodified: 'Prašom pirmiau pasirinkti steką arba spauskite Atšaukti',
      fieldStackAllAdded: 'Visi įmanomi stekai jau pridėti',
      fieldState: 'Būsena',
      fieldStateSwitching: 'Perjungiama būsena...',
      fieldStateStop: 'Stabdyti visus servisus',
      fieldStateRun: 'Paleisti visus servisus'
    },
    gearspecs: {
      versionMismatch: 'Nepavyko rasti tikslios nurodytos versijos (v.{0}), vietoje to bus naudojama artimiausia suderinama versija (v.{1}).'
    },
    services: {
      doNotRun: 'Nenaudoti šio serviso',
      select: 'Pasirinkite servisą...',
      readonlyWhileRunning: 'Negalima keisti serviso, kol projektas paleistas!',
      change: 'Keisti servisą',
      open: 'Atverti {0}{1}',
      notRunning: '(sustabdyta)',
      frontend: 'Svetainę',
      dashboard: 'Administracinę dalį'
    },
    stacks: {
      hide: 'Slėpti servisus',
      show: 'Rodyti servisus',
      moreActions: 'Daugiau veiksmų',
      moreStackActions: 'Daugiau veiksmų',
      remove: 'Šalinti steką',
      readonlyWhileRunning: 'Negalima šalinti steko, kol projektas paleistas!'
    },
    button: {
      cancel: 'Atšaukti',
      addNotes: 'Pridėti užrašą',
      addStack: 'Pridėti steką',
      viewLocation: 'Žiūrėti projekto vietą',
      saveChanges: 'Išsaugoti pakeitimus',
      makeChanges: 'Pirmiau kažką pakeiskite',
      close: 'Užverti'
    },
    process: {
      loading: 'Parsiunčiama...',
      updating: 'Saugoma...',
      deleting: 'Šalinama...',
      restoring: 'Atstatoma...'
    }
  }
}

const i18n = new VueI18n({
  locale: 'en', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages // set locale messages
})

export default i18n
