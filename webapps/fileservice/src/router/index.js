import Vue from 'vue'
import Router from 'vue-router'
import LogIn from '@/pages/LogIn'
import Main from '@/pages/Main'
import FileList from "@/pages/FileList.vue";
import SysSetting from "@/pages/SysSetting.vue";
import Preview from '@/pages/Preview'
import Preview4Audio from '@/components/preview/audio'
import Preview4Picture from '@/components/preview/picture'
import Preview4Video from '@/components/preview/video'

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
    },
    {
      path: '/preview',
      name: 'Preview',
      component: Preview,
      children: [
        { name: 'audio', path: 'audio', component: Preview4Audio },
        { name: 'picture', path: 'picture', component: Preview4Picture },
        { name: 'video', path: 'video', component: Preview4Video },
      ]
    }
  ]
})
