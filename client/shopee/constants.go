package shopee

// API 路径常量
const (
	// API 基础路径
	BaseSellerURL  = "https://seller.shopee.cn"
	BaseSellerHost = "seller.shopee.cn"

	// 账户相关接口
	APIPathLogin              = "/api/cnsc/selleraccount/login/"
	APIPathSwitchMerchantShop = "/api/cnsc/selleraccount/switch_merchant_shop/"
	APIPathGetSession         = "/api/cnsc/selleraccount/get_session/"

	// 商品相关接口
	APIPathUpdateProductInfo   = "/api/v3/product/update_product_info"
	APIPathProductList         = "/api/v3/mpsku/list/v2/get_product_list"
	APIPathGetMerchantShopList = "/api/cnsc/selleraccount/get_merchant_shop_list/"
	APIPathGetOrSetShop        = "/api/cnsc/selleraccount/get_or_set_shop/"
	APIPathProductDetailList   = "/api/v3/product/search_product_list_v2/"
	APIPathDeleteProduct       = "/api/v3/product/delete_product/"

	// 折扣相关接口
	APIPathGetDiscountList    = "/api/marketing/v3/public/discount/list/"
	APIPathDeleteDiscount     = "/api/marketing/v4/discount/delete_stop_discount/"
	APIPathCreateDiscount     = "/api/marketing/v4/discount/create_discount/"
	APIPathGetDiscountItem    = "/api/marketing/v4/discount/get_discount_items_aggregated/"
	APIPathUpdateDiscountItem = "/api/marketing/v4/discount/update_seller_discount_items/"

	// 台湾地址
	BaseSellerURLForTw = "https://partner.shopeemobile.com"
	// 测试环境地址
	BaseTestSellerURLForTw = "https://partner.test-stable.shopeemobile.com"
	// 正式环境地址
	BaseLiveSellerURLForTw                = "https://partner.shopeemobile.com"
	APIPathProductListForTw               = "/api/v2/product/get_item_list"
	APIPathProductUpdateForTw             = "/api/v2/product/update_item"
	APIPathSignForTw                      = "/api/v2/shop/auth_partner"
	APIPathAuthTokenForTw                 = "/api/v2/auth/token/get"
	APIPathAccessTokenForTw               = "/api/v2/auth/access_token/get"
	APIPathGetBaseProductInfo             = "/api/v2/product/get_item_base_info"
	APIPathBatchUpdateProductInfo         = "/api/v3/product/update_product/"
	APIPathBatchUpdateProductInfoWithFile = "/api/mass/mpsku/upload_edit_template/"
)

// API 请求方法
const (
	HTTPMethodGet  = "GET"
	HTTPMethodPost = "POST"
	HTTPMethodPut  = "PUT"
)

// API 响应码
const (
	ResponseCodeSuccess = 0
	ResponseCodeError   = 1
	ProcessCode         = 1000101612
)

// 媒体侧 Error
const (
	RateLimitError    = "requests too frequent"
	RateLimitCode     = 429
	TokenNotFoundCode = 2
)

// 成功 message
const SuccessMessage = "success"

// 国家地区
const (
	SgRegion = "SG"
	MyRegion = "MY"
	ThRegion = "TH"
	VNRegion = "VN"

	SGRegionSmall = "sg"
)

const (
	ListTypeAll      = "all"
	ListTypeLive     = "live_all"
	ListTypeDelisted = "delisted"
)

// 上下架状态
const (
	ListedStatus   = true
	UnlistedStatus = false
)

// 批量更新的来源
const (
	SourceAttributeTool = "attribute_tool"
	SourceSellerCenter  = "seller_center"
)

const (
	DefaultOptType = -1
	EmailOptType   = 4
)

// 登陆类型
const (
	LoginTypeEmail = "email"
)

// 操作可选类型
const (
	DeleteDiscountAction = 1
	StopDiscountAction   = 2
)
