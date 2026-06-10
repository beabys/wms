import{k as x,l as C,C as E,B as I,y as V}from"./vue.esm-bundler-ar8GxPRD.js";import{_ as $}from"./_plugin-vue_export-helper-DlAUqK2U.js";const N=x({__name:"WmsBadge",props:{variant:{default:"neutral"},size:{default:"sm"}},setup(o){return(k,q)=>(V(),C("span",{class:I(["wms-badge",[`wms-badge--${o.variant}`,`wms-badge--${o.size}`]])},[E(k.$slots,"default",{},void 0,!0)],2))}}),h=$(N,[["__scopeId","data-v-bf35ff8f"]]);N.__docgenInfo={exportName:"default",displayName:"WmsBadge",description:"",tags:{},props:[{name:"variant",required:!1,type:{name:"union",elements:[{name:'"success"'},{name:'"warning"'},{name:'"error"'},{name:'"info"'},{name:'"neutral"'}]},defaultValue:{func:!1,value:"'neutral'"}},{name:"size",required:!1,type:{name:"union",elements:[{name:'"sm"'},{name:'"md"'}]},defaultValue:{func:!1,value:"'sm'"}}],slots:[{name:"default"}],sourceFiles:["/Users/beabys/go/src/github.com/beabys/wms/customer-ui/src/components/common/WmsBadge.vue"]};const M={title:"Common/WmsBadge",component:h,tags:["autodocs"],argTypes:{variant:{control:"select",options:["success","warning","error","info","neutral"]},size:{control:"select",options:["sm","md"]}}},e={args:{variant:"success",default:"Active"}},a={args:{variant:"warning",default:"Pending"}},r={args:{variant:"error",default:"Blocked"}},s={args:{variant:"info",default:"New"}},n={args:{variant:"neutral",default:"Draft"}},t={args:{variant:"success",size:"md",default:"Verified"}};var c,m,u;e.parameters={...e.parameters,docs:{...(c=e.parameters)==null?void 0:c.docs,source:{originalSource:`{
  args: {
    variant: 'success',
    default: 'Active'
  }
}`,...(u=(m=e.parameters)==null?void 0:m.docs)==null?void 0:u.source}}};var i,d,l;a.parameters={...a.parameters,docs:{...(i=a.parameters)==null?void 0:i.docs,source:{originalSource:`{
  args: {
    variant: 'warning',
    default: 'Pending'
  }
}`,...(l=(d=a.parameters)==null?void 0:d.docs)==null?void 0:l.source}}};var f,p,g;r.parameters={...r.parameters,docs:{...(f=r.parameters)==null?void 0:f.docs,source:{originalSource:`{
  args: {
    variant: 'error',
    default: 'Blocked'
  }
}`,...(g=(p=r.parameters)==null?void 0:p.docs)==null?void 0:g.source}}};var v,_,S;s.parameters={...s.parameters,docs:{...(v=s.parameters)==null?void 0:v.docs,source:{originalSource:`{
  args: {
    variant: 'info',
    default: 'New'
  }
}`,...(S=(_=s.parameters)==null?void 0:_.docs)==null?void 0:S.source}}};var w,B,b;n.parameters={...n.parameters,docs:{...(w=n.parameters)==null?void 0:w.docs,source:{originalSource:`{
  args: {
    variant: 'neutral',
    default: 'Draft'
  }
}`,...(b=(B=n.parameters)==null?void 0:B.docs)==null?void 0:b.source}}};var z,y,W;t.parameters={...t.parameters,docs:{...(z=t.parameters)==null?void 0:z.docs,source:{originalSource:`{
  args: {
    variant: 'success',
    size: 'md',
    default: 'Verified'
  }
}`,...(W=(y=t.parameters)==null?void 0:y.docs)==null?void 0:W.source}}};const P=["Success","Warning","Error","Info","Neutral","MediumSize"];export{r as Error,s as Info,t as MediumSize,n as Neutral,e as Success,a as Warning,P as __namedExportsOrder,M as default};
