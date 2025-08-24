// shopee 请求参数
package shopee

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ----------------------- 公共请求参数 ----------------------
type CommomParam struct {
	ShopId string `json:"cnsc_shop_id"`
	Region string `json:"cbsc_shop_region"`
}

func (p *CommomParam) ToFormValues() url.Values {
	form := url.Values{}
	form.Set("SPC_CDS", uuid.New().String())
	form.Set("SPC_CDS_VER", "2")
	if p.ShopId != "" {
		form.Set("cnsc_shop_id", p.ShopId)
	}
	if p.Region != "" {
		form.Set("cbsc_shop_region", p.Region)
	}
	return form
}

// ----------------------- 登录请求参数 ----------------------
// LoginParam 登录请求参数
type LoginParam struct {
	PasswordHash    string `json:"password_hash"`
	Remember        bool   `json:"remember"`
	OtpType         int    `json:"otp_type"`
	Subaccount      string `json:"subaccount"`
	SubaccountPhone string `json:"subaccount_phone"`
	SubaccountEmail string `json:"subaccount_email"`
	Vcode           string `json:"vcode,omitempty"`
}

func (p *LoginParam) ToFormValues() url.Values {
	form := url.Values{}
	form.Set("password_hash", p.PasswordHash)
	form.Set("remember", strconv.FormatBool(p.Remember))
	form.Set("otp_type", strconv.Itoa(p.OtpType))

	if p.Subaccount != "" {
		form.Set("subaccount", p.Subaccount)
	}
	if p.SubaccountPhone != "" {
		form.Set("subaccount_phone", p.SubaccountPhone)
	}
	if p.SubaccountEmail != "" {
		form.Set("subaccount_email", p.SubaccountEmail)
	}
	if p.Vcode != "" {
		form.Set("vcode", p.Vcode)
	}
	return form
}

// ----------------------- GetOrSetShop --------------------
type GetOrSetShopReq struct {
	ShopId int64 `json:"shop_id"`
}

// ----------------------- 更新商品信息 ----------------------
type ProductStatusInfo struct {
	Unlisted bool `json:"unlisted"`
}

// UpdateProductInfoReq 更新商品信息请求
type UpdateProductInfoReq struct {
	ProductId     int64             `json:"product_id"`
	DaysToShip    int               `json:"days_to_ship,omitempty"`
	Cookies       string            `json:"cookies"`
	ShopID        string            `json:"shop_id"`
	Region        string            `json:"region"`
	ProductStatus ProductStatusInfo `json:"product_info,omitempty"`
}

// PreOrderInfo 商品预售信息
type PreOrderInfo struct {
	PreOrder   bool `json:"pre_order"`
	DaysToShip int  `json:"days_to_ship"`
}

// ProductInfo 商品基础信息
type ProductInfo struct {
	EnableModelLevelDts bool         `json:"enable_model_level_dts"`
	PreOrderInfo        PreOrderInfo `json:"pre_order_info"`
	Unlisted            bool         `json:"unlisted"`
}

// UpdateProductInfoRequest 更新商品信息请求参数
type UpdateProductInfoRequest struct {
	ProductID   int64       `json:"product_id"`
	ProductInfo ProductInfo `json:"product_info"`
	IsDraft     bool        `json:"is_draft"`
}

// ----------------------- 获取token信息 ----------------------
// GetAccessTokenReq 获取 accessToken 请求参数
type GetAccessTokenReq struct {
	Code         string `json:"code"`
	ShopId       int64  `json:"shop_id"`
	PartnerId    int64  `json:"partner_id"`
	RefreshToken string `json:"refresh_token"`
}

// ----------------------- 商品列表信息 ----------------------

// ProductListRequest 商品列表请求
type ProductListRequest struct {
	PageSize   int    `json:"page_size"`
	PageNo     int    `json:"page_no"`
	SearchType string `json:"search_type,omitempty"`
	Keyword    string `json:"keyword,omitempty"`
	SortBy     string `json:"sort_by,omitempty"`
	SortType   int    `json:"sort_type,omitempty"`
}

type OngoingCampaigns struct {
	ProductID int `json:"product_id"`
}

type PromotionDetail struct {
	OngoingCampaigns []OngoingCampaigns `json:"ongoing_campaigns"`
}

// 商品信息
type Product struct {
	ID              int               `json:"id"`
	Name            string            `json:"name"`
	PromotionDetail PromotionDetail   `json:"promotion_detail"`
	ModelList       []Model           `json:"model_list"`
	Statistics      ProductStatistics `json:"statistics"`
	CreateTime      int64             `json:"create_time"`
}

// ProductStatistics 商品销售信息
type ProductStatistics struct {
	LikedCount int64 `json:"liked_count"`
	SoldCount  int64 `json:"sold_count"`
	ViewCount  int64 `json:"view_count"`
}

// DeleteProductReq 删除商品信息
type DeleteProductReq struct {
	ProductIdList []int64 `json:"product_id_list"`
}

// ----------------------- 折扣 ----------------------
// 获取折扣列表 请求参数
type getDiscountListRequest struct {
	DiscountType int    `json:"discount_type"` // 1 表示限时折扣
	TimeStatus   int    `json:"time_status"`   // 0 全部，1 进行中，2 已结束
	Offset       int    `json:"offset"`        // 分页 offset
	Limit        int    `json:"limit"`         // 每页条数
	Keyword      string `json:"keyword"`       // 搜索关键词
	PeriodFrom   int64  `json:"period_from"`   // 开始时间
	PeriodTo     int64  `json:"period_to"`     // 结束时间
}

// 获取折扣列表 请求参数
type getDiscountItemRequest struct {
	PromotionId int64 `json:"promotion_id"` // 折扣 id
	Offset      int   `json:"offset"`       // 分页 offset
	Limit       int   `json:"limit"`        // 每页条数
}

// 更新折扣列表 请求参数
type UpdateDiscountItemRequest struct {
	PromotionId       int                `json:"promotion_id"` // 折扣 id
	DiscountModelList []DiscountItemList `json:"discount_model_list"`
}

// delete_discount 删除请求参数
type DeleteDiscountReq struct {
	PromotionID int64 `json:"promotion_id"` // 折扣活动 ID
	Action      int   `json:"action"`       // 1 = 删除，2 = 停止
}

type CreateDiscountReq struct {
	StartTime        int64    `json:"start_time"`         // Unix 时间戳
	EndTime          int64    `json:"end_time"`           // Unix 时间戳
	Title            string   `json:"title"`              // 折扣标题
	Status           int      `json:"status"`             // 折扣状态
	Source           int      `json:"source"`             // 来源（例如 0：手动创建）
	Images           []string `json:"images"`             // 图片资源 ID 列表
	TotalProduct     int      `json:"total_product"`      // 商品数量
	GlobalDiscountID int64    `json:"global_discount_id"` // 全局折扣 ID
}

func (req *CreateDiscountReq) ConvertFromDiscount(discount Discount) {
	future := time.Now().AddDate(0, 1, 0) // 当前时间加1个月
	req.StartTime = time.Now().Add(10 * time.Minute).Unix()
	req.StartTime = time.Now().Unix()
	req.EndTime = future.Unix()
	req.Status = 1
	req.Source = 0
	req.Images = discount.SellerDiscount.ItemPreview.Images
	req.TotalProduct = discount.SellerDiscount.ItemPreview.ItemCount

	// 原始折扣名
	originalTitle := discount.SellerDiscount.Name
	currentTimeStr := time.Now().Format("20060102150405")

	// 处理 title 逻辑
	baseTitle := originalTitle
	if strings.Contains(originalTitle, "_copy_") {
		parts := strings.Split(originalTitle, "_copy_")
		if len(parts) == 2 {
			baseTitle = parts[0]
		}
	}
	req.Title = fmt.Sprintf("%s_copy_%s", baseTitle, currentTimeStr)
}

type AccountInfo struct {
	AccountId   int64  `json:"account_id"`
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
}

// 获取 session 返回参数
type GetSessionResp struct {
	Code        int         `json:"code"`
	Errcode     int         `json:"errcode"`
	Message     string      `json:"message"`
	AccountInfo AccountInfo `json:"sub_account_info"`
}
