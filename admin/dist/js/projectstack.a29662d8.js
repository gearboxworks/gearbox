(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projectstack"],{"014b":function(t,e,r){"use strict";var n=r("e53d"),o=r("07e3"),c=r("8e60"),i=r("63b6"),a=r("9138"),s=r("ebfd").KEY,f=r("294c"),u=r("dbdb"),p=r("45f2"),l=r("62a0"),b=r("5168"),v=r("ccb9"),g=r("6718"),d=r("47ee"),m=r("9003"),h=r("e4ae"),y=r("f772"),j=r("36c3"),x=r("1bc3"),O=r("aebd"),S=r("a159"),k=r("0395"),w=r("bf0b"),_=r("d9f6"),P=r("c3a1"),N=w.f,E=_.f,I=k.f,F=n.Symbol,A=n.JSON,C=A&&A.stringify,T="prototype",D=b("_hidden"),B=b("toPrimitive"),R={}.propertyIsEnumerable,H=u("symbol-registry"),M=u("symbols"),$=u("op-symbols"),G=Object[T],J="function"==typeof F,q=n.QObject,L=!q||!q[T]||!q[T].findChild,U=c&&f(function(){return 7!=S(E({},"a",{get:function(){return E(this,"a",{value:7}).a}})).a})?function(t,e,r){var n=N(G,e);n&&delete G[e],E(t,e,r),n&&t!==G&&E(G,e,n)}:E,V=function(t){var e=M[t]=S(F[T]);return e._k=t,e},Y=J&&"symbol"==typeof F.iterator?function(t){return"symbol"==typeof t}:function(t){return t instanceof F},K=function(t,e,r){return t===G&&K($,e,r),h(t),e=x(e,!0),h(r),o(M,e)?(r.enumerable?(o(t,D)&&t[D][e]&&(t[D][e]=!1),r=S(r,{enumerable:O(0,!1)})):(o(t,D)||E(t,D,O(1,{})),t[D][e]=!0),U(t,e,r)):E(t,e,r)},W=function(t,e){h(t);var r,n=d(e=j(e)),o=0,c=n.length;while(c>o)K(t,r=n[o++],e[r]);return t},z=function(t,e){return void 0===e?S(t):W(S(t),e)},X=function(t){var e=R.call(this,t=x(t,!0));return!(this===G&&o(M,t)&&!o($,t))&&(!(e||!o(this,t)||!o(M,t)||o(this,D)&&this[D][t])||e)},Q=function(t,e){if(t=j(t),e=x(e,!0),t!==G||!o(M,e)||o($,e)){var r=N(t,e);return!r||!o(M,e)||o(t,D)&&t[D][e]||(r.enumerable=!0),r}},Z=function(t){var e,r=I(j(t)),n=[],c=0;while(r.length>c)o(M,e=r[c++])||e==D||e==s||n.push(e);return n},tt=function(t){var e,r=t===G,n=I(r?$:j(t)),c=[],i=0;while(n.length>i)!o(M,e=n[i++])||r&&!o(G,e)||c.push(M[e]);return c};J||(F=function(){if(this instanceof F)throw TypeError("Symbol is not a constructor!");var t=l(arguments.length>0?arguments[0]:void 0),e=function(r){this===G&&e.call($,r),o(this,D)&&o(this[D],t)&&(this[D][t]=!1),U(this,t,O(1,r))};return c&&L&&U(G,t,{configurable:!0,set:e}),V(t)},a(F[T],"toString",function(){return this._k}),w.f=Q,_.f=K,r("6abf").f=k.f=Z,r("355d").f=X,r("9aa9").f=tt,c&&!r("b8e3")&&a(G,"propertyIsEnumerable",X,!0),v.f=function(t){return V(b(t))}),i(i.G+i.W+i.F*!J,{Symbol:F});for(var et="hasInstance,isConcatSpreadable,iterator,match,replace,search,species,split,toPrimitive,toStringTag,unscopables".split(","),rt=0;et.length>rt;)b(et[rt++]);for(var nt=P(b.store),ot=0;nt.length>ot;)g(nt[ot++]);i(i.S+i.F*!J,"Symbol",{for:function(t){return o(H,t+="")?H[t]:H[t]=F(t)},keyFor:function(t){if(!Y(t))throw TypeError(t+" is not a symbol!");for(var e in H)if(H[e]===t)return e},useSetter:function(){L=!0},useSimple:function(){L=!1}}),i(i.S+i.F*!J,"Object",{create:z,defineProperty:K,defineProperties:W,getOwnPropertyDescriptor:Q,getOwnPropertyNames:Z,getOwnPropertySymbols:tt}),A&&i(i.S+i.F*(!J||f(function(){var t=F();return"[null]"!=C([t])||"{}"!=C({a:t})||"{}"!=C(Object(t))})),"JSON",{stringify:function(t){var e,r,n=[t],o=1;while(arguments.length>o)n.push(arguments[o++]);if(r=e=n[1],(y(e)||void 0!==t)&&!Y(t))return m(e)||(e=function(t,e){if("function"==typeof r&&(e=r.call(this,t,e)),!Y(e))return e}),n[1]=e,C.apply(A,n)}}),F[T][B]||r("35e8")(F[T],B,F[T].valueOf),p(F,"Symbol"),p(Math,"Math",!0),p(n.JSON,"JSON",!0)},"0395":function(t,e,r){var n=r("36c3"),o=r("6abf").f,c={}.toString,i="object"==typeof window&&window&&Object.getOwnPropertyNames?Object.getOwnPropertyNames(window):[],a=function(t){try{return o(t)}catch(e){return i.slice()}};t.exports.f=function(t){return i&&"[object Window]"==c.call(t)?a(t):o(n(t))}},"055b":function(t,e,r){t.exports=r.p+"img/python.51c2eab2.svg"},"090e":function(t,e,r){t.exports=r.p+"img/nodejs.94cafb0d.svg"},"11e9":function(t,e,r){var n=r("52a7"),o=r("4630"),c=r("6821"),i=r("6a99"),a=r("69a8"),s=r("c69a"),f=Object.getOwnPropertyDescriptor;e.f=r("9e1e")?f:function(t,e){if(t=c(t),e=i(e,!0),s)try{return f(t,e)}catch(r){}if(a(t,e))return o(!n.f.call(t,e),t[e])}},1540:function(t,e,r){t.exports=r.p+"img/php.fa78b345.svg"},"1b4d":function(t,e,r){t.exports=r.p+"img/flask.318d58cb.svg"},"268f":function(t,e,r){t.exports=r("fde4")},"319f":function(t,e,r){t.exports=r.p+"img/rails.2db29782.svg"},"31e8":function(t,e,r){var n={"./angular.svg":"a230","./apache.svg":"b77f","./codeigniter.svg":"7939","./django.svg":"c6da","./drupal.svg":"a88c","./elasticsearch.svg":"81bb","./flask.svg":"1b4d","./joomla.svg":"5390","./laravel.svg":"41c8","./logo.svg":"9b19","./mariadb.svg":"613e","./memcached.svg":"a0ba","./mysql.svg":"6c4c","./nginx.svg":"c502","./nodejs.svg":"090e","./perl.svg":"c44f","./php.svg":"1540","./python.svg":"055b","./rails.svg":"319f","./react.svg":"b2e9","./redis.svg":"8bcb","./ruby.svg":"9401","./wordpress.svg":"ee3c"};function o(t){var e=c(t);return r(e)}function c(t){var e=n[t];if(!(e+1)){var r=new Error("Cannot find module '"+t+"'");throw r.code="MODULE_NOT_FOUND",r}return e}o.keys=function(){return Object.keys(n)},o.resolve=c,t.exports=o,o.id="31e8"},"32a6":function(t,e,r){var n=r("241e"),o=r("c3a1");r("ce7e")("keys",function(){return function(t){return o(n(t))}})},3983:function(t,e,r){"use strict";r.r(e);var n=function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"project-stack-list",attrs:{role:"tablist",id:t.projectBase+"stack"}},t._l(t.groupProjectStacks(t.projectStack),function(e,o,c){return n("div",{key:o,staticClass:"project-stack"},[n("h2",{staticClass:"stack-title"},[t._v(t._s(o.replace("gearbox.works/","")))]),n("b-button",{staticClass:"js-remove-stack",attrs:{tabindex:100*t.projectIndex+10*c,size:"sm",variant:"outline-secondary","aria-label":"Remove this stack from project",title:"Remove this stack from project"},on:{click:function(e){return e.preventDefault(),t.removeProjectStack(o)}}},[t._v("×")]),n("ul",{staticClass:"service-list"},t._l(t.stackServices(o),function(i,a,s){return n("li",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],key:t.projectBase+t.escAttr(a),staticClass:"service-item",attrs:{id:t.projectBase+t.escAttr(a),tabindex:100*t.projectIndex+10*c+s+1,title:e[a]?e[a].service_id.replace("gearboxworks/",""):null}},[n("img",{staticClass:"service-program",attrs:{src:r("31e8")("./"+e[a].program+".svg")}}),n("h6",{staticClass:"service-role"},[t._v(t._s(t.stackRoles(o)[a].label))]),n("b-popover",{ref:t.projectBase+t.escAttr(a)+"-popover",refInFor:!0,attrs:{target:t.projectBase+t.escAttr(a),container:t.projectHostname+"stack",triggers:"focus",placement:"bottom"}},[n("template",{slot:"title"},[n("b-button",{staticClass:"close",attrs:{"aria-label":"Close"},on:{click:function(e){t.onClosePopoverFor(t.projectBase+t.escAttr(a))}}},[n("span",{staticClass:"d-inline-block",attrs:{"aria-hidden":"true"}},[t._v("×")])]),t._v("\n            "+t._s(t.stackRoles(o)[a].program)+"\n          ")],1),n("div",[n("b-form-select",{attrs:{id:t.projectBase+t.escAttr(a)+"_input",value:e[a]?e[a].service_id:null,tabindex:100*t.projectIndex+10*c+s+9},on:{change:function(e){return t.changeProjectService(a,e)}}},[n("option",{attrs:{disabled:"",value:""}},[t._v("Select service...")]),t._l(t.optionGroups(i.options),function(e,r){return n("optgroup",{key:r,attrs:{label:r}},t._l(e,function(e){return n("option",{key:e,domProps:{value:i.org+"/"+e}},[t._v(t._s(e))])}),0)})],2)],1)],2)],1)}),0)],1)}),0)},o=[],c=(r("28a5"),r("a481"),r("cebc")),i=(r("c5f6"),r("2f62")),a={name:"ProjectStack",props:{projectHostname:{type:String,required:!0},projectStack:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return{}},components:{},computed:Object(c["a"])({},Object(i["b"])(["groupProjectStacks","stackRoles","stackServices"]),{projectBase:function(){return this.escAttr(this.projectHostname)+"-"},project:function(){return this.$store.getters.projectBy("hostname",this.projectHostname)}}),methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},mapOptions:function(t){var e=[];for(var r in t)e.push({value:r,text:t[r]});return e},optionGroups:function(t){var e={};for(var r in t){var n=t[r].split(":")[0];"undefined"===typeof e[n]&&(e[n]={}),e[n][r]=t[r]}return e},changeProjectService:function(t,e){this.$store.dispatch("changeProjectService",{projectHostname:this.projectHostname,serviceName:t,serviceId:e})},onClosePopoverFor:function(t){console.log("onClosePopoverFor",t),this.$root.$emit("bv::popover::hide",t)},removeProjectStack:function(t){this.$store.dispatch("removeProjectStack",{projectHostname:this.projectHostname,stackName:t})}}},s=a,f=(r("6b94"),r("2877")),u=Object(f["a"])(s,n,o,!1,null,"459a1d8a",null);e["default"]=u.exports},"41c8":function(t,e,r){t.exports=r.p+"img/laravel.1766a461.svg"},"454f":function(t,e,r){r("46a7");var n=r("584a").Object;t.exports=function(t,e,r){return n.defineProperty(t,e,r)}},"46a7":function(t,e,r){var n=r("63b6");n(n.S+n.F*!r("8e60"),"Object",{defineProperty:r("d9f6").f})},"47ee":function(t,e,r){var n=r("c3a1"),o=r("9aa9"),c=r("355d");t.exports=function(t){var e=n(t),r=o.f;if(r){var i,a=r(t),s=c.f,f=0;while(a.length>f)s.call(t,i=a[f++])&&e.push(i)}return e}},5390:function(t,e,r){t.exports=r.p+"img/joomla.d8aa2e45.svg"},"5dbc":function(t,e,r){var n=r("d3f4"),o=r("8b97").set;t.exports=function(t,e,r){var c,i=e.constructor;return i!==r&&"function"==typeof i&&(c=i.prototype)!==r.prototype&&n(c)&&o&&o(t,c),t}},"613e":function(t,e,r){t.exports=r.p+"img/mariadb.e16110bc.svg"},6718:function(t,e,r){var n=r("e53d"),o=r("584a"),c=r("b8e3"),i=r("ccb9"),a=r("d9f6").f;t.exports=function(t){var e=o.Symbol||(o.Symbol=c?{}:n.Symbol||{});"_"==t.charAt(0)||t in e||a(e,t,{value:i.f(t)})}},"6abf":function(t,e,r){var n=r("e6f3"),o=r("1691").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return n(t,o)}},"6b94":function(t,e,r){"use strict";var n=r("eb9b"),o=r.n(n);o.a},"6c4c":function(t,e,r){t.exports=r.p+"img/mysql.dd2a5a35.svg"},7939:function(t,e,r){t.exports=r.p+"img/codeigniter.434bf735.svg"},"81bb":function(t,e,r){t.exports=r.p+"img/elasticsearch.3ecfa530.svg"},"85f2":function(t,e,r){t.exports=r("454f")},"8aae":function(t,e,r){r("32a6"),t.exports=r("584a").Object.keys},"8b97":function(t,e,r){var n=r("d3f4"),o=r("cb7c"),c=function(t,e){if(o(t),!n(e)&&null!==e)throw TypeError(e+": can't set as prototype!")};t.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(t,e,n){try{n=r("9b43")(Function.call,r("11e9").f(Object.prototype,"__proto__").set,2),n(t,[]),e=!(t instanceof Array)}catch(o){e=!0}return function(t,r){return c(t,r),e?t.__proto__=r:n(t,r),t}}({},!1):void 0),check:c}},"8bcb":function(t,e,r){t.exports=r.p+"img/redis.3c39fafe.svg"},9093:function(t,e,r){var n=r("ce10"),o=r("e11e").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return n(t,o)}},9401:function(t,e,r){t.exports=r.p+"img/ruby.514befa7.svg"},"9aa9":function(t,e){e.f=Object.getOwnPropertySymbols},"9b19":function(t,e,r){t.exports=r.p+"img/logo.63a7d78d.svg"},a0ba:function(t,e,r){t.exports=r.p+"img/memcached.2bcccabf.svg"},a230:function(t,e,r){t.exports=r.p+"img/angular.e224f5ed.svg"},a4bb:function(t,e,r){t.exports=r("8aae")},a88c:function(t,e,r){t.exports=r.p+"img/drupal.66089b06.svg"},aa77:function(t,e,r){var n=r("5ca1"),o=r("be13"),c=r("79e5"),i=r("fdef"),a="["+i+"]",s="​",f=RegExp("^"+a+a+"*"),u=RegExp(a+a+"*$"),p=function(t,e,r){var o={},a=c(function(){return!!i[t]()||s[t]()!=s}),f=o[t]=a?e(l):i[t];r&&(o[r]=f),n(n.P+n.F*a,"String",o)},l=p.trim=function(t,e){return t=String(o(t)),1&e&&(t=t.replace(f,"")),2&e&&(t=t.replace(u,"")),t};t.exports=p},b2e9:function(t,e,r){t.exports=r.p+"img/react.9a28da9f.svg"},b77f:function(t,e,r){t.exports=r.p+"img/apache.12c49354.svg"},bf0b:function(t,e,r){var n=r("355d"),o=r("aebd"),c=r("36c3"),i=r("1bc3"),a=r("07e3"),s=r("794b"),f=Object.getOwnPropertyDescriptor;e.f=r("8e60")?f:function(t,e){if(t=c(t),e=i(e,!0),s)try{return f(t,e)}catch(r){}if(a(t,e))return o(!n.f.call(t,e),t[e])}},bf90:function(t,e,r){var n=r("36c3"),o=r("bf0b").f;r("ce7e")("getOwnPropertyDescriptor",function(){return function(t,e){return o(n(t),e)}})},c44f:function(t,e,r){t.exports=r.p+"img/perl.a025edb4.svg"},c502:function(t,e,r){t.exports=r.p+"img/nginx.eae76401.svg"},c5f6:function(t,e,r){"use strict";var n=r("7726"),o=r("69a8"),c=r("2d95"),i=r("5dbc"),a=r("6a99"),s=r("79e5"),f=r("9093").f,u=r("11e9").f,p=r("86cc").f,l=r("aa77").trim,b="Number",v=n[b],g=v,d=v.prototype,m=c(r("2aeb")(d))==b,h="trim"in String.prototype,y=function(t){var e=a(t,!1);if("string"==typeof e&&e.length>2){e=h?e.trim():l(e,3);var r,n,o,c=e.charCodeAt(0);if(43===c||45===c){if(r=e.charCodeAt(2),88===r||120===r)return NaN}else if(48===c){switch(e.charCodeAt(1)){case 66:case 98:n=2,o=49;break;case 79:case 111:n=8,o=55;break;default:return+e}for(var i,s=e.slice(2),f=0,u=s.length;f<u;f++)if(i=s.charCodeAt(f),i<48||i>o)return NaN;return parseInt(s,n)}}return+e};if(!v(" 0o1")||!v("0b1")||v("+0x1")){v=function(t){var e=arguments.length<1?0:t,r=this;return r instanceof v&&(m?s(function(){d.valueOf.call(r)}):c(r)!=b)?i(new g(y(e)),r,v):y(e)};for(var j,x=r("9e1e")?f(g):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),O=0;x.length>O;O++)o(g,j=x[O])&&!o(v,j)&&p(v,j,u(g,j));v.prototype=d,d.constructor=v,r("2aba")(n,b,v)}},c6da:function(t,e,r){t.exports=r.p+"img/django.28fe09a0.svg"},ccb9:function(t,e,r){e.f=r("5168")},ce7e:function(t,e,r){var n=r("63b6"),o=r("584a"),c=r("294c");t.exports=function(t,e){var r=(o.Object||{})[t]||Object[t],i={};i[t]=e(r),n(n.S+n.F*c(function(){r(1)}),"Object",i)}},cebc:function(t,e,r){"use strict";var n=r("268f"),o=r.n(n),c=r("e265"),i=r.n(c),a=r("a4bb"),s=r.n(a),f=r("85f2"),u=r.n(f);function p(t,e,r){return e in t?u()(t,e,{value:r,enumerable:!0,configurable:!0,writable:!0}):t[e]=r,t}function l(t){for(var e=1;e<arguments.length;e++){var r=null!=arguments[e]?arguments[e]:{},n=s()(r);"function"===typeof i.a&&(n=n.concat(i()(r).filter(function(t){return o()(r,t).enumerable}))),n.forEach(function(e){p(t,e,r[e])})}return t}r.d(e,"a",function(){return l})},e265:function(t,e,r){t.exports=r("ed33")},eb9b:function(t,e,r){},ebfd:function(t,e,r){var n=r("62a0")("meta"),o=r("f772"),c=r("07e3"),i=r("d9f6").f,a=0,s=Object.isExtensible||function(){return!0},f=!r("294c")(function(){return s(Object.preventExtensions({}))}),u=function(t){i(t,n,{value:{i:"O"+ ++a,w:{}}})},p=function(t,e){if(!o(t))return"symbol"==typeof t?t:("string"==typeof t?"S":"P")+t;if(!c(t,n)){if(!s(t))return"F";if(!e)return"E";u(t)}return t[n].i},l=function(t,e){if(!c(t,n)){if(!s(t))return!0;if(!e)return!1;u(t)}return t[n].w},b=function(t){return f&&v.NEED&&s(t)&&!c(t,n)&&u(t),t},v=t.exports={KEY:n,NEED:!1,fastKey:p,getWeak:l,onFreeze:b}},ed33:function(t,e,r){r("014b"),t.exports=r("584a").Object.getOwnPropertySymbols},ee3c:function(t,e,r){t.exports=r.p+"img/wordpress.b08e20e3.svg"},fde4:function(t,e,r){r("bf90");var n=r("584a").Object;t.exports=function(t,e){return n.getOwnPropertyDescriptor(t,e)}},fdef:function(t,e){t.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);
//# sourceMappingURL=projectstack.a29662d8.js.map