import {createRouter, createWebHashHistory} from 'vue-router'

export const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            redirect: "/Dashboard"
        },
        {
            path: '/Dashboard',
            name: 'Dashboard',
            meta: {name: 'Dashboard', displayName: '仪表盘'},
            component: () => import('@/pages/dashboard')
        },
        {
            path: '/auth',
            name: 'Auth',
            meta: {name: 'API Keys', displayName: 'API 密钥'},
            component: () => import('@/pages/auth')
        },
        {
            path: '/apps',
            name: 'Apps',
            meta: {name: 'Apps', displayName: '应用'},
            component: () => import('@/pages/apps')
        },
        {
            path: '/bundle-ids',
            name: 'BundleIds',
            meta: {name: 'Bundle IDs', displayName: 'Bundle ID'},
            component: () => import('@/pages/bundle-ids')
        },
        {
            path: '/certificates',
            name: 'Certificates',
            meta: {name: 'Certificates', displayName: '证书'},
            component: () => import('@/pages/certificates')
        },
        {
            path: '/profiles',
            name: 'Profiles',
            meta: {name: 'Profiles', displayName: '描述文件'},
            component: () => import('@/pages/profiles')
        },
        {
            path: '/devices',
            name: 'Devices',
            meta: {name: 'Devices', displayName: '设备'},
            component: () => import('@/pages/devices')
        },
        {
            path: '/testflight',
            name: 'TestFlight',
            meta: {name: 'TestFlight', displayName: 'TestFlight'},
            component: () => import('@/pages/testflight')
        }
    ]
})

export default router
