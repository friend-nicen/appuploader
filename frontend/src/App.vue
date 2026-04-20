<template>
  <a-config-provider :locale="zhCN">
    <a-layout style="min-height: 100vh">
      <!-- 左侧侧边栏 (Fixed Left Sidebar) -->
      <a-layout-sider
        theme="dark"
        :style="{ overflow: 'auto', height: '100vh', position: 'fixed', left: 0, top: 0, bottom: 0 }"
        :width="256"
      >
        <div class="h-16 flex items-center justify-center border-b border-gray-800">
          <span class="text-xl font-bold text-white tracking-wider">AppStore Connect</span>
        </div>
        
        <a-menu
          v-model:selectedKeys="selectedKeys"
          theme="dark"
          mode="inline"
          @click="handleMenuClick"
          class="mt-4 border-r-0"
        >
          <a-menu-item key="/">
            <template #icon>
              <DashboardOutlined />
            </template>
            Dashboard
          </a-menu-item>
          
          <a-menu-item key="/auth">
            <template #icon>
              <KeyOutlined />
            </template>
            API Keys
          </a-menu-item>

          <a-menu-item key="/apps">
            <template #icon>
              <AppstoreAddOutlined />
            </template>
            Apps
          </a-menu-item>
          
          <a-menu-item key="/bundle-ids">
            <template #icon>
              <AppstoreOutlined />
            </template>
            Bundle IDs
          </a-menu-item>
          
          <a-menu-item key="/certificates">
            <template #icon>
              <SafetyCertificateOutlined />
            </template>
            Certificates
          </a-menu-item>
          
          <a-menu-item key="/profiles">
            <template #icon>
              <ProfileOutlined />
            </template>
            Profiles
          </a-menu-item>
          
          <a-menu-item key="/devices">
            <template #icon>
              <MobileOutlined />
            </template>
            Devices
          </a-menu-item>

          <a-menu-item key="/testflight">
            <template #icon>
              <CloudUploadOutlined />
            </template>
            TestFlight
          </a-menu-item>
        </a-menu>
        
        <!-- 侧边栏底部 -->
        <div class="absolute bottom-0 w-full p-4 border-t border-gray-800 bg-[#001529]">
          <div class="text-xs text-gray-500 text-center">
            &copy; {{ new Date().getFullYear() }} ASC GUI
          </div>
        </div>
      </a-layout-sider>

      <!-- 右侧内容区 -->
      <a-layout :style="{ marginLeft: '256px', transition: 'all 0.2s', minHeight: '100vh', display: 'flex', flexDirection: 'column' }">
        <!-- 顶部导航栏 -->
        <a-layout-header style="background: #fff; padding: 0 24px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #f0f0f0; height: 64px;">
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-800 m-0">{{ currentRouteName }}</h1>
          </div>
          <div class="flex items-center space-x-4">
            <KeySelector />
          </div>
        </a-layout-header>

        <!-- 主视图区域 -->
        <a-layout-content style="margin: 24px; padding: 24px; background: #fff; min-height: 280px; border-radius: 8px; flex: 1; overflow: auto;">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-config-provider>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import KeySelector from '@/components/KeySelector.vue'

// Icons from Ant Design Vue
import {
  DashboardOutlined,
  KeyOutlined,
  AppstoreOutlined,
  AppstoreAddOutlined,
  SafetyCertificateOutlined,
  ProfileOutlined,
  MobileOutlined,
  CloudUploadOutlined
} from '@ant-design/icons-vue'

const route = useRoute()
const router = useRouter()

// 选中的菜单项，默认根据当前路由路径设置
const selectedKeys = ref([route.path])

// 监听路由变化，更新侧边栏选中的菜单项
watch(
  () => route.path,
  (newPath) => {
    selectedKeys.value = [newPath]
  },
  { immediate: true }
)

const handleMenuClick = ({ key }) => {
  router.push(key)
}

// 菜单项配置，用于计算当前页面的标题
const menuItems = [
  { name: 'Dashboard', path: '/' },
  { name: 'API Keys', path: '/auth' },
  { name: 'Apps', path: '/apps' },
  { name: 'Bundle IDs', path: '/bundle-ids' },
  { name: 'Certificates', path: '/certificates' },
  { name: 'Profiles', path: '/profiles' },
  { name: 'Devices', path: '/devices' },
  { name: 'TestFlight', path: '/testflight' }
]

// 计算当前路由名称，用于顶部栏显示
const currentRouteName = computed(() => {
  const match = menuItems.find(item => item.path === route.path)
  return match ? match.name : route.name
})
</script>

<style>
/* 视图切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
