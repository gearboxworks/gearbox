(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projects"],{"018e":function(t,e,s){},"0baa":function(t,e,s){"use strict";const r=s("df7c"),a=s("4845"),n=s("d40f"),o=s("6d1c"),i=100,c=/[\u0000-\u001f\u0080-\u009f]/g,l=/^\.+/,d=(t,e={})=>{if("string"!==typeof t)throw new TypeError("Expected a string");const s=void 0===e.replacement?"!":e.replacement;if(n().test(s)&&c.test(s))throw new Error("Replacement string cannot contain reserved filename characters");return t=t.replace(n(),s),t=t.replace(c,s),t=t.replace(l,s),s.length>0&&(t=a(t,s),t=t.length>1?o(t,s):t),t=n.windowsNames().test(t)?t+s:t,t=t.slice(0,i),t};d.path=((t,e)=>{return t=r.resolve(t),r.join(r.dirname(t),d(r.basename(t),e))}),t.exports=d},"1c80":function(t,e,s){},"1f2e":function(t,e,s){},"254e":function(t,e,s){},"414b":function(t,e,s){"use strict";var r=s("1f2e"),a=s.n(r);a.a},4845:function(t,e,s){"use strict";var r=s("a318");t.exports=function(t,e){if("string"!==typeof t||"string"!==typeof e)throw new TypeError("Expected a string");return t.replace(new RegExp("(?:"+r(e)+"){2,}","g"),e)}},"504c":function(t,e,s){var r=s("0d58"),a=s("6821"),n=s("52a7").f;t.exports=function(t){return function(e){var s,o=a(e),i=r(o),c=i.length,l=0,d=[];while(c>l)n.call(o,s=i[l++])&&d.push(t?[s,o[s]]:o[s]);return d}}},"6cd5":function(t,e,s){"use strict";var r=s("254e"),a=s.n(r);a.a},"6d1c":function(t,e,s){"use strict";var r=s("a318");t.exports=function(t,e){if("string"!==typeof t||"string"!==typeof e)throw new TypeError;return e=r(e),t.replace(new RegExp("^"+e+"|"+e+"$","g"),"")}},7319:function(t,e,s){"use strict";var r=s("1c80"),a=s.n(r);a.a},8760:function(t,e,s){"use strict";var r=s("d9e3"),a=s.n(r);a.a},"98b6":function(t,e,s){},a318:function(t,e,s){"use strict";var r=/[|\\{}()[\]^$+*?.]/g;t.exports=function(t){if("string"!==typeof t)throw new TypeError("Expected a string");return t.replace(r,"\\$&")}},a4e2:function(t,e,s){"use strict";var r=s("018e"),a=s.n(r);a.a},acca:function(t,e,s){"use strict";s.r(e);var r=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"projects-container"},[s("projects-drawer",{attrs:{visible:"false"},on:{"switch-view-mode":t.switchViewMode}}),"cards"===t.viewMode?s("b-card-group",{staticClass:"pl-3 pr-3",attrs:{columns:""}},t._l(t.projects,function(t,e){return s("project-card",{key:t.id,attrs:{project:t,projectIndex:e}})}),1):s("table",[s("thead",[s("tr",[s("th",[t._v("Project ID")]),s("th",[t._v("Status")]),s("th")])]),s("tbody",t._l(t.projects,function(e,r){return s("tr",{key:e.id,attrs:{project:e,projectIndex:r}},[s("td",[t._v(t._s(e.id))]),s("td",[t._v(t._s(e.attributes.enabled?"Running":"Stopped"))])])}),0)])],1)},a=[],n=s("cebc"),o=s("2f62"),i=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"drawer mb-3 clearfix"},[t.expanded?s("div",{staticClass:"drawer-contents clearfix"},[s("div",{staticClass:"left-panel"},[s("b-form",{staticClass:"filter-form"},[s("b-form-checkbox-group",{staticClass:"form-group--states",attrs:{name:"show_states",label:"","label-for":"filter-state",description:"State",switches:"",stack:""},model:{value:t.show_states,callback:function(e){t.show_states=e},expression:"show_states"}},[s("b-form-checkbox",{attrs:{value:"running",title:"Include projects that are currently RUNNING"},on:{change:function(e){return t.toggle_state("running")}}},[t._v("Running")]),s("b-form-checkbox",{attrs:{value:"stopped",title:"Include projects that are currently STOPPED"},on:{change:function(e){return t.toggle_state("stopped")}}},[t._v("Stopped")]),s("b-form-checkbox",{attrs:{value:"candidates",title:"Include projects that are yet to be imported"},on:{change:function(e){return t.toggle_state("candidates")}}},[t._v("Candidates")]),s("small",{staticClass:"form-text text-muted",attrs:{tabindex:"-1"}},[t._v("Project State")])],1),t.hasExtraBasedirs?s("b-form-group",{staticClass:"form-group--location",attrs:{label:"","label-for":"filter-location",description:"Location"}},[s("b-select",{attrs:{id:"filter-location",variant:"secondary",options:t.basedirsAsOptions},model:{value:t.show_locations,callback:function(e){t.show_locations=e},expression:"show_locations"}},[s("template",{slot:"first"},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Show projects from...")]),s("option",{attrs:{value:"all"}},[t._v("All known locations")])])],2)],1):t._e(),s("b-form-group",{staticClass:"form-group--stacks",attrs:{label:"","label-for":"filter-stacks",description:"Used Stacks"}},[s("b-select",{attrs:{id:"filter-stacks",variant:"secondary"},model:{value:t.show_stacks,callback:function(e){t.show_stacks=e},expression:"show_stacks"}},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Filter by stacks...")]),s("option",{attrs:{value:"all"}},[t._v("Any stack")]),s("optgroup",{attrs:{label:"Specific Stacks"}},t._l(t.stacksAsOptions,function(e){return s("option",{key:e.value,domProps:{value:e.value}},[t._v(t._s(e.text.toUpperCase()))])}),0),s("option",{attrs:{value:"none"}},[t._v("No stacks assigned")])])],1)],1)],1),s("div",{staticClass:"right-panel"},[s("b-form",{staticClass:"sort-form"},[s("b-form-group",{attrs:{id:"sort-by-group",label:"","label-for":"sort-by",description:"Sort by"}},[s("b-select",{attrs:{id:"sort-by",variant:"secondary"},model:{value:t.sort_by,callback:function(e){t.sort_by=e},expression:"sort_by"}},[s("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Sort by...")]),s("option",{attrs:{value:"access-date"}},[t._v("Access date")]),s("option",{attrs:{value:"creation-date"}},[t._v("Creation date")]),s("option",{attrs:{value:"project-title"}},[t._v("Project title")])])],1),s("b-form-group",{attrs:{id:"sort-order-group",label:"","label-for":"sort-order-select",description:"Order"}},[s("a",{staticClass:"view-mode view-mode--order",attrs:{target:"_blank",href:"#",title:"Sort Order"},on:{click:function(e){e.preventDefault(),t.sort_ascending=!t.sort_ascending}}},[s("font-awesome-icon",{attrs:{icon:["fa",t.sort_ascending?"sort-alpha-down":"sort-alpha-up"]}})],1)])],1),s("b-form",{staticClass:"view-form"},[s("b-form-group",{attrs:{id:"view-form",label:"","label-for":"view-select",description:"View Mode"}},[s("a",{class:{"view-mode":!0,"view-mode--cards":!0,"is-inactive":"cards"!=t.view_mode},attrs:{target:"_blank",href:"#",title:"Cards View"},on:{click:function(e){e.preventDefault(),t.view_mode="cards",t.$emit("switch-view-mode",e,"cards")}}},[s("font-awesome-icon",{attrs:{icon:["fa","columns"]}})],1),s("a",{class:{"view-mode":!0,"view-mode--table":!0,"is-inactive":"table"!=t.view_mode},attrs:{target:"_blank",href:"#",title:"Table View"},on:{click:function(e){e.preventDefault(),t.view_mode="table",t.$emit("switch-view-mode",e,"table")}}},[s("font-awesome-icon",{attrs:{icon:["fa","th-list"]}})],1)])],1)],1)]):t._e(),s("div",{staticClass:"drawer-handle",on:{click:function(e){t.expanded=!t.expanded}}},[s("div",{staticClass:"label small"},[s("span",[t._v("Viewing Options \n      "),t.expanded?s("font-awesome-icon",{attrs:{icon:["fa","chevron-up"]}}):s("font-awesome-icon",{attrs:{icon:["fa","chevron-down"]}})],1)]),s("div",{staticClass:"current-filter"},[s("b-badge",{attrs:{title:"Project State",variant:t.states_variant}},[t._v(t._s(t.states_label))]),t.hasExtraBasedirs?s("b-badge",{attrs:{title:"Project Locations",variant:"all"==t.show_locations?"secondary":"warning"}},[t._v(t._s(t.locations_label))]):t._e(),s("b-badge",{attrs:{title:"Stacks",variant:"all"==t.show_stacks?"secondary":"warning"}},[t._v(t._s(t.stacks_label))]),s("b-badge",{attrs:{title:"Sorting"}},[t._v(t._s(t.sorting_label))])],1)])])},c=[],l=(s("a481"),{name:"ProjectsDrawer",props:{},computed:Object(n["a"])({},Object(o["b"])(["basedirBy","stackBy","basedirsAsOptions","stacksAsOptions","hasExtraBasedirs"]),{states_label:function(){var t=this.show_states,e=-1!==t.indexOf("running")?"Running projects":"",s=-1!==t.indexOf("stopped")?"Stopped projects":"",r=-1!==t.indexOf("candidates"),a=e&&s?"All projects":e||s?e+s:"";return a?a+(r&&(e||s)?"":" (no candidates)"):r?"Project candidates":""},states_variant:function(){var t=this.show_states,e=-1!==t.indexOf("running")?"Running projects":"",s=-1!==t.indexOf("stopped")?"Stopped projects":"",r=-1!==t.indexOf("candidates");return e&&s&&r?"secondary":"warning"},locations_label:function(){var t="all"!==this.show_locations?this.basedirBy("id",this.show_locations):null;return"From "+(t?t.attributes.basedir:"all known locations")},stacks_label:function(){var t;if("none"===this.show_stacks)t="With no stacks assigned";else{var e="all"!==this.show_stacks?this.stackBy("id",this.show_stacks):null;t="Using "+(e?e.attributes.stackname.toUpperCase()+" stack":"any stack")}return t},sorting_label:function(){return"Sorted by "+this.sort_by.replace("-"," ")+(this.sort_ascending?"":" (reverse)")}}),data:function(){return{expanded:!1,show_states:["running","stopped","candidates"],show_locations:"all",show_stacks:"all",sort_by:"access-date",sort_ascending:!0,view_mode:"cards"}},methods:{toggle_state:function(t){var e=this.show_states,s=-1!==e.indexOf("running"),r=-1!==e.indexOf("stopped"),a=-1!==e.indexOf("candidates");"candidates"!==t||s||r||!a?("running"!==t||!s||r||a)&&("stopped"!==t||s||!r||a)||(this.show_states=["candidates"]):this.show_states=["running","stopped"]}}}),d=l,u=(s("8760"),s("7319"),s("2877")),p=Object(u["a"])(d,i,c,!1,null,"63c691c8",null),f=p.exports,h=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("b-card",{class:{"card--project":!0,"showing-details":t.showingDetails,"not-showing-details":!t.showingDetails},attrs:{to:{path:"/project/"+t.id}}},[s("b-form",{staticClass:"clearfix"},[s("b-form-group",{staticClass:"hostname-group",attrs:{id:"hostname-group-"+t.projectIndex,label:"","label-for":"hostname-input-"+t.projectIndex,description:t.showingDetails?"Hostname":""}},[s("b-form-input",{staticClass:"hostname-input",attrs:{id:"hostname-input-"+t.projectIndex,type:"text",size:"lg",required:"",placeholder:""},on:{change:t.maybeSubmit,click:function(e){t.showingDetails=!0}},model:{value:t.hostname,callback:function(e){t.hostname=e},expression:"hostname"}})],1),s("project-toolbar",{attrs:{project:t.project,projectIndex:t.projectIndex},on:{"run-stop":t.onRunStop}}),t.showingDetails?s("project-details",{attrs:{project:t.project,projectIndex:t.projectIndex},on:{"toggle-details":t.toggleDetails}}):t._e()],1),s("b-alert",{attrs:{show:t.alertShow,dismissible:t.alertDismissible,variant:t.alertVariant,fade:""},on:{dismissed:function(e){t.alertShow=!1}}},[t._v(t._s(t.alertContent))]),t.project.attributes.stack&&t.project.attributes.stack.length?s("div",{attrs:{slot:"footer"},slot:"footer"},[s("project-stack-list",{attrs:{project:t.project,projectIndex:t.projectIndex}})],1):t._e()],1)},b=[],v=(s("c5f6"),function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("ul",{staticClass:"toolbar-list"},[s("li",{staticClass:"toolbar-item"},[s("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link toolbar-link--state",attrs:{target:"_blank",href:"#",title:t.isRunning?"Stop all services":"Run all services"},on:{click:function(e){return e.preventDefault(),t.$emit("run-stop")}}},[s("font-awesome-icon",{attrs:{icon:["fa",t.isRunning?"stop":"play"]}})],1)])])}),g=[],m={name:"ProjectToolbar",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return{id:this.project.id,hostname:this.project.attributes.hostname}},computed:{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},isRunning:function(){return this.project.attributes.enabled}},methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")}}},j=m,w=(s("414b"),Object(u["a"])(j,v,g,!1,null,"262b50d4",null)),_=w.exports,k=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{attrs:{id:t.projectBase+"details",role:"tabpanel"}},[s("b-form-group",{attrs:{id:t.projectBase+"location-group","label-for":t.projectBase+"location-input",label:"",description:"Location"}},[s("b-input-group",[s("b-form-input",{staticClass:"location-input",attrs:{disabled:"",id:t.projectBase+"location-input",value:t.resolveDir(t.currentBasedir,t.path)}}),s("b-input-group-append",[s("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"folder-icon",attrs:{variant:"outline-info",title:"Open project directory",id:t.projectBase+"open-location",href:"#"}},[s("font-awesome-icon",{attrs:{icon:["fa","folder"]}})],1)],1)],1)],1),s("b-form-group",{attrs:{id:"notesGroup",label:"","label-for":"notesInput",description:"(will be visible only here)"}},[s("b-form-textarea",{attrs:{id:"textarea",placeholder:"Notes...",rows:"3","max-rows":"6"},on:{change:t.maybeSubmit},model:{value:t.notes,callback:function(e){t.notes=e},expression:"notes"}})],1),s("b-form-select",{staticClass:"add-stack",attrs:{disabled:!t.hasStacksNotInProject,required:!0},on:{change:t.addProjectStack},model:{value:t.selectedService,callback:function(e){t.selectedService=e},expression:"selectedService"}},[s("option",{attrs:{value:"",disabled:""}},[t._v(t._s(t.hasStacksNotInProject?"Add stack...":"All stacks already added"))]),t._l(t.stacksNotInProject,function(e,r){return s("option",{key:r,domProps:{value:r}},[t._v(t._s(e.attributes.stackname))])})],2),s("a",{staticClass:"hide-details",attrs:{title:"Hide project details"},on:{click:function(e){return t.$emit("toggle-details")}}},[s("font-awesome-icon",{attrs:{icon:["fa","chevron-up"]}}),s("span",[t._v("Hide")])],1)],1)},y=[],x=(s("ac6a"),s("ffc1"),s("0baa")),S=s.n(x),P={name:"ProjectDetails",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(n["a"])({id:this.project.id},this.project.attributes,{selectedService:""})},computed:Object(n["a"])({},Object(o["b"])({basedirBy:"basedirBy",serviceBy:"serviceBy",gearspecBy:"gearspecBy",allGearspecs:"gearspecs/all",allStacks:"stacks/all",hasExtraBasedirs:"hasExtraBasedirs"}),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},currentBasedir:function(){var t=this.basedirBy("id",this.basedir);return t?t.attributes.basedir:""},stacksNotInProject:function(){var t={},e=this.project.attributes.stack?this.groupProjectServicesByStack(this.project.attributes.stack):{};for(var s in this.allStacks){var r=this.allStacks[s];"undefined"===typeof e[r.id]&&(t[r.id]=r)}return t},hasStacksNotInProject:function(){return Object.entries(this.stacksNotInProject).length>0},servicesInProject:function(){for(var t={},e=0;e>this.stack.length;e++){var s=this.serviceBy("id",this.stack[e].service_id);s&&(t[this.stack[e].service_id]=s)}return t},gearsInProject:function(){for(var t={},e=0;e>this.stack.length;e++){var s=this.gearspecBy("id",this.stack[e].gearspec_id);s&&(t[this.stack[e].gearspec_id]=s)}return t}}),methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},groupProjectServicesByStack:function(t){var e=this,s={};return t.forEach(function(t,r){var a=e.gearspecBy("id",t.gearspec_id),n=e.serviceBy("id",t.service_id);a&&n&&("undefined"===typeof s[a.attributes.stack_id]&&(s[a.attributes.stack_id]={}),s[a.attributes.stack_id][a.attributes.role]=n)}),s},resolveDir:function(t,e){return t+(-1!==t.indexOf("/")?"/":"\\")+e},maybeSubmit:function(t){this.$store.dispatch("updateProject",{id:this.id,attributes:this.$data}).then(function(){})},onClosePopoverFor:function(t){this.$root.$emit("bv::hide::popover",t)},addProjectStack:function(t){this.selectedService="",this.$store.dispatch("addProjectStack",{projectId:this.id,stackId:t})},sanitizePath:function(t){var e=S()(t).trim();return e||"project"}}},C=P,O=(s("6cd5"),Object(u["a"])(C,k,y,!1,null,"d1c97864",null)),B=O.exports,I=s("2123"),A={name:"ProjectCard",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(n["a"])({id:this.project.id},this.project.attributes,{showingDetails:!1,alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"warning"})},components:{ProjectToolbar:_,ProjectDetails:B,ProjectStackList:I["default"]},computed:{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},isRunning:function(){return this.project.attributes.enabled}},methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},toggleDetails:function(){this.showingDetails=!this.showingDetails},showAlert:function(t){"string"===typeof t?this.alertContent=t:(this.alertVariant=t.variant||this.alertVariant,this.alertDismissible=t.dismissible||this.alertDismissible,this.alertContent=t.content||this.alertContent),this.alertShow=!0},maybeSubmit:function(t){this.$store.dispatch("updateProject",{id:this.id,attributes:this.$data}).then(function(){})},onRunStop:function(){this.project.attributes.stack&&this.project.attributes.stack.length>0?this.$store.dispatch("changeProjectState",{projectId:this.id,isEnabled:!this.isRunning}):this.showAlert("Please add some stacks first!")}}},$=A,D=(s("a4e2"),Object(u["a"])($,h,b,!1,null,"2ceab748",null)),E=D.exports,N={name:"ProjectList",data:function(){return{viewMode:"cards"}},components:{ProjectsDrawer:f,ProjectCard:E},computed:Object(n["a"])({},Object(o["b"])({projects:"projects/all"})),methods:{switchViewMode:function(t,e){console.log(e),this.viewMode=e}},mounted:function(){var t=this;this.$store.dispatch("basedirs/loadAll").then(function(){}),this.$store.dispatch("stacks/loadAll").then(function(){}),this.$store.dispatch("services/loadAll").then(function(){}),this.$store.dispatch("gearspecs/loadAll").then(function(){}),this.$store.dispatch("projects/loadAll").then(function(){}).then(function(){t.$store.dispatch("loadProjectDetails")})}},R=N,T=(s("d9f1"),Object(u["a"])(R,r,a,!1,null,"1fc1d8ac",null));e["default"]=T.exports},d40f:function(t,e,s){"use strict";t.exports=(()=>/[<>:"\/\\|?*\x00-\x1F]/g),t.exports.windowsNames=(()=>/^(con|prn|aux|nul|com[0-9]|lpt[0-9])$/i)},d9e3:function(t,e,s){},d9f1:function(t,e,s){"use strict";var r=s("98b6"),a=s.n(r);a.a},df7c:function(t,e,s){(function(t){function s(t,e){for(var s=0,r=t.length-1;r>=0;r--){var a=t[r];"."===a?t.splice(r,1):".."===a?(t.splice(r,1),s++):s&&(t.splice(r,1),s--)}if(e)for(;s--;s)t.unshift("..");return t}var r=/^(\/?|)([\s\S]*?)((?:\.{1,2}|[^\/]+?|)(\.[^.\/]*|))(?:[\/]*)$/,a=function(t){return r.exec(t).slice(1)};function n(t,e){if(t.filter)return t.filter(e);for(var s=[],r=0;r<t.length;r++)e(t[r],r,t)&&s.push(t[r]);return s}e.resolve=function(){for(var e="",r=!1,a=arguments.length-1;a>=-1&&!r;a--){var o=a>=0?arguments[a]:t.cwd();if("string"!==typeof o)throw new TypeError("Arguments to path.resolve must be strings");o&&(e=o+"/"+e,r="/"===o.charAt(0))}return e=s(n(e.split("/"),function(t){return!!t}),!r).join("/"),(r?"/":"")+e||"."},e.normalize=function(t){var r=e.isAbsolute(t),a="/"===o(t,-1);return t=s(n(t.split("/"),function(t){return!!t}),!r).join("/"),t||r||(t="."),t&&a&&(t+="/"),(r?"/":"")+t},e.isAbsolute=function(t){return"/"===t.charAt(0)},e.join=function(){var t=Array.prototype.slice.call(arguments,0);return e.normalize(n(t,function(t,e){if("string"!==typeof t)throw new TypeError("Arguments to path.join must be strings");return t}).join("/"))},e.relative=function(t,s){function r(t){for(var e=0;e<t.length;e++)if(""!==t[e])break;for(var s=t.length-1;s>=0;s--)if(""!==t[s])break;return e>s?[]:t.slice(e,s-e+1)}t=e.resolve(t).substr(1),s=e.resolve(s).substr(1);for(var a=r(t.split("/")),n=r(s.split("/")),o=Math.min(a.length,n.length),i=o,c=0;c<o;c++)if(a[c]!==n[c]){i=c;break}var l=[];for(c=i;c<a.length;c++)l.push("..");return l=l.concat(n.slice(i)),l.join("/")},e.sep="/",e.delimiter=":",e.dirname=function(t){var e=a(t),s=e[0],r=e[1];return s||r?(r&&(r=r.substr(0,r.length-1)),s+r):"."},e.basename=function(t,e){var s=a(t)[2];return e&&s.substr(-1*e.length)===e&&(s=s.substr(0,s.length-e.length)),s},e.extname=function(t){return a(t)[3]};var o="b"==="ab".substr(-1)?function(t,e,s){return t.substr(e,s)}:function(t,e,s){return e<0&&(e=t.length+e),t.substr(e,s)}}).call(this,s("f28c"))},ffc1:function(t,e,s){var r=s("5ca1"),a=s("504c")(!0);r(r.S,"Object",{entries:function(t){return a(t)}})}}]);
//# sourceMappingURL=projects.e1a7d903.js.map