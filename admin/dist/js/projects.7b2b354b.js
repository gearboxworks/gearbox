(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projects"],{"00b6":function(t,e,s){"use strict";var i=s("6079"),a=s.n(i);a.a},"1c80":function(t,e,s){},"1c9c":function(t,e,s){"use strict";var i=s("50d0"),a=s.n(i);a.a},"372a":function(t,e,s){"use strict";var i=s("968c"),a=s.n(i);a.a},"43f9":function(t,e,s){},"504c":function(t,e,s){var i=s("9e1e"),a=s("0d58"),o=s("6821"),r=s("52a7").f;t.exports=function(t){return function(e){var s,n=o(e),c=a(n),l=c.length,d=0,p=[];while(l>d)s=c[d++],i&&!r.call(n,s)||p.push(t?[s,n[s]]:n[s]);return p}}},"50d0":function(t,e,s){},"59a5":function(t,e,s){},6079:function(t,e,s){},6347:function(t,e,s){},7319:function(t,e,s){"use strict";var i=s("1c80"),a=s.n(i);a.a},"7d23":function(t,e,s){"use strict";var i=s("43f9"),a=s.n(i);a.a},8477:function(t,e,s){"use strict";var i=s("890c"),a=s.n(i);a.a},8700:function(t,e,s){"use strict";var i=s("baa5"),a=s.n(i);a.a},8889:function(t,e,s){"use strict";var i=s("6347"),a=s.n(i);a.a},"890c":function(t,e,s){},"8c70":function(t,e,s){},"8d25":function(t,e,s){"use strict";var i=s("8c70"),a=s.n(i);a.a},"968c":function(t,e,s){},"9deb":function(t,e,s){},a962:function(t,e,s){},acca:function(t,e,s){"use strict";s.r(e);var i=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"projects-container"},[s("projects-drawer",{attrs:{visible:"false"},on:{"switch-view-mode":t.switchViewMode}}),t.projects.length?s("div",{staticClass:"filtered-projects"},["cards"===t.viewMode?s("b-card-group",{staticClass:"pl-3 pr-3",attrs:{columns:""}},t._l(t.projects,function(t,e){return s("project-card",{key:t.id,attrs:{project:t,projectIndex:e}})}),1):s("table",{staticClass:"projects-table"},[s("thead",[s("tr",[s("th",{staticClass:"th--state"},[t._v("State")]),s("th",{staticClass:"th--hostname"},[t._v("Project Name")]),s("th",{staticClass:"th--location"},[t._v("Location")]),s("th",{staticClass:"th--stack"},[t._v("Stack")]),s("th",{staticClass:"th--notes"},[t._v("Notes")])])]),s("tbody",t._l(t.projects,function(t,e){return s("project-row",{key:t.id,attrs:{project:t,projectIndex:e}})}),1)])],1):s("div",{staticClass:"filtered-projects is-empty"},[s("h5",[t._v("No projects match the current criteria.")])])],1)},a=[],o=s("cebc"),r=s("2f62"),n=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"drawer mb-3 clearfix"},[t.expanded?s("div",{staticClass:"drawer-contents clearfix"},[s("div",{staticClass:"left-panel"},[s("b-form",{staticClass:"filter-form"},[s("b-form-checkbox-group",{staticClass:"form-group--states",attrs:{name:"showStates",label:"","label-for":"filter-state",description:"State",switches:"",stack:""},model:{value:t.showStates,callback:function(e){t.showStates=e},expression:"showStates"}},[s("b-form-checkbox",{attrs:{value:"running",title:"Include projects that are currently RUNNING"},on:{change:function(e){return t.toggleState(e,"running")}}},[t._v("Running")]),s("b-form-checkbox",{attrs:{value:"stopped",title:"Include projects that are currently STOPPED"},on:{change:function(e){return t.toggleState(e,"stopped")}}},[t._v("Stopped")]),s("b-form-checkbox",{attrs:{value:"candidates",title:"Include projects that are yet to be imported"},on:{change:function(e){return t.toggleState(e,"candidates")}}},[t._v("Candidates")]),s("small",{staticClass:"form-text text-muted",attrs:{tabindex:"-1"}},[t._v("Project State")])],1),t.hasExtraBasedirs?s("b-form-group",{staticClass:"form-group--location",attrs:{label:"","label-for":"filter-location",description:"Location"}},[s("b-select",{attrs:{id:"filter-basedirs",variant:"secondary",options:t.basedirsAsOptions},on:{change:function(e){return t.changeFilter(e,"basedir")}},model:{value:t.showBasedirs,callback:function(e){t.showBasedirs=e},expression:"showBasedirs"}},[s("template",{slot:"first"},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Show projects from...")]),s("option",{attrs:{value:"all"}},[t._v("All known locations")])])],2)],1):t._e(),s("b-form-group",{staticClass:"form-group--stacks",attrs:{label:"","label-for":"filter-stacks",description:"Used Stacks"}},[s("b-select",{attrs:{id:"filter-stacks",variant:"secondary"},on:{change:function(e){return t.changeFilter(e,"stacks")}},model:{value:t.showStacks,callback:function(e){t.showStacks=e},expression:"showStacks"}},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Filter by stacks...")]),s("option",{attrs:{value:"all"}},[t._v("Any stack")]),s("optgroup",{attrs:{label:"Specific Stacks"}},t._l(t.stacksAsOptions,function(e){return s("option",{key:e.value,domProps:{value:e.value}},[t._v(t._s(e.text.toUpperCase()))])}),0),s("option",{attrs:{value:"none"}},[t._v("No stacks assigned")])])],1),s("b-form-group",{staticClass:"form-group--programs",attrs:{label:"","label-for":"filter-programs",description:"Used Programs"}},[s("b-select",{attrs:{id:"filter-programs",variant:"secondary"},on:{change:function(e){return t.changeFilter(e,"programs")}},model:{value:t.showPrograms,callback:function(e){t.showPrograms=e},expression:"showPrograms"}},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Filter by programs...")]),s("option",{attrs:{value:"all"}},[t._v("Any program")]),s("optgroup",{attrs:{label:"Specific Program"}},t._l(t.programsAsOptions,function(e){return s("option",{key:e.value,domProps:{value:e.value}},[t._v(t._s(e.text.toUpperCase()))])}),0),s("option",{attrs:{value:"none"}},[t._v("No programs assigned")])])],1)],1)],1),s("div",{staticClass:"right-panel"},[s("b-form",{staticClass:"sort-form"},[s("b-form-group",{attrs:{id:"sort-by-group",label:"","label-for":"sort-by",description:"Sort by"}},[s("b-select",{attrs:{id:"sort-by",variant:"secondary"},model:{value:t.sortBy,callback:function(e){t.sortBy=e},expression:"sortBy"}},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Sort by...")]),s("option",{attrs:{value:"access-date",disabled:""}},[t._v("Access date")]),s("option",{attrs:{value:"creation-date",disabled:""}},[t._v("Creation date")]),s("option",{attrs:{value:"project-title"}},[t._v("Project title")])])],1),s("b-form-group",{attrs:{id:"sort-order-group",label:"","label-for":"sort-order-select",description:"Order"}},[s("a",{staticClass:"view-mode view-mode--order",attrs:{target:"_blank",href:"#",title:"Sort Order"},on:{click:function(e){e.preventDefault(),t.sortAscending=!t.sortAscending}}},[s("font-awesome-icon",{attrs:{icon:["fa",t.sortAscending?"sort-alpha-down":"sort-alpha-up"]}})],1)])],1),s("b-form",{staticClass:"view-form"},[s("b-form-group",{attrs:{id:"view-form",label:"","label-for":"view-select",description:"View Mode"}},[s("a",{class:{"view-mode":!0,"view-mode--cards":!0,"is-inactive":"cards"!=t.viewMode},attrs:{target:"_blank",href:"#",title:"Cards View"},on:{click:function(e){e.preventDefault(),t.viewMode="cards",t.$emit("switch-view-mode",e,"cards")}}},[s("font-awesome-icon",{attrs:{icon:["fa","columns"]}})],1),s("a",{class:{"view-mode":!0,"view-mode--table":!0,"is-inactive":"table"!=t.viewMode},attrs:{target:"_blank",href:"#",title:"Table View"},on:{click:function(e){e.preventDefault(),t.viewMode="table",t.$emit("switch-view-mode",e,"table")}}},[s("font-awesome-icon",{attrs:{icon:["fa","th-list"]}})],1)])],1)],1)]):t._e(),s("div",{staticClass:"drawer-handle",on:{click:function(e){t.expanded=!t.expanded}}},[s("div",{staticClass:"label small"},[s("span",[t._v("Viewing Options \n      "),t.expanded?s("font-awesome-icon",{attrs:{icon:["fa","chevron-up"]}}):s("font-awesome-icon",{attrs:{icon:["fa","chevron-down"]}})],1)]),s("div",{staticClass:"current-filter"},[s("b-badge",{attrs:{title:"Filter by state",variant:t.statesVariant}},[t._v(t._s(t.labelStates))]),t.hasExtraBasedirs?s("b-badge",{attrs:{title:"Filter by location",variant:"all"==t.showBasedirs?"secondary":"warning"}},[t._v(t._s(t.labelBasedirs))]):t._e(),s("b-badge",{attrs:{title:"Filter by used stack",variant:"all"==t.showStacks?"secondary":"warning"}},[t._v(t._s(t.labelStacks))]),s("b-badge",{attrs:{title:"Filter by used program",variant:"all"==t.showPrograms?"secondary":"warning"}},[t._v(t._s(t.labelPrograms))]),s("b-badge",{attrs:{title:"Sorting order"}},[t._v(t._s(t.labelSorting))])],1)])])},c=[],l=(s("a481"),{name:"ProjectsDrawer",props:{},computed:Object(o["a"])({},Object(r["c"])(["basedirBy","stackBy","basedirsAsOptions","stacksAsOptions","programsAsOptions","hasExtraBasedirs"]),{labelStates:function(){var t=this.showStates,e=-1!==t.indexOf("running")?"Running projects":"",s=-1!==t.indexOf("stopped")?"Stopped projects":"",i=-1!==t.indexOf("candidates"),a=e&&s?"All projects":e||s?e+s:"";return a?a+(i&&(e||s)?"":" (no candidates)"):i?"Project candidates":""},statesVariant:function(){var t=this.showStates,e=-1!==t.indexOf("running")?"Running projects":"",s=-1!==t.indexOf("stopped")?"Stopped projects":"",i=-1!==t.indexOf("candidates");return e&&s&&i?"secondary":"warning"},labelBasedirs:function(){var t="all"!==this.showBasedirs?this.basedirBy("id",this.showBasedirs):null;return"From "+(t?t.attributes.basedir:"all known locations")},labelStacks:function(){var t;if("none"===this.showStacks)t="With no stacks assigned";else{var e="all"!==this.showStacks?this.stackBy("id",this.showStacks):null;t="Using "+(e?e.attributes.stackname.toUpperCase()+" stack":"any stack")}return t},labelPrograms:function(){var t;if("none"===this.showPrograms)t="With no programs assigned";else{var e="all"!==this.showPrograms?this.showPrograms:null;t="Using "+(e?e.toUpperCase():"any program")}return t},labelSorting:function(){return"Sorted by "+this.sortBy.replace("-"," ")+(this.sortAscending?"":" (reverse)")}}),data:function(){return{expanded:!1,showStates:["running","stopped","candidates"],showBasedirs:"all",showStacks:"all",showPrograms:"all",sortBy:"project-title",sortAscending:!0,viewMode:"cards"}},methods:{toggleState:function(t,e){var s=this.showStates,i=-1!==s.indexOf("running"),a=-1!==s.indexOf("stopped"),o=-1!==s.indexOf("candidates");"candidates"!==e||i||a||!o?("running"!==e||!i||a||o)&&("stopped"!==e||i||!a||o)||(this.showStates=["candidates"]):this.showStates=["running","stopped"]},changeFilter:function(t,e){this.$store.dispatch("setProjectsFilter",{field:e,values:t})}},watch:{showStates:function(t,e){this.$store.dispatch("setProjectsFilter",{field:"states",values:this.showStates})}}}),d=l,p=(s("8477"),s("7319"),s("2877")),u=Object(p["a"])(d,n,c,!1,null,"33ea99db",null),h=u.exports,f=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("b-card",{staticClass:"card--project"},[s("div",{staticClass:"clearfix"},[s("project-hostname",{attrs:{project:t.project,projectIndex:t.projectIndex,"is-multimodal":!0},on:{"show-alert":t.showAlert}}),s("project-toolbar",{attrs:{project:t.project,projectIndex:t.projectIndex,isUpdating:t.isUpdating},on:{"run-stop":t.onRunStop}})],1),s("b-alert",{attrs:{show:t.alertShow,dismissible:t.alertDismissible,variant:t.alertVariant,fade:""},on:{dismissed:function(e){t.alertShow=!1}}},[t._v(t._s(t.alertContent))]),s("div",{staticClass:"clearfix",attrs:{slot:"footer"},slot:"footer"},[s("project-stack-list",{attrs:{project:t.project,projectIndex:t.projectIndex}}),s("project-stack-add",{attrs:{project:t.project,projectIndex:t.projectIndex},on:{"maybe-hide-alert":t.maybeHideAlert}}),s("project-location",{attrs:{project:t.project,projectIndex:t.projectIndex}}),s("project-note",{attrs:{project:t.project,projectIndex:t.projectIndex}})],1)],1)},b=[],m=(s("c5f6"),function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("ul",{staticClass:"toolbar-list"},[s("li",{staticClass:"toolbar-item"},[s("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link toolbar-link--state",attrs:{target:"_blank",href:"#",title:t.isSwitching?"Switching state...":t.isRunning?"Stop all services":"Run all services"},on:{click:function(e){return e.preventDefault(),t.$emit("run-stop")}}},[t.isSwitching?s("font-awesome-icon",{attrs:{icon:["fa","circle-notch"],spin:""}}):s("font-awesome-icon",{attrs:{icon:["fa",t.isRunning?"stop":"play"]}})],1)])])}),j=[],v={name:"ProjectToolbar",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},isUpdating:{type:Boolean,required:!0}},data:function(){return{id:this.project.id,hostname:this.project.attributes.hostname}},computed:{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},isRunning:function(){return this.project.attributes.enabled},isSwitching:function(){return this.isUpdating}},methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")}}},g=v,k=(s("c3b6"),Object(p["a"])(g,m,j,!1,null,"34920594",null)),w=k.exports,C=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("b-input-group",{class:{"input-group--hostname":!0,"is-editing":t.isEditing,"is-multimodal":t.isMultimodal,"is-modified":t.isModified,"is-updating":t.isUpdating},attrs:{id:t.projectBase+"hostname-group",role:"tabpanel"}},[s("b-form-input",{staticClass:"hostname-input",attrs:{id:"hostname-input-"+t.projectIndex,type:"text",size:t.isMultimodal?"lg":"md",required:"",placeholder:"",readonly:t.isUpdating||t.project.attributes.enabled,autocomplete:"off"},on:{click:function(e){return e.preventDefault(),t.onInputClicked(e)}},model:{value:t.hostname,callback:function(e){t.hostname=e},expression:"hostname"}}),t.isEditing?s("b-input-group-append",[s("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"btn--submit",attrs:{variant:"outline-info",title:t.isModified?"Submit the new hostname":t.isMultimodal?"Cancel":"No changes to submit",disabled:!t.isMultimodal&&!t.isModified||t.isUpdating},on:{click:function(e){return e.preventDefault(),t.onButtonClicked(e)}}},[t.isUpdating?s("font-awesome-icon",{attrs:{icon:"circle-notch",spin:""}}):s("font-awesome-icon",{attrs:{icon:["fa",t.isMultimodal?t.isModified?"check":"times":"check"]}})],1)],1):t._e()],1)},y=[],S={name:"ProjectLocation",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},isMultimodal:{type:Boolean,required:!1,default:!1}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes,{isEditing:!this.isMultimodal,isModified:!1,isUpdating:!1})},computed:Object(o["a"])({},Object(r["c"])({basedirBy:"basedirBy"}),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},currentBasedir:function(){var t=this.basedirBy("id",this.basedir);return t?t.attributes.basedir:""}}),methods:Object(o["a"])({},Object(r["b"])(["updateProjectHostname"]),{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},resolveDir:function(t,e){return t+(-1!==t.indexOf("/")?"/":"\\")+e},onInputClicked:function(){this.isEditing||(this.project.attributes.enabled?this.$emit("show-alert","Hostname cannot be changed while the project is running!"):this.isMultimodal&&(this.isEditing=!0))},onButtonClicked:function(){this.isModified?(console.log("TODO call the API to change project hostname"),this.maybeSubmit(this.hostname)):this.isEditing=!1},maybeSubmit:function(t){var e=this;this.isUpdating=!0,this.updateProjectHostname({projectId:this.id,hostname:t}).then(function(){e.isMultimodal&&(e.isEditing=!1),e.isModified=!1,e.isUpdating=!1})}}),watch:{hostname:function(t,e){this.isModified=!!t&&t!==this.project.id}}},x=S,_=(s("372a"),Object(p["a"])(x,C,y,!1,null,"09d51d36",null)),P=_.exports,O=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("b-input-group",{class:{"input-group--location":!0,"is-collapsed":t.isCollapsed,"is-multimodal":t.isMultimodal},attrs:{id:t.projectBase+"location",role:"tabpanel"}},[s("b-form-input",{directives:[{name:"show",rawName:"v-show",value:!t.isCollapsed,expression:"!isCollapsed"}],ref:t.projectBase+"location",staticClass:"location-input",attrs:{readonly:"",value:t.resolveDir(t.currentBasedir,t.path),autocomplete:"off"},on:{keyup:function(e){if(!e.type.indexOf("key")&&t._k(e.keyCode,"esc",27,e.key,["Esc","Escape"]))return null;t.isCollapsed=!0}}}),s("b-input-group-append",[t.isCollapsed?t._e():s("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"btn--copy-dir",attrs:{variant:"outline-info",title:"Copy to clipboard",href:"#"},on:{click:t.onCopyToClipboard}},[s("font-awesome-icon",{attrs:{icon:["fa","clone"]}})],1),s("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"btn--open-dir",attrs:{variant:"outline-info",title:t.isCollapsed?"View project location":"Open in file manager",href:"#"},on:{click:t.onButtonClicked}},[s("font-awesome-icon",{attrs:{icon:["fa",t.isCollapsed?"folder":"folder-open"]}})],1)],1)],1)},B=[],A={name:"ProjectLocation",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},isMultimodal:{type:Boolean,required:!1,default:!0}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes,{isCollapsed:this.isMultimodal})},computed:Object(o["a"])({},Object(r["c"])({basedirBy:"basedirBy"}),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},currentBasedir:function(){var t=this.basedirBy("id",this.basedir);return t?t.attributes.basedir:""}}),methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},resolveDir:function(t,e){return t+(-1!==t.indexOf("/")?"/":"\\")+e},onButtonClicked:function(){var t=this;this.isMultimodal&&this.isCollapsed?(this.isCollapsed=!1,this.$nextTick(function(){var e=t.$refs["".concat(t.projectBase,"location")].$el;e.focus(),e.setSelectionRange(0,9999)})):(console.log("TODO: call API method to open directory in file manager"),this.isMultimodal&&this.$nextTick(function(){t.isCollapsed=!0}))},onCopyToClipboard:function(){var t=this;console.log("TODO: implement copy to clipboard"),this.isMultimodal&&this.$nextTick(function(){t.isCollapsed=!0})}}},I=A,M=(s("00b6"),Object(p["a"])(I,O,B,!1,null,"82a43bba",null)),U=M.exports,N=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("b-input-group",{class:{"input-group--note":!0,"is-collapsed":t.isCollapsed,"is-modified":t.isModified,"is-updating":t.isUpdating,"is-empty":!!t.notes},attrs:{role:"tabpanel"}},[t.isCollapsed?t._e():s("b-form-input",{ref:t.projectBase+"note-input",staticClass:"notes-input",attrs:{placeholder:"Add note...",readonly:t.isUpdating,autocomplete:"off",autofocus:""},on:{keyup:function(e){if(!e.type.indexOf("key")&&t._k(e.keyCode,"esc",27,e.key,["Esc","Escape"]))return null;t.isCollapsed=!0}},model:{value:t.notes,callback:function(e){t.notes=e},expression:"notes"}}),s("b-input-group-append",[s("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],class:{"btn--submit":!0,"btn--add":t.isCollapsed},attrs:{variant:t.isCollapsed&&t.notes?"outline-warning":"outline-info",title:t.isCollapsed?t.notes?t.notes:"Add a note":t.isModified?"Submit the note":"Please enter some text first or Click to cancel",disabled:t.isUpdating},on:{click:function(e){return e.preventDefault(),t.onButtonClicked(e)}}},[t.isUpdating?s("font-awesome-icon",{attrs:{icon:"circle-notch",spin:""}}):s("font-awesome-icon",{attrs:{icon:["fa",t.isCollapsed?"sticky-note":t.isModified?"check":"times"]}}),t.isUpdating||t.notes?t._e():s("span",[t._v(t._s(t.isCollapsed?"+":""))])],1)],1)],1)},$=[],E={name:"ProjectNote",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes,{isCollapsed:!0,isModified:!1,isUpdating:!1})},computed:{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"}},methods:Object(o["a"])({},Object(r["b"])(["addProjectNote"]),{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},onButtonClicked:function(){this.isCollapsed?this.isCollapsed=!1:this.isModified?this.maybeSubmit():this.isCollapsed=!0},maybeSubmit:function(){var t=this;this.isUpdating=!0,this.addProjectNote({projectId:this.id,text:this.notes}).then(function(){t.isCollapsed=!0,t.isModified=!1,t.isUpdating=!1})}}),watch:{notes:function(t,e){this.isModified=!!t}}},D=E,R=(s("7d23"),Object(p["a"])(D,N,$,!1,null,"03ac824b",null)),q=R.exports,V=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("b-input-group",{class:{"input-group--stack":!0,"is-collapsed":t.isCollapsed,"is-modified":t.isModified,"is-updating":t.isUpdating},attrs:{id:t.projectBase+"stack",role:"tabpanel"}},[s("b-form-select",{directives:[{name:"show",rawName:"v-show",value:!t.isCollapsed,expression:"!isCollapsed"}],ref:t.projectBase+"-select",staticClass:"select-stack",attrs:{disabled:!t.hasStacksNotInProject||t.isUpdating,required:!0,autofocus:""},on:{change:function(e){t.isModified=!0}},model:{value:t.selectedStack,callback:function(e){t.selectedStack=e},expression:"selectedStack"}},[s("option",{attrs:{value:"",disabled:""}},[t._v(t._s(t.hasStacksNotInProject?"Add stack...":"All stacks already added"))]),t._l(t.stacksNotInProject,function(e,i){return s("option",{key:i,domProps:{value:i}},[t._v(t._s(e.stack.attributes.stackname+(e.isRemoved?" (removed)":"")+(e.isDefault?" (default)":"")))])})],2),s("b-input-group-append",[s("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],class:{"btn--submit":!0,"btn--add":t.isCollapsed},attrs:{variant:"outline-info",title:t.isUpdating?"Updating...":t.isCollapsed?"Add a stack":t.isModified?"Add the selected stack":"Please select some stack first or Click to cancel",disabled:t.isUpdating},on:{click:function(e){return e.preventDefault(),t.onButtonClicked(e)}}},[t.isUpdating?s("font-awesome-icon",{attrs:{icon:"circle-notch",spin:""}}):s("font-awesome-icon",{attrs:{icon:["fa",t.isCollapsed?"layer-group":t.isModified?"check":"times"]}}),s("span",[t._v(t._s(t.isCollapsed&&!t.isUpdating?"+":""))])],1)],1)],1)},T=[],F=(s("ac6a"),s("ffc1"),{name:"ProjectStackAdd",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes,{selectedStack:"",isCollapsed:!0,isModified:!1,isUpdating:!1})},computed:Object(o["a"])({},Object(r["c"])({serviceBy:"serviceBy",gearspecBy:"gearspecBy",allGearspecs:"gearspecs/all",allStacks:"stacks/all",hasExtraBasedirs:"hasExtraBasedirs"}),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},stacksNotInProject:function(){var t={},e=this.projectGearsGroupedByStack;for(var s in this.allStacks){var i=this.allStacks[s];"undefined"===typeof e[i.id]?t[i.id]={stack:i,isRemoved:!1}:e[i.id].isRemoved&&(t[i.id]={stack:i,isRemoved:!1,isDefault:!0},t[i.id+"(removed)"]={stack:i,isRemoved:!0})}return t},hasStacksNotInProject:function(){return Object.entries(this.stacksNotInProject).length>0},servicesInProject:function(){for(var t={},e=0;e>this.stack.length;e++)if(!this.stack[e].isRemoved){var s=this.serviceBy("id",this.stack[e].service_id);s&&(t[this.stack[e].service_id]=s)}return t},gearsInProject:function(){for(var t={},e=0;e>this.stack.length;e++){var s=this.gearspecBy("id",this.stack[e].gearspec_id);s&&(t[this.stack[e].gearspec_id]=s)}return t},projectGearsGroupedByStack:function(){var t=this,e={};return this.project.attributes.stack&&this.project.attributes.stack.forEach(function(s,i){var a=t.gearspecBy("id",s.gearspec_id),o=t.serviceBy("id",s.service_id);a&&o&&("undefined"===typeof e[a.attributes.stack_id]&&(e[a.attributes.stack_id]={isRemoved:s.isRemoved||!1}),e[a.attributes.stack_id][a.attributes.role]=o)}),e}}),methods:Object(o["a"])({},Object(r["b"])(["addProjectStack"]),{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},maybeAddProjectStack:function(t){var e=this;t&&(this.isUpdating=!0,this.addProjectStack({projectId:this.id,stackId:t}).then(function(){e.isUpdating=!1,e.isCollapsed=!0,e.selectedStack="",e.isModified=!1,e.$emit("maybe-hide-alert","Please add some stacks first!")}))},onButtonClicked:function(){var t=this;this.isCollapsed?(this.isCollapsed=!1,this.$nextTick(function(){t.$refs["".concat(t.projectBase,"-select")].$el.focus()})):this.isModified?this.maybeAddProjectStack(this.selectedStack):this.isCollapsed=!0}})}),L=F,H=(s("1c9c"),Object(p["a"])(L,V,T,!1,null,"0748f145",null)),G=H.exports,J=s("2232"),W={name:"ProjectCard",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes,{showingDetails:!1,alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"warning",isUpdating:!1})},components:{ProjectToolbar:w,ProjectHostname:P,ProjectLocation:U,ProjectStackList:J["default"],ProjectStackAdd:G,ProjectNote:q},computed:{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},isRunning:function(){return this.project.attributes.enabled}},methods:Object(o["a"])({},Object(r["b"])(["changeProjectState"]),{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},showAlert:function(t){"string"===typeof t?this.alertContent=t:(this.alertVariant=t.variant||this.alertVariant,this.alertDismissible=t.dismissible||this.alertDismissible,this.alertContent=t.content||this.alertContent),this.alertShow=!0},hideAlert:function(){this.alertContent="",this.alertShow=!1},maybeHideAlert:function(t){this.alertContent===t&&this.hideAlert()},onRunStop:function(){var t=this;this.project.attributes.stack&&this.project.attributes.stack.length>0?(this.isUpdating=!0,this.changeProjectState({projectId:this.id,isEnabled:!this.isRunning}).then(function(e){t.isUpdating=!1,t.hideAlert()})):this.showAlert("Please add some stacks first!")}})},z=W,K=(s("ef09"),s("8700"),Object(p["a"])(z,f,b,!1,null,"3c08b214",null)),Q=K.exports,X=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("tr",{staticClass:"row--project"},[s("td",{staticClass:"td--state"},[s("project-toolbar",{attrs:{project:t.project,projectIndex:t.projectIndex,"is-updating":t.isUpdating},on:{"run-stop":t.onRunStop}})],1),s("td",{staticClass:"td--hostname"},[s("project-hostname",{attrs:{project:t.project,projectIndex:t.projectIndex,"is-multimodal":!1},on:{"show-alert":t.showAlert}})],1),s("td",{staticClass:"td--location"},[s("project-location",{attrs:{project:t.project,projectIndex:t.projectIndex,"is-multimodal":!1}})],1),s("td",{staticClass:"td--stack"},[s("project-stack-list",{attrs:{project:t.project,projectIndex:t.projectIndex,"start-collapsed":!0}}),s("project-stack-add",{attrs:{project:t.project,projectIndex:t.projectIndex}})],1),s("td",{staticClass:"td--notes"},[s("project-note",{attrs:{project:t.project,projectIndex:t.projectIndex}})],1)])},Y=[],Z={name:"ProjectRow",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes,{alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"warning",isUpdating:!1})},components:{ProjectHostname:P,ProjectToolbar:w,ProjectLocation:U,ProjectNote:q,ProjectStackAdd:G,ProjectStackList:J["default"]},computed:{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},isRunning:function(){return this.project.attributes.enabled}},methods:Object(o["a"])({},Object(r["b"])(["changeProjectState"]),{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},showAlert:function(t){"string"===typeof t?this.alertContent=t:(this.alertVariant=t.variant||this.alertVariant,this.alertDismissible=t.dismissible||this.alertDismissible,this.alertContent=t.content||this.alertContent),this.alertShow=!0},maybeSubmit:function(t){this.$store.dispatch("updateProject",{id:this.id,attributes:this.$data}).then(function(){})},onRunStop:function(){var t=this;this.project.attributes.stack&&this.project.attributes.stack.length>0?(this.isUpdating=!0,this.changeProjectState({projectId:this.id,isEnabled:!this.isRunning}).then(function(e){t.isUpdating=!1})):this.showAlert("Please add some stacks first!")}})},tt=Z,et=(s("de0a"),s("8889"),Object(p["a"])(tt,X,Y,!1,null,"18af4de8",null)),st=et.exports,it={name:"ProjectList",data:function(){return{viewMode:"cards"}},components:{ProjectsDrawer:h,ProjectCard:Q,ProjectRow:st},computed:Object(o["a"])({},Object(r["c"])({projects:"filteredProjects"})),methods:{switchViewMode:function(t,e){this.viewMode=e}},mounted:function(){var t=this;this.$store.dispatch("basedirs/loadAll").then(function(){}),this.$store.dispatch("stacks/loadAll").then(function(){}),this.$store.dispatch("services/loadAll").then(function(){}),this.$store.dispatch("gearspecs/loadAll").then(function(){}),this.$store.dispatch("projects/loadAll").then(function(){}).then(function(){t.$store.dispatch("loadProjectDetails")})}},at=it,ot=(s("8d25"),Object(p["a"])(at,i,a,!1,null,"76f276e6",null));e["default"]=ot.exports},baa5:function(t,e,s){},c3b6:function(t,e,s){"use strict";var i=s("59a5"),a=s.n(i);a.a},de0a:function(t,e,s){"use strict";var i=s("a962"),a=s.n(i);a.a},ef09:function(t,e,s){"use strict";var i=s("9deb"),a=s.n(i);a.a},ffc1:function(t,e,s){var i=s("5ca1"),a=s("504c")(!0);i(i.S,"Object",{entries:function(t){return a(t)}})}}]);
//# sourceMappingURL=projects.7b2b354b.js.map