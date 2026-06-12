import{W as a}from"./WmsCard--e3bIPwJ.js";import"./vue.esm-bundler-BobL4fDe.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";const y={title:"Common/WmsCard",component:a,tags:["autodocs"],argTypes:{padding:{control:"select",options:["none","sm","md","lg"]}}},r={args:{padding:"md"},render:e=>({components:{WmsCard:a},setup:()=>({args:e}),template:'<WmsCard v-bind="args">This is a basic card with some content.</WmsCard>'})},s={args:{padding:"md"},render:e=>({components:{WmsCard:a},setup:()=>({args:e}),template:`
      <WmsCard v-bind="args">
        <template #header><strong style="font-size: 1.125rem">Card Header</strong></template>
        Main card content goes here.
        <template #footer><span style="color: var(--color-text-muted)">Last updated: today</span></template>
      </WmsCard>
    `})},t={args:{padding:"sm"},render:e=>({components:{WmsCard:a},setup:()=>({args:e}),template:'<WmsCard v-bind="args">Compact card content.</WmsCard>'})},d={args:{padding:"lg"},render:e=>({components:{WmsCard:a},setup:()=>({args:e}),template:'<WmsCard v-bind="args">Spacious card with extra padding.</WmsCard>'})};var n,o,m;r.parameters={...r.parameters,docs:{...(n=r.parameters)==null?void 0:n.docs,source:{originalSource:`{
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
}`,...(m=(o=r.parameters)==null?void 0:o.docs)==null?void 0:m.source}}};var p,c,g;s.parameters={...s.parameters,docs:{...(p=s.parameters)==null?void 0:p.docs,source:{originalSource:`{
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
}`,...(g=(c=s.parameters)==null?void 0:c.docs)==null?void 0:g.source}}};var i,l,u;t.parameters={...t.parameters,docs:{...(i=t.parameters)==null?void 0:i.docs,source:{originalSource:`{
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
}`,...(u=(l=t.parameters)==null?void 0:l.docs)==null?void 0:u.source}}};var C,W,h;d.parameters={...d.parameters,docs:{...(C=d.parameters)==null?void 0:C.docs,source:{originalSource:`{
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
}`,...(h=(W=d.parameters)==null?void 0:W.docs)==null?void 0:h.source}}};const x=["Default","WithHeaderAndFooter","Compact","LargePadding"];export{t as Compact,r as Default,d as LargePadding,s as WithHeaderAndFooter,x as __namedExportsOrder,y as default};
