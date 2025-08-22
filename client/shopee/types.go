package shopee

// IsActiveProduct 判断是否是活跃商品
func (stat ProductStatistics) IsActiveProduct() bool {
	return stat.LikedCount == 0 && stat.SoldCount == 0 && stat.ViewCount == 0
}

type Model struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PageInfo struct {
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	Total      int    `json:"total"`
	Cursor     string `json:"cursor"`
}

type MerchantShopListData struct {
	Shops []MerchantShop `json:"shops"`
}

type UpdateProductInfoData struct {
	ProductID int64 `json:"product_id"`
}

type ProductListData struct {
	Products []Product `json:"products"`
	PageInfo PageInfo  `json:"page_info"`
}

// 商品详细信息列表
type ProductDetailListData struct {
	List     []ProductDetail `json:"list"`
	PageInfo PageInfo        `json:"page_info"`
}

type ProductDetail struct {
	ID            int     `json:"id"`
	DaysToShip    int     `json:"days_to_ship"`
	EstimatedDays int     `json:"estimated_days"`
	PreOrder      bool    `json:"pre_order"`
	ModelList     []Model `json:"model_list"`
}

// ProductListResponse 商品列表响应
type ProductListResponse struct {
	Code        int             `json:"code"`
	Message     string          `json:"message"`
	UserMessage string          `json:"user_message"`
	Data        ProductListData `json:"data"`
}

// ProductListResponse 商品列表响应
type ProductDetailListResponse struct {
	Code        int                   `json:"code"`
	Message     string                `json:"message"`
	UserMessage string                `json:"user_message"`
	Data        ProductDetailListData `json:"data"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	ProductID   int64   `json:"product_id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       int     `json:"stock,omitempty"`
}

// TWGetAccessTokenResp tw 获取 accessToken resp
type TWGetAccessTokenResp struct {
	ErrorCode    string `json:"error"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	ExpireIn     int64  `json:"expire_in"`
	RefreshToken string `json:"refresh_token"`
}

type TWProductItem struct {
	ItemId     int64  `json:"item_id"`
	ItemStatus string `json:"item_status"`
}

type TWProductListData struct {
	Items       []TWProductItem `json:"item"`
	TotalCount  int64           `json:"total_count"`
	HasNextPage bool            `json:"has_next_page"`
	NextOffset  int64           `json:"next_offset"`
}

// TWProductListResponse 台湾商品列表响应
type TWProductListResponse struct {
	Message  string            `json:"message"`
	Error    string            `json:"error"`
	Response TWProductListData `json:"response"`
}

// TWProductUpdatePreOrder 更新商品信息请求
type TWProductUpdatePreOrder struct {
	DaysToShip int  `json:"days_to_ship"`
	IsPreOrder bool `json:"is_pre_order"`
}

type TWProductUpdateData struct {
	ItemId   int64                   `json:"item_id"`
	ItemName string                  `json:"item_name"`
	PreOrder TWProductUpdatePreOrder `json:"pre_order"`
}

// TWProductUpdateResponse 台湾更新商品响应
type TWProductUpdateResponse struct {
	Message  string              `json:"message"`
	Msg      string              `json:"msg"`
	Error    string              `json:"error"`
	Response TWProductUpdateData `json:"response"`
}

// TWUpdateProductRequest 台湾更新商品请求
type TWUpdateProductRequest struct {
	ProductID   int64   `json:"product_id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       int     `json:"stock,omitempty"`
}

type MerchantShop struct {
	Region string `json:"region"`
	ShopID int64  `json:"shop_id"`
}

type MerchantShopList struct {
	Shops []MerchantShop `json:"shops"`
}

type MerchantShopListResponse struct {
	Code      int              `json:"code"`
	ErrorCode int              `json:"errcode"`
	Error     string           `json:"error"`
	Message   string           `json:"message"`
	Data      MerchantShopList `json:"data"`
}

type ProductBaseInfoWithAreaTwResp struct {
	Error    string                        `json:"error"`
	Message  string                        `json:"message"`
	Response ProductBaseInfoListWithAreaTw `json:"response"`
}

type ProductBaseInfoListWithAreaTw struct {
	ItemList []ProductBaseInfoWithAreaTw `json:"item_list"`
}

// ProductPreOrderInfoWithAreaTwItem 更新商品信息请求
type ProductPreOrderInfoWithAreaTwItem struct {
	DaysToShip int  `json:"days_to_ship"`
	IsPreOrder bool `json:"is_pre_order"`
}

// ProductBaseInfoWithAreaTw 商品基础信息返回参数
type ProductBaseInfoWithAreaTw struct {
	ItemId   int64                             `json:"item_id"`
	PreOrder ProductPreOrderInfoWithAreaTwItem `json:"pre_order"`
}

// ProductBaseInfoWithAreaTw tw基本商品信息--平铺展示
type ProductBaseInfoWithAreaTwComplate struct {
	ItemId     int64 `json:"item_id"`
	DaysToShip int   `json:"days_to_ship"`
	IsPreOrder bool  `json:"is_pre_order"`
}

// BatchUpdateProductInfoItem 批量更新商品信息
type BatchUpdateProductInfoItem struct {
	ID         int64 `json:"id"`
	PreOrder   bool  `json:"pre_order,omitempty"`
	DaysToShip int   `json:"days_to_ship,omitempty"`
	Unlisted   bool  `json:"unlisted"`
}
