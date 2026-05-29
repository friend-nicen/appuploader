package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

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

// UploadIPA is a placeholder for uploading IPA files
func (a *App) UploadIPA(path string) error {
	return fmt.Errorf("Not implemented")
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
	p12Data, err := pkcs12.Encode(rand.Reader, privKey, cert, nil, password)
	if err != nil {
		return "", fmt.Errorf("生成 .p12 失败: %w", err)
	}
	/* 返回 base64 编码的 .p12 内容，前端直接下载 */
	return base64.StdEncoding.EncodeToString(p12Data), nil
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
	p12Data, err := pkcs12.Encode(rand.Reader, privKey, cert, nil, local.Password)
	if err != nil {
		return "", fmt.Errorf("生成 .p12 失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(p12Data), nil
}
