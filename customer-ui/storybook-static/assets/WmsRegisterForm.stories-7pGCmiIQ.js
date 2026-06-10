import{W as b}from"./WmsRegisterForm-D9howU6F.js";import"./vue.esm-bundler-ar8GxPRD.js";import"./WmsInput-DkBnqTky.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";import"./WmsButton-UxMbN2n3.js";import"./WmsPasswordInput-DBNzoznx.js";const F={title:"Auth/WmsRegisterForm",component:b,tags:["autodocs"]},r={args:{inviteToken:"invite-abc-123"}},e={},a={args:{loading:!0}},o={args:{error:"Registration failed. The invite token may be invalid."}},s={play:async({canvasElement:W})=>{const t=W.querySelector("form");t==null||t.requestSubmit()}};var n,i,c;r.parameters={...r.parameters,docs:{...(n=r.parameters)==null?void 0:n.docs,source:{originalSource:`{
  args: {
    inviteToken: 'invite-abc-123'
  }
}`,...(c=(i=r.parameters)==null?void 0:i.docs)==null?void 0:c.source}}};var m,p,d;e.parameters={...e.parameters,docs:{...(m=e.parameters)==null?void 0:m.docs,source:{originalSource:"{}",...(d=(p=e.parameters)==null?void 0:p.docs)==null?void 0:d.source}}};var u,l,g;a.parameters={...a.parameters,docs:{...(u=a.parameters)==null?void 0:u.docs,source:{originalSource:`{
  args: {
    loading: true
  }
}`,...(g=(l=a.parameters)==null?void 0:l.docs)==null?void 0:g.source}}};var v,f,S;o.parameters={...o.parameters,docs:{...(v=o.parameters)==null?void 0:v.docs,source:{originalSource:`{
  args: {
    error: 'Registration failed. The invite token may be invalid.'
  }
}`,...(S=(f=o.parameters)==null?void 0:f.docs)==null?void 0:S.source}}};var y,h,E;s.parameters={...s.parameters,docs:{...(y=s.parameters)==null?void 0:y.docs,source:{originalSource:`{
  play: async ({
    canvasElement
  }) => {
    const form = canvasElement.querySelector('form');
    form?.requestSubmit();
  }
}`,...(E=(h=s.parameters)==null?void 0:h.docs)==null?void 0:E.source}}};const L=["Default","WithNoToken","Loading","WithError","ValidationErrors"];export{r as Default,a as Loading,s as ValidationErrors,o as WithError,e as WithNoToken,L as __namedExportsOrder,F as default};
