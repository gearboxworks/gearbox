(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projectstack"],{"014b":function(t,e,r){"use strict";var n=r("e53d"),o=r("07e3"),c=r("8e60"),a=r("63b6"),i=r("9138"),u=r("ebfd").KEY,s=r("294c"),f=r("dbdb"),l=r("45f2"),p=r("62a0"),b=r("5168"),d=r("ccb9"),v=r("6718"),h=r("47ee"),g=r("9003"),y=r("e4ae"),m=r("f772"),x=r("36c3"),k=r("1bc3"),S=r("aebd"),j=r("a159"),_=r("0395"),w=r("bf0b"),O=r("d9f6"),P=r("c3a1"),E=w.f,A=O.f,R=_.f,$=n.Symbol,N=n.JSON,I=N&&N.stringify,C="prototype",T=b("_hidden"),F=b("toPrimitive"),D={}.propertyIsEnumerable,J=f("symbol-registry"),M=f("symbols"),H=f("op-symbols"),q=Object[C],G="function"==typeof $,K=n.QObject,W=!K||!K[C]||!K[C].findChild,V=c&&s(function(){return 7!=j(A({},"a",{get:function(){return A(this,"a",{value:7}).a}})).a})?function(t,e,r){var n=E(q,e);n&&delete q[e],A(t,e,r),n&&t!==q&&A(q,e,n)}:A,Y=function(t){var e=M[t]=j($[C]);return e._k=t,e},z=G&&"symbol"==typeof $.iterator?function(t){return"symbol"==typeof t}:function(t){return t instanceof $},B=function(t,e,r){return t===q&&B(H,e,r),y(t),e=k(e,!0),y(r),o(M,e)?(r.enumerable?(o(t,T)&&t[T][e]&&(t[T][e]=!1),r=j(r,{enumerable:S(0,!1)})):(o(t,T)||A(t,T,S(1,{})),t[T][e]=!0),V(t,e,r)):A(t,e,r)},Q=function(t,e){y(t);var r,n=h(e=x(e)),o=0,c=n.length;while(c>o)B(t,r=n[o++],e[r]);return t},L=function(t,e){return void 0===e?j(t):Q(j(t),e)},U=function(t){var e=D.call(this,t=k(t,!0));return!(this===q&&o(M,t)&&!o(H,t))&&(!(e||!o(this,t)||!o(M,t)||o(this,T)&&this[T][t])||e)},X=function(t,e){if(t=x(t),e=k(e,!0),t!==q||!o(M,e)||o(H,e)){var r=E(t,e);return!r||!o(M,e)||o(t,T)&&t[T][e]||(r.enumerable=!0),r}},Z=function(t){var e,r=R(x(t)),n=[],c=0;while(r.length>c)o(M,e=r[c++])||e==T||e==u||n.push(e);return n},tt=function(t){var e,r=t===q,n=R(r?H:x(t)),c=[],a=0;while(n.length>a)!o(M,e=n[a++])||r&&!o(q,e)||c.push(M[e]);return c};G||($=function(){if(this instanceof $)throw TypeError("Symbol is not a constructor!");var t=p(arguments.length>0?arguments[0]:void 0),e=function(r){this===q&&e.call(H,r),o(this,T)&&o(this[T],t)&&(this[T][t]=!1),V(this,t,S(1,r))};return c&&W&&V(q,t,{configurable:!0,set:e}),Y(t)},i($[C],"toString",function(){return this._k}),w.f=X,O.f=B,r("6abf").f=_.f=Z,r("355d").f=U,r("9aa9").f=tt,c&&!r("b8e3")&&i(q,"propertyIsEnumerable",U,!0),d.f=function(t){return Y(b(t))}),a(a.G+a.W+a.F*!G,{Symbol:$});for(var et="hasInstance,isConcatSpreadable,iterator,match,replace,search,species,split,toPrimitive,toStringTag,unscopables".split(","),rt=0;et.length>rt;)b(et[rt++]);for(var nt=P(b.store),ot=0;nt.length>ot;)v(nt[ot++]);a(a.S+a.F*!G,"Symbol",{for:function(t){return o(J,t+="")?J[t]:J[t]=$(t)},keyFor:function(t){if(!z(t))throw TypeError(t+" is not a symbol!");for(var e in J)if(J[e]===t)return e},useSetter:function(){W=!0},useSimple:function(){W=!1}}),a(a.S+a.F*!G,"Object",{create:L,defineProperty:B,defineProperties:Q,getOwnPropertyDescriptor:X,getOwnPropertyNames:Z,getOwnPropertySymbols:tt}),N&&a(a.S+a.F*(!G||s(function(){var t=$();return"[null]"!=I([t])||"{}"!=I({a:t})||"{}"!=I(Object(t))})),"JSON",{stringify:function(t){var e,r,n=[t],o=1;while(arguments.length>o)n.push(arguments[o++]);if(r=e=n[1],(m(e)||void 0!==t)&&!z(t))return g(e)||(e=function(t,e){if("function"==typeof r&&(e=r.call(this,t,e)),!z(e))return e}),n[1]=e,I.apply(N,n)}}),$[C][F]||r("35e8")($[C],F,$[C].valueOf),l($,"Symbol"),l(Math,"Math",!0),l(n.JSON,"JSON",!0)},"02f4":function(t,e,r){var n=r("4588"),o=r("be13");t.exports=function(t){return function(e,r){var c,a,i=String(o(e)),u=n(r),s=i.length;return u<0||u>=s?t?"":void 0:(c=i.charCodeAt(u),c<55296||c>56319||u+1===s||(a=i.charCodeAt(u+1))<56320||a>57343?t?i.charAt(u):c:t?i.slice(u,u+2):a-56320+(c-55296<<10)+65536)}}},"0390":function(t,e,r){"use strict";var n=r("02f4")(!0);t.exports=function(t,e,r){return e+(r?n(t,e).length:1)}},"0395":function(t,e,r){var n=r("36c3"),o=r("6abf").f,c={}.toString,a="object"==typeof window&&window&&Object.getOwnPropertyNames?Object.getOwnPropertyNames(window):[],i=function(t){try{return o(t)}catch(e){return a.slice()}};t.exports.f=function(t){return a&&"[object Window]"==c.call(t)?i(t):o(n(t))}},"0bfb":function(t,e,r){"use strict";var n=r("cb7c");t.exports=function(){var t=n(this),e="";return t.global&&(e+="g"),t.ignoreCase&&(e+="i"),t.multiline&&(e+="m"),t.unicode&&(e+="u"),t.sticky&&(e+="y"),e}},"214f":function(t,e,r){"use strict";r("b0c5");var n=r("2aba"),o=r("32e9"),c=r("79e5"),a=r("be13"),i=r("2b4c"),u=r("520a"),s=i("species"),f=!c(function(){var t=/./;return t.exec=function(){var t=[];return t.groups={a:"7"},t},"7"!=="".replace(t,"$<a>")}),l=function(){var t=/(?:)/,e=t.exec;t.exec=function(){return e.apply(this,arguments)};var r="ab".split(t);return 2===r.length&&"a"===r[0]&&"b"===r[1]}();t.exports=function(t,e,r){var p=i(t),b=!c(function(){var e={};return e[p]=function(){return 7},7!=""[t](e)}),d=b?!c(function(){var e=!1,r=/a/;return r.exec=function(){return e=!0,null},"split"===t&&(r.constructor={},r.constructor[s]=function(){return r}),r[p](""),!e}):void 0;if(!b||!d||"replace"===t&&!f||"split"===t&&!l){var v=/./[p],h=r(a,p,""[t],function(t,e,r,n,o){return e.exec===u?b&&!o?{done:!0,value:v.call(e,r,n)}:{done:!0,value:t.call(r,e,n)}:{done:!1}}),g=h[0],y=h[1];n(String.prototype,t,g),o(RegExp.prototype,p,2==e?function(t,e){return y.call(t,this,e)}:function(t){return y.call(t,this)})}}},"268f":function(t,e,r){t.exports=r("fde4")},"28a5":function(t,e,r){"use strict";var n=r("aae3"),o=r("cb7c"),c=r("ebd6"),a=r("0390"),i=r("9def"),u=r("5f1b"),s=r("520a"),f=r("79e5"),l=Math.min,p=[].push,b="split",d="length",v="lastIndex",h=4294967295,g=!f(function(){RegExp(h,"y")});r("214f")("split",2,function(t,e,r,f){var y;return y="c"=="abbc"[b](/(b)*/)[1]||4!="test"[b](/(?:)/,-1)[d]||2!="ab"[b](/(?:ab)*/)[d]||4!="."[b](/(.?)(.?)/)[d]||"."[b](/()()/)[d]>1||""[b](/.?/)[d]?function(t,e){var o=String(this);if(void 0===t&&0===e)return[];if(!n(t))return r.call(o,t,e);var c,a,i,u=[],f=(t.ignoreCase?"i":"")+(t.multiline?"m":"")+(t.unicode?"u":"")+(t.sticky?"y":""),l=0,b=void 0===e?h:e>>>0,g=new RegExp(t.source,f+"g");while(c=s.call(g,o)){if(a=g[v],a>l&&(u.push(o.slice(l,c.index)),c[d]>1&&c.index<o[d]&&p.apply(u,c.slice(1)),i=c[0][d],l=a,u[d]>=b))break;g[v]===c.index&&g[v]++}return l===o[d]?!i&&g.test("")||u.push(""):u.push(o.slice(l)),u[d]>b?u.slice(0,b):u}:"0"[b](void 0,0)[d]?function(t,e){return void 0===t&&0===e?[]:r.call(this,t,e)}:r,[function(r,n){var o=t(this),c=void 0==r?void 0:r[e];return void 0!==c?c.call(r,o,n):y.call(String(o),r,n)},function(t,e){var n=f(y,t,this,e,y!==r);if(n.done)return n.value;var s=o(t),p=String(this),b=c(s,RegExp),d=s.unicode,v=(s.ignoreCase?"i":"")+(s.multiline?"m":"")+(s.unicode?"u":"")+(g?"y":"g"),m=new b(g?s:"^(?:"+s.source+")",v),x=void 0===e?h:e>>>0;if(0===x)return[];if(0===p.length)return null===u(m,p)?[p]:[];var k=0,S=0,j=[];while(S<p.length){m.lastIndex=g?S:0;var _,w=u(m,g?p:p.slice(S));if(null===w||(_=l(i(m.lastIndex+(g?0:S)),p.length))===k)S=a(p,S,d);else{if(j.push(p.slice(k,S)),j.length===x)return j;for(var O=1;O<=w.length-1;O++)if(j.push(w[O]),j.length===x)return j;S=k=_}}return j.push(p.slice(k)),j}]})},"32a6":function(t,e,r){var n=r("241e"),o=r("c3a1");r("ce7e")("keys",function(){return function(t){return o(n(t))}})},3983:function(t,e,r){"use strict";r.r(e);var n=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{staticClass:"project-stack",attrs:{role:"tablist"}},[t._l(t.groupProjectStacks(t.projectStack),function(e,n,o){return r("b-card",{key:n,staticClass:"mb-1",attrs:{"no-body":""}},[r("b-card-header",{staticClass:"p-1",attrs:{"header-tag":"header",role:"tab"}},[r("b-button",{directives:[{name:"b-toggle",rawName:"v-b-toggle",value:t.project_base+"_accordion_"+o,expression:"project_base + '_accordion_' + index"}],attrs:{block:"",href:"#",variant:"info"}},[r("project-stack-header",{attrs:{stackName:n,stackRoles:e}})],1)],1),r("b-collapse",{attrs:{id:t.project_base+"_accordion_"+o,accordion:t.project_base+"_accordion",role:"tabpanel"}},[r("b-card-body",[r("b-form",t._l(t.stackServices(n),function(o,c){return r("b-form-group",{key:t.project_base+t.escAttr(c),attrs:{label:t.stackRoles(n)[c].short_label,"label-for":t.project_base+t.escAttr(c)+"_input",description:t.stackRoles(n)[c].label,"label-cols-sm":"4","label-cols-lg":"3"}},[r("b-form-select",{attrs:{id:t.project_base+t.escAttr(c)+"_input"}},t._l(t.optionGroups(o.options),function(n,c){return r("optgroup",{key:c,attrs:{label:c}},t._l(n,function(n){return r("option",{key:n,domProps:{value:n,selected:t.stackIncludesService(e,o.org+"/"+n)}},[t._v(t._s(n))])}),0)}),0)],1)}),1)],1)],1)],1)}),r("br"),r("b-form-select",{on:{change:t.addProjectStack},model:{value:t.stackToAdd,callback:function(e){t.stackToAdd=e},expression:"stackToAdd"}},[r("option",{domProps:{value:null}},[t._v("Add Stack...")]),t._l(t.unusedProjectStacks(t.projectStack),function(e,n){return r("option",{key:n},[t._v(t._s(n.replace("gearbox.works/","")))])})],2)],2)},o=[],c=(r("28a5"),r("a481"),r("cebc")),a=(r("cadf"),r("551c"),r("f751"),r("097d"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{staticClass:"project-stack-header"},[r("strong",[t._v(t._s(t.stackName.replace("gearbox.works/","")))]),t._v(" (\n  "),t._l(t.stackRoles,function(e,n,o){return r("span",{key:n,attrs:{title:n.replace(e.stack+"/","")}},[o?r("span",[t._v(",")]):t._e(),t._v(" "+t._s(e.program)+" "),r("small",[t._v(t._s(t.serviceVersion(e.version)))])])}),t._v("\n  )\n")],2)}),i=[],u={name:"ProjectStackHeader",props:{stackName:{type:String,required:!0},stackRoles:{type:Object,required:!0}},methods:{serviceVersion:function(t){var e="";return t.major&&(e+=t.major,t.minor&&(e+="."+t.minor,t.patch&&(e+="."+t.patch))),e}}},s=u,f=r("2877"),l=Object(f["a"])(s,a,i,!1,null,"5b26d090",null),p=l.exports,b=r("2f62"),d={name:"ProjectStack",props:{projectHostname:{type:String,required:!0},projectStack:{type:Object,required:!0}},data:function(){return{stackToAdd:null}},components:{ProjectStackHeader:p},computed:Object(c["a"])({},Object(b["b"])(["groupProjectStacks","stackRoles","stackServices"]),{project_base:function(){return this.escAttr(this.projectHostname)},project:function(){return this.$store.getters.projectBy("hostname",this.projectHostname)}}),methods:{escAttr:function(t){return t.replace(/\//g,"_").replace(/\./g,"_")},stackIncludesService:function(t,e){var r=!1;for(var n in t)if(e===t[n].service_id){r=!0;break}return r},unusedProjectStacks:function(t){var e={},r=this.groupProjectStacks(t);for(var n in this.$store.state.gearStacks){var o=this.$store.state.gearStacks[n];"undefined"===typeof r[o]&&(e[o]=this.$store.state.gearStacks[o])}return e},mapOptions:function(t){var e=[];for(var r in t)e.push({value:r,text:t[r]});return e},optionGroups:function(t){var e={};for(var r in t){var n=t[r].split(":")[0];"undefined"===typeof e[n]&&(e[n]={}),e[n][r]=t[r]}return e},addProjectStack:function(t){console.log("Selected",this.stackToAdd,t),this.$store.dispatch("addProjectStack",t),this.stackToAdd=null}}},v=d,h=Object(f["a"])(v,n,o,!1,null,"14889c9b",null);e["default"]=h.exports},"454f":function(t,e,r){r("46a7");var n=r("584a").Object;t.exports=function(t,e,r){return n.defineProperty(t,e,r)}},"46a7":function(t,e,r){var n=r("63b6");n(n.S+n.F*!r("8e60"),"Object",{defineProperty:r("d9f6").f})},"47ee":function(t,e,r){var n=r("c3a1"),o=r("9aa9"),c=r("355d");t.exports=function(t){var e=n(t),r=o.f;if(r){var a,i=r(t),u=c.f,s=0;while(i.length>s)u.call(t,a=i[s++])&&e.push(a)}return e}},"520a":function(t,e,r){"use strict";var n=r("0bfb"),o=RegExp.prototype.exec,c=String.prototype.replace,a=o,i="lastIndex",u=function(){var t=/a/,e=/b*/g;return o.call(t,"a"),o.call(e,"a"),0!==t[i]||0!==e[i]}(),s=void 0!==/()??/.exec("")[1],f=u||s;f&&(a=function(t){var e,r,a,f,l=this;return s&&(r=new RegExp("^"+l.source+"$(?!\\s)",n.call(l))),u&&(e=l[i]),a=o.call(l,t),u&&a&&(l[i]=l.global?a.index+a[0].length:e),s&&a&&a.length>1&&c.call(a[0],r,function(){for(f=1;f<arguments.length-2;f++)void 0===arguments[f]&&(a[f]=void 0)}),a}),t.exports=a},"5f1b":function(t,e,r){"use strict";var n=r("23c6"),o=RegExp.prototype.exec;t.exports=function(t,e){var r=t.exec;if("function"===typeof r){var c=r.call(t,e);if("object"!==typeof c)throw new TypeError("RegExp exec method returned something other than an Object or null");return c}if("RegExp"!==n(t))throw new TypeError("RegExp#exec called on incompatible receiver");return o.call(t,e)}},6718:function(t,e,r){var n=r("e53d"),o=r("584a"),c=r("b8e3"),a=r("ccb9"),i=r("d9f6").f;t.exports=function(t){var e=o.Symbol||(o.Symbol=c?{}:n.Symbol||{});"_"==t.charAt(0)||t in e||i(e,t,{value:a.f(t)})}},"6abf":function(t,e,r){var n=r("e6f3"),o=r("1691").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return n(t,o)}},"85f2":function(t,e,r){t.exports=r("454f")},"8aae":function(t,e,r){r("32a6"),t.exports=r("584a").Object.keys},"9aa9":function(t,e){e.f=Object.getOwnPropertySymbols},a481:function(t,e,r){"use strict";var n=r("cb7c"),o=r("4bf8"),c=r("9def"),a=r("4588"),i=r("0390"),u=r("5f1b"),s=Math.max,f=Math.min,l=Math.floor,p=/\$([$&`']|\d\d?|<[^>]*>)/g,b=/\$([$&`']|\d\d?)/g,d=function(t){return void 0===t?t:String(t)};r("214f")("replace",2,function(t,e,r,v){return[function(n,o){var c=t(this),a=void 0==n?void 0:n[e];return void 0!==a?a.call(n,c,o):r.call(String(c),n,o)},function(t,e){var o=v(r,t,this,e);if(o.done)return o.value;var l=n(t),p=String(this),b="function"===typeof e;b||(e=String(e));var g=l.global;if(g){var y=l.unicode;l.lastIndex=0}var m=[];while(1){var x=u(l,p);if(null===x)break;if(m.push(x),!g)break;var k=String(x[0]);""===k&&(l.lastIndex=i(p,c(l.lastIndex),y))}for(var S="",j=0,_=0;_<m.length;_++){x=m[_];for(var w=String(x[0]),O=s(f(a(x.index),p.length),0),P=[],E=1;E<x.length;E++)P.push(d(x[E]));var A=x.groups;if(b){var R=[w].concat(P,O,p);void 0!==A&&R.push(A);var $=String(e.apply(void 0,R))}else $=h(w,p,O,P,A,e);O>=j&&(S+=p.slice(j,O)+$,j=O+w.length)}return S+p.slice(j)}];function h(t,e,n,c,a,i){var u=n+t.length,s=c.length,f=b;return void 0!==a&&(a=o(a),f=p),r.call(i,f,function(r,o){var i;switch(o.charAt(0)){case"$":return"$";case"&":return t;case"`":return e.slice(0,n);case"'":return e.slice(u);case"<":i=a[o.slice(1,-1)];break;default:var f=+o;if(0===f)return r;if(f>s){var p=l(f/10);return 0===p?r:p<=s?void 0===c[p-1]?o.charAt(1):c[p-1]+o.charAt(1):r}i=c[f-1]}return void 0===i?"":i})}})},a4bb:function(t,e,r){t.exports=r("8aae")},aae3:function(t,e,r){var n=r("d3f4"),o=r("2d95"),c=r("2b4c")("match");t.exports=function(t){var e;return n(t)&&(void 0!==(e=t[c])?!!e:"RegExp"==o(t))}},b0c5:function(t,e,r){"use strict";var n=r("520a");r("5ca1")({target:"RegExp",proto:!0,forced:n!==/./.exec},{exec:n})},bf0b:function(t,e,r){var n=r("355d"),o=r("aebd"),c=r("36c3"),a=r("1bc3"),i=r("07e3"),u=r("794b"),s=Object.getOwnPropertyDescriptor;e.f=r("8e60")?s:function(t,e){if(t=c(t),e=a(e,!0),u)try{return s(t,e)}catch(r){}if(i(t,e))return o(!n.f.call(t,e),t[e])}},bf90:function(t,e,r){var n=r("36c3"),o=r("bf0b").f;r("ce7e")("getOwnPropertyDescriptor",function(){return function(t,e){return o(n(t),e)}})},ccb9:function(t,e,r){e.f=r("5168")},ce7e:function(t,e,r){var n=r("63b6"),o=r("584a"),c=r("294c");t.exports=function(t,e){var r=(o.Object||{})[t]||Object[t],a={};a[t]=e(r),n(n.S+n.F*c(function(){r(1)}),"Object",a)}},cebc:function(t,e,r){"use strict";var n=r("268f"),o=r.n(n),c=r("e265"),a=r.n(c),i=r("a4bb"),u=r.n(i),s=r("85f2"),f=r.n(s);function l(t,e,r){return e in t?f()(t,e,{value:r,enumerable:!0,configurable:!0,writable:!0}):t[e]=r,t}function p(t){for(var e=1;e<arguments.length;e++){var r=null!=arguments[e]?arguments[e]:{},n=u()(r);"function"===typeof a.a&&(n=n.concat(a()(r).filter(function(t){return o()(r,t).enumerable}))),n.forEach(function(e){l(t,e,r[e])})}return t}r.d(e,"a",function(){return p})},e265:function(t,e,r){t.exports=r("ed33")},ebfd:function(t,e,r){var n=r("62a0")("meta"),o=r("f772"),c=r("07e3"),a=r("d9f6").f,i=0,u=Object.isExtensible||function(){return!0},s=!r("294c")(function(){return u(Object.preventExtensions({}))}),f=function(t){a(t,n,{value:{i:"O"+ ++i,w:{}}})},l=function(t,e){if(!o(t))return"symbol"==typeof t?t:("string"==typeof t?"S":"P")+t;if(!c(t,n)){if(!u(t))return"F";if(!e)return"E";f(t)}return t[n].i},p=function(t,e){if(!c(t,n)){if(!u(t))return!0;if(!e)return!1;f(t)}return t[n].w},b=function(t){return s&&d.NEED&&u(t)&&!c(t,n)&&f(t),t},d=t.exports={KEY:n,NEED:!1,fastKey:l,getWeak:p,onFreeze:b}},ed33:function(t,e,r){r("014b"),t.exports=r("584a").Object.getOwnPropertySymbols},fde4:function(t,e,r){r("bf90");var n=r("584a").Object;t.exports=function(t,e){return n.getOwnPropertyDescriptor(t,e)}}}]);
//# sourceMappingURL=projectstack.2b7185e0.js.map