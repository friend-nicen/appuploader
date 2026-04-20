<template>
  <div class="space-y-6">
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-8 text-center">
      <h2 class="text-2xl font-bold text-gray-800 mb-4">欢迎使用 AppStore Connect GUI</h2>
      <p class="text-gray-600 mb-8 max-w-2xl mx-auto">
        这是一个用于管理 Apple 开发者账户资源的现代化桌面应用。您可以轻松管理您的 API 密钥、Bundle IDs、证书、描述文件和设备。
      </p>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 text-left mt-8">
        <!-- 快捷入口卡片 -->
        <router-link to="/auth" class="bg-blue-50 hover:bg-blue-100 p-6 rounded-xl border border-blue-100 transition-colors group">
          <div class="w-12 h-12 bg-blue-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" /></svg>
          </div>
          <h3 class="font-semibold text-gray-900 mb-2">配置 API 密钥</h3>
          <p class="text-sm text-gray-600">添加和管理您的 App Store Connect API 密钥以进行认证。</p>
        </router-link>

        <router-link to="/bundle-ids" class="bg-indigo-50 hover:bg-indigo-100 p-6 rounded-xl border border-indigo-100 transition-colors group">
          <div class="w-12 h-12 bg-indigo-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" /></svg>
          </div>
          <h3 class="font-semibold text-gray-900 mb-2">管理 Bundle IDs</h3>
          <p class="text-sm text-gray-600">查看和管理您应用的所有应用标识符。</p>
        </router-link>

        <router-link to="/certificates" class="bg-purple-50 hover:bg-purple-100 p-6 rounded-xl border border-purple-100 transition-colors group">
          <div class="w-12 h-12 bg-purple-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" /></svg>
          </div>
          <h3 class="font-semibold text-gray-900 mb-2">查看证书</h3>
          <p class="text-sm text-gray-600">列出您账户下所有可用的开发和发布证书。</p>
        </router-link>

        <router-link to="/profiles" class="bg-green-50 hover:bg-green-100 p-6 rounded-xl border border-green-100 transition-colors group">
          <div class="w-12 h-12 bg-green-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" /></svg>
          </div>
          <h3 class="font-semibold text-gray-900 mb-2">描述文件</h3>
          <p class="text-sm text-gray-600">管理所有的 Provisioning Profiles 以进行构建和分发。</p>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
// Dashboard 组件，用于展示欢迎页面和快捷入口
</script>
