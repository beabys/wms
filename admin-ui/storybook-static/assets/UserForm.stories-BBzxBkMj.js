import{r}from"./vue.esm-bundler-BobL4fDe.js";import{U as e}from"./UserForm-Cbmn80OA.js";import"./WmsInput-Br2ujW9m.js";import"./_plugin-vue_export-helper-DlAUqK2U.js";import"./WmsButton-CrMUKyHj.js";const h={email:"",password:"",name:"",role:""},x={title:"Users/UserForm",component:e,tags:["autodocs"]},t={render:()=>({components:{UserForm:e},setup:()=>({data:r({...h})}),template:'<UserForm v-model="data" />'})},s={render:()=>({components:{UserForm:e},setup:()=>({data:r({email:"bad",password:"short",name:"",role:""})}),template:'<UserForm v-model="data" />'})},o={render:()=>({components:{UserForm:e},setup:()=>({data:r({email:"test@test.com",password:"password123",name:"Test User",role:"admin"})}),template:'<UserForm v-model="data" :loading="true" />'})},m={render:()=>({components:{UserForm:e},setup:()=>({data:r({email:"test@test.com",password:"password123",name:"Test User",role:"admin"})}),template:'<UserForm v-model="data" error="Email already exists" />'})},d={render:()=>({components:{UserForm:e},setup:()=>({data:r({...h})}),template:'<UserForm v-model="data" />'}),parameters:{themes:{theme:"dark"}}};var n,p,c;t.parameters={...t.parameters,docs:{...(n=t.parameters)==null?void 0:n.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserForm
    },
    setup: () => {
      const data = ref<CreateUserRequest>({
        ...defaultData
      });
      return {
        data
      };
    },
    template: '<UserForm v-model="data" />'
  })
}`,...(c=(p=t.parameters)==null?void 0:p.docs)==null?void 0:c.source}}};var l,u,i;s.parameters={...s.parameters,docs:{...(l=s.parameters)==null?void 0:l.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserForm
    },
    setup: () => {
      const data = ref<CreateUserRequest>({
        email: 'bad',
        password: 'short',
        name: '',
        role: ''
      });
      return {
        data
      };
    },
    template: '<UserForm v-model="data" />'
  })
}`,...(i=(u=s.parameters)==null?void 0:u.docs)==null?void 0:i.source}}};var U,F,f;o.parameters={...o.parameters,docs:{...(U=o.parameters)==null?void 0:U.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserForm
    },
    setup: () => {
      const data = ref<CreateUserRequest>({
        email: 'test@test.com',
        password: 'password123',
        name: 'Test User',
        role: 'admin'
      });
      return {
        data
      };
    },
    template: '<UserForm v-model="data" :loading="true" />'
  })
}`,...(f=(F=o.parameters)==null?void 0:F.docs)==null?void 0:f.source}}};var v,w,g;m.parameters={...m.parameters,docs:{...(v=m.parameters)==null?void 0:v.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserForm
    },
    setup: () => {
      const data = ref<CreateUserRequest>({
        email: 'test@test.com',
        password: 'password123',
        name: 'Test User',
        role: 'admin'
      });
      return {
        data
      };
    },
    template: '<UserForm v-model="data" error="Email already exists" />'
  })
}`,...(g=(w=m.parameters)==null?void 0:w.docs)==null?void 0:g.source}}};var S,D,E;d.parameters={...d.parameters,docs:{...(S=d.parameters)==null?void 0:S.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserForm
    },
    setup: () => {
      const data = ref<CreateUserRequest>({
        ...defaultData
      });
      return {
        data
      };
    },
    template: '<UserForm v-model="data" />'
  }),
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(E=(D=d.parameters)==null?void 0:D.docs)==null?void 0:E.source}}};const T=["Default","ValidationErrors","Submitting","ServerError","DarkMode"];export{d as DarkMode,t as Default,m as ServerError,o as Submitting,s as ValidationErrors,T as __namedExportsOrder,x as default};
