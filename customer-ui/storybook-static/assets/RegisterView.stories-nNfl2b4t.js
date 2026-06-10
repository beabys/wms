import{k as N,E as S,D as W,y as q,l as B,p as i,q as c,z as l,s as C,j as u}from"./vue.esm-bundler-ar8GxPRD.js";import{a as D,b as E,u as M}from"./useAuth-CebID0S7.js";import{W as A}from"./WmsRegisterForm-D9howU6F.js";import{W as F}from"./WmsThemeToggle-n41KqdBx.js";import{_ as I}from"./_plugin-vue_export-helper-DlAUqK2U.js";import"./pinia-Cx913iQk.js";import"./WmsInput-DkBnqTky.js";import"./WmsButton-UxMbN2n3.js";import"./WmsPasswordInput-DBNzoznx.js";const j={class:"register-view"},z={class:"register-view__toggle"},O={class:"register-view__card"},P={class:"register-view__login"},y=N({__name:"RegisterView",setup(G){const b=D(),R=E(),{register:V}=M(),m=u(""),a=u(!1),n=u("");S(()=>{const t=R.query.token;t&&(m.value=t)});async function T(t){a.value=!0,n.value="";try{await V(t),b.push("/login?registered=true")}catch(e){n.value=(e==null?void 0:e.message)||"Registration failed. Please try again."}finally{a.value=!1}}return(t,e)=>{const x=W("router-link");return q(),B("div",j,[i("div",z,[c(F)]),i("div",O,[c(A,{"invite-token":m.value,loading:a.value,error:n.value,onSubmit:T},null,8,["invite-token","loading","error"]),i("p",P,[e[1]||(e[1]=l(" Already have an account? ",-1)),c(x,{to:"/login"},{default:C(()=>[...e[0]||(e[0]=[l("Sign in",-1)])]),_:1})])])])}}}),U=I(y,[["__scopeId","data-v-f2b0c870"]]);y.__docgenInfo={exportName:"default",displayName:"RegisterView",description:"",tags:{},sourceFiles:["/Users/beabys/go/src/github.com/beabys/wms/customer-ui/src/views/RegisterView.vue"]};const ee={title:"Views/RegisterView",component:U,tags:["autodocs"]},s={parameters:{vueRouter:{query:{token:"invite-abc-123"}}}},r={},o={parameters:{themes:{theme:"dark"}}};var d,p,g;s.parameters={...s.parameters,docs:{...(d=s.parameters)==null?void 0:d.docs,source:{originalSource:`{
  parameters: {
    vueRouter: {
      query: {
        token: 'invite-abc-123'
      }
    }
  }
}`,...(g=(p=s.parameters)==null?void 0:p.docs)==null?void 0:g.source}}};var _,v,f;r.parameters={...r.parameters,docs:{...(_=r.parameters)==null?void 0:_.docs,source:{originalSource:"{}",...(f=(v=r.parameters)==null?void 0:v.docs)==null?void 0:f.source}}};var k,h,w;o.parameters={...o.parameters,docs:{...(k=o.parameters)==null?void 0:k.docs,source:{originalSource:`{
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(w=(h=o.parameters)==null?void 0:h.docs)==null?void 0:w.source}}};const te=["WithToken","NoToken","DarkMode"];export{o as DarkMode,r as NoToken,s as WithToken,te as __namedExportsOrder,ee as default};
