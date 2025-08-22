package shopee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"github.com/shopee_tool_base/pkg/logger"
)

// RequestOption 请求配置选项
type RequestOption func(*RequestConfig)

// RequestConfig 请求配置
type RequestConfig struct {
	Headers       map[string]string
	Timeout       time.Duration
	RetryTimes    int
	RetryDelay    time.Duration
	UseProxy      bool
	SkipErrorLog  bool
}

// WithHeaders 设置请求头
func WithHeaders(headers map[string]string) RequestOption {
	return func(config *RequestConfig) {
		if config.Headers == nil {
			config.Headers = make(map[string]string)
		}
		for k, v := range headers {
			config.Headers[k] = v
		}
	}
}

// WithRequestTimeout 设置请求超时时间
func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(config *RequestConfig) {
		config.Timeout = timeout
	}
}

// WithRequestRetry 设置重试配置
func WithRequestRetry(times int, delay time.Duration) RequestOption {
	return func(config *RequestConfig) {
		config.RetryTimes = times
		config.RetryDelay = delay
	}
}

// WithProxy 启用代理
func WithProxy(useProxy bool) RequestOption {
	return func(config *RequestConfig) {
		config.UseProxy = useProxy
	}
}

// WithSkipErrorLog 跳过错误日志
func WithSkipErrorLog(skip bool) RequestOption {
	return func(config *RequestConfig) {
		config.SkipErrorLog = skip
	}
}

// RequestManager 请求管理器
type RequestManager struct {
	client *Client
}

// NewRequestManager 创建请求管理器
func NewRequestManager(client *Client) *RequestManager {
	return &RequestManager{client: client}
}

// DoRequest 执行请求并返回原始响应
func (rm *RequestManager) DoRequest(ctx context.Context, method, path string, reqBody interface{}, cookies string, opts ...RequestOption) (*http.Response, *ShopeeError) {
	config := &RequestConfig{
		Timeout:    rm.client.timeout,
		RetryTimes: rm.client.retryTimes,
		RetryDelay: rm.client.retryDelay,
	}
	
	for _, opt := range opts {
		opt(config)
	}

	var bodyReader io.Reader
	if reqBody != nil {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, NewParsingError("marshal request body failed", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	url := rm.client.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, NewNetworkError("create request failed", err)
	}

	// 设置默认headers
	rm.client.setCommonHeaders(req)
	req.Header.Set("Cookie", cookies)

	// 设置自定义headers
	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	// 执行请求
	var resp *http.Response
	var lastErr error

	for i := 0; i <= config.RetryTimes; i++ {
		if config.UseProxy {
			resp, lastErr = rm.client.executeWithProxy(req)
		} else {
			resp, lastErr = rm.client.executeWithRetry(req)
		}

		if lastErr == nil {
			return resp, nil
		}

		// 检查是否可重试
		if shopeeErr, ok := lastErr.(*ShopeeError); ok && !shopeeErr.IsRetryable() {
			break
		}

		if i == config.RetryTimes {
			break
		}

		if !config.SkipErrorLog {
			logger.Warn("request failed, retrying",
				zap.String("url", url),
				zap.Int("attempt", i+1),
				zap.Error(lastErr),
			)
		}

		time.Sleep(config.RetryDelay)
	}

	return nil, NewNetworkError(fmt.Sprintf("request failed after %d retries", config.RetryTimes), lastErr)
}

// DoRequestWithResponse 执行请求并解析响应
func (rm *RequestManager) DoRequestWithResponse(ctx context.Context, method, path string, reqBody interface{}, cookies string, opts ...RequestOption) ([]byte, *ShopeeError) {
	resp, err := rm.DoRequest(ctx, method, path, reqBody, cookies, opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, NewNetworkError("read response body failed", readErr)
	}

	// 检查HTTP状态码
	if httpErr := HandleHTTPError(resp, body); httpErr != nil {
		return body, httpErr
	}

	return body, nil
}

// DoRequestWithJSONResponse 执行请求并解析JSON响应
func (rm *RequestManager) DoRequestWithJSONResponse(ctx context.Context, method, path string, reqBody interface{}, cookies string, target interface{}, opts ...RequestOption) *ShopeeError {
	body, err := rm.DoRequestWithResponse(ctx, method, path, reqBody, cookies, opts...)
	if err != nil {
		return err
	}

	if unmarshalErr := json.Unmarshal(body, target); unmarshalErr != nil {
		return NewParsingError("unmarshal response failed", unmarshalErr)
	}

	return nil
}

// DoRequestWithCommonResponse 执行请求并解析通用响应格式
func DoRequestWithCommonResponse[T any](rm *RequestManager, ctx context.Context, method, path string, reqBody interface{}, cookies string, opts ...RequestOption) (*T, *ShopeeError) {
	body, err := rm.DoRequestWithResponse(ctx, method, path, reqBody, cookies, opts...)
	if err != nil {
		return nil, err
	}

	var resp CommonResponse[T]
	if unmarshalErr := json.Unmarshal(body, &resp); unmarshalErr != nil {
		return nil, NewParsingError("unmarshal response failed", unmarshalErr)
	}

	// 检查业务错误
	if businessErr := CheckBusinessError(resp.Code, resp.Message); businessErr != nil {
		return nil, businessErr
	}

	// 检查错误码
	if businessErr := CheckBusinessError(resp.ErrCode, resp.Message); businessErr != nil {
		return nil, businessErr
	}

	return &resp.Data, nil
}

// 为了向后兼容，提供简化的方法
func (c *Client) DoSimpleRequest(method, path string, reqBody interface{}, cookies string) (*http.Response, error) {
	rm := NewRequestManager(c)
	resp, err := rm.DoRequest(context.Background(), method, path, reqBody, cookies)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DoSimpleRequestWithResponse(method, path string, reqBody interface{}, cookies string) ([]byte, error) {
	rm := NewRequestManager(c)
	body, err := rm.DoRequestWithResponse(context.Background(), method, path, reqBody, cookies)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func DoSimpleRequestWithCommonResponse[T any](c *Client, method, path string, reqBody interface{}, cookies string) (*T, error) {
	rm := NewRequestManager(c)
	result, err := DoRequestWithCommonResponse[T](rm, context.Background(), method, path, reqBody, cookies)
	if err != nil {
		return nil, err
	}
	return result, nil
}