(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projectstack"],{"014b":function(t,e,r){"use strict";var n=r("e53d"),o=r("07e3"),c=r("8e60"),a=r("63b6"),s=r("9138"),i=r("ebfd").KEY,f=r("294c"),u=r("dbdb"),p=r("45f2"),l=r("62a0"),b=r("5168"),v=r("ccb9"),d=r("6718"),g=r("47ee"),m=r("9003"),h=r("e4ae"),y=r("f772"),j=r("36c3"),k=r("1bc3"),S=r("aebd"),x=r("a159"),_=r("0395"),O=r("bf0b"),P=r("d9f6"),w=r("c3a1"),N=O.f,E=P.f,I=_.f,A=n.Symbol,C=n.JSON,F=C&&C.stringify,T="prototype",$=b("_hidden"),D=b("toPrimitive"),H={}.propertyIsEnumerable,R=u("symbol-registry"),U=u("symbols"),M=u("op-symbols"),G=Object[T],J="function"==typeof A,q=n.QObject,L=!q||!q[T]||!q[T].findChild,V=c&&f(function(){return 7!=x(E({},"a",{get:function(){return E(this,"a",{value:7}).a}})).a})?function(t,e,r){var n=N(G,e);n&&delete G[e],E(t,e,r),n&&t!==G&&E(G,e,n)}:E,Y=function(t){var e=U[t]=x(A[T]);return e._k=t,e},K=J&&"symbol"==typeof A.iterator?function(t){return"symbol"==typeof t}:function(t){return t instanceof A},W=function(t,e,r){return t===G&&W(M,e,r),h(t),e=k(e,!0),h(r),o(U,e)?(r.enumerable?(o(t,$)&&t[$][e]&&(t[$][e]=!1),r=x(r,{enumerable:S(0,!1)})):(o(t,$)||E(t,$,S(1,{})),t[$][e]=!0),V(t,e,r)):E(t,e,r)},z=function(t,e){h(t);var r,n=g(e=j(e)),o=0,c=n.length;while(c>o)W(t,r=n[o++],e[r]);return t},X=function(t,e){return void 0===e?x(t):z(x(t),e)},B=function(t){var e=H.call(this,t=k(t,!0));return!(this===G&&o(U,t)&&!o(M,t))&&(!(e||!o(this,t)||!o(U,t)||o(this,$)&&this[$][t])||e)},Q=function(t,e){if(t=j(t),e=k(e,!0),t!==G||!o(U,e)||o(M,e)){var r=N(t,e);return!r||!o(U,e)||o(t,$)&&t[$][e]||(r.enumerable=!0),r}},Z=function(t){var e,r=I(j(t)),n=[],c=0;while(r.length>c)o(U,e=r[c++])||e==$||e==i||n.push(e);return n},tt=function(t){var e,r=t===G,n=I(r?M:j(t)),c=[],a=0;while(n.length>a)!o(U,e=n[a++])||r&&!o(G,e)||c.push(U[e]);return c};J||(A=function(){if(this instanceof A)throw TypeError("Symbol is not a constructor!");var t=l(arguments.length>0?arguments[0]:void 0),e=function(r){this===G&&e.call(M,r),o(this,$)&&o(this[$],t)&&(this[$][t]=!1),V(this,t,S(1,r))};return c&&L&&V(G,t,{configurable:!0,set:e}),Y(t)},s(A[T],"toString",function(){return this._k}),O.f=Q,P.f=W,r("6abf").f=_.f=Z,r("355d").f=B,r("9aa9").f=tt,c&&!r("b8e3")&&s(G,"propertyIsEnumerable",B,!0),v.f=function(t){return Y(b(t))}),a(a.G+a.W+a.F*!J,{Symbol:A});for(var et="hasInstance,isConcatSpreadable,iterator,match,replace,search,species,split,toPrimitive,toStringTag,unscopables".split(","),rt=0;et.length>rt;)b(et[rt++]);for(var nt=w(b.store),ot=0;nt.length>ot;)d(nt[ot++]);a(a.S+a.F*!J,"Symbol",{for:function(t){return o(R,t+="")?R[t]:R[t]=A(t)},keyFor:function(t){if(!K(t))throw TypeError(t+" is not a symbol!");for(var e in R)if(R[e]===t)return e},useSetter:function(){L=!0},useSimple:function(){L=!1}}),a(a.S+a.F*!J,"Object",{create:X,defineProperty:W,defineProperties:z,getOwnPropertyDescriptor:Q,getOwnPropertyNames:Z,getOwnPropertySymbols:tt}),C&&a(a.S+a.F*(!J||f(function(){var t=A();return"[null]"!=F([t])||"{}"!=F({a:t})||"{}"!=F(Object(t))})),"JSON",{stringify:function(t){var e,r,n=[t],o=1;while(arguments.length>o)n.push(arguments[o++]);if(r=e=n[1],(y(e)||void 0!==t)&&!K(t))return m(e)||(e=function(t,e){if("function"==typeof r&&(e=r.call(this,t,e)),!K(e))return e}),n[1]=e,F.apply(C,n)}}),A[T][D]||r("35e8")(A[T],D,A[T].valueOf),p(A,"Symbol"),p(Math,"Math",!0),p(n.JSON,"JSON",!0)},"0395":function(t,e,r){var n=r("36c3"),o=r("6abf").f,c={}.toString,a="object"==typeof window&&window&&Object.getOwnPropertyNames?Object.getOwnPropertyNames(window):[],s=function(t){try{return o(t)}catch(e){return a.slice()}};t.exports.f=function(t){return a&&"[object Window]"==c.call(t)?s(t):o(n(t))}},"055b":function(t,e,r){t.exports=r.p+"img/python.51c2eab2.svg"},"090e":function(t,e,r){t.exports=r.p+"img/nodejs.94cafb0d.svg"},"11e9":function(t,e,r){var n=r("52a7"),o=r("4630"),c=r("6821"),a=r("6a99"),s=r("69a8"),i=r("c69a"),f=Object.getOwnPropertyDescriptor;e.f=r("9e1e")?f:function(t,e){if(t=c(t),e=a(e,!0),i)try{return f(t,e)}catch(r){}if(s(t,e))return o(!n.f.call(t,e),t[e])}},1540:function(t,e,r){t.exports=r.p+"img/php.fa78b345.svg"},"1b4d":function(t,e,r){t.exports=r.p+"img/flask.318d58cb.svg"},"268f":function(t,e,r){t.exports=r("fde4")},"2d1f":function(t,e,r){t.exports=r("b606")},"319f":function(t,e,r){t.exports=r.p+"img/rails.2db29782.svg"},"31e8":function(t,e,r){var n={"./angular.svg":"a230","./apache.svg":"b77f","./codeigniter.svg":"7939","./django.svg":"c6da","./drupal.svg":"a88c","./elasticsearch.svg":"81bb","./flask.svg":"1b4d","./joomla.svg":"5390","./laravel.svg":"41c8","./logo.svg":"9b19","./mariadb.svg":"613e","./memcached.svg":"a0ba","./mysql.svg":"6c4c","./nginx.svg":"c502","./nodejs.svg":"090e","./perl.svg":"c44f","./php.svg":"1540","./python.svg":"055b","./rails.svg":"319f","./react.svg":"b2e9","./redis.svg":"8bcb","./ruby.svg":"9401","./wordpress.svg":"ee3c"};function o(t){var e=c(t);return r(e)}function c(t){var e=n[t];if(!(e+1)){var r=new Error("Cannot find module '"+t+"'");throw r.code="MODULE_NOT_FOUND",r}return e}o.keys=function(){return Object.keys(n)},o.resolve=c,t.exports=o,o.id="31e8"},"32a6":function(t,e,r){var n=r("241e"),o=r("c3a1");r("ce7e")("keys",function(){return function(t){return o(n(t))}})},3983:function(t,e,r){"use strict";r.r(e);var n=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"project-stack-list",attrs:{role:"tablist",id:t.project_base+"-stack"}},[t._l(t.groupProjectStacks(t.projectStack),function(e,o,c){return n("div",{key:o,staticClass:"project-stack"},[n("h2",{staticClass:"stack-title"},[t._v(t._s(o.replace("gearbox.works/","")))]),n("b-button",{staticClass:"js-remove-stack",attrs:{tabindex:100*t.projectIndex+10*c,size:"sm",variant:"outline-secondary","aria-label":"Remove this stack from project",title:"Remove this stack from project"},on:{click:function(e){return e.preventDefault(),t.removeProjectStack(o)}}},[t._v("×")]),n("ul",{staticClass:"service-list"},t._l(t.stackServices(o),function(a,s,i){return n("li",{key:t.project_base+t.escAttr(s),staticClass:"service-item",attrs:{id:t.project_base+t.escAttr(s),tabindex:100*t.projectIndex+10*c+i+1}},[n("img",{staticClass:"service-program",attrs:{src:r("31e8")("./"+e[s].program+".svg")}}),n("h6",{staticClass:"service-role"},[t._v(t._s(t.stackRoles(o)[s].label))]),n("b-popover",{ref:t.project_base+t.escAttr(s)+"_popover",refInFor:!0,attrs:{target:t.project_base+t.escAttr(s),container:t.projectHostname+"-stack",triggers:"focus",placement:"bottom"}},[n("template",{slot:"title"},[n("b-button",{staticClass:"close",attrs:{"aria-label":"Close"},on:{click:function(e){t.onClosePopoverFor(t.project_base+t.escAttr(s))}}},[n("span",{staticClass:"d-inline-block",attrs:{"aria-hidden":"true"}},[t._v("×")])]),t._v("\n            "+t._s(t.stackRoles(o)[s].name)+"\n          ")],1),n("div",[n("b-form-select",{attrs:{id:t.project_base+t.escAttr(s)+"_input",value:e[s]?e[s].service_id:null,tabindex:100*t.projectIndex+10*c+i+9},on:{change:function(e){return t.changeProjectService(s,e)}}},[n("option",{attrs:{disabled:"",value:""}},[t._v("Please select one...")]),t._l(t.optionGroups(a.options),function(e,r){return n("optgroup",{key:r,attrs:{label:r}},t._l(e,function(e){return n("option",{key:e,domProps:{value:a.org+"/"+e}},[t._v(t._s(e))])}),0)})],2)],1)],2)],1)}),0)],1)}),t.hasUnusedStacks?n("b-form-select",{staticClass:"add-stack",on:{change:t.addProjectStack},model:{value:t.stackToAdd,callback:function(e){t.stackToAdd=e},expression:"stackToAdd"}},[n("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Add Stack...")]),t._l(t.stacksNotUnusedInProject,function(e,r){return n("option",{key:r,domProps:{value:r}},[t._v(t._s(r.replace("gearbox.works/","")))])})],2):t._e()],2)},o=[],c=(r("28a5"),r("a481"),r("2d1f")),a=r.n(c),s=r("cebc"),i=(r("c5f6"),r("2f62")),f={name:"ProjectStack",props:{projectHostname:{type:String,required:!0},projectStack:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return{stackToAdd:null}},components:{},computed:Object(s["a"])({},Object(i["b"])(["groupProjectStacks","stackRoles","stackServices"]),{project_base:function(){return this.escAttr(this.projectHostname)},project:function(){return this.$store.getters.projectBy("hostname",this.projectHostname)},hasUnusedStacks:function(){return a()(this.stacksNotUnusedInProject).length>0},stacksNotUnusedInProject:function(){var t={},e=this.groupProjectStacks(this.projectStack);for(var r in this.$store.state.gearStacks){var n=this.$store.state.gearStacks[r];"undefined"===typeof e[n]&&(t[n]=this.$store.state.gearStacks[n])}return t}}),methods:{escAttr:function(t){return t.replace(/\//g,"_").replace(/\./g,"_")},stackIncludesService:function(t,e){var r=!1;for(var n in t)if(e===t[n].service_id){r=!0;break}return r&&console.log("found",e),r},mapOptions:function(t){var e=[];for(var r in t)e.push({value:r,text:t[r]});return e},optionGroups:function(t){var e={};for(var r in t){var n=t[r].split(":")[0];"undefined"===typeof e[n]&&(e[n]={}),e[n][r]=t[r]}return e},addProjectStack:function(t){console.log("Selected",this.stackToAdd,t),this.$store.dispatch("addProjectStack",{projectHostname:this.projectHostname,stackName:t}),this.stackToAdd=null},changeProjectService:function(t,e){this.$store.dispatch("changeProjectService",{projectHostname:this.projectHostname,serviceName:t,serviceId:e})},onClosePopoverFor:function(t){console.log("onClosePopoverFor",t),this.$root.$emit("bv::popover::hide",t)},removeProjectStack:function(t){this.$store.dispatch("removeProjectStack",{projectHostname:this.projectHostname,stackName:t})}}},u=f,p=(r("a705"),r("2877")),l=Object(p["a"])(u,n,o,!1,null,"e0b24eca",null);e["default"]=l.exports},"41c8":function(t,e,r){t.exports=r.p+"img/laravel.1766a461.svg"},"454f":function(t,e,r){r("46a7");var n=r("584a").Object;t.exports=function(t,e,r){return n.defineProperty(t,e,r)}},"46a7":function(t,e,r){var n=r("63b6");n(n.S+n.F*!r("8e60"),"Object",{defineProperty:r("d9f6").f})},"47ee":function(t,e,r){var n=r("c3a1"),o=r("9aa9"),c=r("355d");t.exports=function(t){var e=n(t),r=o.f;if(r){var a,s=r(t),i=c.f,f=0;while(s.length>f)i.call(t,a=s[f++])&&e.push(a)}return e}},5390:function(t,e,r){t.exports=r.p+"img/joomla.d8aa2e45.svg"},"5dbc":function(t,e,r){var n=r("d3f4"),o=r("8b97").set;t.exports=function(t,e,r){var c,a=e.constructor;return a!==r&&"function"==typeof a&&(c=a.prototype)!==r.prototype&&n(c)&&o&&o(t,c),t}},"613e":function(t,e,r){t.exports=r.p+"img/mariadb.e16110bc.svg"},6718:function(t,e,r){var n=r("e53d"),o=r("584a"),c=r("b8e3"),a=r("ccb9"),s=r("d9f6").f;t.exports=function(t){var e=o.Symbol||(o.Symbol=c?{}:n.Symbol||{});"_"==t.charAt(0)||t in e||s(e,t,{value:a.f(t)})}},"6abf":function(t,e,r){var n=r("e6f3"),o=r("1691").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return n(t,o)}},"6c4c":function(t,e,r){t.exports=r.p+"img/mysql.dd2a5a35.svg"},7939:function(t,e,r){t.exports=r.p+"img/codeigniter.434bf735.svg"},"81bb":function(t,e,r){t.exports=r.p+"img/elasticsearch.3ecfa530.svg"},"85f2":function(t,e,r){t.exports=r("454f")},"8aae":function(t,e,r){r("32a6"),t.exports=r("584a").Object.keys},"8b97":function(t,e,r){var n=r("d3f4"),o=r("cb7c"),c=function(t,e){if(o(t),!n(e)&&null!==e)throw TypeError(e+": can't set as prototype!")};t.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(t,e,n){try{n=r("9b43")(Function.call,r("11e9").f(Object.prototype,"__proto__").set,2),n(t,[]),e=!(t instanceof Array)}catch(o){e=!0}return function(t,r){return c(t,r),e?t.__proto__=r:n(t,r),t}}({},!1):void 0),check:c}},"8bcb":function(t,e,r){t.exports=r.p+"img/redis.3c39fafe.svg"},9093:function(t,e,r){var n=r("ce10"),o=r("e11e").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return n(t,o)}},9401:function(t,e,r){t.exports=r.p+"img/ruby.514befa7.svg"},"9aa9":function(t,e){e.f=Object.getOwnPropertySymbols},"9b19":function(t,e,r){t.exports=r.p+"img/logo.63a7d78d.svg"},"9c60":function(t,e,r){var n=r("63b6"),o=r("13c8")(!0);n(n.S,"Object",{entries:function(t){return o(t)}})},a0ba:function(t,e,r){t.exports=r.p+"img/memcached.2bcccabf.svg"},a230:function(t,e,r){t.exports=r.p+"img/angular.e224f5ed.svg"},a4bb:function(t,e,r){t.exports=r("8aae")},a705:function(t,e,r){"use strict";var n=r("e06e"),o=r.n(n);o.a},a88c:function(t,e,r){t.exports=r.p+"img/drupal.66089b06.svg"},aa77:function(t,e,r){var n=r("5ca1"),o=r("be13"),c=r("79e5"),a=r("fdef"),s="["+a+"]",i="​",f=RegExp("^"+s+s+"*"),u=RegExp(s+s+"*$"),p=function(t,e,r){var o={},s=c(function(){return!!a[t]()||i[t]()!=i}),f=o[t]=s?e(l):a[t];r&&(o[r]=f),n(n.P+n.F*s,"String",o)},l=p.trim=function(t,e){return t=String(o(t)),1&e&&(t=t.replace(f,"")),2&e&&(t=t.replace(u,"")),t};t.exports=p},b2e9:function(t,e,r){t.exports=r.p+"img/react.9a28da9f.svg"},b606:function(t,e,r){r("9c60"),t.exports=r("584a").Object.entries},b77f:function(t,e,r){t.exports=r.p+"img/apache.12c49354.svg"},bf0b:function(t,e,r){var n=r("355d"),o=r("aebd"),c=r("36c3"),a=r("1bc3"),s=r("07e3"),i=r("794b"),f=Object.getOwnPropertyDescriptor;e.f=r("8e60")?f:function(t,e){if(t=c(t),e=a(e,!0),i)try{return f(t,e)}catch(r){}if(s(t,e))return o(!n.f.call(t,e),t[e])}},bf90:function(t,e,r){var n=r("36c3"),o=r("bf0b").f;r("ce7e")("getOwnPropertyDescriptor",function(){return function(t,e){return o(n(t),e)}})},c44f:function(t,e,r){t.exports=r.p+"img/perl.a025edb4.svg"},c502:function(t,e,r){t.exports=r.p+"img/nginx.eae76401.svg"},c5f6:function(t,e,r){"use strict";var n=r("7726"),o=r("69a8"),c=r("2d95"),a=r("5dbc"),s=r("6a99"),i=r("79e5"),f=r("9093").f,u=r("11e9").f,p=r("86cc").f,l=r("aa77").trim,b="Number",v=n[b],d=v,g=v.prototype,m=c(r("2aeb")(g))==b,h="trim"in String.prototype,y=function(t){var e=s(t,!1);if("string"==typeof e&&e.length>2){e=h?e.trim():l(e,3);var r,n,o,c=e.charCodeAt(0);if(43===c||45===c){if(r=e.charCodeAt(2),88===r||120===r)return NaN}else if(48===c){switch(e.charCodeAt(1)){case 66:case 98:n=2,o=49;break;case 79:case 111:n=8,o=55;break;default:return+e}for(var a,i=e.slice(2),f=0,u=i.length;f<u;f++)if(a=i.charCodeAt(f),a<48||a>o)return NaN;return parseInt(i,n)}}return+e};if(!v(" 0o1")||!v("0b1")||v("+0x1")){v=function(t){var e=arguments.length<1?0:t,r=this;return r instanceof v&&(m?i(function(){g.valueOf.call(r)}):c(r)!=b)?a(new d(y(e)),r,v):y(e)};for(var j,k=r("9e1e")?f(d):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),S=0;k.length>S;S++)o(d,j=k[S])&&!o(v,j)&&p(v,j,u(d,j));v.prototype=g,g.constructor=v,r("2aba")(n,b,v)}},c6da:function(t,e,r){t.exports=r.p+"img/django.28fe09a0.svg"},ccb9:function(t,e,r){e.f=r("5168")},ce7e:function(t,e,r){var n=r("63b6"),o=r("584a"),c=r("294c");t.exports=function(t,e){var r=(o.Object||{})[t]||Object[t],a={};a[t]=e(r),n(n.S+n.F*c(function(){r(1)}),"Object",a)}},cebc:function(t,e,r){"use strict";var n=r("268f"),o=r.n(n),c=r("e265"),a=r.n(c),s=r("a4bb"),i=r.n(s),f=r("85f2"),u=r.n(f);function p(t,e,r){return e in t?u()(t,e,{value:r,enumerable:!0,configurable:!0,writable:!0}):t[e]=r,t}function l(t){for(var e=1;e<arguments.length;e++){var r=null!=arguments[e]?arguments[e]:{},n=i()(r);"function"===typeof a.a&&(n=n.concat(a()(r).filter(function(t){return o()(r,t).enumerable}))),n.forEach(function(e){p(t,e,r[e])})}return t}r.d(e,"a",function(){return l})},e06e:function(t,e,r){},e265:function(t,e,r){t.exports=r("ed33")},ebfd:function(t,e,r){var n=r("62a0")("meta"),o=r("f772"),c=r("07e3"),a=r("d9f6").f,s=0,i=Object.isExtensible||function(){return!0},f=!r("294c")(function(){return i(Object.preventExtensions({}))}),u=function(t){a(t,n,{value:{i:"O"+ ++s,w:{}}})},p=function(t,e){if(!o(t))return"symbol"==typeof t?t:("string"==typeof t?"S":"P")+t;if(!c(t,n)){if(!i(t))return"F";if(!e)return"E";u(t)}return t[n].i},l=function(t,e){if(!c(t,n)){if(!i(t))return!0;if(!e)return!1;u(t)}return t[n].w},b=function(t){return f&&v.NEED&&i(t)&&!c(t,n)&&u(t),t},v=t.exports={KEY:n,NEED:!1,fastKey:p,getWeak:l,onFreeze:b}},ed33:function(t,e,r){r("014b"),t.exports=r("584a").Object.getOwnPropertySymbols},ee3c:function(t,e,r){t.exports=r.p+"img/wordpress.b08e20e3.svg"},fde4:function(t,e,r){r("bf90");var n=r("584a").Object;t.exports=function(t,e){return n.getOwnPropertyDescriptor(t,e)}},fdef:function(t,e){t.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);
//# sourceMappingURL=projectstack.fcb38643.js.map