package nacos

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	authUrl = "/v1/auth/login"
	tokenExpireBuffer = 300 // token过期前5分钟自动刷新
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrAuthFailed   = errors.New("authentication failed")
)

// getCacheDir 获取缓存目录
func getCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".nacosctl")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", err
	}

	return cacheDir, nil
}

// getCacheFilePath 获取token缓存文件路径
func getCacheFilePath(addr string) (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	// 使用地址的hash作为文件名，避免特殊字符问题
	h := md5.New()
	h.Write([]byte(addr))
	hash := hex.EncodeToString(h.Sum(nil))

	return filepath.Join(cacheDir, "token_"+hash+".json"), nil
}

// loadToken 从缓存加载token
func loadToken(addr string) (*TokenCache, error) {
	cacheFile, err := getCacheFilePath(addr)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // 没有缓存文件
		}
		return nil, err
	}

	var token TokenCache
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// saveToken 保存token到缓存
func saveToken(addr string, token *TokenCache) error {
	cacheFile, err := getCacheFilePath(addr)
	if err != nil {
		return err
	}

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0600)
}

// clearToken 清除token缓存
func clearToken(addr string) error {
	cacheFile, err := getCacheFilePath(addr)
	if err != nil {
		return err
	}

	if err := os.Remove(cacheFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// isTokenValid 检查token是否有效
func isTokenValid(token *TokenCache) bool {
	if token == nil || token.AccessToken == "" {
		return false
	}

	// 提前过期检查，留出缓冲时间
	now := time.Now().Unix()
	return token.ExpireTime > (now + tokenExpireBuffer)
}

// Login 登录获取accessToken
func Login(addr, username, password string) (*AuthResponse, error) {
	if addr == "" {
		return nil, errors.New("address is required")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// 确保地址以/nacos结尾
	if !strings.HasSuffix(addr, "/nacos") {
		if strings.HasSuffix(addr, "/") {
			addr = addr + "nacos"
		} else {
			addr = addr + "/nacos"
		}
	}

	// 构建登录URL
	loginURL, err := url.JoinPath(addr, authUrl)
	if err != nil {
		return nil, err
	}

	// 构建请求体
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)

	resp, err := http.PostForm(loginURL, formData)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d, response: %s", ErrAuthFailed, resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, err
	}

	if authResp.AccessToken == "" {
		return nil, ErrAuthFailed
	}

	return &authResp, nil
}

// GetAccessToken 获取有效的accessToken，优先从缓存获取，过期则重新登录
func GetAccessToken(config *NacosConfig) (string, error) {
	if config.Username == "" || config.Password == "" {
		// 没有配置用户名密码，返回空token（无需认证）
		return "", nil
	}

	// 尝试从缓存加载
	cachedToken, err := loadToken(config.Addr)
	if err == nil && isTokenValid(cachedToken) {
		return cachedToken.AccessToken, nil
	}

	// token过期或无效，重新登录
	authResp, err := Login(config.Addr, config.Username, config.Password)
	if err != nil {
		return "", err
	}

	// 保存到缓存
	expireTime := time.Now().Unix() + authResp.TokenTTL
	tokenCache := &TokenCache{
		AccessToken: authResp.AccessToken,
		ExpireTime:  expireTime,
	}

	if err := saveToken(config.Addr, tokenCache); err != nil {
		// 保存失败不影响使用，只记录错误
		fmt.Fprintf(os.Stderr, "Warning: failed to save token cache: %v\n", err)
	}

	return authResp.AccessToken, nil
}

// ClearAccessToken 清除缓存的accessToken
func ClearAccessToken(addr string) error {
	return clearToken(addr)
}
