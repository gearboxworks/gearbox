(function(e){function t(t){for(var r,a,s=t[0],i=t[1],u=t[2],f=0,l=[];f<s.length;f++)a=s[f],o[a]&&l.push(o[a][0]),o[a]=0;for(r in i)Object.prototype.hasOwnProperty.call(i,r)&&(e[r]=i[r]);p&&p(t);while(l.length)l.shift()();return c.push.apply(c,u||[]),n()}function n(){for(var e,t=0;t<c.length;t++){for(var n=c[t],r=!0,a=1;a<n.length;a++){var s=n[a];0!==o[s]&&(r=!1)}r&&(c.splice(t--,1),e=i(i.s=n[0]))}return e}var r={},a={app:0},o={app:0},c=[];function s(e){return i.p+"js/"+({about:"about",gear:"gear",preferences:"preferences",projects:"projects",projectstack:"projectstack",stack:"stack"}[e]||e)+"."+{about:"10d8c6ca",gear:"7f2aab7b",preferences:"d999ac7d",projects:"64eea0a1",projectstack:"9839b10f",stack:"5afb343d"}[e]+".js"}function i(t){if(r[t])return r[t].exports;var n=r[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,i),n.l=!0,n.exports}i.e=function(e){var t=[],n={projects:1,projectstack:1};a[e]?t.push(a[e]):0!==a[e]&&n[e]&&t.push(a[e]=new Promise(function(t,n){for(var r="css/"+({about:"about",gear:"gear",preferences:"preferences",projects:"projects",projectstack:"projectstack",stack:"stack"}[e]||e)+"."+{about:"31d6cfe0",gear:"31d6cfe0",preferences:"31d6cfe0",projects:"eb4d7850",projectstack:"0a353f03",stack:"31d6cfe0"}[e]+".css",o=i.p+r,c=document.getElementsByTagName("link"),s=0;s<c.length;s++){var u=c[s],f=u.getAttribute("data-href")||u.getAttribute("href");if("stylesheet"===u.rel&&(f===r||f===o))return t()}var l=document.getElementsByTagName("style");for(s=0;s<l.length;s++){u=l[s],f=u.getAttribute("data-href");if(f===r||f===o)return t()}var p=document.createElement("link");p.rel="stylesheet",p.type="text/css",p.onload=t,p.onerror=function(t){var r=t&&t.target&&t.target.src||o,c=new Error("Loading CSS chunk "+e+" failed.\n("+r+")");c.request=r,delete a[e],p.parentNode.removeChild(p),n(c)},p.href=o;var d=document.getElementsByTagName("head")[0];d.appendChild(p)}).then(function(){a[e]=0}));var r=o[e];if(0!==r)if(r)t.push(r[2]);else{var c=new Promise(function(t,n){r=o[e]=[t,n]});t.push(r[2]=c);var u,f=document.createElement("script");f.charset="utf-8",f.timeout=120,i.nc&&f.setAttribute("nonce",i.nc),f.src=s(e),u=function(t){f.onerror=f.onload=null,clearTimeout(l);var n=o[e];if(0!==n){if(n){var r=t&&("load"===t.type?"missing":t.type),a=t&&t.target&&t.target.src,c=new Error("Loading chunk "+e+" failed.\n("+r+": "+a+")");c.type=r,c.request=a,n[1](c)}o[e]=void 0}};var l=setTimeout(function(){u({type:"timeout",target:f})},12e4);f.onerror=f.onload=u,document.head.appendChild(f)}return Promise.all(t)},i.m=e,i.c=r,i.d=function(e,t,n){i.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},i.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},i.t=function(e,t){if(1&t&&(e=i(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(i.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var r in e)i.d(n,r,function(t){return e[t]}.bind(null,r));return n},i.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return i.d(t,"a",t),t},i.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},i.p="",i.oe=function(e){throw console.error(e),e};var u=window["webpackJsonp"]=window["webpackJsonp"]||[],f=u.push.bind(u);u.push=t,u=u.slice();for(var l=0;l<u.length;l++)t(u[l]);var p=f;c.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},"034f":function(e,t,n){"use strict";var r=n("64a9"),a=n.n(r);a.a},"56d7":function(e,t,n){"use strict";n.r(t);n("cadf"),n("551c"),n("f751"),n("097d");var r=n("2b0e"),a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app"}},[n("b-alert",{attrs:{show:e.isConnectionProblem,variant:"warning"}},[n("h4",[e._v("Connection Problem")]),n("p",[e._v("It seems that Gearbox Server is not running. Remaining connection attempts: "+e._s(e.remainingRetries))])]),e.isUnrecoverableConnectionProblem?n("b-alert",{attrs:{show:"",variant:"danger"}},[n("h4",[e._v("Connection Problem")]),n("p",[e._v("Failed to connect to Gearbox Server.")])]):e._e(),n("top-bar"),n("br"),n("router-view")],1)},o=[],c=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("b-navbar",{attrs:{toggleable:"lg",type:"dark",variant:"info"}},[n("b-navbar-brand",{attrs:{to:"/"}},[e._v("Gearbox")]),n("b-navbar-toggle",{attrs:{target:"nav_collapse"}}),n("b-collapse",{attrs:{"is-nav":"",id:"nav_collapse"}},[n("b-navbar-nav",[n("b-nav-item",{attrs:{to:"/projects"}},[e._v("Projects")]),n("b-nav-item",{attrs:{to:"/preferences"}},[e._v("Preferences")])],1),n("b-navbar-nav",{staticClass:"ml-auto"},[n("b-nav-form",[n("b-form-input",{staticClass:"mr-sm-2",attrs:{size:"sm",type:"text",placeholder:"Search"}}),n("b-button",{staticClass:"my-2 my-sm-0",attrs:{size:"sm",type:"submit"}},[e._v("Search")])],1),n("b-nav-item-dropdown",{attrs:{text:"Lang",right:""}},[n("b-dropdown-item",{attrs:{href:"#"}},[e._v("EN")]),n("b-dropdown-item",{attrs:{href:"#"}},[e._v("ES")]),n("b-dropdown-item",{attrs:{href:"#"}},[e._v("RU")]),n("b-dropdown-item",{attrs:{href:"#"}},[e._v("FR")])],1),n("b-nav-item-dropdown",{attrs:{right:""}},[n("template",{slot:"button-content"},[n("em",[e._v("User")])]),n("b-dropdown-item",{attrs:{href:"#"}},[e._v("Profile")]),n("b-dropdown-item",{attrs:{href:"#"}},[e._v("Signout")])],2)],1)],1)],1)],1)},s=[],i={name:"TopBar.vue"},u=i,f=n("2877"),l=Object(f["a"])(u,c,s,!1,null,"fb66cdee",null),p=l.exports,d={name:"App",components:{TopBar:p},computed:{isConnectionProblem:function(){return this.$store.state.connectionStatus.networkError&&this.$store.state.connectionStatus.remainingRetries>0},remainingRetries:function(){return this.$store.state.connectionStatus.remainingRetries},isUnrecoverableConnectionProblem:function(){return this.$store.state.connectionStatus.networkError?0===this.$store.state.connectionStatus.remainingRetries:""}}},m=d,h=(n("034f"),Object(f["a"])(m,a,o,!1,null,null,null)),v=h.exports,b=n("8c4f"),g=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("h1",[e._v(e._s(e.title))])},E=[],S={props:{title:{type:String,default:"Welcome!"}}},R=S,_=Object(f["a"])(R,g,E,!1,null,null,null),y=_.exports;r["a"].use(b["a"]);var j=new b["a"]({routes:[{path:"/",name:"welcome",component:y},{path:"/about",name:"about",component:function(){return n.e("about").then(n.bind(null,"f820"))}},{path:"/preferences",name:"preferences",component:function(){return n.e("preferences").then(n.bind(null,"a55d"))}},{path:"/projects",name:"projects",component:function(){return n.e("projects").then(n.bind(null,"acca"))},children:[{path:":hostname/stack",component:function(){return n.e("projectstack").then(n.bind(null,"3983"))}}]},{path:"/stack/:stackName",name:"stack",component:function(){return n.e("stack").then(n.bind(null,"9e3c"))}},{path:"/gear/:gearName",name:"gear",component:function(){return n.e("gear").then(n.bind(null,"babd"))}}]}),T=(n("7f7f"),n("ac6a"),n("a481"),n("28a5"),n("7514"),n("cebc")),O=n("2f62"),k=n("bc3a"),C=n.n(k),P=n("a7fe"),A=n.n(P),w=n("6916"),D=n("768b"),x=n("5d73"),N=n.n(x),I=n("db0c"),G=n.n(I),B=C.a.create({baseURL:"http://127.0.0.1:9999/",headers:{"Content-Type":"application/vnd.api+json"}});B.defaults.raxConfig={instance:B,retry:5,noResponseRetries:5,retryDelay:3e3,httpMethodsToRetry:["GET","HEAD","OPTIONS","DELETE","PUT"],statusCodesToRetry:[[100,199],[429,429],[500,599]],shouldRetry:function(e){var t=Object(w["getConfig"])(e);if(!t||0===t.retry)return!1;if(!e.response&&(t.currentRetryAttempt||0)>=t.noResponseRetries)return!1;if(!e.config.method||G()(t.httpMethodsToRetry).indexOf(e.config.method.toUpperCase())<0)return!1;if(e.response&&e.response.status){var n=!1,r=!0,a=!1,o=void 0;try{for(var c,s=N()(t.statusCodesToRetry);!(r=(c=s.next()).done);r=!0){var i=Object(D["a"])(c.value,2),u=i[0],f=i[1],l=e.response.status;if(l>=u&&l<=f){n=!0;break}}}catch(e){a=!0,o=e}finally{try{r||null==s.return||s.return()}finally{if(a)throw o}}if(!n)return!1}return t.currentRetryAttempt=t.currentRetryAttempt||0,!(t.currentRetryAttempt>=t.retry)}},Object(w["attach"])(B);var J=B,K=n("a592");r["a"].use(O["a"]),r["a"].use(A.a,C.a);var M=new O["a"].Store({strict:!0,modules:Object(T["a"])({},Object(K["a"])({names:["baseDirs","stacks","projects"],httpClient:J})),state:{baseDirs:{},projects:[],stacks:[],stack_members:[],gearStacks:{},gearRoles:{},gearServices:{},connectionStatus:{networkError:null,remainingRetries:5}},getters:{projectBy:function(e){return function(t,n){return e.projects.find(function(e){return e[t]===n})}},groupProjectStacks:function(e){return function(e){var t={};for(var n in e)if(e.hasOwnProperty(n)){var r=e[n].authority+"/"+e[n].stack,a=e[n].authority+"/"+n;"undefined"===typeof t[r]&&(t[r]={}),t[r][a]=e[n]}return t}},projectStackServices:function(e){return function(t){var n,r={};for(n in e.gearServices)e.gearServices.hasOwnProperty(n)&&-1!==n.indexOf(t)&&(r[n]=e.gearServices[n]);return r}},stackRoles:function(e){return function(t){var n,r={};for(n in e.gearRoles)e.gearRoles.hasOwnProperty(n)&&-1!==n.indexOf(t)&&(r[n]=e.gearRoles[n]);return r}},stackServices:function(e){return function(t){var n,r={};for(n in e.gearServices)e.gearServices.hasOwnProperty(n)&&-1!==n.indexOf(t)&&(r[n]=e.gearServices[n]);return r}},baseDirsAsOptions:function(e){var t=[];return t},projectServiceDefaults:function(e){return function(t,n){var r=t.substring(0,t.indexOf("/")),a=e.gearServices[t],o=-1,c=-1;if(a.default)for(var s=a.options.length;s--;)if(-1!==a.options[s].indexOf(a.default)&&(-1===o&&(o=s),a.options[s]===a.default)){c=s;break}var i=-1!==o?a.options[-1!==c?c:o]:"",u=i?i.split(":")[1].split("."):"",f={authority:r,org:r.replace(".",""),stack:n.substring(n.indexOf("/")+1),service_id:i?a.org+"/"+i:"",program:i?i.substring(0,i.indexOf(":")):"",version:{}};return u.length>0&&(f.version.major=u[0]),u.length>1&&(f.version.minor=u[1]),u.length>2&&(f.version.patch=u[2]),f}}},actions:{loadBaseDirs:function(e){var t=e.commit;try{J.get("basedirs",{crossDomain:!0,raxConfig:{onRetryAttempt:function(e){var n=Object(w["getConfig"])(e);t("SET_NETWORK_ERROR",e.message),t("SET_REMAINING_RETRIES",n.retry-n.currentRetryAttempt)}}}).catch(function(e,t){console.log("rejected",e)}).then(function(e){return e?e.data.data:null}).then(function(e){if(e){var n={};for(var r in e)e.hasOwnProperty(r)&&(n[r]={value:r,text:e[r]});t("SET_BASEDIRS",n)}})}catch(n){console.log(n)}},loadProjectHeaders:function(e){var t=e.commit;try{J.get("projects/with-details",{crossDomain:!0,raxConfig:{onRetryAttempt:function(e){var n=Object(w["getConfig"])(e);t("SET_NETWORK_ERROR",e.message),t("SET_REMAINING_RETRIES",n.retry-n.currentRetryAttempt)}}}).catch(function(e,t){console.log("rejected",e)}).then(function(e){return e?e.data.data:null}).then(function(e){if(e){var n=[];for(var r in e)if(e.hasOwnProperty(r)){var a=e[r],o=a.data;n.push({baseDir:o.basedir,path:o.path,hostname:o.hostname,fullPath:o.project_dir,enabled:o.enabled,notes:o.notes,aliases:o.aliases,stack:o.stack})}t("SET_PROJECTS",n)}})}catch(n){console.log(n)}},loadProjectDetails:function(e){var t=e.commit;this.state.projects.forEach(function(e,n){try{J.get("projects/"+e.hostname,{crossDomain:!0,raxConfig:{onRetryAttempt:function(e){var n=Object(w["getConfig"])(e);t("SET_NETWORK_ERROR",e.message),t("SET_REMAINING_RETRIES",n.retry-n.currentRetryAttempt)}}}).catch(function(e,t){console.log("rejected",e)}).then(function(e){return e?e.data.data:null}).then(function(e){var n={path:e.hostname,enabled:e.enabled};t("SET_PROJECT_DETAILS",n)})}catch(r){console.log(r)}})},loadStacks:function(e){var t=e.commit;J.get("stacks",{crossDomain:!0}).catch(function(e){console.log("rejected",e)}).then(function(e){return e?e.data.data:null}).then(function(e){if(e){var n=[];for(var r in e)if(e.hasOwnProperty(r)){var a=e[r];n.push({code:r,name:a.name,label:a.label,examples:a.examples,stack:a.stack,optional:a.optional,memberType:a.member_type})}t("SET_STACKS",e)}})},loadGears:function(e){var t=e.commit;J.get("stacks",{crossDomain:!0}).catch(function(e){}).then(function(e){return e.data}).then(function(e){t("SET_GEAR_STACKS",e.stacks),t("SET_GEAR_ROLES",e.roles),t("SET_GEAR_SERVICES",e.services)})},updateProject:function(e,t){var n=e.commit,r=t.hostname,a=t.project;n("UPDATE_PROJECT",{hostname:r,project:a}),J({method:"post",url:"project/"+r,data:a}).then(function(e){return e.data}).then(function(e){}).catch(function(e){console.log("rejected",e)})},addBaseDir:function(e,t){var n=e.commit,r=t.name,a=t.path;n("ADD_BASEDIR",{value:r,text:a}),J({method:"post",url:"basedirs/new",data:t}).then(function(e){return e.data}).then(function(e){}).catch(function(e){console.log("rejected",e)})},addProjectStack:function(e,t){var n=e.commit;n("ADD_PROJECT_STACK",t)},removeProjectStack:function(e,t){var n=e.commit;n("REMOVE_PROJECT_STACK",t)},changeProjectService:function(e,t){var n=e.commit;n("CHANGE_PROJECT_SERVICE",t)},changeProjectState:function(e,t){var n=e.commit;n("CHANGE_PROJECT_STATE",t)}},mutations:{SET_PROJECTS:function(e,t){e.projects=t},SET_PROJECT_DETAILS:function(e,t){var n=this.getters.projectBy("path",t.path);n.enabled=t.enabled},SET_STACKS:function(e,t){e.stacks=t},SET_GEAR_STACKS:function(e,t){e.gearStacks=t},SET_GEAR_ROLES:function(e,t){e.gearRoles=t},SET_GEAR_SERVICES:function(e,t){e.gearServices=t},UPDATE_PROJECT:function(e,t){var n=t.hostname,r=t.project;console.log(t);var a=this.getters.projectBy("hostname",n);a.hostname=r.hostname,a.notes=r.notes,a.baseDir=r.baseDir,a.path=r.path,a.fullPath=r.fullPath,a.enabled=r.enabled},SET_NETWORK_ERROR:function(e,t){e.connectionStatus.networkError=t},CLEAR_NETWORK_ERROR:function(e){e.connectionStatus.networkError=""},SET_REMAINING_RETRIES:function(e,t){e.connectionStatus.remainingRetries=t},SET_BASEDIRS:function(e,t){e.baseDirs=t},ADD_BASEDIR:function(e,t){e.baseDirs[t.value]=t},ADD_PROJECT_STACK:function(e,t){var n=t.projectHostname,a=t.stackName,o=this.getters.projectBy("hostname",n),c=serviceName.substring(serviceName.indexOf("/")+1);if(o)for(var s in this.getters.stackServices(a))r["a"].set(o.stack,c,this.projectServiceDefaults(a,s))},REMOVE_PROJECT_STACK:function(e,t){var n=t.projectHostname,a=t.stackName,o=this.getters.projectBy("hostname",n);if(o){var c=a.split("/")[1];for(var s in o.stack)o.stack.hasOwnProperty(s)&&-1!==s.indexOf(c)&&r["a"].delete(o.stack,s)}},CHANGE_PROJECT_SERVICE:function(e,t){var n=t.projectHostname,a=t.serviceName,o=t.serviceId,c=this.getters.projectBy("hostname",n),s=e.gearServices[a];if(c&&s){var i=a.substring(a.indexOf("/")+1);if(o)o.split("/")[1],o.split(":")[1].split(".");else r["a"].delete(c.stack,i)}},CHANGE_PROJECT_STATE:function(e,t){var n=t.projectHostname,a=t.isEnabled,o=this.getters.projectBy("hostname",n);o&&r["a"].set(o,"enabled",!!a)}}}),H=n("9f7b"),L=n.n(H);n("f9e3"),n("2dd8");r["a"].use(L.a);var U=n("ecee"),$=n("c074"),V=n("ad3d");U["c"].add($["d"]),U["c"].add($["e"]),U["c"].add($["f"]),U["c"].add($["c"]),U["c"].add($["a"]),U["c"].add($["b"]),r["a"].component("font-awesome-icon",V["a"]),r["a"].config.productionTip=!1,new r["a"]({router:j,store:M,render:function(e){return e(v)}}).$mount("#app")},"64a9":function(e,t,n){}});
//# sourceMappingURL=app.17c7b375.js.map