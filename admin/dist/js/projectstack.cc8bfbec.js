(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projectstack"],{"0368":function(e,t,r){},"055b":function(e,t,r){e.exports=r.p+"img/python.51c2eab2.svg"},"090e":function(e,t,r){e.exports=r.p+"img/nodejs.94cafb0d.svg"},"11e9":function(e,t,r){var i=r("52a7"),s=r("4630"),a=r("6821"),n=r("6a99"),o=r("69a8"),c=r("c69a"),l=Object.getOwnPropertyDescriptor;t.f=r("9e1e")?l:function(e,t){if(e=a(e),t=n(t,!0),c)try{return l(e,t)}catch(r){}if(o(e,t))return s(!i.f.call(e,t),e[t])}},1540:function(e,t,r){e.exports=r.p+"img/php.fa78b345.svg"},"1b4d":function(e,t,r){e.exports=r.p+"img/flask.318d58cb.svg"},2232:function(e,t,r){"use strict";r.r(t);var i=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{class:{"project-stack-list":!0,"start-collapsed":e.startCollapsed,"is-loading":e.isLoading},attrs:{id:e.projectBase+"stack",role:"tablist"}},[e.isLoading?r("font-awesome-icon",{attrs:{icon:"circle-notch",spin:"",title:"Loading project details..."}}):r("div",{staticClass:"project-stack-list-wrap"},e._l(e.groupedStackItems,function(t,i,s){return r("stack-card",{key:i,attrs:{stackId:i,stackIndex:s,stackItems:t,project:e.project,projectIndex:e.projectIndex,"start-collapsed":e.startCollapsed||!e.startCollapsed&&Object.entries(e.groupedStackItems).length>1}})}),1)],1)},s=[],a=(r("a481"),r("ac6a"),r("cebc")),n=(r("c5f6"),function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{class:{"stack-card":!0,"is-collapsible":e.isCollapsible,"is-collapsed":e.isCollapsed}},[r("h2",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"stack-title",attrs:{title:e.isCollapsed?"Show services":"Hide services"},on:{click:e.onTitleClicked}},[r("font-awesome-icon",{staticStyle:{color:"silver"},attrs:{icon:["fa","layer-group"]}}),e._v(" \n    "+e._s(e.stackId.replace("gearbox.works/",""))+"\n  ")],1),e.isCollapsed?e._e():r("stack-toolbar",{attrs:{project:e.project,projectIndex:e.projectIndex,stackId:e.stackId}}),r("div",{staticClass:"stack-content"},[e.isCollapsible&&e.isCollapsed?e._e():r("ul",{staticClass:"service-list"},e._l(e.stackItems,function(t,i){return r("li",{key:e.id+t.gearspec.attributes.role,staticClass:"service-item"},[r("stack-gear",{attrs:{project:e.project,stackItem:t,projectIndex:e.projectIndex,stackIndex:e.stackIndex,itemIndex:i}})],1)}),0),r("b-alert",{key:e.stackId,attrs:{show:e.alertShow,dismissible:e.alertDismissible,variant:e.alertVariant,fade:""},on:{dismissed:function(t){e.alertShow=!1}}},[e._v(e._s(e.alertContent))])],1)],1)}),o=[],c=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"toolbar-list"},[r("div",{key:e.projectBase+e.stackId+"more",staticClass:"toolbar-item toolbar-item--more"},[r("b-dropdown",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],attrs:{variant:"link","no-caret":"",title:"More stack actions"}},[r("template",{slot:"button-content"},[r("font-awesome-icon",{attrs:{icon:["fa","ellipsis-v"]}}),r("span",{staticClass:"sr-only"},[e._v("More actions")])],1),r("b-dropdown-item",{attrs:{href:"#"},on:{click:function(t){return t.preventDefault(),e.removeProjectStack(t)}}},[e._v("Remove stack")])],2)],1),e.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[e.isRunning?r("li",{key:e.projectBase+e.stackId+"frontend",class:["toolbar-item","toolbar-item--frontend"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+e.hostname+"/",title:"Open Frontend"+(e.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","home"]}})],1)]):e._e()]):e._e(),e.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[e.isRunning?r("li",{key:e.projectBase+e.stackId+"dashboard",class:["toolbar-item","toolbar-item--dashboard"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+e.hostname+"/wp-admin/",title:"Open Dashboard"+(e.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","tachometer-alt"]}})],1)]):e._e()]):e._e(),e.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[e.isRunning?r("li",{key:e.projectBase+e.stackId+"adminer",class:["toolbar-item","toolbar-item--adminer"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+e.hostname+"/wp-admin/",title:"Open Adminer"+(e.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","database"]}})],1)]):e._e()]):e._e(),e.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[e.isRunning?r("li",{key:e.projectBase+e.stackId+"mailhog",class:["toolbar-item","toolbar-item--frontend"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+e.hostname+":4003",title:"Open Mailhog"+(e.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","envelope"]}})],1)]):e._e()]):e._e()],1)},l=[],p={name:"StackToolbar",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackId:{type:String,required:!0}},data:function(){return{isWordPress:-1!==this.stackId.indexOf("/wordpress")}},computed:{projectBase:function(){return"gb-"+this.escAttr(this.project.id)+"-"},hostname:function(){return this.project.attributes.hostname},isRunning:function(){return this.project.attributes.enabled}},methods:{escAttr:function(e){return e.replace(/\//g,"-").replace(/\./g,"-")},removeProjectStack:function(e){this.$store.dispatch("removeProjectStack",{projectId:this.project.id,stackId:this.stackId})}}},d=p,u=(r("d423"),r("2877")),f=Object(u["a"])(d,c,l,!1,null,"a9d62bc0",null),g=f.exports,v=function(){var e=this,t=e.$createElement,i=e._self._c||t;return i("div",{staticClass:"project-gear",attrs:{id:e.gearControlId,tabindex:100*e.projectIndex+10*e.stackIndex+e.itemIndex+1}},[e.service?i("img",{class:{"service-program":!0,"is-loaded":e.isLoaded,"is-switching":e.isSwitching,"is-switching-same":e.isSwitchingSame,"is-switching-same-again":e.isSwitchingSameAgain},attrs:{src:r("31e8")("./"+e.service.attributes.program+".svg")},on:{load:e.onImageLoaded}}):i("font-awesome-icon",{attrs:{icon:["fa","expand"]}}),i("h6",{staticClass:"gear-role"},[e._v(e._s(e.gearspec.attributes.role))]),i("b-tooltip",{key:e.gearControlId+"-"+(e.service?e.service.id:"unselected"),attrs:{triggers:"hover",target:e.gearControlId,title:e.programTooltip}}),i("b-popover",{ref:e.gearControlId+"-popover",attrs:{target:e.gearControlId,container:e.projectBase+"stack",triggers:"focus",placement:"bottom"}},[i("template",{slot:"title"},[i("b-button",{staticClass:"close",attrs:{"aria-label":"Close"},on:{click:e.closePopover}},[i("span",{staticClass:"d-inline-block",attrs:{"aria-hidden":"true"}},[e._v("×")])]),e._v("\n      Change service\n    ")],1),i("b-form-group",[i("label",{attrs:{for:e.gearControlId+"-input"}},[e._v(e._s(e.gearspec.attributes.role)+":")]),i("b-form-select",{ref:e.gearControlId+"-select",attrs:{value:e.preselectClosestGearServiceId,tabindex:100*e.projectIndex+10*e.stackIndex+e.itemIndex+9},on:{change:function(t){return e.onChangeService(t)}}},[e.defaultService?e._e():i("option",{attrs:{value:""}},[e._v("Do not run this service")]),i("option",{attrs:{disabled:""},domProps:{value:null}},[e._v("Select service...")]),e._l(e.servicesGroupedByRole,function(t,r){return i("optgroup",{key:r,attrs:{label:r}},e._l(t,function(t){return i("option",{key:t,attrs:{disabled:e.project.attributes.enabled},domProps:{value:t}},[e._v(e._s(t.replace("gearboxworks/","")))])}),0)})],2),i("b-alert",{attrs:{show:!e.stackItem.service,variant:"warning"}},[e._v("Note, the currently selected version of the service is different from what is in project specification!")]),i("b-alert",{attrs:{show:e.project.attributes.enabled}},[e._v("Note, you cannot change this service while the project is running!")])],1)],2)],1)},h=[],b=(r("28a5"),r("2f62")),m={name:"StackGear",props:{project:{type:Object,required:!0},stackItem:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackIndex:{type:Number,required:!0},itemIndex:{type:Number,required:!0}},data:function(){return{isLoaded:!1,isSwitching:!0,isSwitchingSame:!1,isSwitchingSameAgain:!1}},computed:Object(a["a"])({},Object(b["c"])(["serviceBy","gearspecBy","stackBy","stackDefaultServiceByRole","stackServicesByRole","preselectServiceId"]),{projectBase:function(){return"gb-"+this.escAttr(this.project.id)+"-"},gearspec:function(){return this.stackItem.gearspec},service:function(){var e=null;if(this.stackItem.service)e=this.stackItem.service;else if(this.stackItem.serviceId){var t=this.preselectClosestGearServiceId;t&&(e=this.serviceBy("id",t))}return e},stack:function(){return this.stackBy("id",this.gearspec.attributes.stack_id)},gearControlId:function(){return this.projectBase+(this.stack?this.stack.attributes.stackname+"-":"")+this.gearspec.attributes.role},defaultService:function(){return this.stackDefaultServiceByRole(this.stack,this.stackItem.gearspecId)},preselectClosestGearServiceId:function(){return this.preselectServiceId(this.stackServicesByRole(this.stack,this.stackItem.gearspecId),this.defaultService,this.stackItem.serviceId)},servicesGroupedByRole:function(){var e=this.stackServicesByRole(this.stack,this.gearspec.id),t={};for(var r in e){var i=e[r].split(":")[0].replace("gearboxworks/","");"undefined"===typeof t[i]&&(t[i]={}),t[i][r]=e[r]}return t},programTooltip:function(){var e=this.stackItem.serviceId,t=e&&this.service?this.service.attributes:null,r=t?t.program:"",i=t?t.version:"";return e&&(!t||this.service&&this.service.id!==e)&&(r=e.split("/")[1].split(":")[0],i=e.split("/")[1].split(":")[1]),r&&i?r+" "+i:"Service not selected"}}),methods:{escAttr:function(e){return e.replace(/\//g,"-").replace(/\./g,"-")},onChangeService:function(e){var t=this.service?this.service.id:"",r=t?t.split("/")[1].split(":")[0]:"",i=e?e.split("/")[1].split(":")[0]:"";r!==i?(this.isLoaded=!1,this.isSwitching=!0,this.isSwitchingSame=!1,this.isSwitchingSameAgain=!1):t!==e&&(this.isSwitchingSame||this.isSwitchingSameAgain?(this.isSwitchingSame=!this.isSwitchingSame,this.isSwitchingSameAgain=!this.isSwitchingSameAgain):(this.isSwitchingSame=!0,this.isSwitchingSameAgain=!1)),this.$store.dispatch("changeProjectService",{projectId:this.project.id,gearspecId:this.gearspec.id,serviceId:e}),this.closePopover()},closePopover:function(){this.$root.$emit("bv::hide::popover",this.gearControlId)},onImageLoaded:function(e){this.isSwitching=!1,this.isLoaded=!0}}},k=m,I=(r("5e01"),Object(u["a"])(k,v,h,!1,null,"6f9591f6",null)),j=I.exports,S={name:"StackCard",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackId:{type:String,required:!0},stackIndex:{type:Number,required:!0},stackItems:{type:Array,required:!0},isCollapsible:{type:Boolean,default:!0,required:!1},startCollapsed:{type:Boolean,required:!1,default:!0}},data:function(){return{id:this.project.id,alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"info",isCollapsed:this.startCollapsed}},components:{StackToolbar:g,StackGear:j},computed:Object(a["a"])({},Object(b["c"])(["serviceBy","gearspecBy"]),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"}}),methods:{escAttr:function(e){return e.replace(/\//g,"-").replace(/\./g,"-")},onTitleClicked:function(){this.isCollapsed=!this.isCollapsed},showAlert:function(e){"string"===typeof e?this.alertContent=e:(this.alertVariant=e.variant||this.alertVariant,this.alertDismissible=e.dismissible||this.alertDismissible,this.alertContent=e.content||this.alertContent),this.alertShow=!0}}},y=S,_=(r("844e"),r("f843"),Object(u["a"])(y,n,o,!1,null,"086a489d",null)),w=_.exports,x={name:"StackCardList",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},startCollapsed:{type:Boolean,required:!1,default:!1}},data:function(){return{id:this.project.id}},components:{StackCard:w},computed:Object(a["a"])({},Object(b["c"])(["serviceBy","gearspecBy"]),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"},isLoading:function(){return"undefined"===typeof this.project.attributes.stack},groupedStackItems:function(){var e=this,t={},r=this.project.attributes.stack||[];return r.forEach(function(r){var i=e.gearspecBy("id",r.gearspec_id);if(i){"undefined"===typeof t[i.attributes.stack_id]&&(t[i.attributes.stack_id]=[]);var s=r.service_id?e.serviceBy("id",r.service_id):null;t[i.attributes.stack_id].push({gearspecId:r.gearspec_id,gearspec:i,serviceId:r.service_id,service:s})}}),t}}),methods:{escAttr:function(e){return e.replace(/\//g,"-").replace(/\./g,"-")}}},C=x,N=(r("e5eb"),Object(u["a"])(C,i,s,!1,null,"e9b5085c",null));t["default"]=N.exports},"319f":function(e,t,r){e.exports=r.p+"img/rails.2db29782.svg"},"31e8":function(e,t,r){var i={"./angular.svg":"a230","./apache.svg":"b77f","./codeigniter.svg":"7939","./django.svg":"c6da","./drupal.svg":"a88c","./elasticsearch.svg":"81bb","./flask.svg":"1b4d","./joomla.svg":"5390","./laravel.svg":"41c8","./logo.svg":"9b19","./mariadb.svg":"613e","./memcached.svg":"a0ba","./mysql.svg":"6c4c","./nginx.svg":"c502","./nodejs.svg":"090e","./perl.svg":"c44f","./php.svg":"1540","./python.svg":"055b","./rails.svg":"319f","./react.svg":"b2e9","./redis.svg":"8bcb","./ruby.svg":"9401","./wordpress.svg":"ee3c"};function s(e){var t=a(e);return r(t)}function a(e){var t=i[e];if(!(t+1)){var r=new Error("Cannot find module '"+e+"'");throw r.code="MODULE_NOT_FOUND",r}return t}s.keys=function(){return Object.keys(i)},s.resolve=a,e.exports=s,s.id="31e8"},"41c8":function(e,t,r){e.exports=r.p+"img/laravel.1766a461.svg"},4661:function(e,t,r){},5390:function(e,t,r){e.exports=r.p+"img/joomla.d8aa2e45.svg"},"5dbc":function(e,t,r){var i=r("d3f4"),s=r("8b97").set;e.exports=function(e,t,r){var a,n=t.constructor;return n!==r&&"function"==typeof n&&(a=n.prototype)!==r.prototype&&i(a)&&s&&s(e,a),e}},"5e01":function(e,t,r){"use strict";var i=r("ccb9b"),s=r.n(i);s.a},"5e19":function(e,t,r){},"613e":function(e,t,r){e.exports=r.p+"img/mariadb.e16110bc.svg"},"6c4c":function(e,t,r){e.exports=r.p+"img/mysql.dd2a5a35.svg"},7939:function(e,t,r){e.exports=r.p+"img/codeigniter.434bf735.svg"},"81bb":function(e,t,r){e.exports=r.p+"img/elasticsearch.3ecfa530.svg"},"844e":function(e,t,r){"use strict";var i=r("c133"),s=r.n(i);s.a},"8b97":function(e,t,r){var i=r("d3f4"),s=r("cb7c"),a=function(e,t){if(s(e),!i(t)&&null!==t)throw TypeError(t+": can't set as prototype!")};e.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(e,t,i){try{i=r("9b43")(Function.call,r("11e9").f(Object.prototype,"__proto__").set,2),i(e,[]),t=!(e instanceof Array)}catch(s){t=!0}return function(e,r){return a(e,r),t?e.__proto__=r:i(e,r),e}}({},!1):void 0),check:a}},"8bcb":function(e,t,r){e.exports=r.p+"img/redis.3c39fafe.svg"},9093:function(e,t,r){var i=r("ce10"),s=r("e11e").concat("length","prototype");t.f=Object.getOwnPropertyNames||function(e){return i(e,s)}},9401:function(e,t,r){e.exports=r.p+"img/ruby.514befa7.svg"},"9b19":function(e,t,r){e.exports=r.p+"img/logo.63a7d78d.svg"},a0ba:function(e,t,r){e.exports=r.p+"img/memcached.2bcccabf.svg"},a230:function(e,t,r){e.exports=r.p+"img/angular.e224f5ed.svg"},a481:function(e,t,r){"use strict";var i=r("cb7c"),s=r("4bf8"),a=r("9def"),n=r("4588"),o=r("0390"),c=r("5f1b"),l=Math.max,p=Math.min,d=Math.floor,u=/\$([$&`']|\d\d?|<[^>]*>)/g,f=/\$([$&`']|\d\d?)/g,g=function(e){return void 0===e?e:String(e)};r("214f")("replace",2,function(e,t,r,v){return[function(i,s){var a=e(this),n=void 0==i?void 0:i[t];return void 0!==n?n.call(i,a,s):r.call(String(a),i,s)},function(e,t){var s=v(r,e,this,t);if(s.done)return s.value;var d=i(e),u=String(this),f="function"===typeof t;f||(t=String(t));var b=d.global;if(b){var m=d.unicode;d.lastIndex=0}var k=[];while(1){var I=c(d,u);if(null===I)break;if(k.push(I),!b)break;var j=String(I[0]);""===j&&(d.lastIndex=o(u,a(d.lastIndex),m))}for(var S="",y=0,_=0;_<k.length;_++){I=k[_];for(var w=String(I[0]),x=l(p(n(I.index),u.length),0),C=[],N=1;N<I.length;N++)C.push(g(I[N]));var O=I.groups;if(f){var A=[w].concat(C,x,u);void 0!==O&&A.push(O);var B=String(t.apply(void 0,A))}else B=h(w,u,x,C,O,t);x>=y&&(S+=u.slice(y,x)+B,y=x+w.length)}return S+u.slice(y)}];function h(e,t,i,a,n,o){var c=i+e.length,l=a.length,p=f;return void 0!==n&&(n=s(n),p=u),r.call(o,p,function(r,s){var o;switch(s.charAt(0)){case"$":return"$";case"&":return e;case"`":return t.slice(0,i);case"'":return t.slice(c);case"<":o=n[s.slice(1,-1)];break;default:var p=+s;if(0===p)return r;if(p>l){var u=d(p/10);return 0===u?r:u<=l?void 0===a[u-1]?s.charAt(1):a[u-1]+s.charAt(1):r}o=a[p-1]}return void 0===o?"":o})}})},a88c:function(e,t,r){e.exports=r.p+"img/drupal.66089b06.svg"},aa77:function(e,t,r){var i=r("5ca1"),s=r("be13"),a=r("79e5"),n=r("fdef"),o="["+n+"]",c="​",l=RegExp("^"+o+o+"*"),p=RegExp(o+o+"*$"),d=function(e,t,r){var s={},o=a(function(){return!!n[e]()||c[e]()!=c}),l=s[e]=o?t(u):n[e];r&&(s[r]=l),i(i.P+i.F*o,"String",s)},u=d.trim=function(e,t){return e=String(s(e)),1&t&&(e=e.replace(l,"")),2&t&&(e=e.replace(p,"")),e};e.exports=d},b2e9:function(e,t,r){e.exports=r.p+"img/react.9a28da9f.svg"},b77f:function(e,t,r){e.exports=r.p+"img/apache.12c49354.svg"},c133:function(e,t,r){},c44f:function(e,t,r){e.exports=r.p+"img/perl.a025edb4.svg"},c502:function(e,t,r){e.exports=r.p+"img/nginx.eae76401.svg"},c5f6:function(e,t,r){"use strict";var i=r("7726"),s=r("69a8"),a=r("2d95"),n=r("5dbc"),o=r("6a99"),c=r("79e5"),l=r("9093").f,p=r("11e9").f,d=r("86cc").f,u=r("aa77").trim,f="Number",g=i[f],v=g,h=g.prototype,b=a(r("2aeb")(h))==f,m="trim"in String.prototype,k=function(e){var t=o(e,!1);if("string"==typeof t&&t.length>2){t=m?t.trim():u(t,3);var r,i,s,a=t.charCodeAt(0);if(43===a||45===a){if(r=t.charCodeAt(2),88===r||120===r)return NaN}else if(48===a){switch(t.charCodeAt(1)){case 66:case 98:i=2,s=49;break;case 79:case 111:i=8,s=55;break;default:return+t}for(var n,c=t.slice(2),l=0,p=c.length;l<p;l++)if(n=c.charCodeAt(l),n<48||n>s)return NaN;return parseInt(c,i)}}return+t};if(!g(" 0o1")||!g("0b1")||g("+0x1")){g=function(e){var t=arguments.length<1?0:e,r=this;return r instanceof g&&(b?c(function(){h.valueOf.call(r)}):a(r)!=f)?n(new v(k(t)),r,g):k(t)};for(var I,j=r("9e1e")?l(v):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),S=0;j.length>S;S++)s(v,I=j[S])&&!s(g,I)&&d(g,I,p(v,I));g.prototype=h,h.constructor=g,r("2aba")(i,f,g)}},c6da:function(e,t,r){e.exports=r.p+"img/django.28fe09a0.svg"},ccb9b:function(e,t,r){},d423:function(e,t,r){"use strict";var i=r("0368"),s=r.n(i);s.a},e5eb:function(e,t,r){"use strict";var i=r("5e19"),s=r.n(i);s.a},ee3c:function(e,t,r){e.exports=r.p+"img/wordpress.b08e20e3.svg"},f843:function(e,t,r){"use strict";var i=r("4661"),s=r.n(i);s.a},fdef:function(e,t){e.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);
//# sourceMappingURL=projectstack.cc8bfbec.js.map