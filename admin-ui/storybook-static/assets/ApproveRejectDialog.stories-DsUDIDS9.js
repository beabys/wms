import{A as j}from"./ApproveRejectDialog-CPBnmKqH.js";import"./vue.esm-bundler-BobL4fDe.js";import"./WmsCard--e3bIPwJ.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";import"./WmsButton-CrMUKyHj.js";const a={id:"1",company_name:"Acme Corp",email:"acme@test.com",phone:"123456789",vat_number:"VAT123",status:"pending",created_at:"2024-01-01"},D={title:"Customers/ApproveRejectDialog",component:j,tags:["autodocs"],argTypes:{action:{control:"select",options:["approve","reject","suspend"]}}},e={args:{visible:!0,customer:a,action:"approve"}},r={args:{visible:!0,customer:a,action:"reject"}},s={args:{visible:!0,customer:a,action:"suspend"}},t={args:{visible:!0,customer:a,action:"reject"},parameters:{themes:{theme:"dark"}}};var o,n,c;e.parameters={...e.parameters,docs:{...(o=e.parameters)==null?void 0:o.docs,source:{originalSource:`{
  args: {
    visible: true,
    customer,
    action: 'approve'
  }
}`,...(c=(n=e.parameters)==null?void 0:n.docs)==null?void 0:c.source}}};var p,i,m;r.parameters={...r.parameters,docs:{...(p=r.parameters)==null?void 0:p.docs,source:{originalSource:`{
  args: {
    visible: true,
    customer,
    action: 'reject'
  }
}`,...(m=(i=r.parameters)==null?void 0:i.docs)==null?void 0:m.source}}};var u,d,l;s.parameters={...s.parameters,docs:{...(u=s.parameters)==null?void 0:u.docs,source:{originalSource:`{
  args: {
    visible: true,
    customer,
    action: 'suspend'
  }
}`,...(l=(d=s.parameters)==null?void 0:d.docs)==null?void 0:l.source}}};var g,v,b;t.parameters={...t.parameters,docs:{...(g=t.parameters)==null?void 0:g.docs,source:{originalSource:`{
  args: {
    visible: true,
    customer,
    action: 'reject'
  },
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(b=(v=t.parameters)==null?void 0:v.docs)==null?void 0:b.source}}};const R=["Approve","Reject","Suspend","DarkMode"];export{e as Approve,t as DarkMode,r as Reject,s as Suspend,R as __namedExportsOrder,D as default};
