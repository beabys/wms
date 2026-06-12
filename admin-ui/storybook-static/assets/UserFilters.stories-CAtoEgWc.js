import{r as o}from"./vue.esm-bundler-BobL4fDe.js";import{U as s}from"./UserFilters-CFEvBV_A.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";const h={title:"Users/UserFilters",component:s,tags:["autodocs"]},e={render:()=>({components:{UserFilters:s},setup:()=>({filters:o({})}),template:'<UserFilters v-model="filters" />'})},r={render:()=>({components:{UserFilters:s},setup:()=>({filters:o({role:"admin"})}),template:'<UserFilters v-model="filters" />'})},t={render:()=>({components:{UserFilters:s},setup:()=>({filters:o({})}),template:'<UserFilters v-model="filters" />'}),parameters:{themes:{theme:"dark"}}};var a,i,m;e.parameters={...e.parameters,docs:{...(a=e.parameters)==null?void 0:a.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserFilters
    },
    setup: () => {
      const filters = ref<{
        role?: string;
      }>({});
      return {
        filters
      };
    },
    template: '<UserFilters v-model="filters" />'
  })
}`,...(m=(i=e.parameters)==null?void 0:i.docs)==null?void 0:m.source}}};var n,c,p;r.parameters={...r.parameters,docs:{...(n=r.parameters)==null?void 0:n.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserFilters
    },
    setup: () => {
      const filters = ref<{
        role?: string;
      }>({
        role: 'admin'
      });
      return {
        filters
      };
    },
    template: '<UserFilters v-model="filters" />'
  })
}`,...(p=(c=r.parameters)==null?void 0:c.docs)==null?void 0:p.source}}};var d,f,u;t.parameters={...t.parameters,docs:{...(d=t.parameters)==null?void 0:d.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserFilters
    },
    setup: () => {
      const filters = ref<{
        role?: string;
      }>({});
      return {
        filters
      };
    },
    template: '<UserFilters v-model="filters" />'
  }),
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(u=(f=t.parameters)==null?void 0:f.docs)==null?void 0:u.source}}};const v=["Default","WithRoleSelected","DarkMode"];export{t as DarkMode,e as Default,r as WithRoleSelected,v as __namedExportsOrder,h as default};
