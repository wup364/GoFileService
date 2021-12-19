import Vue from 'vue'
import Router from 'vue-router'
import LogIn from '@/components/LogIn'
import Main from '@/components/Main'

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
      component: Main
    }
  ]
})
