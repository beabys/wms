import{U as b}from"./UserTable-DJsctQ9a.js";import"./vue.esm-bundler-BobL4fDe.js";import"./WmsBadge-B4x3mWOX.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";const t=[{id:"1",email:"alice@example.com",name:"Alice Johnson",role:"admin",created_at:"2024-01-15T10:30:00Z",updated_at:"2024-01-15T10:30:00Z"},{id:"2",email:"bob@example.com",name:"Bob Smith",role:"warehouse_staff",created_at:"2024-02-20T14:00:00Z",updated_at:"2024-02-20T14:00:00Z"},{id:"3",email:"carol@example.com",name:"Carol Davis",role:"billing_manager",created_at:"2024-03-10T09:15:00Z",updated_at:"2024-03-10T09:15:00Z"}],E={title:"Users/UserTable",component:b,tags:["autodocs"]},e={args:{users:t,loading:!1}},a={args:{users:[],loading:!0}},s={args:{users:[],loading:!1}},r={args:{users:[t[0]],loading:!1}},o={args:{users:t,loading:!1},parameters:{themes:{theme:"dark"}}};var n,m,c;e.parameters={...e.parameters,docs:{...(n=e.parameters)==null?void 0:n.docs,source:{originalSource:`{
  args: {
    users: sampleUsers,
    loading: false
  }
}`,...(c=(m=e.parameters)==null?void 0:m.docs)==null?void 0:c.source}}};var d,l,i;a.parameters={...a.parameters,docs:{...(d=a.parameters)==null?void 0:d.docs,source:{originalSource:`{
  args: {
    users: [],
    loading: true
  }
}`,...(i=(l=a.parameters)==null?void 0:l.docs)==null?void 0:i.source}}};var p,u,g;s.parameters={...s.parameters,docs:{...(p=s.parameters)==null?void 0:p.docs,source:{originalSource:`{
  args: {
    users: [],
    loading: false
  }
}`,...(g=(u=s.parameters)==null?void 0:u.docs)==null?void 0:g.source}}};var f,_,h;r.parameters={...r.parameters,docs:{...(f=r.parameters)==null?void 0:f.docs,source:{originalSource:`{
  args: {
    users: [sampleUsers[0]],
    loading: false
  }
}`,...(h=(_=r.parameters)==null?void 0:_.docs)==null?void 0:h.source}}};var T,U,D;o.parameters={...o.parameters,docs:{...(T=o.parameters)==null?void 0:T.docs,source:{originalSource:`{
  args: {
    users: sampleUsers,
    loading: false
  },
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(D=(U=o.parameters)==null?void 0:U.docs)==null?void 0:D.source}}};const y=["Default","Loading","Empty","WithData","DarkMode"];export{o as DarkMode,e as Default,s as Empty,a as Loading,r as WithData,y as __namedExportsOrder,E as default};
