(function(t){function e(e){for(var n,a,c=e[0],i=e[1],u=e[2],d=0,f=[];d<c.length;d++)a=c[d],o[a]&&f.push(o[a][0]),o[a]=0;for(n in i)Object.prototype.hasOwnProperty.call(i,n)&&(t[n]=i[n]);p&&p(e);while(f.length)f.shift()();return s.push.apply(s,u||[]),r()}function r(){for(var t,e=0;e<s.length;e++){for(var r=s[e],n=!0,a=1;a<r.length;a++){var c=r[a];0!==o[c]&&(n=!1)}n&&(s.splice(e--,1),t=i(i.s=r[0]))}return t}var n={},a={app:0},o={app:0},s=[];function c(t){return i.p+"js/"+({about:"about",gear:"gear",preferences:"preferences",projects:"projects",projectstack:"projectstack",stack:"stack"}[t]||t)+"."+{about:"2454b3d1",gear:"d47a5e7a",preferences:"957b4466",projects:"d6cfef6f",projectstack:"43adcce2",stack:"f22cef34"}[t]+".js"}function i(e){if(n[e])return n[e].exports;var r=n[e]={i:e,l:!1,exports:{}};return t[e].call(r.exports,r,r.exports,i),r.l=!0,r.exports}i.e=function(t){var e=[],r={projects:1,projectstack:1};a[t]?e.push(a[t]):0!==a[t]&&r[t]&&e.push(a[t]=new Promise(function(e,r){for(var n="css/"+({about:"about",gear:"gear",preferences:"preferences",projects:"projects",projectstack:"projectstack",stack:"stack"}[t]||t)+"."+{about:"31d6cfe0",gear:"31d6cfe0",preferences:"31d6cfe0",projects:"b54b7b91",projectstack:"4fe95edc",stack:"31d6cfe0"}[t]+".css",o=i.p+n,s=document.getElementsByTagName("link"),c=0;c<s.length;c++){var u=s[c],d=u.getAttribute("data-href")||u.getAttribute("href");if("stylesheet"===u.rel&&(d===n||d===o))return e()}var f=document.getElementsByTagName("style");for(c=0;c<f.length;c++){u=f[c],d=u.getAttribute("data-href");if(d===n||d===o)return e()}var p=document.createElement("link");p.rel="stylesheet",p.type="text/css",p.onload=e,p.onerror=function(e){var n=e&&e.target&&e.target.src||o,s=new Error("Loading CSS chunk "+t+" failed.\n("+n+")");s.request=n,delete a[t],p.parentNode.removeChild(p),r(s)},p.href=o;var l=document.getElementsByTagName("head")[0];l.appendChild(p)}).then(function(){a[t]=0}));var n=o[t];if(0!==n)if(n)e.push(n[2]);else{var s=new Promise(function(e,r){n=o[t]=[e,r]});e.push(n[2]=s);var u,d=document.createElement("script");d.charset="utf-8",d.timeout=120,i.nc&&d.setAttribute("nonce",i.nc),d.src=c(t),u=function(e){d.onerror=d.onload=null,clearTimeout(f);var r=o[t];if(0!==r){if(r){var n=e&&("load"===e.type?"missing":e.type),a=e&&e.target&&e.target.src,s=new Error("Loading chunk "+t+" failed.\n("+n+": "+a+")");s.type=n,s.request=a,r[1](s)}o[t]=void 0}};var f=setTimeout(function(){u({type:"timeout",target:d})},12e4);d.onerror=d.onload=u,document.head.appendChild(d)}return Promise.all(e)},i.m=t,i.c=n,i.d=function(t,e,r){i.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:r})},i.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},i.t=function(t,e){if(1&e&&(t=i(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var r=Object.create(null);if(i.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var n in t)i.d(r,n,function(e){return t[e]}.bind(null,n));return r},i.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return i.d(e,"a",e),e},i.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},i.p="",i.oe=function(t){throw console.error(t),t};var u=window["webpackJsonp"]=window["webpackJsonp"]||[],d=u.push.bind(u);u.push=e,u=u.slice();for(var f=0;f<u.length;f++)e(u[f]);var p=d;s.push([0,"chunk-vendors"]),r()})({0:function(t,e,r){t.exports=r("56d7")},"034f":function(t,e,r){"use strict";var n=r("64a9"),a=r.n(n);a.a},"56d7":function(t,e,r){"use strict";r.r(e);r("cadf"),r("551c"),r("f751"),r("097d");var n=r("2b0e"),a=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{attrs:{id:"app"}},[r("b-alert",{attrs:{show:t.isConnectionProblem,variant:"warning"}},[r("h4",[t._v("Connection Problem")]),r("p",[t._v("It seems that Gearbox Server is not running. Remaining connection attempts: "+t._s(t.remainingRetries))])]),t.isUnrecoverableConnectionProblem?r("b-alert",{attrs:{show:"",variant:"danger"}},[r("h4",[t._v("Connection Problem")]),r("p",[t._v("Failed to connect to Gearbox Server.")])]):t._e(),r("top-bar"),r("br"),r("router-view")],1)},o=[],s=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",[r("b-navbar",{attrs:{toggleable:"lg",type:"dark",variant:"info"}},[r("b-navbar-brand",{attrs:{to:"/"}},[t._v("Gearbox")]),r("b-navbar-toggle",{attrs:{target:"nav_collapse"}}),r("b-collapse",{attrs:{"is-nav":"",id:"nav_collapse"}},[r("b-navbar-nav",[r("b-nav-item",{attrs:{to:"/projects"}},[t._v("Projects")]),r("b-nav-item",{attrs:{to:"/preferences"}},[t._v("Preferences")])],1),r("b-navbar-nav",{staticClass:"ml-auto"},[r("b-nav-form",[r("b-form-input",{staticClass:"mr-sm-2",attrs:{size:"sm",type:"text",placeholder:"Search"}}),r("b-button",{staticClass:"my-2 my-sm-0",attrs:{size:"sm",type:"submit"}},[t._v("Search")])],1),r("b-nav-item-dropdown",{attrs:{text:"Lang",right:""}},[r("b-dropdown-item",{attrs:{href:"#"}},[t._v("EN")]),r("b-dropdown-item",{attrs:{href:"#"}},[t._v("ES")]),r("b-dropdown-item",{attrs:{href:"#"}},[t._v("RU")]),r("b-dropdown-item",{attrs:{href:"#"}},[t._v("FR")])],1),r("b-nav-item-dropdown",{attrs:{right:""}},[r("template",{slot:"button-content"},[r("em",[t._v("User")])]),r("b-dropdown-item",{attrs:{href:"#"}},[t._v("Profile")]),r("b-dropdown-item",{attrs:{href:"#"}},[t._v("Signout")])],2)],1)],1)],1)],1)},c=[],i={name:"TopBar.vue"},u=i,d=r("2877"),f=Object(d["a"])(u,s,c,!1,null,"fb66cdee",null),p=f.exports,l={name:"App",components:{TopBar:p},computed:{isConnectionProblem:function(){return this.$store.state.connectionStatus.networkError&&this.$store.state.connectionStatus.remainingRetries>0},remainingRetries:function(){return this.$store.state.connectionStatus.remainingRetries},isUnrecoverableConnectionProblem:function(){return this.$store.state.connectionStatus.networkError?0===this.$store.state.connectionStatus.remainingRetries:""}}},b=l,v=(r("034f"),Object(d["a"])(b,a,o,!1,null,null,null)),m=v.exports,h=r("8c4f"),g=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("h1",[t._v(t._s(t.title))])},E=[],_={props:{title:{type:String,default:"Welcome!"}}},y=_,j=Object(d["a"])(y,g,E,!1,null,null,null),R=j.exports;n["a"].use(h["a"]);var k=new h["a"]({routes:[{path:"/",name:"welcome",component:R},{path:"/about",name:"about",component:function(){return r.e("about").then(r.bind(null,"f820"))}},{path:"/preferences",name:"preferences",component:function(){return r.e("preferences").then(r.bind(null,"a55d"))}},{path:"/projects",name:"projects",component:function(){return r.e("projects").then(r.bind(null,"acca"))},children:[{path:":hostname/stack",component:function(){return r.e("projectstack").then(r.bind(null,"3983"))}}]},{path:"/stack/:stackName",name:"stack",component:function(){return r.e("stack").then(r.bind(null,"9e3c"))}},{path:"/gear/:gearName",name:"gear",component:function(){return r.e("gear").then(r.bind(null,"babd"))}}]}),S=(r("28a5"),r("7f7f"),r("ac6a"),r("7514"),r("cebc")),T=r("2f62"),C=r("bc3a"),O=r.n(C),w=r("a7fe"),P=r.n(w),A=r("6916"),x=r("768b"),B=r("5d73"),I=r.n(B),D=r("db0c"),N=r.n(D),J=O.a.create({baseURL:"http://127.0.0.1:9999/",headers:{"Content-Type":"application/vnd.api+json"}});J.defaults.raxConfig={instance:J,retry:5,noResponseRetries:5,retryDelay:3e3,httpMethodsToRetry:["GET","HEAD","OPTIONS","DELETE","PUT"],statusCodesToRetry:[[100,199],[429,429],[500,599]],shouldRetry:function(t){var e=Object(A["getConfig"])(t);if(!e||0===e.retry)return!1;if(!t.response&&(e.currentRetryAttempt||0)>=e.noResponseRetries)return!1;if(!t.config.method||N()(e.httpMethodsToRetry).indexOf(t.config.method.toUpperCase())<0)return!1;if(t.response&&t.response.status){var r=!1,n=!0,a=!1,o=void 0;try{for(var s,c=I()(e.statusCodesToRetry);!(n=(s=c.next()).done);n=!0){var i=Object(x["a"])(s.value,2),u=i[0],d=i[1],f=t.response.status;if(f>=u&&f<=d){r=!0;break}}}catch(t){a=!0,o=t}finally{try{n||null==c.return||c.return()}finally{if(a)throw o}}if(!r)return!1}return e.currentRetryAttempt=e.currentRetryAttempt||0,!(e.currentRetryAttempt>=e.retry)}},Object(A["attach"])(J);var M=J,G=r("a592");n["a"].use(T["a"]),n["a"].use(P.a,O.a);var K=new T["a"].Store({strict:!0,modules:Object(S["a"])({},Object(G["a"])({names:["stacks","services","gearspecs","projects","basedirs"],httpClient:M})),state:{stacks:[],services:[],gearspecs:[],projects:[],basedirs:[],connectionStatus:{networkError:null,remainingRetries:5}},getters:{basedirBy:function(t){return function(e,r){return"id"===e?t.basedirs.records.find(function(t){return t.id===r}):t.basedirs.records.find(function(t){return t.attributes[e]===r})}},stackBy:function(t){return function(e,r){return"id"===e?t.stacks.records.find(function(t){return t.id===r}):t.stacks.records.find(function(t){return t.attributes[e]===r})}},serviceBy:function(t){return function(e,r){return"id"===e?t.services.records.find(function(t){return t.id===r}):t.services.records.find(function(t){return t.attributes[e]===r})}},gearspecBy:function(t){return function(e,r){return"id"===e?t.gearspecs.records.find(function(t){return t.id===r}):t.gearspecs.records.find(function(t){return t.attributes[e]===r})}},projectBy:function(t){return function(e,r){return"id"===e?t.projects.records.find(function(t){return t.id===r}):t.projects.records.find(function(t){return t.attributes[e]===r})}},projectStackMemberIndexBy:function(t){return function(t,e,r){var n=-1;return t.attributes.stack.find(function(t,a){return t[e]===r&&(n=a,!0)}),n}},baseDirsAsOptions:function(t){var e=[];return t.basedirs.records.forEach(function(t,r){e.push({value:t.id,text:t.host_dir})}),e},preselectService:function(t){return function(t,e){var r=-1,n=-1;if(e)for(var a=t.length;a--;)if(-1!==t[a].indexOf(e)&&(-1===r&&(r=a),t[a]===e)){n=a;break}var o=-1!==r?t[-1!==n?n:r]:"";return o}}},actions:{loadProjectDetails:function(t){var e=t.commit;for(var r in this.state.projects.records){var n=this.state.projects.records[r];try{M.get("projects/"+n.id,{crossDomain:!0,raxConfig:{onRetryAttempt:function(t){var r=Object(A["getConfig"])(t);e("SET_NETWORK_ERROR",t.message),e("SET_REMAINING_RETRIES",r.retry-r.currentRetryAttempt)}}}).catch(function(t,e){console.log("rejected",t)}).then(function(t){return t?t.data:null}).then(function(t){var r=t.data;if(e("SET_PROJECT",r),t.included.length)for(var n in t.included){var a=t.included[n];"service"===a.type&&e("SET_SERVICE",a),"stack"===a.type&&e("SET_STACK",a)}})}catch(a){console.log(a)}}},updateProject:function(t,e){var r=t.commit,n=e.hostname,a=e.project;r("UPDATE_PROJECT",{hostname:n,project:a}),M({method:"post",url:"project/"+n,data:a}).then(function(t){return t.data}).then(function(t){}).catch(function(t){console.log("rejected",t)})},addBaseDir:function(t,e){var r=t.commit,n=e.name,a=e.path;r("ADD_BASEDIR",{value:n,text:a}),M({method:"post",url:"basedirs/new",data:e}).then(function(t){return t.data}).then(function(t){}).catch(function(t){console.log("rejected",t)})},addProjectStack:function(t,e){var r=t.commit;r("ADD_PROJECT_STACK",e)},removeProjectStack:function(t,e){var r=t.commit;r("REMOVE_PROJECT_STACK",e)},changeProjectService:function(t,e){var r=t.commit;r("CHANGE_PROJECT_SERVICE",e)},changeProjectState:function(t,e){var r=t.commit;r("CHANGE_PROJECT_STATE",e)}},mutations:{SET_PROJECT:function(t,e){var r=this.getters.projectBy("id",e.id);r?n["a"].set(r.attributes,"stack",e.attributes.stack):t.projects.records.push(e)},SET_STACK:function(t,e){var r=this.getters.stackBy("id",e.id);r?r.attributes=e.attributes:t.stacks.records.push(e)},SET_SERVICE:function(t,e){var r=this.getters.serviceBy("id",e.id);r?r.attributes=e.attributes:t.services.records.push(e)},SET_GEARSTACK:function(t,e){var r=this.getters.gearstackBy("id",e.id);r?r.attributes=e.attributes:t.gearstacks.records.push(e)},SET_NETWORK_ERROR:function(t,e){t.connectionStatus.networkError=e},CLEAR_NETWORK_ERROR:function(t){t.connectionStatus.networkError=""},SET_REMAINING_RETRIES:function(t,e){t.connectionStatus.remainingRetries=e},ADD_BASEDIR:function(t,e){t.baseDirs[e.value]=e},ADD_PROJECT_STACK:function(t,e){var r=this,a=e.projectId,o=e.stackId,s=this.getters.projectBy("id",a),c=this.getters.stackBy("id",o);s&&c&&c.attributes.members.length&&("undefined"===typeof s.attributes.stack&&n["a"].set(s.attributes,"stack",[]),c.attributes.members.forEach(function(t,e){var n=r.getters.preselectService(t.services,t.default_service);n&&t.gearspec_id&&s.attributes.stack.push({service_id:n,gearspec_id:t.gearspec_id})}))},REMOVE_PROJECT_STACK:function(t,e){var r=e.projectId,a=e.stackId,o=this.getters.projectBy("id",r);if(o)for(var s=a.split("/")[1],c=o.attributes.stack.length-1;c>=0;c--)o.attributes.stack[c].gearspec_id.split("/")[1]===s&&n["a"].delete(o.attributes.stack,c)},CHANGE_PROJECT_SERVICE:function(t,e){var r=e.projectId,a=e.serviceId,o=e.gearspecId,s=this.getters.projectBy("id",r);if(s){var c=this.getters.projectStackMemberIndexBy(s,"gearspec_id",o);a?n["a"].set(s.attributes.stack[c],"service_id",a):n["a"].delete(s.attributes.stack,c)}},CHANGE_PROJECT_STATE:function(t,e){var r=e.projectId,n=e.isEnabled,a=this.getters.projectBy("id",r);a&&(a.attributes.enabled=!!n)}}}),$=r("9f7b"),U=r.n($);r("f9e3"),r("2dd8");n["a"].use(U.a);var L=r("ecee"),V=r("c074"),H=r("ad3d");L["c"].add(V["d"]),L["c"].add(V["e"]),L["c"].add(V["f"]),L["c"].add(V["c"]),L["c"].add(V["a"]),L["c"].add(V["b"]),n["a"].component("font-awesome-icon",H["a"]),n["a"].config.productionTip=!1,new n["a"]({router:k,store:K,render:function(t){return t(m)}}).$mount("#app")},"64a9":function(t,e,r){}});
//# sourceMappingURL=app.bd863796.js.map