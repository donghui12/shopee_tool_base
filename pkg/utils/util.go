package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/shopee_tool_base/pkg/logger"

	"go.uber.org/zap"
)

func SaveToFile(filename string, response interface{}) error {
	// 序列化 JSON
	resultJson, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		logger.Error("序列化结果失败", zap.Error(err))
		return err
	}

	// 添加 UTF-8 BOM 并写入文件
	err = os.WriteFile(filename, append([]byte{0xEF, 0xBB, 0xBF}, resultJson...), 0644)
	if err != nil {
		logger.Error("写入结果文件失败", zap.Error(err), zap.String("filename", filename))
		return err
	}

	logger.Info("结果已写入文件", zap.String("filename", filename))
	return nil
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func GenerateFileName(shopId string) string {
	// 生成带时间戳的文件名
	timestamp := time.Now().Unix()
	return fmt.Sprintf("update_pre_day_ship_%s_%d.xlsx", shopId, timestamp)
}

func StrSanitize(s string) string {
	// 1. 去掉所有前后空白符（空格、换行、制表符等）
	s = strings.TrimSpace(s)

	// 2. 去除所有控制字符（如 \n \r \t）
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, s)

	// 3. 去除前缀的特殊字符，比如 + - = _ 等
	s = regexp.MustCompile(`^[\W_]+`).ReplaceAllString(s, "")
	return s
}

func Contains(key int64, list []int64) bool {
	for _, id := range list {
		if id == key {
			return true
		}
	}
	return false
}
