import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Task from "../components/Task";
import Worker from "../components/Worker";
import Log from "../components/Log";
import Save from "../components/Save";

const routes = [
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/home',
    name: 'Home',
    meta:{
      title: "首页",
      keepalive: true,
    },
    component: Home,
    children:[
      {
        path: 'tasks',
        component: Task,
      },
      {
        path: 'workers',
        component: Worker,
      },
      {
        path: 'log',
        component: Log,
      },
      {
        path: 'save',
        component: Save,
      }
    ]
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
