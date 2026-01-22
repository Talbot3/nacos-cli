package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	baseUrl = "/cs/configs"
	authHeader = "Authorization"
)

// Client Nacos客户端
type Client struct {
	Config *NacosConfig
}

// doRequest 执行带有认证的HTTP请求
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	// 获取accessToken
	token, err := GetAccessToken(c.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 添加认证头
	if token != "" {
		req.Header.Set(authHeader, "Bearer "+token)
	}

	client := &http.Client{}
	return client.Do(req)
}

// createRequest 创建HTTP请求
func (c *Client) createRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	// 获取accessToken
	token, err := GetAccessToken(c.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 添加认证头
	if token != "" {
		req.Header.Set(authHeader, "Bearer "+token)
	}

	return req, nil
}

// Get获取配置
func (c *Client) Get(operation ConfigGetOperation) (*NacosConfigDetail, error) {

	configUrl, err := getUrl(c.Config)

	if err != nil {
		return nil, err
	}

	// Nacos API 中 public 命名空间用空字符串表示
	tenant := operation.Namespace
	if tenant == "public" {
		tenant = ""
	}

	requestUrl := fmt.Sprintf(configUrl+"?dataId=%s&group=%s&tenant=%s", operation.DataId, operation.Group, tenant)

	req, err := c.createRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// 检查是否是认证错误
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			// 清除缓存并重试一次
			_ = ClearAccessToken(c.Config.Addr)
			token, err := GetAccessToken(c.Config)
			if err != nil {
				return nil, fmt.Errorf("authentication failed: %w", err)
			}

			req, err = http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				return nil, err
			}
			if token != "" {
				req.Header.Set(authHeader, "Bearer "+token)
			}

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != http.StatusOK {
				// 配置不存在时返回 403
				if resp.StatusCode == http.StatusForbidden && len(body) == 0 {
					return nil, errors.New("config not exists")
				}
				return nil, fmt.Errorf("response error,status code:%d\n%s", resp.StatusCode, body)
			}
		} else {
			return nil, fmt.Errorf("response error,status code:%d\n%s", resp.StatusCode, body)
		}
	}

	// 检查是否为空响应（配置不存在）
	if len(body) == 0 {
		return nil, errors.New("config not exists")
	}

	// 检查 Content-Type，如果是 text/plain，直接返回内容
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/plain") {
		return &NacosConfigDetail{
			DataID: operation.DataId,
			Group:  operation.Group,
			Tenant: operation.Namespace,
			Content: string(body),
			Md5: resp.Header.Get("Content-MD5"),
		}, nil
	}

	detail := NacosConfigDetail{}

	if err = json.Unmarshal(body, &detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

// AllConfig 获取所有配置
func (c *Client) AllConfig(operation ConfigGetOperation) ([]NacosPageItem, error) {

	configUrl, err := getUrl(c.Config)

	if err != nil {
		return nil, err
	}

	// Nacos API 中 public 命名空间用空字符串表示
	tenant := operation.Namespace
	if tenant == "public" {
		tenant = ""
	}

	requestUrl := fmt.Sprintf(configUrl+"?dataId=&group=%s&tenant=%s&pageNo=1&pageSize=999&search=accurate", operation.Group, tenant)

	req, err := c.createRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// 检查是否是认证错误，尝试重新认证
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			_ = ClearAccessToken(c.Config.Addr)
			token, err := GetAccessToken(c.Config)
			if err != nil {
				return nil, fmt.Errorf("authentication failed: %w", err)
			}

			req, err = http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				return nil, err
			}
			if token != "" {
				req.Header.Set(authHeader, "Bearer "+token)
			}

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("response error,status code:%d", resp.StatusCode)
			}
		} else {
			return nil, fmt.Errorf("response error,status code:%d", resp.StatusCode)
		}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		body = []byte{}
	}

	result := NacosPageResult{}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.PageItems, nil
}

// Edit 更新配置
func (c *Client) Edit(operation ConfigEditOperation) error {

	configUrl, err := getUrl(c.Config)

	if err != nil {
		return err
	}

	// 获取token
	token, err := GetAccessToken(c.Config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Nacos API 中 public 命名空间用空字符串表示
	tenant := operation.Namespace
	if tenant == "public" {
		tenant = ""
	}

	formData := url.Values{
		"dataId":  []string{operation.DataId},
		"group":   []string{operation.Group},
		"content": []string{operation.Content},
		"tenant":  []string{tenant},
		"type":    []string{operation.Type},
	}

	req, err := http.NewRequest(http.MethodPost, configUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		req.Header.Set(authHeader, "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 检查是否是认证错误
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			_ = ClearAccessToken(c.Config.Addr)
			token, err := GetAccessToken(c.Config)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			req, err = http.NewRequest(http.MethodPost, configUrl, strings.NewReader(formData.Encode()))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if token != "" {
				req.Header.Set(authHeader, "Bearer "+token)
			}

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("response error,status code:%d", resp.StatusCode)
			}
		} else {
			return fmt.Errorf("response error,status code:%d", resp.StatusCode)
		}
	}

	return nil
}

// DeleteConfig 删除配置
func (c *Client) DeleteConfig(operation ConfigDeleteOperation) error {
	configUrl, err := getUrl(c.Config)
	if err != nil {
		return err
	}

	// 获取token
	token, err := GetAccessToken(c.Config)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Nacos API 中 public 命名空间用空字符串表示
	// 但删除时不应该包含 tenant 参数（而不是传空字符串）
	var deleteUrl string
	if operation.Namespace == "public" || operation.Namespace == "" {
		deleteUrl = fmt.Sprintf("%s?dataId=%s&group=%s", configUrl, operation.DataId, operation.Group)
	} else {
		deleteUrl = fmt.Sprintf("%s?dataId=%s&group=%s&tenant=%s", configUrl, operation.DataId, operation.Group, operation.Namespace)
	}

	req, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		return err
	}

	if token != "" {
		req.Header.Set(authHeader, "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 检查是否是认证错误
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			_ = ClearAccessToken(c.Config.Addr)
			token, err := GetAccessToken(c.Config)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			req, err = http.NewRequest(http.MethodDelete, deleteUrl, nil)
			if err != nil {
				return err
			}
			if token != "" {
				req.Header.Set(authHeader, "Bearer "+token)
			}

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("response error,status code:%d", resp.StatusCode)
			}
		} else {
			return fmt.Errorf("response error,status code:%d", resp.StatusCode)
		}
	}

	return nil
}

func getUrl(config *NacosConfig) (string, error) {
	return url.JoinPath(config.Addr, config.ApiVersion, baseUrl)
}

func NewDefaultClient() *Client {

	addr := os.Getenv("NACOS_ADDR")
	apiVersion := os.Getenv("NACOS_API_VERSION")
	username := os.Getenv("NACOS_USERNAME")
	password := os.Getenv("NACOS_PASSWORD")

	if addr == "" {
		addr = "http://127.0.0.1:8848/nacos"
	}

	if apiVersion == "" {
		apiVersion = "v1"
	}

	return &Client{
		Config: &NacosConfig{
			Addr:       addr,
			ApiVersion: apiVersion,
			Username:   username,
			Password:   password,
		},
	}
}

// NewClient 创建自定义配置的客户端
func NewClient(addr, apiVersion, username, password string) *Client {
	if apiVersion == "" {
		apiVersion = "v1"
	}

	return &Client{
		Config: &NacosConfig{
			Addr:       addr,
			ApiVersion: apiVersion,
			Username:   username,
			Password:   password,
		},
	}
}
