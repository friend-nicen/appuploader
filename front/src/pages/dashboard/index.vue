<template>
  <a-card :bordered="false" class="card" :title="$route.meta.name">
    <div class="space-y-6">
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-8 text-center">
        <h2 class="text-2xl font-bold text-gray-800 mb-4">欢迎使用 AppStore Connect GUI</h2>
        <p class="text-gray-600 mb-8 max-w-2xl mx-auto">
          这是一个用于管理 Apple 开发者账户资源的现代化桌面应用。您可以轻松管理您的 API 密钥、Bundle IDs、证书、描述文件和设备。
        </p>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 text-left mt-8">
          <!-- 快捷入口卡片 -->
          <router-link to="/auth" class="bg-blue-50 hover:bg-blue-100 p-6 rounded-xl border border-blue-100 transition-colors group block">
            <div class="w-12 h-12 bg-blue-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <v-icon name="key" />
            </div>
            <h3 class="font-semibold text-gray-900 mb-2">配置 API 密钥</h3>
            <p class="text-sm text-gray-600">添加和管理您的 App Store Connect API 密钥以进行认证。</p>
          </router-link>

          <router-link to="/bundle-ids" class="bg-indigo-50 hover:bg-indigo-100 p-6 rounded-xl border border-indigo-100 transition-colors group block">
            <div class="w-12 h-12 bg-indigo-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <v-icon name="app" />
            </div>
            <h3 class="font-semibold text-gray-900 mb-2">管理 Bundle IDs</h3>
            <p class="text-sm text-gray-600">查看和管理您应用的所有应用标识符。</p>
          </router-link>

          <router-link to="/certificates" class="bg-purple-50 hover:bg-purple-100 p-6 rounded-xl border border-purple-100 transition-colors group block">
            <div class="w-12 h-12 bg-purple-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <v-icon name="safety-certificate" />
            </div>
            <h3 class="font-semibold text-gray-900 mb-2">查看证书</h3>
            <p class="text-sm text-gray-600">列出您账户下所有可用的开发和发布证书。</p>
          </router-link>

          <router-link to="/profiles" class="bg-green-50 hover:bg-green-100 p-6 rounded-xl border border-green-100 transition-colors group block">
            <div class="w-12 h-12 bg-green-600 text-white rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <v-icon name="profile" />
            </div>
            <h3 class="font-semibold text-gray-900 mb-2">描述文件</h3>
            <p class="text-sm text-gray-600">管理所有的 Provisioning Profiles 以进行构建和分发。</p>
          </router-link>
        </div>
      </div>
    </div>
  </a-card>
</template>

<script setup>
</script>

<style lang="scss" scoped>
.space-y-6 {
  .bg-white {
    padding: 32px;
    text-align: center;
    .grid {
      display: grid;
      grid-template-columns: repeat(1, minmax(0, 1fr));
      gap: 24px;
      text-align: left;
      margin-top: 32px;
      @media (min-width: 768px) {
        grid-template-columns: repeat(2, minmax(0, 1fr));
      }
      @media (min-width: 1024px) {
        grid-template-columns: repeat(4, minmax(0, 1fr));
      }
      .group {
        display: block;
        padding: 24px;
        border-radius: 12px;
        border: 1px solid #e5e7eb;
        transition: background-color 0.3s, border-color 0.3s;
        text-decoration: none;
        .w-12 {
          width: 48px;
          height: 48px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 16px;
          color: white;
          transition: transform 0.3s;
        }
        h3 {
          font-weight: 600;
          color: #111827;
          margin-bottom: 8px;
        }
        p {
          font-size: 14px;
          color: #4b5563;
        }
        &:hover {
          .w-12 {
            transform: scale(1.1);
          }
        }
      }
      .bg-blue-50 {
        background-color: #eff6ff;
        border-color: #dbeafe;
        .w-12 { background-color: #2563eb; }
        &:hover { background-color: #dbeafe; }
      }
      .bg-indigo-50 {
        background-color: #eef2ff;
        border-color: #e0e7ff;
        .w-12 { background-color: #4f46e5; }
        &:hover { background-color: #e0e7ff; }
      }
      .bg-purple-50 {
        background-color: #faf5ff;
        border-color: #f3e8ff;
        .w-12 { background-color: #9333ea; }
        &:hover { background-color: #f3e8ff; }
      }
      .bg-green-50 {
        background-color: #f0fdf4;
        border-color: #dcfce7;
        .w-12 { background-color: #16a34a; }
        &:hover { background-color: #dcfce7; }
      }
    }
  }
}
</style>
