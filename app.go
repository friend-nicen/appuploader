package main

import (
	"archive/zip"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"howett.net/plist"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
	"software.sslmate.com/src/go-pkcs12"

	"appstore-connect-gui/backend"
)

// App struct
type App struct {
	ctx context.Context
	db  *gorm.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Setup SQLite Database
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "AppStoreConnectGUI")
	os.MkdirAll(appDir, 0755)

	dbPath := filepath.Join(appDir, "app_data.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database:", err)
	}

	// Auto Migrate
	db.AutoMigrate(&backend.ApiKey{}, &backend.LocalCertificate{})
	a.db = db
}

// GetApiKeys retrieves all API keys from the database
func (a *App) GetApiKeys() []backend.ApiKey {
	var keys []backend.ApiKey
	a.db.Find(&keys)
	// Mask private keys for security when sending to frontend
	for i := range keys {
		keys[i].PrivateKey = "*** MASKED ***"
	}
	return keys
}

// AddApiKey adds a new API key
func (a *App) AddApiKey(name, issuerID, keyID, privateKey string) error {
	key := backend.ApiKey{
		Name:       name,
		IssuerID:   issuerID,
		KeyID:      keyID,
		PrivateKey: privateKey,
		IsActive:   false,
	}
	
	// If it's the first key, make it active
	var count int64
	a.db.Model(&backend.ApiKey{}).Count(&count)
	if count == 0 {
		key.IsActive = true
	}

	return a.db.Create(&key).Error
}

// SetCurrentKey sets the active API key
func (a *App) SetCurrentKey(id uint) error {
	// Deactivate all
	a.db.Model(&backend.ApiKey{}).Where("1 = 1").Update("is_active", false)
	// Activate selected
	return a.db.Model(&backend.ApiKey{}).Where("id = ?", id).Update("is_active", true).Error
}

// UpdateApiKey updates an existing API key
func (a *App) UpdateApiKey(id uint, name, issuerID, keyID, privateKey string) error {
	updates := map[string]interface{}{
		"name":      name,
		"issuer_id": issuerID,
		"key_id":    keyID,
	}
	if privateKey != "" {
		updates["private_key"] = privateKey
	}
	return a.db.Model(&backend.ApiKey{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteApiKey deletes an API key
func (a *App) DeleteApiKey(id uint) error {
	return a.db.Delete(&backend.ApiKey{}, id).Error
}

func (a *App) getActiveClient() (*backend.AppStoreClient, error) {
	var activeKey backend.ApiKey
	if err := a.db.Where("is_active = ?", true).First(&activeKey).Error; err != nil {
		return nil, fmt.Errorf("no active API key found")
	}
	return backend.NewAppStoreClient(activeKey.IssuerID, activeKey.KeyID, activeKey.PrivateKey), nil
}

// ListBundleIds retrieves App IDs
func (a *App) ListBundleIds() (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get("/v1/bundleIds?limit=200")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// ListCertificates retrieves certificates from Apple
func (a *App) ListCertificates() (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get("/v1/certificates?limit=200")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// ListCertificatesWithLocal retrieves certificates from Apple and overrides the name
// with the locally stored name (the name set during creation).
func (a *App) ListCertificatesWithLocal() (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get("/v1/certificates?limit=200")
	if err != nil {
		return "", err
	}
	/* 获取本地存储的证书信息 */
	var locals []backend.LocalCertificate
	a.db.Find(&locals)
	type localInfo struct {
		Name     string
		Password string
	}
	localMap := make(map[string]localInfo, len(locals))
	for _, l := range locals {
		localMap[l.AppleCertID] = localInfo{Name: l.Name, Password: l.Password}
	}
	if len(localMap) == 0 {
		return string(res), nil
	}
	/* 解析 Apple 响应，注入本地数据 */
	var response map[string]interface{}
	if err := json.Unmarshal(res, &response); err != nil {
		return "", fmt.Errorf("解析 Apple 响应失败: %w", err)
	}
	data, ok := response["data"].([]interface{})
	if !ok {
		return string(res), nil
	}
	for _, item := range data {
		cert, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := cert["id"].(string)
		if id == "" {
			continue
		}
		info, exists := localMap[id]
		if !exists {
			continue
		}
		attrs, ok := cert["attributes"].(map[string]interface{})
		if !ok {
			continue
		}
		attrs["name"] = info.Name
		attrs["password"] = info.Password
	}
	result, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("JSON marshal 失败: %w", err)
	}
	return string(result), nil
}

// ListProfiles retrieves provisioning profiles
func (a *App) ListProfiles() (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get("/v1/profiles?limit=200")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// ListDevices retrieves devices
func (a *App) ListDevices() (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get("/v1/devices?limit=200")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// TestAuth tests if the active key is valid
func (a *App) TestAuth() (bool, error) {
	_, err := a.ListBundleIds()
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListApps retrieves Apps
func (a *App) ListApps() (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get("/v1/apps?limit=200")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// DeleteBundleId deletes a Bundle ID
func (a *App) DeleteBundleId(id string) error {
	client, err := a.getActiveClient()
	if err != nil {
		return err
	}
	return client.Delete(fmt.Sprintf("/v1/bundleIds/%s", id))
}

// RevokeCertificate revokes a Certificate on Apple and cleans up local record
func (a *App) RevokeCertificate(id string) error {
	client, err := a.getActiveClient()
	if err != nil {
		return err
	}
	if err := client.Delete(fmt.Sprintf("/v1/certificates/%s", id)); err != nil {
		return err
	}
	/* 删除本地存储的证书数据 */
	a.db.Where("apple_cert_id = ?", id).Delete(&backend.LocalCertificate{})
	return nil
}

// DownloadCertificate retrieves certificate content for download
func (a *App) DownloadCertificate(id string) (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get(fmt.Sprintf("/v1/certificates/%s", id))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// DeleteProfile deletes a Profile
func (a *App) DeleteProfile(id string) error {
	client, err := a.getActiveClient()
	if err != nil {
		return err
	}
	return client.Delete(fmt.Sprintf("/v1/profiles/%s", id))
}

// DeleteDevice deletes a Device
func (a *App) DeleteDevice(id string) error {
	client, err := a.getActiveClient()
	if err != nil {
		return err
	}
	return client.Delete(fmt.Sprintf("/v1/devices/%s", id))
}

// CreateBundleId creates a Bundle ID via App Store Connect API
func (a *App) CreateBundleId(name, identifier, platform string) (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}

	body := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "bundleIds",
			"attributes": map[string]interface{}{
				"name":       name,
				"identifier": identifier,
				"platform":   platform,
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("JSON marshal failed: %w", err)
	}

	res, err := client.Post("/v1/bundleIds", payload)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// CreateDevice creates a Device (Stub)
func (a *App) CreateDevice(name, udid string) (string, error) {
	return "", fmt.Errorf("Not implemented yet")
}

// CreateProfile creates a Profile via App Store Connect API
func (a *App) CreateProfile(name, profileType, bundleId, certIds, deviceIds string) (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	/* 解析证书 ID 列表 */
	var certData []map[string]string
	if err := json.Unmarshal([]byte(certIds), &certData); err != nil {
		return "", fmt.Errorf("解析证书 ID 列表失败: %w", err)
	}
	/* 构造 relationships */
	relationships := map[string]interface{}{
		"bundleId": map[string]interface{}{
			"data": map[string]string{
				"id":   bundleId,
				"type": "bundleIds",
			},
		},
		"certificates": map[string]interface{}{
			"data": certData,
		},
	}
	/* 如果传入了设备 ID，添加 devices relationship */
	if deviceIds != "" && deviceIds != "[]" {
		var deviceData []map[string]string
		if err := json.Unmarshal([]byte(deviceIds), &deviceData); err != nil {
			return "", fmt.Errorf("解析设备 ID 列表失败: %w", err)
		}
		relationships["devices"] = map[string]interface{}{
			"data": deviceData,
		}
	}
	/* 构造请求体 */
	body := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "profiles",
			"attributes": map[string]string{
				"name":        name,
				"profileType": profileType,
			},
			"relationships": relationships,
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("JSON marshal failed: %w", err)
	}
	res, err := client.Post("/v1/profiles", payload)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// DownloadProfile retrieves profile content from App Store Connect
func (a *App) DownloadProfile(id string) (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	res, err := client.Get(fmt.Sprintf("/v1/profiles/%s", id))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// SelectFile uses Wails runtime dialog to select an .ipa file
func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 IPA 文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "IPA Files (*.ipa)",
				Pattern:     "*.ipa",
			},
		},
	})
}

// extractIPAInfo reads an IPA file (ZIP) and extracts version, build number, and platform
// from the embedded Info.plist. Supports both XML and binary plist formats.
type ipaInfo struct {
	VersionString string `json:"versionString"`
	BuildNumber   string `json:"buildNumber"`
	Platform      string `json:"platform"`
	BundleID      string `json:"bundleId,omitempty"`
}

func extractIPAInfo(path string) *ipaInfo {
	info := &ipaInfo{}

	reader, err := zip.OpenReader(path)
	if err != nil {
		return info
	}
	defer reader.Close()

	/* 在 Payload/*.app/Info.plist 中查找 */
	var plistFile *zip.File
	for _, f := range reader.File {
		if strings.HasPrefix(f.Name, "Payload/") && strings.HasSuffix(f.Name, ".app/Info.plist") {
			plistFile = f
			break
		}
	}
	if plistFile == nil {
		return info
	}

	rc, err := plistFile.Open()
	if err != nil {
		return info
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return info
	}

	/* 解析 plist（同时支持 XML 和二进制格式） */
	var plistData map[string]interface{}
	if _, err := plist.Unmarshal(content, &plistData); err != nil {
		return info
	}

	if v, ok := plistData["CFBundleShortVersionString"].(string); ok {
		info.VersionString = v
	}

	/* CFBundleVersion 可能是字符串也可能是数字 */
	if v, ok := plistData["CFBundleVersion"].(string); ok {
		info.BuildNumber = v
	} else if v, ok := plistData["CFBundleVersion"].(uint64); ok {
		info.BuildNumber = fmt.Sprintf("%d", v)
	} else if v, ok := plistData["CFBundleVersion"].(int); ok {
		info.BuildNumber = fmt.Sprintf("%d", v)
	} else if v, ok := plistData["CFBundleVersion"].(float64); ok {
		info.BuildNumber = fmt.Sprintf("%.0f", v)
	}

	/* 读取 Bundle ID */
	if v, ok := plistData["CFBundleIdentifier"].(string); ok {
		info.BundleID = v
	}

	platformRaw, _ := plistData["DTPlatformName"].(string)
	switch platformRaw {
	case "iphoneos":
		info.Platform = "IOS"
	case "macosx":
		info.Platform = "MAC_OS"
	case "appletvos":
		info.Platform = "TV_OS"
	case "xros":
		info.Platform = "VISION_OS"
	default:
		info.Platform = "IOS"
	}

	return info
}

// GetIPAInfo extracts version, build number, and platform from an IPA file
// and returns them as JSON for the frontend to display.
func (a *App) GetIPAInfo(path string) (string, error) {
	info := extractIPAInfo(path)
	result, err := json.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %w", err)
	}
	return string(result), nil
}

// UploadIPA uploads an IPA file to App Store Connect for TestFlight distribution.
func (a *App) UploadIPA(appID, path, versionString, buildNumber, platform string) error {
	client, err := a.getActiveClient()
	if err != nil {
		return err
	}

	/* 获取文件信息 */
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}
	fileSize := fileStat.Size()

	/* 计算文件 MD5 */
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return fmt.Errorf("计算 MD5 失败: %w", err)
	}
	/* Step 1: 创建 build upload 预约 */
	uploadBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "buildUploads",
			"attributes": map[string]interface{}{
				"cfBundleShortVersionString": versionString,
				"cfBundleVersion":             buildNumber,
				"platform":                   platform,
			},
			"relationships": map[string]interface{}{
				"app": map[string]interface{}{
					"data": map[string]string{
						"id":   appID,
						"type": "apps",
					},
				},
			},
		},
	}
	payload, _ := json.Marshal(uploadBody)
	res, err := client.Post("/v1/buildUploads", payload)
	if err != nil {
		return fmt.Errorf("创建上传任务失败: %w", err)
	}

	/* 解析 buildUpload 响应，获取 buildUpload ID */
	var buildUploadResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res, &buildUploadResp); err != nil {
		return fmt.Errorf("解析 buildUpload 响应失败: %w", err)
	}
	buildUploadID := buildUploadResp.Data.ID
	if buildUploadID == "" {
		return fmt.Errorf("创建上传任务失败: 未获取到 upload ID")
	}

	/* 提取文件名 */
	fileName := filepath.Base(path)

	/* Step 2: 创建 buildUploadFile，获取上传地址 */
	fileBody := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "buildUploadFiles",
			"attributes": map[string]interface{}{
				"fileName":  fileName,
				"fileSize":  fileSize,
				"assetType": "ASSET",
				"uti":       "com.apple.ipa",
			},
			"relationships": map[string]interface{}{
				"buildUpload": map[string]interface{}{
					"data": map[string]string{
						"id":   buildUploadID,
						"type": "buildUploads",
					},
				},
			},
		},
	}
	filePayload, _ := json.Marshal(fileBody)
	res, err = client.Post("/v1/buildUploadFiles", filePayload)
	if err != nil {
		return fmt.Errorf("创建文件上传任务失败: %w", err)
	}

	/* 解析 buildUploadFile 响应 */
	var fileResp struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				UploadOperations []struct {
					RequestURL    string `json:"url"`
					RequestMethod string `json:"method"`
					RequestHeaders []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"requestHeaders"`
					MD5    string `json:"md5"`
					Offset int64  `json:"offset"`
					Length int64  `json:"length"`
				} `json:"uploadOperations"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res, &fileResp); err != nil {
		return fmt.Errorf("解析文件上传响应失败: %w", err)
	}

	uploadFileID := fileResp.Data.ID
	operations := fileResp.Data.Attributes.UploadOperations
	if len(operations) == 0 {
		return fmt.Errorf("未获取到上传地址")
	}

	/* Step 3: 逐个上传分片到预签名 URL */
	file.Seek(0, 0)

	for i, op := range operations {
		chunk := make([]byte, op.Length)
		n, err := file.ReadAt(chunk, op.Offset)
		if err != nil && err != io.EOF {
			return fmt.Errorf("读取文件分片 %d 失败: %w", i, err)
		}
		chunk = chunk[:n]

		/* 构建请求头 */
		headers := make(map[string]string)
		for _, h := range op.RequestHeaders {
			headers[h.Name] = h.Value
		}
		if op.MD5 != "" {
			if _, exists := headers["Content-MD5"]; !exists {
				headers["Content-MD5"] = op.MD5
			}
		}

		if err := client.PutUpload(op.RequestURL, chunk, headers); err != nil {
			return fmt.Errorf("上传分片 %d 失败: %w", i, err)
		}
	}

	/* Step 4: 标记 build upload file 为已上传 */
	if uploadFileID != "" {
		patchBody := map[string]interface{}{
			"data": map[string]interface{}{
				"type": "buildUploadFiles",
				"id":   uploadFileID,
				"attributes": map[string]interface{}{
					"uploaded": true,
				},
			},
		}
		patchPayload, _ := json.Marshal(patchBody)
		if _, err := client.Patch("/v1/buildUploadFiles/"+uploadFileID, patchPayload); err != nil {
			return fmt.Errorf("标记上传完成失败: %w", err)
		}
	}

	return nil
}

/* 生成 RSA 2048 密钥对和 PEM 格式的 CSR（Apple 要求使用 RSA 2048） */
func generateKeyPair(commonName string) (*rsa.PrivateKey, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("生成密钥对失败: %w", err)
	}
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: commonName,
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("生成 CSR 失败: %w", err)
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
	return priv, csrPEM, nil
}

// CreateCertificate generates a key pair + CSR, submits it to App Store Connect,
// stores the result locally, and returns a JSON with the signed certificate, private key, and metadata.
func (a *App) CreateCertificate(name, certificateType, password string) (string, error) {
	client, err := a.getActiveClient()
	if err != nil {
		return "", err
	}
	/* 生成密钥对和 CSR */
	priv, csrPEM, err := generateKeyPair(name)
	if err != nil {
		return "", err
	}
	/* 构造 POST 请求体提交 CSR 到 Apple */
	body := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "certificates",
			"attributes": map[string]interface{}{
				"csrContent":      string(csrPEM),
				"certificateType": certificateType,
			},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("JSON marshal failed: %w", err)
	}
	res, err := client.Post("/v1/certificates", payload)
	if err != nil {
		return "", err
	}
	/* 解析 Apple 返回的证书内容 */
	var certResp struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				CertificateContent string `json:"certificateContent"`
				DisplayName        string `json:"displayName"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res, &certResp); err != nil {
		return "", fmt.Errorf("解析 Apple 响应失败: %w", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("序列化私钥失败: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	/* 存储到本地数据库 */
	localCert := backend.LocalCertificate{
		AppleCertID:     certResp.Data.ID,
		Name:            name,
		CertificateType: certificateType,
		CertPem:         certResp.Data.Attributes.CertificateContent,
		KeyPem:          string(keyPEM),
		CsrPem:          string(csrPEM),
		Password:        password,
	}
	if err := a.db.Create(&localCert).Error; err != nil {
		return "", fmt.Errorf("存储证书到本地失败: %w", err)
	}
	/* 构造返回结果 */
	result := map[string]interface{}{
		"certPem":       certResp.Data.Attributes.CertificateContent,
		"keyPem":        string(keyPEM),
		"name":          name,
		"certificateId": certResp.Data.ID,
		"displayName":   certResp.Data.Attributes.DisplayName,
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("JSON marshal failed: %w", err)
	}
	return string(resultJSON), nil
}

// ExportP12 combines a certificate and private key into a password-protected .p12 file
// and returns the content as base64.
// The certContent is base64-encoded DER (Apple's certificateContent format).
func (a *App) ExportP12(certContent, keyPem, password string) (string, error) {
	/* 解析证书内容（base64 DER 格式） */
	certDER, err := base64.StdEncoding.DecodeString(certContent)
	if err != nil {
		return "", fmt.Errorf("解码证书内容失败: %w", err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return "", fmt.Errorf("解析证书失败: %w", err)
	}
	/* 解析 PEM 格式的私钥 */
	keyBlock, _ := pem.Decode([]byte(keyPem))
	if keyBlock == nil {
		return "", fmt.Errorf("解析私钥 PEM 失败")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}
	/* 将证书和私钥打包为 PKCS#12 格式 */
	p12Data, err := pkcs12.LegacyDES.Encode(privKey, cert, nil, password)
	if err != nil {
		return "", fmt.Errorf("生成 .p12 失败: %w", err)
	}
	/* 返回 base64 编码的 .p12 内容，前端直接下载 */
	return base64.StdEncoding.EncodeToString(p12Data), nil
}

// UpdateCertificatePassword updates the stored password for a local certificate.
func (a *App) UpdateCertificatePassword(appleCertID, newPassword string) (string, error) {
	result := a.db.Model(&backend.LocalCertificate{}).
		Where("apple_cert_id = ?", appleCertID).
		Update("password", newPassword)
	if result.Error != nil {
		return "", fmt.Errorf("更新密码失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return "", fmt.Errorf("未找到本地证书记录")
	}
	return `{"code":1,"errMsg":"密码修改成功"}`, nil
}

// ExportP12WithOpenSSL exports a .p12 file using the system's OpenSSL command.
// This provides better compatibility with tools like HBuilder.
func (a *App) ExportP12WithOpenSSL(appleCertID string) (string, error) {
	var local backend.LocalCertificate
	if err := a.db.Where("apple_cert_id = ?", appleCertID).First(&local).Error; err != nil {
		return "", fmt.Errorf("未找到本地证书记录: %w", err)
	}
	/* 创建临时目录存放证书和私钥 */
	tmpDir, err := os.MkdirTemp("", "p12export-*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")
	p12File := filepath.Join(tmpDir, "output.p12")

	/* 写入证书文件（Apple 的 cert_pem 是 base64 DER，需要先转成 PEM） */
	certDER, err := base64.StdEncoding.DecodeString(local.CertPem)
	if err != nil {
		return "", fmt.Errorf("解码证书失败: %w", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err := os.WriteFile(certFile, certPEM, 0644); err != nil {
		return "", fmt.Errorf("写入证书文件失败: %w", err)
	}
	/* 写入私钥文件 */
	if err := os.WriteFile(keyFile, []byte(local.KeyPem), 0644); err != nil {
		return "", fmt.Errorf("写入私钥文件失败: %w", err)
	}
	/* 调用 openssl 命令生成 .p12 */
	cmd := exec.Command("openssl", "pkcs12", "-export",
		"-in", certFile,
		"-inkey", keyFile,
		"-out", p12File,
		"-passout", "pass:"+local.Password,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("openssl 执行失败: %w\n输出: %s", err, string(output))
	}
	/* 读取生成的 .p12 文件 */
	p12Data, err := os.ReadFile(p12File)
	if err != nil {
		return "", fmt.Errorf("读取 .p12 文件失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(p12Data), nil
}

// ExportPEM exports the private key and certificate as PEM text content.
// Returns a JSON with "keyPem" and "certPem" fields.
func (a *App) ExportPEM(appleCertID string) (string, error) {
	var local backend.LocalCertificate
	if err := a.db.Where("apple_cert_id = ?", appleCertID).First(&local).Error; err != nil {
		return "", fmt.Errorf("未找到本地证书记录: %w", err)
	}
	/* 将 Apple 的 base64 DER 证书转为 PEM */
	certDER, err := base64.StdEncoding.DecodeString(local.CertPem)
	if err != nil {
		return "", fmt.Errorf("解码证书失败: %w", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	result := map[string]string{
		"keyPem":  local.KeyPem,
		"certPem": string(certPEM),
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("JSON marshal 失败: %w", err)
	}
	return string(resultJSON), nil
}

// VerifyP12Password checks if the stored password can correctly encode and decode a .p12 file.
// This is a debug function to verify the password pipeline.
func (a *App) VerifyP12Password(appleCertID string) (string, error) {
	var local backend.LocalCertificate
	if err := a.db.Where("apple_cert_id = ?", appleCertID).First(&local).Error; err != nil {
		return "", fmt.Errorf("未找到本地证书记录: %w", err)
	}
	/* 解析证书 */
	certDER, err := base64.StdEncoding.DecodeString(local.CertPem)
	if err != nil {
		return "", fmt.Errorf("解码证书失败: %w", err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return "", fmt.Errorf("解析证书失败: %w", err)
	}
	/* 解析私钥 */
	keyBlock, _ := pem.Decode([]byte(local.KeyPem))
	if keyBlock == nil {
		return "", fmt.Errorf("解析私钥 PEM 失败")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}
	/* 用存储的密码编码 .p12 */
	p12Data, err := pkcs12.LegacyDES.Encode(privKey, cert, nil, local.Password)
	if err != nil {
		return "", fmt.Errorf("编码 .p12 失败: %w", err)
	}
	/* 尝试用同样的密码解码，验证密码正确性 */
	_, _, err = pkcs12.Decode(p12Data, local.Password)
	if err != nil {
		return fmt.Sprintf("密码校验失败，存储的密码为: [%s]，长度: %d，解码错误: %v", local.Password, len(local.Password), err), nil
	}
	return fmt.Sprintf("密码校验成功，存储的密码为: [%s]，长度: %d", local.Password, len(local.Password)), nil
}

// SaveTextFile shows a save file dialog and writes text content to the selected path.
func (a *App) SaveTextFile(content, defaultName, filterPattern string) (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:          "保存文件",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{
				DisplayName: defaultName,
				Pattern:     filterPattern,
			},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", fmt.Errorf("用户取消了保存")
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}
	return path, nil
}

// SaveBase64File shows a save file dialog, decodes base64 content, and writes to the selected path.
func (a *App) SaveBase64File(content, defaultName, filterPattern string) (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:          "保存文件",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{
				DisplayName: defaultName,
				Pattern:     filterPattern,
			},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", fmt.Errorf("用户取消了保存")
	}
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", fmt.Errorf("解码文件内容失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}
	return path, nil
}

/* 导出数据结构 */
type exportData struct {
	ApiKeys      []backend.ApiKey            `json:"apiKeys"`
	Certificates []backend.LocalCertificate  `json:"certificates"`
}

// ExportLocalData exports all locally stored data (API keys and certificates) as JSON.
func (a *App) ExportLocalData() (string, error) {
	var data exportData
	a.db.Find(&data.ApiKeys)
	a.db.Find(&data.Certificates)
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("导出数据序列化失败: %w", err)
	}
	return string(result), nil
}

// ImportLocalData imports data from a JSON string into the local database.
// Existing records are skipped; new records are inserted.
func (a *App) ImportLocalData(dataStr string) error {
	var data exportData
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return fmt.Errorf("导入数据解析失败: %w", err)
	}
	/* 导入 API 密钥 */
	for i := range data.ApiKeys {
		data.ApiKeys[i].ID = 0
		if err := a.db.Create(&data.ApiKeys[i]).Error; err != nil {
			return fmt.Errorf("导入 API 密钥失败: %w", err)
		}
	}
	/* 导入本地证书（按 AppleCertID 去重） */
	for i := range data.Certificates {
		var count int64
		a.db.Model(&backend.LocalCertificate{}).
			Where("apple_cert_id = ?", data.Certificates[i].AppleCertID).
			Count(&count)
		if count > 0 {
			continue
		}
		data.Certificates[i].ID = 0
		if err := a.db.Create(&data.Certificates[i]).Error; err != nil {
			return fmt.Errorf("导入证书失败: %w", err)
		}
	}
	return nil
}

// ExportCSR retrieves the stored CSR for a certificate by its Apple certificate ID.
func (a *App) ExportCSR(appleCertID string) (string, error) {
	var local backend.LocalCertificate
	if err := a.db.Where("apple_cert_id = ?", appleCertID).First(&local).Error; err != nil {
		return "", fmt.Errorf("未找到本地证书记录: %w", err)
	}
	return local.CsrPem, nil
}

// ExportLocalP12 generates a .p12 file from a locally stored certificate using its stored password.
func (a *App) ExportLocalP12(appleCertID string) (string, error) {
	var local backend.LocalCertificate
	if err := a.db.Where("apple_cert_id = ?", appleCertID).First(&local).Error; err != nil {
		return "", fmt.Errorf("未找到本地证书记录: %w", err)
	}
	/* 解析证书内容（Apple 的 certificateContent 是 base64 DER 格式） */
	certDER, err := base64.StdEncoding.DecodeString(local.CertPem)
	if err != nil {
		return "", fmt.Errorf("解码证书内容失败: %w", err)
	}
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return "", fmt.Errorf("解析证书失败: %w", err)
	}
	/* 解析 PEM 格式的私钥 */
	keyBlock, _ := pem.Decode([]byte(local.KeyPem))
	if keyBlock == nil {
		return "", fmt.Errorf("解析私钥 PEM 失败")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}
	/* 使用存储的密码生成 .p12 */
	p12Data, err := pkcs12.LegacyDES.Encode(privKey, cert, nil, local.Password)
	if err != nil {
		return "", fmt.Errorf("生成 .p12 失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(p12Data), nil
}
