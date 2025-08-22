package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/donghui12/shopee_tool_base/pkg/constant"
	"github.com/donghui12/shopee_tool_base/pkg/logger"

	"go.uber.org/zap"
)

type ProxyData struct {
	Server   string `json:"server"`
	Deadline string `json:"deadline"` // 格式如 "2025-07-20 17:15:37"
}

type ProxyResponse struct {
	Code string      `json:"code"`
	Data []ProxyData `json:"data"`
}

// 缓存结构
var (
	mu         sync.Mutex
	cachedIP   string
	expireTime time.Time
)

// GetProxyIP 获取代理 IP（自动缓存并判断是否过期）
func GetProxyIP() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	// 如果有缓存且未过期，直接返回
	if cachedIP != "" && time.Now().Before(expireTime) {
		logger.Info("当前IP:",
			zap.String("IP:", cachedIP),
			zap.Any("过期时间:", expireTime),
			zap.Any("当前时间:", time.Now()),
		)
		return cachedIP, nil
	}

	// 请求新 IP
	url := fmt.Sprintf("%s/get?key=%s", constant.ProxyHost, constant.ProxyAuthKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 状态码错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var result ProxyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %v", err)
	}

	if result.Code != "SUCCESS" || len(result.Data) == 0 {
		return "", fmt.Errorf("无有效代理 IP 返回")
	}

	// 解析过期时间
	ip := result.Data[0].Server
	deadlineStr := result.Data[0].Deadline
	loc, _ := time.LoadLocation("Asia/Shanghai") // 或 "Local"
	expireAt, err := time.ParseInLocation("2006-01-02 15:04:05", deadlineStr, loc)
	if err != nil {
		return "", fmt.Errorf("解析过期时间失败: %v", err)
	}

	// 更新缓存
	cachedIP = ip
	expireTime = expireAt
	logger.Info("当前IP:",
		zap.String("IP:", cachedIP),
		zap.Any("过期时间:", expireAt))

	return cachedIP, nil
}
