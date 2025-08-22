package shopee

import (
	"fmt"
	"net/http"
)

// ErrType 定义错误类型
type ErrType string

const (
	ErrTypeAuth         ErrType = "auth"         // 认证错误
	ErrTypeNetwork      ErrType = "network"      // 网络错误
	ErrTypeValidation   ErrType = "validation"   // 参数验证错误
	ErrTypeAPI          ErrType = "api"          // API错误
	ErrTypeRateLimit    ErrType = "rate_limit"   // 限流错误
	ErrTypeParsing      ErrType = "parsing"      // 解析错误
	ErrTypeUnknown      ErrType = "unknown"      // 未知错误
)

// ShopeeError 统一错误类型
type ShopeeError struct {
	Type       ErrType
	Code       int
	Message    string
	Err        error
	StatusCode int // HTTP状态码
}

func (e *ShopeeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s:%d] %s: %v", e.Type, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s:%d] %s", e.Type, e.Code, e.Message)
}

func (e *ShopeeError) Unwrap() error {
	return e.Err
}

// IsType 检查错误类型
func (e *ShopeeError) IsType(errType ErrType) bool {
	return e.Type == errType
}

// IsRetryable 判断错误是否可重试
func (e *ShopeeError) IsRetryable() bool {
	if e.Type == ErrTypeRateLimit {
		return true
	}
	if e.StatusCode >= 500 {
		return true
	}
	if e.Type == ErrTypeNetwork {
		return true
	}
	return false
}

// 预定义的错误创建函数
func NewAuthError(code int, message string, err error) *ShopeeError {
	return &ShopeeError{
		Type:    ErrTypeAuth,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewNetworkError(message string, err error) *ShopeeError {
	return &ShopeeError{
		Type:    ErrTypeNetwork,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string) *ShopeeError {
	return &ShopeeError{
		Type:    ErrTypeValidation,
		Message: message,
	}
}

func NewAPIError(code int, message string, statusCode int) *ShopeeError {
	return &ShopeeError{
		Type:       ErrTypeAPI,
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewRateLimitError(message string) *ShopeeError {
	return &ShopeeError{
		Type:    ErrTypeRateLimit,
		Message: message,
	}
}

func NewParsingError(message string, err error) *ShopeeError {
	return &ShopeeError{
		Type:    ErrTypeParsing,
		Message: message,
		Err:     err,
	}
}

// 常见错误处理函数
func HandleHTTPError(resp *http.Response, body []byte) *ShopeeError {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return NewAuthError(resp.StatusCode, "unauthorized access", nil)
	case http.StatusForbidden:
		return NewAuthError(resp.StatusCode, "forbidden access", nil)
	case http.StatusTooManyRequests:
		return NewRateLimitError("rate limit exceeded")
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return NewAPIError(resp.StatusCode, "server error", resp.StatusCode)
	default:
		if resp.StatusCode >= 400 {
			return NewAPIError(resp.StatusCode, string(body), resp.StatusCode)
		}
	}
	return nil
}

// 检查常见的业务错误
func CheckBusinessError(code int, message string) *ShopeeError {
	switch message {
	case "error_server":
		return NewAPIError(code, "请联系管理员", 0)
	case "error_need_vcode":
		return NewValidationError("需要验证码")
	case "error_invalid_vcode":
		return NewValidationError("验证码错误")
	case "error_name_or_password_incorrect":
		return NewAuthError(code, "账号或密码错误", nil)
	case RateLimitError:
		return NewRateLimitError("rate limit exceeded")
	}
	
	if code == TokenNotFoundCode {
		return NewAuthError(code, "cookies 失效，请重新登陆", nil)
	}
	
	if code != ResponseCodeSuccess && code != 0 {
		return NewAPIError(code, message, 0)
	}
	
	return nil
}