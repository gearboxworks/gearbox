(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["preferences"],{"5dbc":function(t,e,r){var i=r("d3f4"),n=r("8b97").set;t.exports=function(t,e,r){var o,a=e.constructor;return a!==r&&"function"==typeof a&&(o=a.prototype)!==r.prototype&&i(o)&&n&&n(t,o),t}},"778f":function(t,e,r){"use strict";var i=r("79f9"),n=r.n(i);n.a},"79f9":function(t,e,r){},"8b97":function(t,e,r){var i=r("d3f4"),n=r("cb7c"),o=function(t,e){if(n(t),!i(e)&&null!==e)throw TypeError(e+": can't set as prototype!")};t.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(t,e,i){try{i=r("9b43")(Function.call,r("11e9").f(Object.prototype,"__proto__").set,2),i(t,[]),e=!(t instanceof Array)}catch(n){e=!0}return function(t,r){return o(t,r),e?t.__proto__=r:i(t,r),t}}({},!1):void 0),check:o}},9283:function(t,e,r){"use strict";var i=r("a240"),n=r.n(i);n.a},a240:function(t,e,r){},a4fc:function(t,e,r){},a55d:function(t,e,r){"use strict";r.r(e);var i=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("b-form",{class:{"form--preferences":!0,"is-loading":t.isLoading}},[r("h1",[t._v("Preferences")]),r("h2",[t._v("Project Directories")]),t.isLoading?r("div",{key:"basedirs-content",staticClass:"basedirs-loading"},[r("font-awesome-icon",{attrs:{icon:"circle-notch",spin:""}}),t._v("\n     "),r("span",[t._v("Loading directories...")])],1):r("div",{key:"basedirs-content",staticClass:"basedirs-wrap"},t._l(t.basedirs,function(t){return r("basedir-row-edit",{key:t.id,attrs:{basedir:t,"is-deletable":"default"!==t.id}})}),1),r("basedir-row-add",{attrs:{"tab-offset":3*t.basedirs.length}})],1)},n=[],o=(r("8e6e"),r("ac6a"),r("456d"),r("bd86")),a=r("2f62"),s=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("b-form-row",{class:{"form-row--basedir-edit":!0,"is-updating":t.isUpdating,"is-modified":t.isModified,"is-deleting":t.isDeleting}},[r("b-input-group",{class:{"input-group--basedir-edit":!0},attrs:{role:"tabpanel"}},[r("b-form-input",{staticClass:"basedir",attrs:{state:t.isValid,type:"text",required:"",placeholder:"Path",tabindex:"0"},on:{keyup:t.touch,change:t.touch},model:{value:t.currentValue,callback:function(e){t.currentValue=e},expression:"currentValue"}}),t.isModified?r("b-input-group-append",[r("b-button",{staticClass:"btn--update",attrs:{type:"submit.prevent",variant:t.isModified?"outline-info":"outline-secondary",disabled:!t.isModified,title:"Update directory reference",tabindex:0},on:{click:function(e){return e.preventDefault(),t.onUpdateBasedir(e)}}},[t.isUpdating?r("font-awesome-icon",{key:"status-icon",attrs:{spin:"",icon:["fa","circle-notch"]}}):r("font-awesome-icon",{key:"status-icon",attrs:{icon:["fa","check"]}})],1)],1):t._e()],1),r("b-button-group",{staticClass:"button-group--extras"},[r("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}},{name:"clipboard",rawName:"v-clipboard:copy",value:t.currentValue,expression:"currentValue",arg:"copy"},{name:"clipboard",rawName:"v-clipboard:success",value:t.onCopyToClipboard,expression:"onCopyToClipboard",arg:"success"}],staticClass:"btn--copy-dir",attrs:{variant:"outline-info",title:"Copy to clipboard"}},[r("font-awesome-icon",{attrs:{icon:["fa","clone"]}})],1),r("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"btn--open-dir",attrs:{variant:"outline-info",title:"Open in file manager"},on:{click:function(e){return e.preventDefault(),t.onOpenDirectory(e)}}},[r("font-awesome-icon",{attrs:{icon:["fa","folder-open"]}})],1)],1),r("b-button",{staticClass:"btn--delete",attrs:{type:"submit.prevent",variant:t.basedir.id===t.isDeletable?"outline-secondary":"outline-warning",disabled:!t.isDeletable,title:t.isDeletable?"Delete this directory":"Cannot delete the default directory",tabindex:0},on:{click:function(e){return e.preventDefault(),t.onDeleteBasedir(e)}}},[t.isDeleting?r("font-awesome-icon",{key:"status-icon",attrs:{spin:"",icon:["fa","circle-notch"]}}):r("font-awesome-icon",{key:"status-icon",attrs:{icon:["fa","trash-alt"]}})],1),t.errors?t._e():r("div",{class:{confirmation:!0,visible:t.notfound[t.currentValue]}},[t._v("\n    This dir does not exist. Would you like to create it?\n    "),r("a",{staticClass:"yes",attrs:{title:"Create directory"},on:{click:function(e){return e.stopPropagation(),t.updateDir(t.currentValue)}}},[t._v("Yes")]),t._v("\n    |\n    "),r("a",{staticClass:"no",attrs:{title:"Try a different dir"},on:{click:function(e){e.stopPropagation(),t.notfound[t.currentValue]=0}}},[t._v("No")])]),t.errors?r("div",{staticClass:"invalid-feedback d-block"},[t._v("\n    "+t._s(t.errors)+"\n  ")]):t._e()],1)},c=[];function d(t,e){var r=Object.keys(t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(t);e&&(i=i.filter(function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable})),r.push.apply(r,i)}return r}function u(t){for(var e=1;e<arguments.length;e++){var r=null!=arguments[e]?arguments[e]:{};e%2?d(r,!0).forEach(function(e){Object(o["a"])(t,e,r[e])}):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(r)):d(r).forEach(function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(r,e))})}return t}var l={name:"BasedirRowEdit",props:{basedir:{type:Object,required:!0},isDeletable:{type:Boolean,required:!0,default:!0}},data:function(){return{id:this.basedir.id,currentValue:this.basedir.attributes.basedir,errors:"",notfound:{},isModified:!1,isUpdating:!1,isDeleting:!1,alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"warning"}},computed:u({},Object(a["d"])(["basedirBy"]),{isValid:function(){return""===this.errors&&null}}),methods:u({},Object(a["c"])({getDirectory:"getDirectory",doUpdateBasedir:"basedirs/update",doDeleteBasedir:"basedirs/delete"}),{showAlert:function(t){"string"===typeof t?this.alertContent=t:(this.alertVariant=t.variant||this.alertVariant,this.alertDismissible=t.dismissible||this.alertDismissible,this.alertContent=t.content||this.alertContent),this.alertShow=!0},touch:function(){var t=this.currentValue;this.isModified=t&&t!==this.basedir.attributes.basedir,this.errors=""},onUpdateBasedir:function(){var t=this,e=this.currentValue;this.id&&e&&this.getDirectory({dir:e}).then(function(t){return t?t.data:null}).then(function(r){t.updateDir(e)}).catch(function(r){404===r.response.status?t.$set(t.notfound,e,1):t.$delete(t.notfound,e)})},updateDir:function(t){var e=this;if(this.isModified){var r={id:this.id,type:"basedirs",attributes:{basedir:t,nickname:this.id}};this.isUpdating=!0,this.doUpdateBasedir(r).then(function(t){e.errors="",e.isModified=!1,e.isUpdating=!1,e.$delete(e.notfound,e.currentValue)}).catch(function(t){e.isUpdating=!1,e.$nextTick(function(){this.errors=t.data.errors[0].title||t.statusText}),console.log("fail",t.response)})}},onDeleteBasedir:function(){var t=this;this.isDeleting=!0,this.doDeleteBasedir({id:this.id}).then(function(e){t.errors="",t.isModified=!1,t.isDeleting=!1}).catch(function(e){t.errors=e,t.isDeleting=!1})},onOpenDirectory:function(){console.log("TODO: call API method to open directory in file manager")}})},f=l,p=(r("9283"),r("2877")),b=Object(p["a"])(f,s,c,!1,null,"28d8562a",null),h=b.exports,y=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("b-form-row",{class:{"form-row--basedir-add":!0,"is-updating":t.isUpdating}},[r("b-input-group",{class:{"input-group--basedir-edit":!0},attrs:{role:"tabpanel"}},[r("b-form-input",{staticClass:"basedir",attrs:{type:"text",placeholder:"Input directory...",state:t.inputState("add"),tabindex:"0"},on:{keyup:function(e){return t.touch("add")},change:function(e){return t.touch("add")}},model:{value:t.currentValue,callback:function(e){t.currentValue=e},expression:"currentValue"}}),r("b-input-group-append",[r("b-button",{staticClass:"btn--add",attrs:{type:"submit.prevent",variant:t.touched["add"]?"outline-info":"outline-secondary",disabled:!t.touched["add"],title:t.currentValue?"Add directory":"First, input some directory!",tabindex:"0"},on:{click:function(e){return e.preventDefault(),t.onAddBasedir(e)}}},[r("font-awesome-icon",{attrs:{icon:["fa","plus"]}})],1)],1)],1),r("div",{class:{confirmation:!0,visible:t.notfound[t.currentValue]}},[t._v("\n    This dir does not exist. Would you like to create it?\n    "),r("a",{staticClass:"yes",attrs:{title:"Create directory"},on:{click:function(e){return e.stopPropagation(),t.createDir(t.currentValue)}}},[t._v("Yes")]),t._v("\n    |\n    "),r("a",{staticClass:"no",attrs:{title:"Try a different dir"},on:{click:function(e){e.stopPropagation(),t.notfound[t.currentValue]=0}}},[t._v("No")])]),r("div",{staticClass:"invalid-feedback d-block"},[t._v("\n    "+t._s(t.errors["add"]||" ")+"\n  ")])],1)},v=[];r("c5f6");function g(t,e){var r=Object.keys(t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(t);e&&(i=i.filter(function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable})),r.push.apply(r,i)}return r}function m(t){for(var e=1;e<arguments.length;e++){var r=null!=arguments[e]?arguments[e]:{};e%2?g(r,!0).forEach(function(e){Object(o["a"])(t,e,r[e])}):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(r)):g(r).forEach(function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(r,e))})}return t}var O={name:"BasedirRowEdit",props:{tabOffset:{type:Number,required:!0}},data:function(){return{isUpdating:!1,currentValue:"",errors:{},touched:{},notfound:{},alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"warning"}},computed:m({},Object(a["d"])(["basedirBy"])),methods:m({},Object(a["c"])({doCreateBasedir:"basedirs/create",getDirectory:"getDirectory"}),{showAlert:function(t){"string"===typeof t?this.alertContent=t:(this.alertVariant=t.variant||this.alertVariant,this.alertDismissible=t.dismissible||this.alertDismissible,this.alertContent=t.content||this.alertContent),this.alertShow=!0},touch:function(t){this.currentValue?(this.$set(this.touched,t,!0),this.$delete(this.errors,t)):(this.$delete(this.touched,t),this.$delete(this.errors,t))},inputState:function(t){return"undefined"===typeof this.errors[t]?null:"no error"===this.errors[t]},onAddBasedir:function(){var t=this,e=this.currentValue;e&&this.getDirectory({dir:e}).then(function(t){return t?t.data:null}).then(function(r){t.createDir(e)}).catch(function(r){404===r.response.status?t.$set(t.notfound,e,1):t.$delete(t.notfound,e)})},createDir:function(t){var e=this,r={attributes:{basedir:t}};this.doCreateBasedir(r).then(function(t){e.$set(e.touched,"add",!0),e.$delete(e.errors,"add"),e.$delete(e.notfound,e.currentValue),e.currentValue=""}).catch(function(t){e.$set(e.errors,"add",t.data.errors[0].title||t.statusText),e.$delete(e.touched,"add")})}})},w=O,D=(r("e165"),Object(p["a"])(w,y,v,!1,null,"6b21ec3e",null)),_=D.exports;function j(t,e){var r=Object.keys(t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(t);e&&(i=i.filter(function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable})),r.push.apply(r,i)}return r}function C(t){for(var e=1;e<arguments.length;e++){var r=null!=arguments[e]?arguments[e]:{};e%2?j(r,!0).forEach(function(e){Object(o["a"])(t,e,r[e])}):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(r)):j(r).forEach(function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(r,e))})}return t}var P={name:"Preferences",components:{BasedirRowEdit:h,BasedirRowAdd:_},data:function(){return{isLoading:!0,errors:{},touched:{}}},computed:C({},Object(a["d"])({basedirs:"basedirs/all"})),mounted:function(){var t=this;this.doLoadBasedirs().then(function(){t.isLoading=!1})},methods:C({},Object(a["c"])({doLoadBasedirs:"basedirs/loadAll"}))},k=P,V=(r("778f"),Object(p["a"])(k,i,n,!1,null,"5b5ad31a",null));e["default"]=V.exports},aa77:function(t,e,r){var i=r("5ca1"),n=r("be13"),o=r("79e5"),a=r("fdef"),s="["+a+"]",c="​",d=RegExp("^"+s+s+"*"),u=RegExp(s+s+"*$"),l=function(t,e,r){var n={},s=o(function(){return!!a[t]()||c[t]()!=c}),d=n[t]=s?e(f):a[t];r&&(n[r]=d),i(i.P+i.F*s,"String",n)},f=l.trim=function(t,e){return t=String(n(t)),1&e&&(t=t.replace(d,"")),2&e&&(t=t.replace(u,"")),t};t.exports=l},c5f6:function(t,e,r){"use strict";var i=r("7726"),n=r("69a8"),o=r("2d95"),a=r("5dbc"),s=r("6a99"),c=r("79e5"),d=r("9093").f,u=r("11e9").f,l=r("86cc").f,f=r("aa77").trim,p="Number",b=i[p],h=b,y=b.prototype,v=o(r("2aeb")(y))==p,g="trim"in String.prototype,m=function(t){var e=s(t,!1);if("string"==typeof e&&e.length>2){e=g?e.trim():f(e,3);var r,i,n,o=e.charCodeAt(0);if(43===o||45===o){if(r=e.charCodeAt(2),88===r||120===r)return NaN}else if(48===o){switch(e.charCodeAt(1)){case 66:case 98:i=2,n=49;break;case 79:case 111:i=8,n=55;break;default:return+e}for(var a,c=e.slice(2),d=0,u=c.length;d<u;d++)if(a=c.charCodeAt(d),a<48||a>n)return NaN;return parseInt(c,i)}}return+e};if(!b(" 0o1")||!b("0b1")||b("+0x1")){b=function(t){var e=arguments.length<1?0:t,r=this;return r instanceof b&&(v?c(function(){y.valueOf.call(r)}):o(r)!=p)?a(new h(m(e)),r,b):m(e)};for(var O,w=r("9e1e")?d(h):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),D=0;w.length>D;D++)n(h,O=w[D])&&!n(b,O)&&l(b,O,u(h,O));b.prototype=y,y.constructor=b,r("2aba")(i,p,b)}},e165:function(t,e,r){"use strict";var i=r("a4fc"),n=r.n(i);n.a},fdef:function(t,e){t.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);
//# sourceMappingURL=preferences.804ca4ee.js.map