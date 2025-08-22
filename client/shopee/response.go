// shopee 响应参数
package shopee

import (
	"encoding/json"
	"fmt"
)

// CommonResponse 公共响应参数
// CommonResponse 是统一响应结构（用于泛型解析）
type CommonResponse[T any] struct {
	Code        int    `json:"code"`
	ErrCode     int    `json:"errcode"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
	Data        T      `json:"data"`
}

// ParseCommonResponse 解析响应 JSON 为 CommonResponse，并返回其中的 Data 字段
func ParseCommonResponse[T any](respBody []byte) (*T, error) {
	var resp CommonResponse[T]
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.ErrCode == TokenNotFoundCode {
		return nil, fmt.Errorf("cookies 失效，请重新登陆")
	}

	// 可选：处理返回码
	if resp.Code != 0 {
		return nil, fmt.Errorf("请求失败：code=%d, msg=%s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}

// ---------------------- 获取或设置店铺响应 -----------------------
// LoginData 登录响应
type LoginData struct {
	Token string `json:"token"`
}

// ---------------------- 获取或设置店铺响应 -----------------------
// GetOrSetShopData 获取或设置店铺响应
type GetOrSetShopData struct {
	ShopId string `json:"shop_id"`
}

// ---------------------- 折扣 -----------------------
type ItemPreview struct {
	ItemCount int      `json:"item_count"`
	Images    []string `json:"images"`
}

type SellerDiscount struct {
	DiscountID       int64       `json:"discount_id"`
	Name             string      `json:"name"`
	TimeStatus       int         `json:"time_status"`
	StartTime        int64       `json:"start_time"`
	EndTime          int64       `json:"end_time"`
	ItemPreview      ItemPreview `json:"item_preview"`
	Source           int         `json:"source"`
	GlobalDiscountID int64       `json:"global_discount_id"`
}

type Discount struct {
	DiscountType   int            `json:"discount_type"`
	SellerDiscount SellerDiscount `json:"seller_discount"`
}

type DiscountList struct {
	Discounts  []Discount `json:"discounts"`
	TotalCount int        `json:"total_count"`
}

type DiscountItemList struct {
	ItemID         int64 `json:"item_id"`         // 商品 ID
	ModelID        int64 `json:"model_id"`        // SKU 模型 ID
	PromotionPrice int64 `json:"promotion_price"` // 活动价（单位：分）
	UserItemLimit  int   `json:"user_item_limit"` // 每人限购
	Status         int   `json:"status"`          // 状态（例如 1 表示启用）
	PromotionStock int   `json:"promotion_stock"` // 活动库存
}

type ItemInfo struct {
	ItemID  int64    `json:"item_id"`  // 商品 ID
	Status  int      `json:"status"`   // 状态
	Name    string   `json:"name"`     // 商品名称
	Images  []string `json:"images"`   // 商品图片 ID 列表
	Price   int64    `json:"price"`    // 商品价格（单位：分）
	Stock   int      `json:"stock"`    // 商品库存
	Sold    int      `json:"sold"`     // 已售数量
	ExtInfo ExtInfo  `json:"ext_info"` // 扩展信息
}

type ExtInfo struct {
	HasWholesale     bool   `json:"has_wholesale"`      // 是否支持批发
	MinPurchaseLimit int    `json:"min_purchase_limit"` // 最小起购数量
	LogisticsInfo    string `json:"logistics_info"`     // 物流信息（JSON 字符串）
}

type ModelInfo struct {
	ItemID int64       `json:"item_id"` // 商品 ID
	Models []ModelItem `json:"models"`  // 商品规格列表
}

type ModelItem struct {
	ItemID    int64  `json:"item_id"`    // 商品 ID
	ModelID   int64  `json:"model_id"`   // 规格 ID
	Name      string `json:"name"`       // 规格名称
	Price     int64  `json:"price"`      // 价格（单位：分）
	Stock     int    `json:"stock"`      // 库存
	Sold      int    `json:"sold"`       // 已售数量
	IsDefault bool   `json:"is_default"` // 是否默认规格
	PffTag    bool   `json:"pff_tag"`    // 是否有 PFF 标签
}

type PriceStockInfo struct {
	ItemID              int64                `json:"item_id"`
	AggregatedPriceInfo AggregatedPriceInfo  `json:"aggregated_price_info"`
	AggregatedStockInfo AggregatedStockInfo  `json:"aggregated_stock_info"`
	SkuStockPriceList   []SkuStockPriceEntry `json:"sku_stock_price_list"`
}

type AggregatedPriceInfo struct {
	MinPrice int64 `json:"min_price"`
	MaxPrice int64 `json:"max_price"`
}

type AggregatedStockInfo struct {
	WMSTotalStock    int `json:"wms_total_stock"`
	SellerTotalStock int `json:"seller_total_stock"`
}

type SkuStockPriceEntry struct {
	ModelID         int64     `json:"model_id"`
	PffTag          bool      `json:"pff_tag"`
	SellerStockInfo StockInfo `json:"seller_stock_info"`
	PriceInfo       PriceInfo `json:"price_info"`
	WmsStockInfo    StockInfo `json:"wms_stock_info"`
}

type StockInfo struct {
	StockPromotionType int `json:"stock_promotion_type"`
	NormalStock        int `json:"normal_stock"`
	PromotionStock     int `json:"promotion_stock"`
	ReservedStock      int `json:"reserved_stock"`
}

type PriceInfo struct {
	PricePromotionType  int   `json:"price_promotion_type"`
	InputPromotionPrice int64 `json:"input_promotion_price"`
	PromotionPrice      int64 `json:"promotion_price"`
	InputNormalPrice    int64 `json:"input_normal_price"`
	NormalPrice         int64 `json:"normal_price"`
}

// DiscountItemData 折扣列表 data
type DiscountItemData struct {
	DiscountItemList []DiscountItemList `json:"discount_item_list"`
	ItemInfo         []ItemInfo         `json:"item_info"`
	ModelInfo        []ModelInfo        `json:"model_info"`
	PriceStockInfo   []PriceStockInfo   `json:"price_stock_info"`
	TotalCount       int                `json:"total_count"`
	ErrorList        string             `json:"error_list"`
}

// 删除折扣错误列表
type DeleteDiscountData struct {
	ErrorList []ItemError `json:"error_list"` // 错误详情，通常为 null
}

// 创建折扣结果
type CreateDiscountData struct {
	ErrorList   []string `json:"error_list"` // 错误详情，通常为 null
	PromationId int64    `json:"promotion_id"`
}

type UpdateSellerDiscountItemsResp struct {
	TotalCount     int         `json:"total_count"`
	SuccessCount   int         `json:"success_count"`
	FailedItemList []int64     `json:"failed_item_list"` // 商品 ID 列表
	ErrorList      []ItemError `json:"error_list"`
	PriceRatio     float64     `json:"price_ratio"` // 如为整数也建议 float64 做兼容
}

type ItemError struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
