(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["project"],{"014b":function(t,e,n){"use strict";var r=n("e53d"),a=n("07e3"),o=n("8e60"),i=n("63b6"),c=n("9138"),s=n("ebfd").KEY,u=n("294c"),f=n("dbdb"),l=n("45f2"),b=n("62a0"),p=n("5168"),d=n("ccb9"),h=n("6718"),m=n("47ee"),v=n("9003"),y=n("e4ae"),j=n("f772"),g=n("36c3"),O=n("1bc3"),w=n("aebd"),P=n("a159"),_=n("0395"),x=n("bf0b"),S=n("d9f6"),k=n("c3a1"),E=x.f,D=S.f,I=_.f,N=r.Symbol,F=r.JSON,$=F&&F.stringify,G="prototype",B=p("_hidden"),J=p("toPrimitive"),C={}.propertyIsEnumerable,q=f("symbol-registry"),T=f("symbols"),K=f("op-symbols"),W=Object[G],H="function"==typeof N,M=r.QObject,Y=!M||!M[G]||!M[G].findChild,z=o&&u(function(){return 7!=P(D({},"a",{get:function(){return D(this,"a",{value:7}).a}})).a})?function(t,e,n){var r=E(W,e);r&&delete W[e],D(t,e,n),r&&t!==W&&D(W,e,r)}:D,A=function(t){var e=T[t]=P(N[G]);return e._k=t,e},L=H&&"symbol"==typeof N.iterator?function(t){return"symbol"==typeof t}:function(t){return t instanceof N},Q=function(t,e,n){return t===W&&Q(K,e,n),y(t),e=O(e,!0),y(n),a(T,e)?(n.enumerable?(a(t,B)&&t[B][e]&&(t[B][e]=!1),n=P(n,{enumerable:w(0,!1)})):(a(t,B)||D(t,B,w(1,{})),t[B][e]=!0),z(t,e,n)):D(t,e,n)},R=function(t,e){y(t);var n,r=m(e=g(e)),a=0,o=r.length;while(o>a)Q(t,n=r[a++],e[n]);return t},U=function(t,e){return void 0===e?P(t):R(P(t),e)},V=function(t){var e=C.call(this,t=O(t,!0));return!(this===W&&a(T,t)&&!a(K,t))&&(!(e||!a(this,t)||!a(T,t)||a(this,B)&&this[B][t])||e)},X=function(t,e){if(t=g(t),e=O(e,!0),t!==W||!a(T,e)||a(K,e)){var n=E(t,e);return!n||!a(T,e)||a(t,B)&&t[B][e]||(n.enumerable=!0),n}},Z=function(t){var e,n=I(g(t)),r=[],o=0;while(n.length>o)a(T,e=n[o++])||e==B||e==s||r.push(e);return r},tt=function(t){var e,n=t===W,r=I(n?K:g(t)),o=[],i=0;while(r.length>i)!a(T,e=r[i++])||n&&!a(W,e)||o.push(T[e]);return o};H||(N=function(){if(this instanceof N)throw TypeError("Symbol is not a constructor!");var t=b(arguments.length>0?arguments[0]:void 0),e=function(n){this===W&&e.call(K,n),a(this,B)&&a(this[B],t)&&(this[B][t]=!1),z(this,t,w(1,n))};return o&&Y&&z(W,t,{configurable:!0,set:e}),A(t)},c(N[G],"toString",function(){return this._k}),x.f=X,S.f=Q,n("6abf").f=_.f=Z,n("355d").f=V,n("9aa9").f=tt,o&&!n("b8e3")&&c(W,"propertyIsEnumerable",V,!0),d.f=function(t){return A(p(t))}),i(i.G+i.W+i.F*!H,{Symbol:N});for(var et="hasInstance,isConcatSpreadable,iterator,match,replace,search,species,split,toPrimitive,toStringTag,unscopables".split(","),nt=0;et.length>nt;)p(et[nt++]);for(var rt=k(p.store),at=0;rt.length>at;)h(rt[at++]);i(i.S+i.F*!H,"Symbol",{for:function(t){return a(q,t+="")?q[t]:q[t]=N(t)},keyFor:function(t){if(!L(t))throw TypeError(t+" is not a symbol!");for(var e in q)if(q[e]===t)return e},useSetter:function(){Y=!0},useSimple:function(){Y=!1}}),i(i.S+i.F*!H,"Object",{create:U,defineProperty:Q,defineProperties:R,getOwnPropertyDescriptor:X,getOwnPropertyNames:Z,getOwnPropertySymbols:tt}),F&&i(i.S+i.F*(!H||u(function(){var t=N();return"[null]"!=$([t])||"{}"!=$({a:t})||"{}"!=$(Object(t))})),"JSON",{stringify:function(t){var e,n,r=[t],a=1;while(arguments.length>a)r.push(arguments[a++]);if(n=e=r[1],(j(e)||void 0!==t)&&!L(t))return v(e)||(e=function(t,e){if("function"==typeof n&&(e=n.call(this,t,e)),!L(e))return e}),r[1]=e,$.apply(F,r)}}),N[G][J]||n("35e8")(N[G],J,N[G].valueOf),l(N,"Symbol"),l(Math,"Math",!0),l(r.JSON,"JSON",!0)},"0395":function(t,e,n){var r=n("36c3"),a=n("6abf").f,o={}.toString,i="object"==typeof window&&window&&Object.getOwnPropertyNames?Object.getOwnPropertyNames(window):[],c=function(t){try{return a(t)}catch(e){return i.slice()}};t.exports.f=function(t){return i&&"[object Window]"==o.call(t)?c(t):a(r(t))}},"1aee":function(t,e,n){},"1ea0":function(t,e,n){"use strict";var r=n("d53f"),a=n.n(r);a.a},"268f":function(t,e,n){t.exports=n("fde4")},"32a6":function(t,e,n){var r=n("241e"),a=n("c3a1");n("ce7e")("keys",function(){return function(t){return a(r(t))}})},"454f":function(t,e,n){n("46a7");var r=n("584a").Object;t.exports=function(t,e,n){return r.defineProperty(t,e,n)}},"46a7":function(t,e,n){var r=n("63b6");r(r.S+r.F*!n("8e60"),"Object",{defineProperty:n("d9f6").f})},"47ee":function(t,e,n){var r=n("c3a1"),a=n("9aa9"),o=n("355d");t.exports=function(t){var e=r(t),n=a.f;if(n){var i,c=n(t),s=o.f,u=0;while(c.length>u)s.call(t,i=c[u++])&&e.push(i)}return e}},6718:function(t,e,n){var r=n("e53d"),a=n("584a"),o=n("b8e3"),i=n("ccb9"),c=n("d9f6").f;t.exports=function(t){var e=a.Symbol||(a.Symbol=o?{}:r.Symbol||{});"_"==t.charAt(0)||t in e||c(e,t,{value:i.f(t)})}},"6abf":function(t,e,n){var r=n("e6f3"),a=n("1691").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return r(t,a)}},7601:function(t,e,n){"use strict";n.r(e);var r=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",[n("h1",[t._v("Project Details")]),t.project?n("b-form",[n("h2",[t._v(t._s(t.project.hostname))]),n("b-form-group",{attrs:{id:"basedirGroup1",label:"Base Dir:","label-for":"basedirInput",description:"Base dir"}},[n("b-form-input",{attrs:{id:"basedirInput",type:"text",required:"",placeholder:"Enter base dir"},model:{value:t.baseDir,callback:function(e){t.baseDir=e},expression:"baseDir"}})],1),n("b-form-group",{attrs:{id:"pathGroup",label:"Dir Name:","label-for":"dirNameInput",description:""}},[n("b-form-input",{attrs:{id:"dirNameInput",type:"text",required:"",placeholder:""},model:{value:t.path,callback:function(e){t.path=e},expression:"path"}})],1),n("b-form-group",{attrs:{id:"fullPathGroup",label:"Full path:","label-for":"fullPathInput",description:""}},[n("b-form-input",{attrs:{id:"fullPathInput",type:"text",required:"",placeholder:""},model:{value:t.fullPath,callback:function(e){t.fullPath=e},expression:"fullPath"}})],1),n("b-form-group",{attrs:{id:"hostnameGroup",label:"Hostname:","label-for":"hostnameInput",description:""}},[n("b-form-input",{attrs:{id:"hostnameInput",type:"text",required:"",placeholder:""},model:{value:t.hostname,callback:function(e){t.hostname=e},expression:"hostname"}})],1),n("b-form-group",{attrs:{id:"notesGroup",label:"Notes:","label-for":"notesInput",description:""}},[n("b-form-textarea",{attrs:{id:"textarea",placeholder:"Enter something...",rows:"3","max-rows":"6"},model:{value:t.notes,callback:function(e){t.notes=e},expression:"notes"}})],1),n("b-form-group",{attrs:{id:"enabledGroup",label:"Status:","label-for":"enabledInput",description:""}},[n("b-form-radio",{attrs:{value:"true",name:"enabledInput"},model:{value:t.enabled,callback:function(e){t.enabled=e},expression:"enabled"}},[t._v("Enabled")]),n("b-form-radio",{attrs:{value:"false",name:"enabledInput"},model:{value:t.enabled,callback:function(e){t.enabled=e},expression:"enabled"}},[t._v("Disabled")])],1)],1):n("div",{staticClass:"project-details"},[n("h2",[t._v(t._s(this.$route.params.hostname))]),n("p",[t._v("This is a dummy project with no actual data!")])])],1)},a=[],o=n("cebc"),i=n("2f62"),c={name:"ProjectDetails",data:function(){return{hostname:"",notes:"",baseDir:"",path:"",fullPath:"",enabled:null,stack:{}}},watch:{"$route.params.hostname":{handler:function(t){var e=this.projectBy("hostname",t);e&&(this.hostname=e.hostname,this.notes=e.notes,this.baseDir=e.baseDir,this.path=e.path,this.fullPath=e.fullPath,this.enabled=e.enabled,this.stack=e.stack)},deep:!0,immediate:!0}},computed:Object(o["a"])({},Object(i["b"])(["projectBy","projectByName"]),{project:function(){return this.projectBy("hostname",this.$route.params.hostname)}}),methods:{onSubmit:function(t){var e=this;this.$store.dispatch("updateProject",{hostname:this.project.hostname,project:{hostname:this.hostname,notes:this.notes,baseDir:this.baseDir,path:this.path,enabled:this.enabled,fullPath:this.fullPath}}).then(function(){e.$router.push("/project/"+e.hostname)})}}},s=c,u=(n("8a41"),n("2877")),f=Object(u["a"])(s,r,a,!1,null,"db21c074",null);e["default"]=f.exports},"85f2":function(t,e,n){t.exports=n("454f")},"8a41":function(t,e,n){"use strict";var r=n("1aee"),a=n.n(r);a.a},"8aae":function(t,e,n){n("32a6"),t.exports=n("584a").Object.keys},"9aa9":function(t,e){e.f=Object.getOwnPropertySymbols},a4bb:function(t,e,n){t.exports=n("8aae")},acca:function(t,e,n){"use strict";n.r(e);var r=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("b-card-group",{attrs:{deck:""}},t._l(t.projects,function(e){return n("b-card",{key:e.path,attrs:{title:e.path,"sub-title":e.hostname,to:{path:"/project/"+e.hostname}}},[n("b-card-text",[n("dl",[n("dt",[t._v("Enabled:")]),n("dd",[t._v(t._s(e.enabled))]),n("dt",[t._v("Directory:")]),n("dd",[t._v(t._s(e.fullPath))])])]),n("div",{attrs:{slot:"footer"},slot:"footer"},[t._v("\n      Stack: "),n("strong",[t._v(t._s(e.stack[Object.keys(e.stack)[0]].named_stack))]),n("ul",{staticClass:"service-list"},t._l(e.stack,function(e){return n("li",{key:e.stack_role},[t._v("\n          "+t._s(e.program)+" "+t._s(e.version.major+"."+e.version.minor+"."+e.version.patch)+" "),n("small",{staticClass:"text-muted"},[t._v("("+t._s(e.role)+")")])])}),0),n("b-button",{attrs:{to:"/projects/"+e.hostname,variant:"primary"}},[t._v("Edit")])],1)],1)}),1)},a=[],o=n("cebc"),i=n("2f62"),c={name:"ProjectList",computed:Object(o["a"])({},Object(i["c"])(["projects"])),mounted:function(){this.$store.dispatch("loadProjectHeaders")}},s=c,u=(n("1ea0"),n("2877")),f=Object(u["a"])(s,r,a,!1,null,"90a66bae",null);e["default"]=f.exports},bf0b:function(t,e,n){var r=n("355d"),a=n("aebd"),o=n("36c3"),i=n("1bc3"),c=n("07e3"),s=n("794b"),u=Object.getOwnPropertyDescriptor;e.f=n("8e60")?u:function(t,e){if(t=o(t),e=i(e,!0),s)try{return u(t,e)}catch(n){}if(c(t,e))return a(!r.f.call(t,e),t[e])}},bf90:function(t,e,n){var r=n("36c3"),a=n("bf0b").f;n("ce7e")("getOwnPropertyDescriptor",function(){return function(t,e){return a(r(t),e)}})},ccb9:function(t,e,n){e.f=n("5168")},ce7e:function(t,e,n){var r=n("63b6"),a=n("584a"),o=n("294c");t.exports=function(t,e){var n=(a.Object||{})[t]||Object[t],i={};i[t]=e(n),r(r.S+r.F*o(function(){n(1)}),"Object",i)}},cebc:function(t,e,n){"use strict";var r=n("268f"),a=n.n(r),o=n("e265"),i=n.n(o),c=n("a4bb"),s=n.n(c),u=n("85f2"),f=n.n(u);function l(t,e,n){return e in t?f()(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}function b(t){for(var e=1;e<arguments.length;e++){var n=null!=arguments[e]?arguments[e]:{},r=s()(n);"function"===typeof i.a&&(r=r.concat(i()(n).filter(function(t){return a()(n,t).enumerable}))),r.forEach(function(e){l(t,e,n[e])})}return t}n.d(e,"a",function(){return b})},d53f:function(t,e,n){},e265:function(t,e,n){t.exports=n("ed33")},ebfd:function(t,e,n){var r=n("62a0")("meta"),a=n("f772"),o=n("07e3"),i=n("d9f6").f,c=0,s=Object.isExtensible||function(){return!0},u=!n("294c")(function(){return s(Object.preventExtensions({}))}),f=function(t){i(t,r,{value:{i:"O"+ ++c,w:{}}})},l=function(t,e){if(!a(t))return"symbol"==typeof t?t:("string"==typeof t?"S":"P")+t;if(!o(t,r)){if(!s(t))return"F";if(!e)return"E";f(t)}return t[r].i},b=function(t,e){if(!o(t,r)){if(!s(t))return!0;if(!e)return!1;f(t)}return t[r].w},p=function(t){return u&&d.NEED&&s(t)&&!o(t,r)&&f(t),t},d=t.exports={KEY:r,NEED:!1,fastKey:l,getWeak:b,onFreeze:p}},ed33:function(t,e,n){n("014b"),t.exports=n("584a").Object.getOwnPropertySymbols},fde4:function(t,e,n){n("bf90");var r=n("584a").Object;t.exports=function(t,e){return r.getOwnPropertyDescriptor(t,e)}}}]);
//# sourceMappingURL=project.b0825e80.js.map