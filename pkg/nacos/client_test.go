package nacos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.Config)
}

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8848/nacos", "v1", "nacos", "nacos")
	assert.NotNil(t, client)
	assert.NotNil(t, client.Config)
	assert.Equal(t, "http://localhost:8848/nacos", client.Config.Addr)
	assert.Equal(t, "v1", client.Config.ApiVersion)
	assert.Equal(t, "nacos", client.Config.Username)
	assert.Equal(t, "nacos", client.Config.Password)
}

func TestGetUrl(t *testing.T) {
	config := &NacosConfig{
		Addr:       "http://localhost:8848/nacos",
		ApiVersion: "v1",
	}
	url, err := getUrl(config)
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost:8848/nacos/v1/cs/configs", url)
}

func TestIsTokenValid(t *testing.T) {
	tests := []struct {
		name  string
		token *TokenCache
		valid bool
	}{
		{"nil token", nil, false},
		{"empty token", &TokenCache{}, false},
		{"valid token", &TokenCache{AccessToken: "test", ExpireTime: 9999999999}, true},
		{"expired token", &TokenCache{AccessToken: "test", ExpireTime: 1000000000}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTokenValid(tt.token)
			assert.Equal(t, tt.valid, result)
		})
	}
}
