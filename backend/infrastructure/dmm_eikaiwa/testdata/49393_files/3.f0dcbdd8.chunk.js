(this["webpackJsonp@biboglobal/eikaiwa-core-front"]=this["webpackJsonp@biboglobal/eikaiwa-core-front"]||[]).push([[3],{222:function(t,e,n){"use strict";n.d(e,"n",(function(){return r})),n.d(e,"d",(function(){return a})),n.d(e,"e",(function(){return c})),n.d(e,"f",(function(){return i})),n.d(e,"g",(function(){return u})),n.d(e,"h",(function(){return l})),n.d(e,"i",(function(){return o})),n.d(e,"j",(function(){return d})),n.d(e,"k",(function(){return f})),n.d(e,"l",(function(){return s})),n.d(e,"o",(function(){return b})),n.d(e,"a",(function(){return j})),n.d(e,"b",(function(){return m})),n.d(e,"c",(function(){return O})),n.d(e,"m",(function(){return g})),n.d(e,"p",(function(){return h})),n.d(e,"q",(function(){return p})),n.d(e,"r",(function(){return x}));var r="#fff",a="#808080",c="#b3b3b3",i="#ccc",u="#d8d8d8",l="#eee",o="#f2f2f2",d="#333",f="#787878",s="#404040",b="#e3463d",j="#4bbcff",m="#43aae6",O="#009DFF",g="#ef871a",h="#ebc56a",p="#FFE539",x="#FFF9D5"},234:function(t,e,n){"use strict";n.r(e);n(0);var r=n(49),a=n(222),c=n(74),i=n(73),u=n(261),l=n(1),o={paragraph:Object(i.a)(Object(i.a)({},u.defaultRules.paragraph),{},{react:function(t,e,n){return Object(l.jsx)("p",{children:e(t.content,n)},n.key)}}),url:u.defaultRules.url,link:Object(i.a)(Object(i.a)({},u.defaultRules.link),{},{react:function(t,e,n){return Object(l.jsx)("a",{href:Object(u.sanitizeUrl)(t.target),target:"_blank",children:e(t.content,n)},n.key)}}),image:u.defaultRules.image,simpleItalic:Object(i.a)(Object(i.a)({},u.defaultRules.em),{},{match:Object(u.inlineRegex)(/^_([\s\S]+?)_(?!_)/)}),simpleBold:Object(i.a)(Object(i.a)({},u.defaultRules.strong),{},{match:Object(u.inlineRegex)(/^\*([\s\S]+?)\*(?!\*)/)}),del:u.defaultRules.del,list:u.defaultRules.list,br:Object(i.a)(Object(i.a)({},u.defaultRules.br),{},{match:function(t,e){return/^\n/.exec(t)},react:function(t,e,n){return Object(l.jsx)("br",{},n.key)}}),text:Object(i.a)(Object(i.a)({},u.defaultRules.text),{},{match:Object(u.inlineRegex)(/^[^\n]+?(?=[^0-9A-Za-z\s\u00c0-\uffff]|\n|\w+:\S|$)/)})},d=Object(u.parserFor)(o),f=Object(u.reactFor)(Object(u.ruleOutput)(o,"react")),s={default:{"& p":{margin:"1em 0"},"& img":{maxWidth:"100%",height:"auto",margin:"1em 0",verticalAlign:"bottom"},"& em":{fontWeight:"normal",fontStyle:"italic"},"& strong":{fontWeight:"bold",fontStyle:"normal"},"& a":{color:a.a},"& a:hover":{color:a.b},"& ul":{paddingLeft:"1.5em",margin:"1em 0",listStyleType:"disc"},"& ol":{paddingLeft:"1.5em",margin:"1em 0",listStyleType:"decimal"}}},b=function(t){var e=t.text,n=t.classes,r=t.showError;try{var a=d("".concat(e,"\n\n"),{inline:!1}),i=f(a);return Object(l.jsx)("div",{className:n.default,"data-testid":"formatted-block-text",children:i})}catch(o){if(r){var u="RenderFormattedText received invalid rules: ".concat(o);return Object(c.a)(o),Object(l.jsx)("div",{className:n.default,"data-testid":"formatted-block-text",children:u})}return Object(l.jsx)("div",{className:n.default,"data-testid":"formatted-block-text",children:e})}};b.defaultProps={showError:!1};e.default=Object(r.b)(s)(b)}}]);
//# sourceMappingURL=3.f0dcbdd8.chunk.js.map