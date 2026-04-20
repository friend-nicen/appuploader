package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"

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

// CreateBundleId creates a Bundle ID (Stub)
func (a *App) CreateBundleId(name, identifier string) (string, error) {
	return "", fmt.Errorf("Not implemented yet")
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
