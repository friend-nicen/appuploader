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
	db.AutoMigrate(&backend.ApiKey{})
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

// ListCertificates retrieves certificates
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

// RevokeCertificate revokes a Certificate
func (a *App) RevokeCertificate(id string) error {
	client, err := a.getActiveClient()
	if err != nil {
		return err
	}
	return client.Delete(fmt.Sprintf("/v1/certificates/%s", id))
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

// CreateProfile creates a Profile (Stub)
func (a *App) CreateProfile(name, profileType, bundleId, certId string) (string, error) {
	return "", fmt.Errorf("Not implemented yet")
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
// and returns a JSON with the signed certificate, private key, and metadata.
func (a *App) CreateCertificate(name, certificateType string) (string, error) {
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

// ExportP12 combines a PEM certificate and PEM private key into a password-protected .p12 file
// and returns the content as base64.
func (a *App) ExportP12(certPem, keyPem, password string) (string, error) {
	/* 解析 PEM 格式的证书 */
	certBlock, _ := pem.Decode([]byte(certPem))
	if certBlock == nil {
		return "", fmt.Errorf("解析证书 PEM 失败")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
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
