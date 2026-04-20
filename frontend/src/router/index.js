import { createRouter, createWebHashHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Auth from '../views/Auth.vue'
import BundleIds from '../views/BundleIds.vue'
import Certificates from '../views/Certificates.vue'
import Profiles from '../views/Profiles.vue'
import Devices from '../views/Devices.vue'
import Apps from '../views/Apps.vue'
import TestFlight from '../views/TestFlight.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/auth',
    name: 'Auth',
    component: Auth
  },
  {
    path: '/apps',
    name: 'Apps',
    component: Apps
  },
  {
    path: '/bundle-ids',
    name: 'BundleIds',
    component: BundleIds
  },
  {
    path: '/certificates',
    name: 'Certificates',
    component: Certificates
  },
  {
    path: '/profiles',
    name: 'Profiles',
    component: Profiles
  },
  {
    path: '/devices',
    name: 'Devices',
    component: Devices
  },
  {
    path: '/testflight',
    name: 'TestFlight',
    component: TestFlight
  }
]

// 使用 hash 模式，因为 Wails 在桌面端运行，没有服务器处理 history 模式的路由
const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
