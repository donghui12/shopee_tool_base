package main

import (
	"fmt"
	"log"

	"github.com/donghui12/shopee_tool_base/client/shopee"
	"github.com/donghui12/shopee_tool_base/global"
)

func main() {
	// 初始化数据库连接
	dsn := "user:password@tcp(localhost:3306)/shopee_tool?charset=utf8mb4&parseTime=True&loc=Local"
	if err := global.SetGlobalDB(dsn); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		// 继续运行，但数据库功能将不可用
	}

	// 初始化 Shopee 客户端
	shopee.InitShopeeClient()

	fmt.Println("Shopee Tool Base 已启动")
	fmt.Println("使用以下功能:")
	fmt.Println("1. 商品管理")
	fmt.Println("2. 店铺管理")
	fmt.Println("3. 折扣管理")
	fmt.Println("4. 账号管理")

	// 示例用法
	client := shopee.GetShopeeClient()
	if client != nil {
		fmt.Println("Shopee 客户端初始化成功")
	}

	twClient := shopee.GetTwShopeeClient()
	if twClient != nil {
		fmt.Println("Shopee TW 客户端初始化成功")
	}
}
