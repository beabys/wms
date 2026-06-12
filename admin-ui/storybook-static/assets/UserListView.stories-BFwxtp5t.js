import{c as f,s as U}from"./pinia-B-feHutF.js";import{k as H,H as O,w as R,q as x,v as P,x as c,y as r,p as i,s as v,E as t,u as l,z as G,r as K,j as Q}from"./vue.esm-bundler-BobL4fDe.js";import{u as X}from"./vue-router-kwq0R5qN.js";import{u as w}from"./users-Dc90cD1P.js";import{U as Y}from"./UserTable-DJsctQ9a.js";import{U as $}from"./UserFilters-CFEvBV_A.js";import{W as ee}from"./WmsCard--e3bIPwJ.js";import{W as V}from"./WmsButton-CrMUKyHj.js";import{_ as se}from"./_plugin-vue_export-helper-DlAUqK2U.js";import"./authClient-DToioyDJ.js";import"./WmsBadge-B4x3mWOX.js";const te={class:"user-list-view"},ae={class:"user-list-view__header"},ne={class:"user-list-view__pagination"},re={class:"user-list-view__page-info"},ie={class:"user-list-view__page-controls"},oe={key:0,class:"user-list-view__error"},J=H({__name:"UserListView",setup(n){const s=X(),e=w(),p=K({});O(()=>{e.fetchUsers()}),R(p,L=>{e.fetchUsers({role:L.role,page:1,page_size:10})},{deep:!0});function F(){s.push("/users/create")}function I(){e.fetchUsers({role:p.value.role,page:e.pagination.page+1,page_size:e.pagination.page_size})}function j(){e.pagination.page>1&&e.fetchUsers({role:p.value.role,page:e.pagination.page-1,page_size:e.pagination.page_size})}const h=Q(()=>Math.max(1,Math.ceil(e.pagination.total_items/e.pagination.page_size)));return(L,a)=>(x(),P("div",te,[c("header",ae,[a[2]||(a[2]=c("h1",{class:"user-list-view__title"},"Users",-1)),r(V,{variant:"primary",size:"md",onClick:F},{default:i(()=>[...a[1]||(a[1]=[v(" Create User ",-1)])]),_:1})]),r(ee,{padding:"md"},{header:i(()=>[r($,{modelValue:p.value,"onUpdate:modelValue":a[0]||(a[0]=q=>p.value=q),loading:t(e).loading},null,8,["modelValue","loading"])]),footer:i(()=>[c("div",ne,[c("span",re," Page "+l(t(e).pagination.page)+" of "+l(h.value)+" ("+l(t(e).pagination.total_items)+" total) ",1),c("div",ie,[r(V,{variant:"secondary",size:"sm",disabled:t(e).pagination.page<=1||t(e).loading,onClick:j},{default:i(()=>[...a[3]||(a[3]=[v(" Previous ",-1)])]),_:1},8,["disabled"]),r(V,{variant:"secondary",size:"sm",disabled:t(e).pagination.page>=h.value||t(e).loading,onClick:I},{default:i(()=>[...a[4]||(a[4]=[v(" Next ",-1)])]),_:1},8,["disabled"])])])]),default:i(()=>[r(Y,{users:t(e).users,loading:t(e).loading},null,8,["users","loading"])]),_:1}),t(e).error?(x(),P("p",oe,l(t(e).error),1)):G("",!0)]))}}),o=se(J,[["__scopeId","data-v-c5c92622"]]);J.__docgenInfo={exportName:"default",displayName:"UserListView",description:"",tags:{},sourceFiles:["/Users/beabys/go/src/github.com/beabys/wms/admin-ui/src/views/users/UserListView.vue"]};const ve={title:"Views/UserListView",component:o,tags:["autodocs"]},m={render:()=>({components:{UserListView:o},template:"<UserListView />"})},d={render:()=>({components:{UserListView:o},setup:()=>{const n=f();U(n);const s=w();return s.loading=!0,{}},template:"<UserListView />"})},u={render:()=>({components:{UserListView:o},setup:()=>{const n=f();return U(n),{}},template:"<UserListView />"})},g={render:()=>({components:{UserListView:o},setup:()=>{const n=f();U(n);const s=w();return s.users=[{id:"1",email:"alice@example.com",name:"Alice Johnson",role:"admin",created_at:"2024-01-15T10:30:00Z",updated_at:""},{id:"2",email:"bob@example.com",name:"Bob Smith",role:"warehouse_staff",created_at:"2024-02-20T14:00:00Z",updated_at:""}],s.pagination={page:1,page_size:10,total_items:25},s.loading=!1,{}},template:"<UserListView />"})},_={render:()=>({components:{UserListView:o},setup:()=>{const n=f();U(n);const s=w();return s.users=[{id:"1",email:"alice@example.com",name:"Alice Johnson",role:"admin",created_at:"2024-01-15T10:30:00Z",updated_at:""}],s.pagination={page:1,page_size:10,total_items:1},s.loading=!1,{}},template:"<UserListView />"}),parameters:{themes:{theme:"dark"}}};var b,z,S;m.parameters={...m.parameters,docs:{...(b=m.parameters)==null?void 0:b.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserListView
    },
    template: '<UserListView />'
  })
}`,...(S=(z=m.parameters)==null?void 0:z.docs)==null?void 0:S.source}}};var k,y,A;d.parameters={...d.parameters,docs:{...(k=d.parameters)==null?void 0:k.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserListView
    },
    setup: () => {
      const pinia = createPinia();
      setActivePinia(pinia);
      const store = useUserStore();
      store.loading = true;
      return {};
    },
    template: '<UserListView />'
  })
}`,...(A=(y=d.parameters)==null?void 0:y.docs)==null?void 0:A.source}}};var C,T,N;u.parameters={...u.parameters,docs:{...(C=u.parameters)==null?void 0:C.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserListView
    },
    setup: () => {
      const pinia = createPinia();
      setActivePinia(pinia);
      return {};
    },
    template: '<UserListView />'
  })
}`,...(N=(T=u.parameters)==null?void 0:T.docs)==null?void 0:N.source}}};var B,W,Z;g.parameters={...g.parameters,docs:{...(B=g.parameters)==null?void 0:B.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserListView
    },
    setup: () => {
      const pinia = createPinia();
      setActivePinia(pinia);
      const store = useUserStore();
      store.users = [{
        id: '1',
        email: 'alice@example.com',
        name: 'Alice Johnson',
        role: 'admin',
        created_at: '2024-01-15T10:30:00Z',
        updated_at: ''
      }, {
        id: '2',
        email: 'bob@example.com',
        name: 'Bob Smith',
        role: 'warehouse_staff',
        created_at: '2024-02-20T14:00:00Z',
        updated_at: ''
      }];
      store.pagination = {
        page: 1,
        page_size: 10,
        total_items: 25
      };
      store.loading = false;
      return {};
    },
    template: '<UserListView />'
  })
}`,...(Z=(W=g.parameters)==null?void 0:W.docs)==null?void 0:Z.source}}};var D,E,M;_.parameters={..._.parameters,docs:{...(D=_.parameters)==null?void 0:D.docs,source:{originalSource:`{
  render: () => ({
    components: {
      UserListView
    },
    setup: () => {
      const pinia = createPinia();
      setActivePinia(pinia);
      const store = useUserStore();
      store.users = [{
        id: '1',
        email: 'alice@example.com',
        name: 'Alice Johnson',
        role: 'admin',
        created_at: '2024-01-15T10:30:00Z',
        updated_at: ''
      }];
      store.pagination = {
        page: 1,
        page_size: 10,
        total_items: 1
      };
      store.loading = false;
      return {};
    },
    template: '<UserListView />'
  }),
  parameters: {
    themes: {
      theme: 'dark'
    }
  }
}`,...(M=(E=_.parameters)==null?void 0:E.docs)==null?void 0:M.source}}};const Ve=["Default","Loading","Empty","WithPagination","DarkMode"];export{_ as DarkMode,m as Default,u as Empty,d as Loading,g as WithPagination,Ve as __namedExportsOrder,ve as default};
