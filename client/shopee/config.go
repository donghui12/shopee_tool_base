package shopee

import (
	"net/http"
	"time"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	BaseURL    string
	UserAgent  string
	RetryTimes int
	RetryDelay time.Duration
	Timeout    time.Duration
	
	// HTTP传输配置
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	
	// 代理配置
	UseProxy      bool
	ProxyRotation bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		BaseURL:             BaseSellerURL,
		UserAgent:          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		RetryTimes:         3,
		RetryDelay:         2 * time.Second,
		Timeout:            30 * time.Second,
		MaxIdleConns:       100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:    90 * time.Second,
		UseProxy:           false,
		ProxyRotation:      false,
	}
}

// TaiwanConfig 返回台湾地区配置
func TaiwanConfig() *ClientConfig {
	config := DefaultConfig()
	config.BaseURL = BaseSellerURLForTw
	config.RetryDelay = 5 * time.Second
	return config
}

// HighPerformanceConfig 返回高性能配置
func HighPerformanceConfig() *ClientConfig {
	config := DefaultConfig()
	config.MaxIdleConns = 200
	config.MaxIdleConnsPerHost = 20
	config.RetryTimes = 1
	config.RetryDelay = 1 * time.Second
	return config
}

// ProxyConfig 返回代理配置
func ProxyConfig() *ClientConfig {
	config := DefaultConfig()
	config.UseProxy = true
	config.ProxyRotation = true
	config.RetryTimes = 5
	config.RetryDelay = 3 * time.Second
	return config
}

// ApplyConfig 应用配置到客户端
func (config *ClientConfig) ApplyConfig(client *Client) {
	client.baseURL = config.BaseURL
	client.userAgent = config.UserAgent
	client.retryTimes = config.RetryTimes
	client.retryDelay = config.RetryDelay
	client.timeout = config.Timeout
	
	// 配置HTTP传输
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
	}
	
	client.httpClient = &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}
}

// NewClientWithConfig 使用配置创建客户端
func NewClientWithConfig(config *ClientConfig) *Client {
	client := &Client{}
	config.ApplyConfig(client)
	return client
}

// 预定义的客户端工厂方法
func NewDefaultClient() *Client {
	return NewClientWithConfig(DefaultConfig())
}

func NewTaiwanClient() *Client {
	return NewClientWithConfig(TaiwanConfig())
}

func NewHighPerformanceClient() *Client {
	return NewClientWithConfig(HighPerformanceConfig())
}

func NewProxyClient() *Client {
	return NewClientWithConfig(ProxyConfig())
}