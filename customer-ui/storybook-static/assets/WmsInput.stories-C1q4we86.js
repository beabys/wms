import{W as E}from"./WmsInput-DkBnqTky.js";import"./vue.esm-bundler-ar8GxPRD.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";const V={title:"Common/WmsInput",component:E,tags:["autodocs"],argTypes:{type:{control:"select",options:["text","email","password"]}}},e={args:{label:"Email",placeholder:"you@example.com",type:"email"}},a={args:{label:"Email",modelValue:"user@example.com",type:"email"}},r={args:{label:"Email",modelValue:"invalid",type:"email",error:"Invalid email format"}},s={args:{label:"Email",placeholder:"you@example.com",disabled:!0}},o={args:{label:"Search",placeholder:"Search..."},render:f=>({components:{WmsInput:E},setup:()=>({args:f}),template:'<WmsInput v-bind="args"><template #prefix><span style="color: var(--color-text-muted)">🔍</span></template></WmsInput>'})};var l,t,m;e.parameters={...e.parameters,docs:{...(l=e.parameters)==null?void 0:l.docs,source:{originalSource:`{
  args: {
    label: 'Email',
    placeholder: 'you@example.com',
    type: 'email'
  }
}`,...(m=(t=e.parameters)==null?void 0:t.docs)==null?void 0:m.source}}};var p,n,c;a.parameters={...a.parameters,docs:{...(p=a.parameters)==null?void 0:p.docs,source:{originalSource:`{
  args: {
    label: 'Email',
    modelValue: 'user@example.com',
    type: 'email'
  }
}`,...(c=(n=a.parameters)==null?void 0:n.docs)==null?void 0:c.source}}};var i,d,u;r.parameters={...r.parameters,docs:{...(i=r.parameters)==null?void 0:i.docs,source:{originalSource:`{
  args: {
    label: 'Email',
    modelValue: 'invalid',
    type: 'email',
    error: 'Invalid email format'
  }
}`,...(u=(d=r.parameters)==null?void 0:d.docs)==null?void 0:u.source}}};var g,b,h;s.parameters={...s.parameters,docs:{...(g=s.parameters)==null?void 0:g.docs,source:{originalSource:`{
  args: {
    label: 'Email',
    placeholder: 'you@example.com',
    disabled: true
  }
}`,...(h=(b=s.parameters)==null?void 0:b.docs)==null?void 0:h.source}}};var x,y,W;o.parameters={...o.parameters,docs:{...(x=o.parameters)==null?void 0:x.docs,source:{originalSource:`{
  args: {
    label: 'Search',
    placeholder: 'Search...'
  },
  render: args => ({
    components: {
      WmsInput
    },
    setup: () => ({
      args
    }),
    template: '<WmsInput v-bind="args"><template #prefix><span style="color: var(--color-text-muted)">🔍</span></template></WmsInput>'
  })
}`,...(W=(y=o.parameters)==null?void 0:y.docs)==null?void 0:W.source}}};const D=["Default","WithValue","WithError","Disabled","WithPrefix"];export{e as Default,s as Disabled,r as WithError,o as WithPrefix,a as WithValue,D as __namedExportsOrder,V as default};
