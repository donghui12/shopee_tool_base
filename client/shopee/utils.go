package shopee

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/donghui12/shopee_tool_base/pkg/constant"
	"github.com/donghui12/shopee_tool_base/pkg/logger"
	"github.com/donghui12/shopee_tool_base/pkg/proxy"

	"go.uber.org/zap"
)

// MD5Hash 计算MD5哈希值
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// Shop API sign
func GetShopSign(partnerId, path, timestamp, accessToken, shopId string) string {
	tmpBaseString := fmt.Sprintf("%s%s%s%s%s", partnerId,
		path, timestamp, accessToken, shopId)
	mac := hmac.New(sha256.New, []byte(constant.PartnerKey))
	// 写入待加密的 base_string
	mac.Write([]byte(tmpBaseString))

	// 获取签名结果并转换为十六进制字符串
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature
}

// Public API sign
func GetPublicSign(partnerId, path, timestamp string) string {
	tmpBaseString := fmt.Sprintf("%s%s%s", partnerId,
		path, timestamp)
	mac := hmac.New(sha256.New, []byte(constant.PartnerKey))
	// 写入待加密的 base_string
	mac.Write([]byte(tmpBaseString))

	// 获取签名结果并转换为十六进制字符串
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature
}

// formatPhone 格式化手机号
func formatPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if !strings.HasPrefix(phone, "86") {
		return "86" + phone
	}
	return phone
}

// isPhone 判断是否是手机号
func isPhone(phone string) bool {
	if len(phone) != 13 {
		return false
	}
	return strings.HasPrefix(phone, "86")
}

// VerifyEmailFormat 验证邮箱格式
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// setCommonHeaders 设置通用请求头
func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Origin", BaseSellerURL)
	req.Header.Set("Referer", BaseSellerURL)
	req.Header.Set("Host", BaseSellerHost)

	// 添加 cookies
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
}

// executeWithRetry 带重试的请求执行
func (c *Client) executeWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.retryTimes; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil {
			return resp, nil
		}
		if i == c.retryTimes {
			return resp, fmt.Errorf("request failed after %d retries: %w", c.retryTimes, err)
		}
		time.Sleep(c.retryDelay)
	}

	return resp, err
}

func (c *Client) executeWithProxy(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var lastErr error

	for i := 0; i <= c.retryTimes; i++ {

		proxyIP, err := proxy.GetProxyIP()

		logger.Info("获取代理IP", zap.String("IP:", proxyIP))
		if err != nil {
			lastErr = err
			logger.Error("获取代理IP失败:", zap.Error(err))
			continue
		}

		proxyURL, err := url.Parse("http://" + proxyIP)
		logger.Info("获取代理url", zap.Any("url:", proxyURL))
		if err != nil {
			lastErr = err
			logger.Error("解析代理IP失败:", zap.Error(err))
			continue
		}

		// 使用动态代理构造新的 HTTP 客户端
		c.httpClient = &http.Client{
			Timeout: c.httpClient.Timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
		resp, err := c.httpClient.Do(req)
		if err == nil {
			return resp, nil
		}
		if i == c.retryTimes {
			return resp, fmt.Errorf("request failed after %d retries: %w", c.retryTimes, err)
		}
		logger.Error("请求失败:", zap.Error(err))
		lastErr = err
		time.Sleep(c.retryDelay)
	}

	return resp, lastErr
}
