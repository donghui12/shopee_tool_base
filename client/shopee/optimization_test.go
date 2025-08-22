package shopee

import (
	"context"
	"testing"
	"time"
)

func TestErrorHandling(t *testing.T) {
	// 测试错误类型创建
	authErr := NewAuthError(401, "unauthorized", nil)
	if !authErr.IsType(ErrTypeAuth) {
		t.Errorf("Expected auth error type, got %s", authErr.Type)
	}

	if !authErr.IsRetryable() == false {
		t.Error("Auth errors should not be retryable")
	}

	// 测试网络错误
	networkErr := NewNetworkError("connection failed", nil)
	if !networkErr.IsRetryable() {
		t.Error("Network errors should be retryable")
	}

	// 测试限流错误
	rateLimitErr := NewRateLimitError("rate limit exceeded")
	if !rateLimitErr.IsRetryable() {
		t.Error("Rate limit errors should be retryable")
	}
}

func TestRequestManager(t *testing.T) {
	// 创建测试客户端
	config := DefaultConfig()
	config.Timeout = 5 * time.Second
	client := NewClientWithConfig(config)
	
	rm := NewRequestManager(client)
	
	// 测试请求配置
	ctx := context.Background()
	
	// 这里可以添加模拟的HTTP服务器来测试实际请求
	// 由于需要模拟服务器，这里只测试配置
	
	// 测试请求选项
	opts := []RequestOption{
		WithRequestTimeout(10 * time.Second),
		WithRequestRetry(2, 1*time.Second),
		WithHeaders(map[string]string{"Custom-Header": "test"}),
		WithProxy(false),
		WithSkipErrorLog(true),
	}
	
	reqConfig := &RequestConfig{}
	for _, opt := range opts {
		opt(reqConfig)
	}
	
	if reqConfig.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", reqConfig.Timeout)
	}
	
	if reqConfig.RetryTimes != 2 {
		t.Errorf("Expected retry times 2, got %d", reqConfig.RetryTimes)
	}
	
	if reqConfig.Headers["Custom-Header"] != "test" {
		t.Error("Custom header not set correctly")
	}
	
	// 避免未使用变量的警告
	_ = rm
	_ = ctx
}

func TestClientConfig(t *testing.T) {
	// 测试默认配置
	defaultConfig := DefaultConfig()
	if defaultConfig.BaseURL != BaseSellerURL {
		t.Errorf("Expected base URL %s, got %s", BaseSellerURL, defaultConfig.BaseURL)
	}
	
	// 测试台湾配置
	twConfig := TaiwanConfig()
	if twConfig.BaseURL != BaseSellerURLForTw {
		t.Errorf("Expected Taiwan base URL %s, got %s", BaseSellerURLForTw, twConfig.BaseURL)
	}
	
	// 测试高性能配置
	hpConfig := HighPerformanceConfig()
	if hpConfig.MaxIdleConns != 200 {
		t.Errorf("Expected max idle conns 200, got %d", hpConfig.MaxIdleConns)
	}
	
	// 测试代理配置
	proxyConfig := ProxyConfig()
	if !proxyConfig.UseProxy {
		t.Error("Proxy config should enable proxy")
	}
}

func TestClientFactory(t *testing.T) {
	// 测试工厂方法
	defaultClient := NewDefaultClient()
	if defaultClient.baseURL != BaseSellerURL {
		t.Error("Default client should use default base URL")
	}
	
	twClient := NewTaiwanClient()
	if twClient.baseURL != BaseSellerURLForTw {
		t.Error("Taiwan client should use Taiwan base URL")
	}
	
	hpClient := NewHighPerformanceClient()
	if hpClient.retryTimes != 1 {
		t.Error("High performance client should have retry times 1")
	}
	
	proxyClient := NewProxyClient()
	if proxyClient.retryTimes != 5 {
		t.Error("Proxy client should have retry times 5")
	}
}

// 基准测试
func BenchmarkErrorCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewAPIError(500, "server error", 500)
	}
}

func BenchmarkRequestConfigCreation(b *testing.B) {
	opts := []RequestOption{
		WithRequestTimeout(10 * time.Second),
		WithRequestRetry(3, 2*time.Second),
		WithHeaders(map[string]string{"Test": "value"}),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &RequestConfig{}
		for _, opt := range opts {
			opt(config)
		}
	}
}

// 示例测试，展示如何使用新的API
func ExampleNewRequestManager() {
	// 创建客户端
	client := NewDefaultClient()
	
	// 创建请求管理器
	rm := NewRequestManager(client)
	
	// 使用请求管理器发送请求
	ctx := context.Background()
	_, err := rm.DoRequestWithResponse(
		ctx,
		"GET",
		"/api/test",
		nil,
		"cookies=test",
		WithRequestTimeout(5*time.Second),
		WithRequestRetry(2, 1*time.Second),
	)
	
	if err != nil {
		// 错误处理
		if err.IsType(ErrTypeAuth) {
			// 处理认证错误
		} else if err.IsRetryable() {
			// 处理可重试错误
		}
	}
}