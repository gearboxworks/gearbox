(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["projectstack"],{"0201":function(t,e,r){"use strict";var i=r("bf6c"),s=r.n(i);s.a},"0368":function(t,e,r){},"055b":function(t,e,r){t.exports=r.p+"img/python.51c2eab2.svg"},"090e":function(t,e,r){t.exports=r.p+"img/nodejs.94cafb0d.svg"},"11e9":function(t,e,r){var i=r("52a7"),s=r("4630"),a=r("6821"),n=r("6a99"),o=r("69a8"),c=r("c69a"),l=Object.getOwnPropertyDescriptor;e.f=r("9e1e")?l:function(t,e){if(t=a(t),e=n(e,!0),c)try{return l(t,e)}catch(r){}if(o(t,e))return s(!i.f.call(t,e),t[e])}},1540:function(t,e,r){t.exports=r.p+"img/php.fa78b345.svg"},"1b4d":function(t,e,r){t.exports=r.p+"img/flask.318d58cb.svg"},2232:function(t,e,r){"use strict";r.r(e);var i=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{class:{"project-stack-list":!0,"cards-mode":t.cardsMode},attrs:{role:"tablist",id:t.projectBase+"stack"}},t._l(t.groupedStackItems(t.projectStackItems),function(e,i,s){return r("project-stack",{key:i,attrs:{stackId:i,stackIndex:s,stackItems:e,project:t.project,projectIndex:t.projectIndex,"is-collapsible":t.cardsMode}})}),1)},s=[],a=(r("ac6a"),r("a481"),r("cebc")),n=(r("c5f6"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{class:{"project-stack":!0,"is-collapsible":t.isCollapsible,"is-collapsed":t.isCollapsed}},[r("h2",{staticClass:"stack-title",on:{click:t.onTitleClicked}},[t._v(t._s(t.stackId.replace("gearbox.works/","")))]),t.isCollapsed?t._e():r("stack-toolbar",{attrs:{project:t.project,projectIndex:t.projectIndex,stackId:t.stackId}}),r("div",{staticClass:"stack-content"},[t.isCollapsible&&t.isCollapsed?t._e():r("ul",{staticClass:"service-list"},t._l(t.stackItems,function(e,i){return r("li",{key:t.id+e.gearspec.attributes.role,staticClass:"service-item"},[r("stack-gear",{attrs:{projectId:t.project.id,stackItem:e,projectIndex:t.projectIndex,stackIndex:t.stackIndex,itemIndex:i}})],1)}),0),r("b-alert",{key:t.stackId,attrs:{show:t.alertShow,dismissible:t.alertDismissible,variant:t.alertVariant,fade:""},on:{dismissed:function(e){t.alertShow=!1}}},[t._v(t._s(t.alertContent))])],1)],1)}),o=[],c=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",{staticClass:"toolbar-list"},[r("div",{key:t.projectBase+t.stackId+"more",staticClass:"toolbar-item toolbar-item--more"},[r("b-dropdown",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],attrs:{variant:"link","no-caret":"",title:"More stack actions"}},[r("template",{slot:"button-content"},[r("font-awesome-icon",{attrs:{icon:["fa","ellipsis-v"]}}),r("span",{staticClass:"sr-only"},[t._v("More actions")])],1),r("b-dropdown-item",{attrs:{href:"#"},on:{click:function(e){return e.preventDefault(),t.removeProjectStack(e)}}},[t._v("Remove stack")])],2)],1),t.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[t.isRunning?r("li",{key:t.projectBase+t.stackId+"frontend",class:["toolbar-item","toolbar-item--frontend"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+t.hostname+"/",title:"Open Frontend"+(t.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","home"]}})],1)]):t._e()]):t._e(),t.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[t.isRunning?r("li",{key:t.projectBase+t.stackId+"dashboard",class:["toolbar-item","toolbar-item--dashboard"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+t.hostname+"/wp-admin/",title:"Open Dashboard"+(t.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","tachometer-alt"]}})],1)]):t._e()]):t._e(),t.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[t.isRunning?r("li",{key:t.projectBase+t.stackId+"adminer",class:["toolbar-item","toolbar-item--adminer"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+t.hostname+"/wp-admin/",title:"Open Adminer"+(t.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","database"]}})],1)]):t._e()]):t._e(),t.isWordPress?r("transition",{attrs:{name:"icons",tag:"ul"}},[t.isRunning?r("li",{key:t.projectBase+t.stackId+"mailhog",class:["toolbar-item","toolbar-item--frontend"]},[r("a",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover",modifiers:{hover:!0}}],staticClass:"toolbar-link",attrs:{target:"_blank",href:"//"+t.hostname+":4003",title:"Open Mailhog"+(t.isRunning?"":" (not running)")}},[r("font-awesome-icon",{attrs:{icon:["fa","envelope"]}})],1)]):t._e()]):t._e()],1)},l=[],p={name:"StackToolbar",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackId:{type:String,required:!0}},data:function(){return{isWordPress:-1!==this.stackId.indexOf("/wordpress")}},computed:{projectBase:function(){return"gb-"+this.escAttr(this.project.id)+"-"},hostname:function(){return this.project.attributes.hostname},isRunning:function(){return this.project.attributes.enabled}},methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},removeProjectStack:function(t){this.$store.dispatch("removeProjectStack",{projectId:this.project.id,stackId:this.stackId})}}},d=p,u=(r("d423"),r("2877")),f=Object(u["a"])(d,c,l,!1,null,"a9d62bc0",null),g=f.exports,v=function(){var t=this,e=t.$createElement,i=t._self._c||e;return i("div",{staticClass:"project-gear",attrs:{id:t.gearControlId,tabindex:100*t.projectIndex+10*t.stackIndex+t.itemIndex+1}},[t.service?i("img",{class:{"service-program":!0,"is-loaded":t.isLoaded,"is-switching":t.isSwitching,"is-switching-same":t.isSwitchingSame,"is-switching-same-again":t.isSwitchingSameAgain},attrs:{src:r("31e8")("./"+t.service.attributes.program+".svg")},on:{load:t.onImageLoaded}}):i("font-awesome-icon",{attrs:{icon:["fa","expand"]}}),i("h6",{staticClass:"gear-role"},[t._v(t._s(t.gearspec.attributes.role))]),i("b-tooltip",{key:t.gearControlId+"-"+(t.service?t.service.id:"unselected"),attrs:{triggers:"hover",target:t.gearControlId,title:t.programTooltip}}),i("b-popover",{ref:t.gearControlId+"-popover",attrs:{target:t.gearControlId,container:t.projectBase+"stack",triggers:"focus",placement:"bottom"}},[i("template",{slot:"title"},[i("b-button",{staticClass:"close",attrs:{"aria-label":"Close"},on:{click:t.closePopover}},[i("span",{staticClass:"d-inline-block",attrs:{"aria-hidden":"true"}},[t._v("×")])]),t._v("\n      Change service\n    ")],1),i("b-form-group",[i("label",{attrs:{for:t.gearControlId+"-input"}},[t._v(t._s(t.gearspec.attributes.role)+":")]),i("b-form-select",{attrs:{id:t.gearControlId+"-input",value:t.preselectClosestGearServiceId,tabindex:100*t.projectIndex+10*t.stackIndex+t.itemIndex+9},on:{change:function(e){return t.onChangeService(e)}}},[t.defaultService?t._e():i("option",{attrs:{value:""}},[t._v("Do not run this service")]),i("option",{attrs:{disabled:""},domProps:{value:null}},[t._v("Select service...")]),t._l(t.servicesGroupedByRole,function(e,r){return i("optgroup",{key:r,attrs:{label:r}},t._l(e,function(e){return i("option",{key:e,domProps:{value:e}},[t._v(t._s(e.replace("gearboxworks/","")))])}),0)})],2)],1)],2)],1)},h=[],m=(r("28a5"),r("2f62")),b={name:"StackGear",props:{projectId:{type:String,required:!0},stackItem:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackIndex:{type:Number,required:!0},itemIndex:{type:Number,required:!0}},data:function(){return{isLoaded:!1,isSwitching:!0,isSwitchingSame:!1,isSwitchingSameAgain:!1}},computed:Object(a["a"])({},Object(m["b"])(["serviceBy","gearspecBy","stackBy","stackDefaultServiceByRole","stackServicesByRole","preselectServiceId"]),{projectBase:function(){return"gb-"+this.escAttr(this.projectId)+"-"},gearspec:function(){return this.stackItem.gearspec},service:function(){var t=null;if(this.stackItem.service)t=this.stackItem.service;else if(this.stackItem.serviceId){var e=this.preselectClosestGearServiceId;e&&(t=this.serviceBy("id",e))}return t},stack:function(){return this.stackBy("id",this.gearspec.attributes.stack_id)},gearControlId:function(){return this.projectBase+(this.stack?this.stack.attributes.stackname+"-":"")+this.gearspec.attributes.role},defaultService:function(){return this.stackDefaultServiceByRole(this.stack,this.stackItem.gearspecId)},preselectClosestGearServiceId:function(){return this.preselectServiceId(this.stackServicesByRole(this.stack,this.stackItem.gearspecId),this.defaultService,this.stackItem.serviceId)},servicesGroupedByRole:function(){var t=this.stackServicesByRole(this.stack,this.gearspec.id),e={};for(var r in t){var i=t[r].split(":")[0].replace("gearboxworks/","");"undefined"===typeof e[i]&&(e[i]={}),e[i][r]=t[r]}return e},programTooltip:function(){var t=this.stackItem.serviceId,e=t&&this.service?this.service.attributes:null,r=e?e.program:"",i=e?e.version:"";return t&&(!e||this.service&&this.service.id!==t)&&(r=t.split("/")[1].split(":")[0],i=t.split("/")[1].split(":")[1]),r&&i?r+" "+i:"Service not selected"}}),methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},onChangeService:function(t){var e=this.service?this.service.id:"",r=e?e.split("/")[1].split(":")[0]:"",i=t?t.split("/")[1].split(":")[0]:"";r!==i?(this.isLoaded=!1,this.isSwitching=!0,this.isSwitchingSame=!1,this.isSwitchingSameAgain=!1):e!==t&&(this.isSwitchingSame||this.isSwitchingSameAgain?(this.isSwitchingSame=!this.isSwitchingSame,this.isSwitchingSameAgain=!this.isSwitchingSameAgain):(this.isSwitchingSame=!0,this.isSwitchingSameAgain=!1)),this.$store.dispatch("changeProjectService",{projectId:this.projectId,gearspecId:this.gearspec.id,serviceId:t}),this.closePopover()},closePopover:function(){this.$root.$emit("bv::hide::popover",this.gearControlId)},onImageLoaded:function(t){this.isSwitching=!1,this.isLoaded=!0}}},k=b,I=(r("7998"),Object(u["a"])(k,v,h,!1,null,"7447962e",null)),S=I.exports,j={name:"StackCard",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},stackId:{type:String,required:!0},stackIndex:{type:Number,required:!0},stackItems:{type:Array,required:!0},isCollapsible:{type:Boolean,default:!1,required:!1}},data:function(){return{id:this.project.id,alertShow:!1,alertContent:"content",alertDismissible:!0,alertVariant:"info",isCollapsed:!0}},components:{StackToolbar:g,StackGear:S},computed:Object(a["a"])({},Object(m["b"])(["serviceBy","gearspecBy"]),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"}}),methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},onTitleClicked:function(){this.isCollapsed=!this.isCollapsed},showAlert:function(t){"string"===typeof t?this.alertContent=t:(this.alertVariant=t.variant||this.alertVariant,this.alertDismissible=t.dismissible||this.alertDismissible,this.alertContent=t.content||this.alertContent),this.alertShow=!0}}},y=j,_=(r("0201"),r("f843"),Object(u["a"])(y,n,o,!1,null,"39c16f19",null)),x=_.exports,w={name:"ProjectStackList",props:{project:{type:Object,required:!0},projectIndex:{type:Number,required:!0},cardsMode:{type:Boolean,required:!1,default:!1}},data:function(){return{id:this.project.id,projectStackItems:this.project.attributes.stack}},components:{ProjectStack:x},computed:Object(a["a"])({},Object(m["b"])(["serviceBy","gearspecBy"]),{projectBase:function(){return"gb-"+this.escAttr(this.id)+"-"}}),methods:{escAttr:function(t){return t.replace(/\//g,"-").replace(/\./g,"-")},groupedStackItems:function(t){var e=this,r={};return t.forEach(function(t){var i=e.gearspecBy("id",t.gearspec_id);if(i){"undefined"===typeof r[i.attributes.stack_id]&&(r[i.attributes.stack_id]=[]);var s=t.service_id?e.serviceBy("id",t.service_id):null;r[i.attributes.stack_id].push({gearspecId:t.gearspec_id,gearspec:i,serviceId:t.service_id,service:s})}}),r}}},C=w,N=(r("fdb0"),Object(u["a"])(C,i,s,!1,null,"26657353",null));e["default"]=N.exports},"319f":function(t,e,r){t.exports=r.p+"img/rails.2db29782.svg"},"31e8":function(t,e,r){var i={"./angular.svg":"a230","./apache.svg":"b77f","./codeigniter.svg":"7939","./django.svg":"c6da","./drupal.svg":"a88c","./elasticsearch.svg":"81bb","./flask.svg":"1b4d","./joomla.svg":"5390","./laravel.svg":"41c8","./logo.svg":"9b19","./mariadb.svg":"613e","./memcached.svg":"a0ba","./mysql.svg":"6c4c","./nginx.svg":"c502","./nodejs.svg":"090e","./perl.svg":"c44f","./php.svg":"1540","./python.svg":"055b","./rails.svg":"319f","./react.svg":"b2e9","./redis.svg":"8bcb","./ruby.svg":"9401","./wordpress.svg":"ee3c"};function s(t){var e=a(t);return r(e)}function a(t){var e=i[t];if(!(e+1)){var r=new Error("Cannot find module '"+t+"'");throw r.code="MODULE_NOT_FOUND",r}return e}s.keys=function(){return Object.keys(i)},s.resolve=a,t.exports=s,s.id="31e8"},"41c8":function(t,e,r){t.exports=r.p+"img/laravel.1766a461.svg"},4661:function(t,e,r){},5390:function(t,e,r){t.exports=r.p+"img/joomla.d8aa2e45.svg"},"5dbc":function(t,e,r){var i=r("d3f4"),s=r("8b97").set;t.exports=function(t,e,r){var a,n=e.constructor;return n!==r&&"function"==typeof n&&(a=n.prototype)!==r.prototype&&i(a)&&s&&s(t,a),t}},"613e":function(t,e,r){t.exports=r.p+"img/mariadb.e16110bc.svg"},"6c4c":function(t,e,r){t.exports=r.p+"img/mysql.dd2a5a35.svg"},7939:function(t,e,r){t.exports=r.p+"img/codeigniter.434bf735.svg"},7998:function(t,e,r){"use strict";var i=r("a367"),s=r.n(i);s.a},"81bb":function(t,e,r){t.exports=r.p+"img/elasticsearch.3ecfa530.svg"},"8b97":function(t,e,r){var i=r("d3f4"),s=r("cb7c"),a=function(t,e){if(s(t),!i(e)&&null!==e)throw TypeError(e+": can't set as prototype!")};t.exports={set:Object.setPrototypeOf||("__proto__"in{}?function(t,e,i){try{i=r("9b43")(Function.call,r("11e9").f(Object.prototype,"__proto__").set,2),i(t,[]),e=!(t instanceof Array)}catch(s){e=!0}return function(t,r){return a(t,r),e?t.__proto__=r:i(t,r),t}}({},!1):void 0),check:a}},"8bcb":function(t,e,r){t.exports=r.p+"img/redis.3c39fafe.svg"},9093:function(t,e,r){var i=r("ce10"),s=r("e11e").concat("length","prototype");e.f=Object.getOwnPropertyNames||function(t){return i(t,s)}},9401:function(t,e,r){t.exports=r.p+"img/ruby.514befa7.svg"},"9b19":function(t,e,r){t.exports=r.p+"img/logo.63a7d78d.svg"},a0ba:function(t,e,r){t.exports=r.p+"img/memcached.2bcccabf.svg"},a230:function(t,e,r){t.exports=r.p+"img/angular.e224f5ed.svg"},a367:function(t,e,r){},a481:function(t,e,r){"use strict";var i=r("cb7c"),s=r("4bf8"),a=r("9def"),n=r("4588"),o=r("0390"),c=r("5f1b"),l=Math.max,p=Math.min,d=Math.floor,u=/\$([$&`']|\d\d?|<[^>]*>)/g,f=/\$([$&`']|\d\d?)/g,g=function(t){return void 0===t?t:String(t)};r("214f")("replace",2,function(t,e,r,v){return[function(i,s){var a=t(this),n=void 0==i?void 0:i[e];return void 0!==n?n.call(i,a,s):r.call(String(a),i,s)},function(t,e){var s=v(r,t,this,e);if(s.done)return s.value;var d=i(t),u=String(this),f="function"===typeof e;f||(e=String(e));var m=d.global;if(m){var b=d.unicode;d.lastIndex=0}var k=[];while(1){var I=c(d,u);if(null===I)break;if(k.push(I),!m)break;var S=String(I[0]);""===S&&(d.lastIndex=o(u,a(d.lastIndex),b))}for(var j="",y=0,_=0;_<k.length;_++){I=k[_];for(var x=String(I[0]),w=l(p(n(I.index),u.length),0),C=[],N=1;N<I.length;N++)C.push(g(I[N]));var A=I.groups;if(f){var O=[x].concat(C,w,u);void 0!==A&&O.push(A);var B=String(e.apply(void 0,O))}else B=h(x,u,w,C,A,e);w>=y&&(j+=u.slice(y,w)+B,y=w+x.length)}return j+u.slice(y)}];function h(t,e,i,a,n,o){var c=i+t.length,l=a.length,p=f;return void 0!==n&&(n=s(n),p=u),r.call(o,p,function(r,s){var o;switch(s.charAt(0)){case"$":return"$";case"&":return t;case"`":return e.slice(0,i);case"'":return e.slice(c);case"<":o=n[s.slice(1,-1)];break;default:var p=+s;if(0===p)return r;if(p>l){var u=d(p/10);return 0===u?r:u<=l?void 0===a[u-1]?s.charAt(1):a[u-1]+s.charAt(1):r}o=a[p-1]}return void 0===o?"":o})}})},a88c:function(t,e,r){t.exports=r.p+"img/drupal.66089b06.svg"},aa77:function(t,e,r){var i=r("5ca1"),s=r("be13"),a=r("79e5"),n=r("fdef"),o="["+n+"]",c="​",l=RegExp("^"+o+o+"*"),p=RegExp(o+o+"*$"),d=function(t,e,r){var s={},o=a(function(){return!!n[t]()||c[t]()!=c}),l=s[t]=o?e(u):n[t];r&&(s[r]=l),i(i.P+i.F*o,"String",s)},u=d.trim=function(t,e){return t=String(s(t)),1&e&&(t=t.replace(l,"")),2&e&&(t=t.replace(p,"")),t};t.exports=d},b2e9:function(t,e,r){t.exports=r.p+"img/react.9a28da9f.svg"},b77f:function(t,e,r){t.exports=r.p+"img/apache.12c49354.svg"},bf6c:function(t,e,r){},c44f:function(t,e,r){t.exports=r.p+"img/perl.a025edb4.svg"},c502:function(t,e,r){t.exports=r.p+"img/nginx.eae76401.svg"},c5f6:function(t,e,r){"use strict";var i=r("7726"),s=r("69a8"),a=r("2d95"),n=r("5dbc"),o=r("6a99"),c=r("79e5"),l=r("9093").f,p=r("11e9").f,d=r("86cc").f,u=r("aa77").trim,f="Number",g=i[f],v=g,h=g.prototype,m=a(r("2aeb")(h))==f,b="trim"in String.prototype,k=function(t){var e=o(t,!1);if("string"==typeof e&&e.length>2){e=b?e.trim():u(e,3);var r,i,s,a=e.charCodeAt(0);if(43===a||45===a){if(r=e.charCodeAt(2),88===r||120===r)return NaN}else if(48===a){switch(e.charCodeAt(1)){case 66:case 98:i=2,s=49;break;case 79:case 111:i=8,s=55;break;default:return+e}for(var n,c=e.slice(2),l=0,p=c.length;l<p;l++)if(n=c.charCodeAt(l),n<48||n>s)return NaN;return parseInt(c,i)}}return+e};if(!g(" 0o1")||!g("0b1")||g("+0x1")){g=function(t){var e=arguments.length<1?0:t,r=this;return r instanceof g&&(m?c(function(){h.valueOf.call(r)}):a(r)!=f)?n(new v(k(e)),r,g):k(e)};for(var I,S=r("9e1e")?l(v):"MAX_VALUE,MIN_VALUE,NaN,NEGATIVE_INFINITY,POSITIVE_INFINITY,EPSILON,isFinite,isInteger,isNaN,isSafeInteger,MAX_SAFE_INTEGER,MIN_SAFE_INTEGER,parseFloat,parseInt,isInteger".split(","),j=0;S.length>j;j++)s(v,I=S[j])&&!s(g,I)&&d(g,I,p(v,I));g.prototype=h,h.constructor=g,r("2aba")(i,f,g)}},c6da:function(t,e,r){t.exports=r.p+"img/django.28fe09a0.svg"},d423:function(t,e,r){"use strict";var i=r("0368"),s=r.n(i);s.a},ee3c:function(t,e,r){t.exports=r.p+"img/wordpress.b08e20e3.svg"},f332:function(t,e,r){},f843:function(t,e,r){"use strict";var i=r("4661"),s=r.n(i);s.a},fdb0:function(t,e,r){"use strict";var i=r("f332"),s=r.n(i);s.a},fdef:function(t,e){t.exports="\t\n\v\f\r   ᠎             　\u2028\u2029\ufeff"}}]);
//# sourceMappingURL=projectstack.9d8f664b.js.map