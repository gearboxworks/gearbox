(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projectstack"],{"055b":function(e,t,r){e.exports=r.p+"img/python.51c2eab2.svg"},"0839":function(e,t,r){"use strict";var c=r("2fa9"),s=r.n(c);s.a},"090e":function(e,t,r){e.exports=r.p+"img/nodejs.94cafb0d.svg"},"11e9":function(e,t,r){var c=r("52a7"),s=r("4630"),o=r("6821"),i=r("6a99"),a=r("69a8"),n=r("c69a"),p=Object.getOwnPropertyDescriptor;t.f=r("9e1e")?p:function(e,t){if(e=o(e),t=i(t,!0),n)try{return p(e,t)}catch(r){}if(a(e,t))return s(!c.f.call(e,t),e[t])}},1540:function(e,t,r){e.exports=r.p+"img/php.fa78b345.svg"},"1b4d":function(e,t,r){e.exports=r.p+"img/flask.318d58cb.svg"},"2fa9":function(e,t,r){},"319f":function(e,t,r){e.exports=r.p+"img/rails.2db29782.svg"},"31e8":function(e,t,r){var c={"./angular.svg":"a230","./apache.svg":"b77f","./codeigniter.svg":"7939","./django.svg":"c6da","./drupal.svg":"a88c","./elasticsearch.svg":"81bb","./flask.svg":"1b4d","./joomla.svg":"5390","./laravel.svg":"41c8","./logo.svg":"9b19","./mariadb.svg":"613e","./memcached.svg":"a0ba","./mysql.svg":"6c4c","./nginx.svg":"c502","./nodejs.svg":"090e","./perl.svg":"c44f","./php.svg":"1540","./python.svg":"055b","./rails.svg":"319f","./react.svg":"b2e9","./redis.svg":"8bcb","./ruby.svg":"9401","./wordpress.svg":"ee3c"};function s(e){var t=o(e);return r(t)}function o(e){var t=c[e];if(!(t+1)){var r=new Error("Cannot find module '"+e+"'");throw r.code="MODULE_NOT_FOUND",r}return t}s.keys=function(){return Object.keys(c)},s.resolve=o,e.exports=s,s.id="31e8"},3983:function(e,t,r){"use strict";r.r(t);var c=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"project-stack-list",attrs:{role:"tablist",id:e.projectBase+"stack"}},e._l(e.groupProjectServicesByStack(e.stack),function(t,c,s){return r("div",{key:c,staticClass:"project-stack"},[r("h2",{staticClass:"stack-title"},[e._v(e._s(c.replace("gearbox.works/","")))]),r("b-button",{staticClass:"js-remove-stack",attrs:{tabindex:100*e.projectIndex+10*s,size:"sm",variant:"outline-secondary","aria-label":"Remove this stack from project",title:"Remove this stack from project"},on:{click:function(t){return t.preventDefault(),e.removeProjectStack(c)}}},[e._v("×")]),r("ul",{staticClass:"service-list"},e._l(t,function(t,c,o){return r("li",{key:e.projectBase+t.id,staticClass:"service-item",attrs:{id:e.projectBase+c,tabindex:100*e.projectIndex+10*s+o+1}},[r("project-service",{attrs:{projectId:e.project.id,service:t,projectIndex:e.projectIndex,stackIndex:s,serviceIndex:o}})],1)}),0)],1)}),0)},s=[],o=(r("ac6a"),r("a481"),r("cebc")),i=(r("c5f6"),function(){var e=this,t=e.$createElement,c=e._self._c||t;return c("div",{staticClass:"project-service"},[e.program?c("img",{staticClass:"service-program",attrs:{src:r("31e8")("./"+e.program+".svg")}}):c("font-awesome-icon",{attrs:{icon:["fa","expand"]}}),c("h6",{staticClass:"service-role"},[e._v(e._s(e.role))]),c("b-tooltip",{key:e.id,attrs:{triggers:"hover",target:e.projectBase+e.role,title:e.programTooltip(e.program+" "+e.version)}}),c("b-popover",{ref:e.projectBase+e.role+"-popover",attrs:{target:e.projectBase+e.role,container:e.projectId+"stack",triggers:"focus",placement:"bottom"}},[c("template",{slot:"title"},[c("b-button",{staticClass:"close",attrs:{"aria-label":"Close"},on:{click:function(t){return e.onClosePopoverFor(e.projectBase+e.role)}}},[c("span",{staticClass:"d-inline-block",attrs:{"aria-hidden":"true"}},[e._v("×")])]),e._v("\n      Change service\n    ")],1),c("div",[c("label",{attrs:{for:e.projectBase+e.role+"_input"}},[e._v("\n        "+e._s(e.gear.attributes.role)+":\n      ")]),c("b-form-select",{attrs:{id:e.projectBase+e.role+"_input",value:e.id,tabindex:100*e.projectIndex+10*e.stackIndex+e.serviceIndex+9},on:{change:function(t){return e.changeProjectService(e.role,t)}}},[e.stackRoleDefaultService(e.role)?e._e():c("option",{attrs:{value:""}},[e._v("Do not run this service")]),c("option",{attrs:{disabled:"",value:""}},[e._v("Select service...")]),e._l(e.serviceGroups(e.stackRoleServices(e.role)),function(t,r){return c("optgroup",{key:r,attrs:{label:r}},e._l(t,function(t){return c("option",{key:t,domProps:{value:t}},[e._v(e._s(t.replace("gearboxworks/","")))])}),0)})],2)],1)],2)],1)}),a=[],n=(r("28a5"),r("7514"),r("2f62")),p={name:"ProjectService",props:{projectId:{type:String,required:!0},service:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackIndex:{type:Number,required:!0},serviceIndex:{type:Number,required:!0}},data:function(){return Object(o["a"])({id:this.service.id},this.service.attributes)},computed:Object(o["a"])({},Object(n["b"])(["gearspecBy","stackBy"]),{projectBase:function(){return this.escAttr(this.projectId)+"-"},gear:function(){return this.gearspecBy("id",this.gearspec_id)},role:function(){var e=this.gear;return e?e.attributes.role:""},stack:function(){var e=this.gear;return e?this.stackBy("id",e.attributes.stack_id):null}}),methods:{escAttr:function(e){return e.replace(/\//g,"-").replace(/\./g,"-")},serviceVersion:function(e){var t="";return e.major&&(t+=e.major,e.minor&&(t+="."+e.minor,e.patch&&(t+="."+e.patch))),t},stackRoleObject:function(e){var t=this.gear,r=t?this.stackBy("id",t.attributes.stack_id):null,c=r?r.attributes.members:[];return c.find(function(t){return t.role===e})},stackRoleDefaultService:function(e){var t=this.stackRoleObject(e);return t&&"undefined"!==typeof t.default_service?t.default_service:""},stackRoleServices:function(e){var t=this.stackRoleObject(e);return t?t.services:[]},stackRoleSmartDefaultService:function(e){var t=this.stackRoleObject(e),r=this.stackRoleDefaultService(e),c=-1,s=-1;if(r)for(var o=t.services.length;o--;)if(-1!==t.services[o].indexOf(r)&&(-1===c&&(c=o),t.services[o]===r)){s=o;break}var i=-1!==c?t.services[-1!==s?s:c]:"";return i},serviceGroups:function(e){var t={};for(var r in e){var c=e[r].split(":")[0].replace("gearboxworks/","");"undefined"===typeof t[c]&&(t[c]={}),t[c][r]=e[r]}return t},removeProjectStack:function(){this.$store.dispatch("removeProjectStack",{projectId:this.projectId,stackName:this.stackName})},programTooltip:function(e){return e||"Not selected"},changeProjectService:function(e,t){this.$store.dispatch("changeProjectService",{id:this.projectId,serviceName:e,serviceId:t}),this.onClosePopoverFor(this.projectBase+this.escAttr(e))},onClosePopoverFor:function(e){this.$root.$emit("bv::hide::popover",e)}}},u=p,f=(r("0839"),r("2877")),l=Object(f["a"])(u,i,a,!1,null,"35c8d391",null),v=l.exports,d={name:"ProjectStack",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0}},data:function(){return Object(o["a"])({id:this.project.id},this.project.attributes)},components:{ProjectService:v},computed:Object(o["a"])({},Object(n["b"])(["serviceBy","gearspecBy"]),{projectBase:function(){return this.escAttr(this.id)+"-"}}),methods:{escAttr:function(e){return e.replace(/\//g,"-").replace(/\./g,"-")},groupProjectServicesByStack:function(e){var t=this,r={};return e.forEach(function(e,c){var s=t.gearspecBy("id",e.gearspec_id),o=t.serviceBy("id",e.service_id);s&&o&&("undefined"===typeof r[s.attributes.stack_id]&&(r[s.attributes.stack_id]={}),r[s.attributes.stack_id][s.attributes.role]=o)}),r},removeProjectStack:function(e){this.$store.dispatch("removeProjectStack",{projectHostname:this.projectHostname,stackName:e})}}},g=d,b=(r("f1e3"),Object(f["a"])(g,c,s,!1,null,"b2c89160",null));t["default"]=b.exports},"3f11":function(e,t,r){},"41c8":function(e,t,r){e.exports=r.p+"img/laravel.1766a461.svg"},5390:function(e,t,r){e.exports=r.p+"img/joomla.d8aa2e45.svg"},"5dbc":function(e,t,r){var c=r("d3f4"),s=r("8b97").set;e.exports=function(e,t,r){var o,i=t.constructor;return i!==r&&"function"==typeof i&&(o=i.prototype)!==r.prototype&&c(o)&&s&&s(e,o),e}},"613e":function(e,t,r){e.exports=r.p+"img/mariadb.e16110bc.svg"},"6c4c":function(e,t,r){e.exports=r.p+"img/mysql.dd2a5a35.svg"},7939:function(e,t,r){e.exports=r.p+"img/codeigniter.434bf735.svg"},"81bb":function(e,t,r){e.exports=r.p+"img/elasticsearch.3ecfa530.svg"},"8b97":function(e,t,r){var c=r("d3f4"),s=r("cb7c"),o=function(e,t){if(s(e),!c(t)&&null!==t)throw TypeError(t+": can't set as prototype!")};e.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(e,t,c){try{c=r("9b43")(Function.call,r("11e9").f(Object.prototype,"__proto__").set,2),c(e,[]),t=!(e instanceof Array)}catch(s){t=!0}return function(e,r){return o(e,r),t?e.__proto__=r:c(e,r),e}}({},!1):void 0),check:o}},"8bcb":function(e,t,r){e.exports=r.p+"img/redis.3c39fafe.svg"},9093:function(e,t,r){var c=r("ce10"),s=r("e11e").concat("length","prototype");t.f=Object.getOwnPropertyNames||function(e){return c(e,s)}},9401:function(e,t,r){e.exports=r.p+"img/ruby.514befa7.svg"},"9b19":function(e,t,r){e.exports=r.p+"img/logo.63a7d78d.svg"},a0ba:function(e,t,r){e.exports=r.p+"img/memcached.2bcccabf.svg"},a230:function(e,t,r){e.exports=r.p+"img/angular.e224f5ed.svg"},a88c:function(e,t,r){e.exports=r.p+"img/drupal.66089b06.svg"},aa77:function(e,t,r){var c=r("5ca1"),s=r("be13"),o=r("79e5"),i=r("fdef"),a="["+i+"]",n="​",p=RegExp("^"+a+a+"*"),u=RegExp(a+a+"*$"),f=function(e,t,r){var s={},a=o(function(){return!!i[e]()||n[e]()!=n}),p=s[e]=a?t(l):i[e];r&&(s[r]=p),c(c.P+c.F*a,"String",s)},l=f.trim=function(e,t){return e=String(s(e)),1&t&&(e=e.replace(p,"")),2&t&&(e=e.replace(u,"")),e};e.exports=f},ac6a:function(e,t,r){for(var c=r("cadf"),s=r("0d58"),o=r("2aba"),i=r("7726"),a=r("32e9"),n=r("84f2"),p=r("2b4c"),u=p("iterator"),f=p("toStringTag"),l=n.Array,v={CSSRuleList:!0,CSSStyleDeclaration:!1,CSSValueList:!1,ClientRectList:!1,DOMRectList:!1,DOMStringList:!1,DOMTokenList:!0,DataTransferItemList:!1,FileList:!1,HTMLAllCollection:!1,HTMLCollection:!1,HTMLFormElement:!1,HTMLSelectElement:!1,MediaList:!0,MimeTypeArray:!1,NamedNodeMap:!1,NodeList:!0,PaintRequestList:!1,Plugin:!1,PluginArray:!1,SVGLengthList:!1,SVGNumberList:!1,SVGPathSegList:!1,SVGPointList:!1,SVGStringList:!1,SVGTransformList:!1,SourceBufferList:!1,StyleSheetList:!0,TextTrackCueList:!1,TextTrackList:!1,TouchList:!1},d=s(v),g=0;g<d.length;g++){var b,m=d[g],h=v[m],j=i[m],k=j&&j.prototype;if(k&&(k[u]||a(k,u,l),k[f]||a(k,f,m),n[m]=l,h))for(b in c)k[b]||o(k,b,c[b],!0)}},b2e9:function(e,t,r){e.exports=r.p+"img/react.9a28da9f.svg"},b77f:function(e,t,r){e.exports=r.p+"img/apache.12c49354.svg"},c44f:function(e,t,r){e.exports=r.p+"img/perl.a025edb4.svg"},c502:function(e,t,r){e.exports=r.p+"img/nginx.eae76401.svg"},c5f6:function(e,t,r){"use strict";var c=r("7726"),s=r("69a8"),o=r("2d95"),i=r("5dbc"),a=r("6a99"),n=r("79e5"),p=r("9093").f,u=r("11e9").f,f=r("86cc").f,l=r("aa77").trim,v="Number",d=c[v],g=d,b=d.prototype,m=o(r("2aeb")(b))==v,h="trim"in String.prototype,j=function(e){var t=a(e,!1);if("string"==typeof t&&t.length>2){t=h?t.trim():l(t,3);var r,c,s,o=t.charCodeAt(0);if(43===o||45===o){if(r=t.charCodeAt(2),88===r||120===r)return NaN}else if(48===o){switch(t.charCodeAt(1)){case 66:case 98:c=2,s=49;break;case 79:case 111:c=8,s=55;break;default:return+t}for(var i,n=t.slice(2),p=0,u=n.length;p<u;p++)if(i=n.charCodeAt(p),i<48||i>s)return NaN;return parseInt(n,c)}}return+t};if(!d(" 0o1")||!d("0b1")||d("+0x1")){d=function(e){var t=arguments.length<1?0:e,r=this;return r instanceof d&&(m?n(function(){b.valueOf.call(r)}):o(r)!=v)?i(new g(j(t)),r,d):j(t)};for(var k,y=r("9e1e")?p(g):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),x=0;y.length>x;x++)s(g,k=y[x])&&!s(d,k)&&f(d,k,u(g,k));d.prototype=b,b.constructor=d,r("2aba")(c,v,d)}},c6da:function(e,t,r){e.exports=r.p+"img/django.28fe09a0.svg"},ee3c:function(e,t,r){e.exports=r.p+"img/wordpress.b08e20e3.svg"},f1e3:function(e,t,r){"use strict";var c=r("3f11"),s=r.n(c);s.a},fdef:function(e,t){e.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);
//# sourceMappingURL=projectstack.4fed23d4.js.map