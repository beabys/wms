import{W as E}from"./WmsLoginForm-CwTr9Vmp.js";import"./vue.esm-bundler-BobL4fDe.js";import"./WmsInput-Br2ujW9m.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";import"./WmsButton-CrMUKyHj.js";import"./WmsPasswordInput-SWa_wJtI.js";const w={title:"Auth/WmsLoginForm",component:E,tags:["autodocs"]},r={args:{}},a={args:{loading:!0}},e={args:{error:"Invalid email or password. Please try again."}},o={args:{},play:async({canvasElement:f})=>{const s=f.querySelector("form");s==null||s.requestSubmit()}};var t,n,m;r.parameters={...r.parameters,docs:{...(t=r.parameters)==null?void 0:t.docs,source:{originalSource:`{
  args: {}
}`,...(m=(n=r.parameters)==null?void 0:n.docs)==null?void 0:m.source}}};var c,i,p;a.parameters={...a.parameters,docs:{...(c=a.parameters)==null?void 0:c.docs,source:{originalSource:`{
  args: {
    loading: true
  }
}`,...(p=(i=a.parameters)==null?void 0:i.docs)==null?void 0:p.source}}};var d,u,l;e.parameters={...e.parameters,docs:{...(d=e.parameters)==null?void 0:d.docs,source:{originalSource:`{
  args: {
    error: 'Invalid email or password. Please try again.'
  }
}`,...(l=(u=e.parameters)==null?void 0:u.docs)==null?void 0:l.source}}};var g,y,S;o.parameters={...o.parameters,docs:{...(g=o.parameters)==null?void 0:g.docs,source:{originalSource:`{
  args: {},
  play: async ({
    canvasElement
  }) => {
    // Submit empty form to trigger validation
    const form = canvasElement.querySelector('form');
    form?.requestSubmit();
  }
}`,...(S=(y=o.parameters)==null?void 0:y.docs)==null?void 0:S.source}}};const x=["Default","Loading","WithError","ValidationErrors"];export{r as Default,a as Loading,o as ValidationErrors,e as WithError,x as __namedExportsOrder,w as default};
