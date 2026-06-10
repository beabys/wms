import{c as F,s as A}from"./pinia-B-feHutF.js";import{k as I,v as M,x as p,y as d,p as l,q,s as O,r as w}from"./vue.esm-bundler-BobL4fDe.js";import{u as R}from"./vue-router-kwq0R5qN.js";import{u as B}from"./users-Dc90cD1P.js";import{U as T}from"./UserForm-Cbmn80OA.js";import{W as j}from"./WmsCard--e3bIPwJ.js";import{W as z}from"./WmsButton-CrMUKyHj.js";import{_ as G}from"./_plugin-vue_export-helper-DlAUqK2U.js";import"./authClient-DToioyDJ.js";import"./WmsInput-Br2ujW9m.js";const H={class:"create-user-view"},J={class:"create-user-view__actions"},D=I({__name:"CreateUserView",setup(f){const s=R(),E=B(),c=w({email:"",password:"",name:"",role:""}),m=w(!1),u=w("");async function N(){m.value=!0,u.value="";try{await E.createUser(c.value),s.push("/users")}catch(t){u.value=(t==null?void 0:t.message)||"Failed to create user"}finally{m.value=!1}}function P(){s.push("/users")}return(t,e)=>(q(),M("div",H,[e[2]||(e[2]=p("header",{class:"create-user-view__header"},[p("h1",{class:"create-user-view__title"},"Create User")],-1)),d(j,{padding:"md"},{footer:l(()=>[p("div",J,[d(z,{variant:"ghost",onClick:P},{default:l(()=>[...e[1]||(e[1]=[O(" Cancel ",-1)])]),_:1})])]),default:l(()=>[d(T,{modelValue:c.value,"onUpdate:modelValue":e[0]||(e[0]=W=>c.value=W),loading:m.value,error:u.value,onSubmit:N},null,8,["modelValue","loading","error"])]),_:1})]))}}),r=G(D,[["__scopeId","data-v-cd117943"]]);D.__docgenInfo={exportName:"default",displayName:"CreateUserView",description:"",tags:{},sourceFiles:["/Users/beabys/go/src/github.com/beabys/wms/admin-ui/src/views/users/CreateUserView.vue"]};const te={title:"Views/CreateUserView",component:r,tags:["autodocs"]},a={render:()=>({components:{CreateUserView:r},template:"<CreateUserView />"})},o={render:()=>({components:{CreateUserView:r},template:"<CreateUserView />"})},n={render:()=>({components:{CreateUserView:r},setup:()=>{const f=F();A(f);const s=B();return s.loading=!0,{}},template:"<CreateUserView />"})},i={render:()=>({components:{CreateUserView:r},template:"<CreateUserView />"}),parameters:{themes:{theme:"dark"}}};var V,U,v;a.parameters={...a.parameters,docs:{...(V=a.parameters)==null?void 0:V.docs,source:{originalSource:`{
  render: () => ({
    components: {
      CreateUserView
    },
    template: '<CreateUserView />'
  })
}`,...(v=(U=a.parameters)==null?void 0:U.docs)==null?void 0:v.source}}};var C,_,g;o.parameters={...o.parameters,docs:{...(C=o.parameters)==null?void 0:C.docs,source:{originalSource:`{
  render: () => ({
    components: {
      CreateUserView
    },
    template: '<CreateUserView />'
  })
}`,...(g=(_=o.parameters)==null?void 0:_.docs)==null?void 0:g.source}}};var h,b,k;n.parameters={...n.parameters,docs:{...(h=n.parameters)==null?void 0:h.docs,source:{originalSource:`{
  render: () => ({
    components: {
      CreateUserView
    },
    setup: () => {
      const pinia = createPinia();
      setActivePinia(pinia);
      const store = useUserStore();
      store.loading = true;
      return {};
    },
    template: '<CreateUserView />'
  })
}`,...(k=(b=n.parameters)==null?void 0:b.docs)==null?void 0:k.source}}};var S,x,y;i.parameters={...i.parameters,docs:{...(S=i.parameters)==null?void 0:S.docs,source:{originalSource:`{
  render: () => ({
    components: {
      CreateUserView
    },
    template: '<CreateUserView />'
  }),
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(y=(x=i.parameters)==null?void 0:x.docs)==null?void 0:y.source}}};const ae=["Default","ValidationErrors","Submitting","DarkMode"];export{i as DarkMode,a as Default,n as Submitting,o as ValidationErrors,ae as __namedExportsOrder,te as default};
