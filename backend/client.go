package backend

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppStoreClient struct {
	IssuerID   string
	KeyID      string
	PrivateKey string
}

func NewAppStoreClient(issuerID, keyID, privateKey string) *AppStoreClient {
	return &AppStoreClient{
		IssuerID:   issuerID,
		KeyID:      keyID,
		PrivateKey: privateKey,
	}
}

func (c *AppStoreClient) generateToken() (string, error) {
	pk := c.PrivateKey
	// 清理可能由于复制粘贴导致的格式问题
	pk = strings.TrimSpace(pk)
	pk = strings.ReplaceAll(pk, "\\n", "\n")
	pk = strings.ReplaceAll(pk, "\\r", "")

	block, _ := pem.Decode([]byte(pk))
	if block == nil {
		return "", fmt.Errorf("解析包含私钥的 PEM 块失败 (请检查格式是否为 -----BEGIN PRIVATE KEY----- 开头)")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %v", err)
	}

	ecdsaKey, ok := privKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("该私钥不是 ECDSA 私钥")
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": c.IssuerID,
		"iat": now.Add(-1 * time.Minute).Unix(), // 减去 1 分钟以防时钟偏移
		"exp": now.Add(10 * time.Minute).Unix(), // 有效期 10 分钟 (Apple 限制最大 20 分钟)
		"aud": "appstoreconnect-v1",
	})
	token.Header["kid"] = c.KeyID

	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *AppStoreClient) Get(endpoint string) ([]byte, error) {
	token, err := c.generateToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.appstoreconnect.apple.com"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	return body, nil
}

func (c *AppStoreClient) Post(endpoint string, payload []byte) ([]byte, error) {
	token, err := c.generateToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.appstoreconnect.apple.com"+endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	return body, nil
}

func (c *AppStoreClient) Patch(endpoint string, payload []byte) ([]byte, error) {
	token, err := c.generateToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", "https://api.appstoreconnect.apple.com"+endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	return body, nil
}

func (c *AppStoreClient) Delete(endpoint string) error {
	token, err := c.generateToken()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", "https://api.appstoreconnect.apple.com"+endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error: %s", string(body))
	}

	return nil
}

// PutUpload uploads raw file data to a presigned URL (used for IPA upload chunks).
// The url is a presigned S3 URL; the headers map contains required upload headers.
func (c *AppStoreClient) PutUpload(url string, data []byte, headers map[string]string) error {
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("创建上传请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("上传分片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("上传分片错误 (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
