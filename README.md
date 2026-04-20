# App Store Connect GUI

一款基于 [Wails](https://wails.io) 和 [Vue 3](https://vuejs.org/) 开发的跨平台（macOS / Windows / Linux）App Store Connect 桌面可视化工具。
它通过直接调用 App Store Connect API 的底层机制，实现了本地化、图形化的证书、设备、Bundle ID 和配置文件的管理，提供了现代化的 SaaS Dashboard 体验。

## 特性
- **现代 UI**: 基于 Vue 3 + TailwindCSS 实现的流畅响应式桌面端界面。
- **本地存储**: 使用无 CGO 依赖的纯 Go SQLite 驱动 `github.com/glebarez/sqlite` 本地加密存储 API Keys 等配置信息。
- **多账户管理**: 支持配置多个 Issuer ID、Key ID 和 Private Key，支持无缝切换。
- **核心功能**:
  - **Bundle ID 管理**: 查看现有的 App IDs
  - **证书管理 (Certificates)**: 查看各类证书信息
  - **描述文件管理 (Profiles)**: 浏览 Provisioning Profiles
  - **设备管理 (Devices)**: 查看和管理测试设备

## 技术栈
- **后端 (Go)**: Go 1.21+, Wails v2, GORM, `golang-jwt/jwt/v5`
- **前端 (Web)**: Vue 3 (Composition API), Vue Router, TailwindCSS 3
- **数据库**: SQLite (`github.com/glebarez/sqlite`)

## 开发指南

本项目代码包含详尽的中英文注释，方便进行二次开发和功能扩展。

### 1. 环境准备
- 安装 [Go 1.26+](https://go.dev/doc/install)
- 安装 [Node.js 18+](https://nodejs.org/)
- 安装 Wails CLI:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### 2. 本地开发 (Dev Mode)
开发模式下，Wails 会启动一个本地 Web 服务器，并提供热重载 (Hot Reload) 支持。
```bash
# 确保在项目根目录运行
wails dev
```
前端代码位于 `frontend/` 目录中，所有的 Go 暴露方法在 `frontend/wailsjs/go/main/App.js` 中自动生成，可以直接在 Vue 组件中以 Promise 的方式调用。

### 3. 编译打包 (Build)
编译生产版本时，Wails 会将前端代码打包并嵌入到最终的二进制执行文件中。
```bash
# 编译当前平台的应用
wails build

# 交叉编译 macOS (如果你在 Windows/Linux 上)
wails build -platform darwin/amd64,darwin/arm64

# 交叉编译 Windows
wails build -platform windows/amd64
```
编译成功后，产物会输出到 `build/bin/` 目录下。

## 项目结构
- `/backend`: Go 后端逻辑，包含 API 客户端和 SQLite 数据库模型。
- `/frontend`: Vue 3 前端代码。
  - `src/views`: 各大功能模块的 Vue 页面组件。
  - `src/router`: 页面路由配置。
- `app.go`: Wails 的主生命周期文件，包含了绑定到前端的 Go 方法。
- `main.go`: Wails 应用启动入口。

## License
MIT License
