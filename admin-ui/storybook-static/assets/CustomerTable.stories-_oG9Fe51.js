import{C as y}from"./CustomerTable-CgGb51Fd.js";import"./vue.esm-bundler-BobL4fDe.js";import"./WmsButton-CrMUKyHj.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";import"./CustomerStatusBadge-BODsa58r.js";import"./WmsBadge-B4x3mWOX.js";const A=[{id:"1",company_name:"Acme Corp",email:"acme@test.com",phone:"",vat_number:"VAT001",status:"pending",created_at:"2024-01-15"},{id:"2",company_name:"Globex Inc",email:"globex@test.com",phone:"",vat_number:"VAT002",status:"pending",created_at:"2024-02-01"}],x=[{id:"1",company_name:"Acme Corp",email:"acme@test.com",phone:"",vat_number:"VAT001",status:"active",created_at:"2024-01-15"},{id:"2",company_name:"Globex Inc",email:"globex@test.com",phone:"",vat_number:"VAT002",status:"pending",created_at:"2024-02-01"},{id:"3",company_name:"Initech",email:"initech@test.com",phone:"",vat_number:"VAT003",status:"rejected",created_at:"2024-03-10"},{id:"4",company_name:"Umbrella Co",email:"umbrella@test.com",phone:"",vat_number:"VAT004",status:"suspended",created_at:"2024-03-15"}],E={title:"Customers/CustomerTable",component:y,tags:["autodocs"]},e={args:{customers:x,loading:!1}},a={args:{customers:[],loading:!0}},s={args:{customers:[],loading:!1}},r={args:{customers:A,loading:!1}},t={args:{customers:x,loading:!1},parameters:{themes:{theme:"dark"}}};var o,n,m;e.parameters={...e.parameters,docs:{...(o=e.parameters)==null?void 0:o.docs,source:{originalSource:`{
  args: {
    customers: mixedCustomers,
    loading: false
  }
}`,...(m=(n=e.parameters)==null?void 0:n.docs)==null?void 0:m.source}}};var c,d,i;a.parameters={...a.parameters,docs:{...(c=a.parameters)==null?void 0:c.docs,source:{originalSource:`{
  args: {
    customers: [],
    loading: true
  }
}`,...(i=(d=a.parameters)==null?void 0:d.docs)==null?void 0:i.source}}};var u,p,l;s.parameters={...s.parameters,docs:{...(u=s.parameters)==null?void 0:u.docs,source:{originalSource:`{
  args: {
    customers: [],
    loading: false
  }
}`,...(l=(p=s.parameters)==null?void 0:p.docs)==null?void 0:l.source}}};var g,_,b;r.parameters={...r.parameters,docs:{...(g=r.parameters)==null?void 0:g.docs,source:{originalSource:`{
  args: {
    customers: pendingCustomers,
    loading: false
  }
}`,...(b=(_=r.parameters)==null?void 0:_.docs)==null?void 0:b.source}}};var h,C,f;t.parameters={...t.parameters,docs:{...(h=t.parameters)==null?void 0:h.docs,source:{originalSource:`{
  args: {
    customers: mixedCustomers,
    loading: false
  },
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(f=(C=t.parameters)==null?void 0:C.docs)==null?void 0:f.source}}};const I=["Default","Loading","Empty","WithPendingCustomers","DarkMode"];export{t as DarkMode,e as Default,s as Empty,a as Loading,r as WithPendingCustomers,I as __namedExportsOrder,E as default};
