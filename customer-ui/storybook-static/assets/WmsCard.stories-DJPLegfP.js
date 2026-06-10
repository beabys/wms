import{k as S,l as n,C as m,v as p,p as $,B as k,y as c}from"./vue.esm-bundler-ar8GxPRD.js";import{_ as x}from"./_plugin-vue_export-helper-DlAUqK2U.js";const B={key:0,class:"wms-card__header"},H={class:"wms-card__body"},L={key:1,class:"wms-card__footer"},w=S({__name:"WmsCard",props:{padding:{default:"md"}},setup(e){return(a,N)=>(c(),n("div",{class:k(["wms-card",`wms-card--${e.padding}`])},[a.$slots.header?(c(),n("div",B,[m(a.$slots,"header",{},void 0,!0)])):p("",!0),$("div",H,[m(a.$slots,"default",{},void 0,!0)]),a.$slots.footer?(c(),n("div",L,[m(a.$slots,"footer",{},void 0,!0)])):p("",!0)],2))}}),s=x(w,[["__scopeId","data-v-be4e0c3a"]]);w.__docgenInfo={exportName:"default",displayName:"WmsCard",description:"",tags:{},props:[{name:"padding",required:!1,type:{name:"union",elements:[{name:'"none"'},{name:'"sm"'},{name:'"md"'},{name:'"lg"'}]},defaultValue:{func:!1,value:"'md'"}}],slots:[{name:"header"},{name:"default"},{name:"footer"}],sourceFiles:["/Users/beabys/go/src/github.com/beabys/wms/customer-ui/src/components/common/WmsCard.vue"]};const T={title:"Common/WmsCard",component:s,tags:["autodocs"],argTypes:{padding:{control:"select",options:["none","sm","md","lg"]}}},r={args:{padding:"md"},render:e=>({components:{WmsCard:s},setup:()=>({args:e}),template:'<WmsCard v-bind="args">This is a basic card with some content.</WmsCard>'})},t={args:{padding:"md"},render:e=>({components:{WmsCard:s},setup:()=>({args:e}),template:`
      <WmsCard v-bind="args">
        <template #header><strong style="font-size: 1.125rem">Card Header</strong></template>
        Main card content goes here.
        <template #footer><span style="color: var(--color-text-muted)">Last updated: today</span></template>
      </WmsCard>
    `})},o={args:{padding:"sm"},render:e=>({components:{WmsCard:s},setup:()=>({args:e}),template:'<WmsCard v-bind="args">Compact card content.</WmsCard>'})},d={args:{padding:"lg"},render:e=>({components:{WmsCard:s},setup:()=>({args:e}),template:'<WmsCard v-bind="args">Spacious card with extra padding.</WmsCard>'})};var i,l,g;r.parameters={...r.parameters,docs:{...(i=r.parameters)==null?void 0:i.docs,source:{originalSource:`{
  args: {
    padding: 'md'
  },
  render: args => ({
    components: {
      WmsCard
    },
    setup: () => ({
      args
    }),
    template: '<WmsCard v-bind="args">This is a basic card with some content.</WmsCard>'
  })
}`,...(g=(l=r.parameters)==null?void 0:l.docs)==null?void 0:g.source}}};var u,C,W;t.parameters={...t.parameters,docs:{...(u=t.parameters)==null?void 0:u.docs,source:{originalSource:`{
  args: {
    padding: 'md'
  },
  render: args => ({
    components: {
      WmsCard
    },
    setup: () => ({
      args
    }),
    template: \`
      <WmsCard v-bind="args">
        <template #header><strong style="font-size: 1.125rem">Card Header</strong></template>
        Main card content goes here.
        <template #footer><span style="color: var(--color-text-muted)">Last updated: today</span></template>
      </WmsCard>
    \`
  })
}`,...(W=(C=t.parameters)==null?void 0:C.docs)==null?void 0:W.source}}};var f,_,h;o.parameters={...o.parameters,docs:{...(f=o.parameters)==null?void 0:f.docs,source:{originalSource:`{
  args: {
    padding: 'sm'
  },
  render: args => ({
    components: {
      WmsCard
    },
    setup: () => ({
      args
    }),
    template: '<WmsCard v-bind="args">Compact card content.</WmsCard>'
  })
}`,...(h=(_=o.parameters)==null?void 0:_.docs)==null?void 0:h.source}}};var v,b,y;d.parameters={...d.parameters,docs:{...(v=d.parameters)==null?void 0:v.docs,source:{originalSource:`{
  args: {
    padding: 'lg'
  },
  render: args => ({
    components: {
      WmsCard
    },
    setup: () => ({
      args
    }),
    template: '<WmsCard v-bind="args">Spacious card with extra padding.</WmsCard>'
  })
}`,...(y=(b=d.parameters)==null?void 0:b.docs)==null?void 0:y.source}}};const V=["Default","WithHeaderAndFooter","Compact","LargePadding"];export{o as Compact,r as Default,d as LargePadding,t as WithHeaderAndFooter,V as __namedExportsOrder,T as default};
