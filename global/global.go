package global

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/shopee_tool_base/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB       *gorm.DB
	ProxyURL *url.URL
)

// SetGlobalDB 设置全局数据库连接
func SetGlobalDB(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      gormLogger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		return err
	}

	// 绑定到 global.DB
	DB = db
	return nil
}

// InitProxyWithURL 使用 URL 字符串初始化代理配置
func InitProxyWithURL(proxyURLStr string) error {
	var err error
	ProxyURL, err = url.Parse(proxyURLStr)
	if err != nil {
		logger.Error("解析代理URL失败:", zap.Error(err))
		return err
	}
	return nil
}
