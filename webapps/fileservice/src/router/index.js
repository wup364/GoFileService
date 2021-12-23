import Vue from 'vue'
import Router from 'vue-router'
import LogIn from '@/pages/LogIn'
import Main from '@/pages/Main'
import FileList from "@/pages/FileList.vue";
import SysSetting from "@/pages/SysSetting.vue";
Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'LogIn',
      component: LogIn
    },
    {
      path: '/main',
      name: 'Main',
      component: Main,
      children: [
        { name: 'files', path: '', component: FileList },
        { name: 'files', path: 'files', component: FileList },
        { name: 'settings', path: 'settings', component: SysSetting }
      ]
    }
  ]
})
