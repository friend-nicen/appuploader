import { createRouter, createWebHashHistory } from 'vue-router'

export const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            name: 'Dashboard',
            meta: { name: 'Dashboard' },
            component: () => import('@/pages/dashboard')
        },
        {
            path: '/auth',
            name: 'Auth',
            meta: { name: 'API Keys' },
            component: () => import('@/pages/auth')
        },
        {
            path: '/apps',
            name: 'Apps',
            meta: { name: 'Apps' },
            component: () => import('@/pages/apps')
        },
        {
            path: '/bundle-ids',
            name: 'BundleIds',
            meta: { name: 'Bundle IDs' },
            component: () => import('@/pages/bundle-ids')
        },
        {
            path: '/certificates',
            name: 'Certificates',
            meta: { name: 'Certificates' },
            component: () => import('@/pages/certificates')
        },
        {
            path: '/profiles',
            name: 'Profiles',
            meta: { name: 'Profiles' },
            component: () => import('@/pages/profiles')
        },
        {
            path: '/devices',
            name: 'Devices',
            meta: { name: 'Devices' },
            component: () => import('@/pages/devices')
        },
        {
            path: '/testflight',
            name: 'TestFlight',
            meta: { name: 'TestFlight' },
            component: () => import('@/pages/testflight')
        }
    ]
})

export default router
