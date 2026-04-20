# App Store Connect API 调用说明

本项目基于 Apple 官方提供的 App Store Connect API 进行开发。以下是当前已集成并使用的核心 API 列表及其说明：

## 1. 鉴权机制 (Authentication)
所有的接口请求都需要使用 JWT (JSON Web Token) 进行鉴权。
- **算法**: ES256
- **Header**: 包含 `alg` (ES256), `kid` (你的 Key ID), `typ` (JWT)
- **Payload**: 包含 `iss` (你的 Issuer ID), `iat` (签发时间), `exp` (过期时间，最长20分钟，本项目默认设置10分钟以防时钟偏移), `aud` ("appstoreconnect-v1")
- **文档参考**: [Generating Tokens for API Requests](https://developer.apple.com/documentation/appstoreconnectapi/generating_tokens_for_api_requests)

## 2. 已使用的 API 接口

### Bundle ID 管理
- **接口路径**: `GET https://api.appstoreconnect.apple.com/v1/bundleIds?limit=200`
- **功能描述**: 获取与当前账号关联的所有应用标识符 (Bundle IDs)。
- **文档参考**: [List Bundle IDs](https://developer.apple.com/documentation/appstoreconnectapi/list_bundle_ids)

### 证书管理 (Certificates)
- **接口路径**: `GET https://api.appstoreconnect.apple.com/v1/certificates?limit=200`
- **功能描述**: 获取所有证书的列表，包括开发证书和发布证书。
- **文档参考**: [List Certificates](https://developer.apple.com/documentation/appstoreconnectapi/list_and_download_certificates)

### 描述文件管理 (Provisioning Profiles)
- **接口路径**: `GET https://api.appstoreconnect.apple.com/v1/profiles?limit=200`
- **功能描述**: 获取所有配置文件（描述文件）的列表。
- **文档参考**: [List Profiles](https://developer.apple.com/documentation/appstoreconnectapi/list_and_download_profiles)

### 设备管理 (Devices)
- **接口路径**: `GET https://api.appstoreconnect.apple.com/v1/devices?limit=200`
- **功能描述**: 获取已注册在开发者账号下的所有测试设备列表（包括 iOS, macOS 等设备）。
- **文档参考**: [List Devices](https://developer.apple.com/documentation/appstoreconnectapi/list_registered_devices)

## 3. 未来可扩展的 API
基于当前的 `backend/client.go` 封装，您可以很容易地通过添加 POST/DELETE 请求来扩展以下功能：
- **注册新设备**: `POST /v1/devices`
- **创建新证书**: `POST /v1/certificates`
- **创建描述文件**: `POST /v1/profiles`
- **TestFlight 管理**: `GET /v1/builds`, `GET /v1/betaTesters` 等。