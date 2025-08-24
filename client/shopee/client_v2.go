package shopee

import (
	"context"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// 重构后的登录方法示例
func (c *Client) LoginV2(account, password, vcode, loginType string) (string, error) {
	// 参数验证
	if account == "" || password == "" {
		return "", NewValidationError("账号或密码不能为空")
	}

	// 构建登录参数
	loginParam := LoginParam{
		PasswordHash: MD5Hash(password),
		Remember:     false,
		OtpType:      DefaultOptType,
		Vcode:        vcode,
	}

	if loginType == LoginTypeEmail {
		loginParam.OtpType = EmailOptType
	}

	if IsPhone(account) {
		loginParam.SubaccountPhone = formatPhone(account)
	} else if VerifyEmailFormat(account) {
		loginParam.SubaccountEmail = account
	} else {
		loginParam.Subaccount = account
	}

	// 使用新的请求管理器
	rm := NewRequestManager(c)

	// 设置表单请求头
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	var commonResp CommonResponse[any]
	err := rm.DoRequestWithJSONResponse(
		context.Background(),
		HTTPMethodPost,
		APIPathLogin,
		loginParam.ToFormValues().Encode(),
		"",
		&commonResp,
		WithHeaders(headers),
	)

	if err != nil {
		return "", err
	}

	// 特殊的登录错误处理
	if businessErr := CheckBusinessError(commonResp.Code, commonResp.Message); businessErr != nil {
		return "", businessErr
	}

	// 从最后一次请求的响应中提取cookies（这里需要改造请求管理器来支持获取response headers）
	// 这是一个需要进一步优化的部分
	cookieString := ""
	return cookieString, nil
}

// 重构后的获取店铺列表方法示例
func (c *Client) GetMerchantShopListV2(cookies string) ([]MerchantShop, error) {
	if cookies == "" {
		return nil, NewValidationError("cookies不能为空")
	}

	rm := NewRequestManager(c)

	// 使用新的通用响应解析
	data, err := DoRequestWithCommonResponse[MerchantShopListData](
		rm,
		context.Background(),
		HTTPMethodGet,
		APIPathGetMerchantShopList,
		nil,
		cookies,
	)

	if err != nil {
		return nil, err
	}

	return data.Shops, nil
}

// 重构后的获取商品列表方法示例（简化版）
func (c *Client) GetProductListV2(cookies, shopID, region, listType string) ([]int64, error) {
	if cookies == "" || shopID == "" || region == "" {
		return nil, NewValidationError("参数不能为空")
	}

	var productIDs []int64
	SPC_CDS := uuid.New().String()
	cookies += "SPC_CDS=" + SPC_CDS + ";"

	// 创建基础参数
	baseParams := url.Values{
		"SPC_CDS":          {SPC_CDS},
		"SPC_CDS_VER":      {"2"},
		"list_type":        {listType},
		"need_ads":         {"true"},
		"cnsc_shop_id":     {shopID},
		"cbsc_shop_region": {region},
		"page_size":        {"48"},
		"page_number":      {"1"},
	}

	rm := NewRequestManager(c)
	apiURL := APIPathProductList + "?" + baseParams.Encode()

	// 获取第一页确定总数
	firstPageResp, err := DoRequestWithCommonResponse[ProductListData](
		rm,
		context.Background(),
		HTTPMethodGet,
		apiURL,
		nil,
		cookies,
	)

	if err != nil {
		return nil, err
	}

	// 处理第一页数据
	for _, product := range firstPageResp.Products {
		productIDs = append(productIDs, int64(product.ID))
	}

	// 计算总页数并处理剩余页面
	pageSize := 48
	totalPages := (firstPageResp.PageInfo.Total + pageSize - 1) / pageSize

	// 这里可以继续使用工作池或者简化为顺序处理
	for pageNumber := 2; pageNumber <= totalPages; pageNumber++ {
		params := copyURLValues(baseParams)
		params.Set("page_number", strconv.Itoa(pageNumber))
		apiURL := APIPathProductList + "?" + params.Encode()

		pageResp, err := DoRequestWithCommonResponse[ProductListData](
			rm,
			context.Background(),
			HTTPMethodGet,
			apiURL,
			nil,
			cookies,
			WithSkipErrorLog(true), // 跳过错误日志避免日志过多
		)

		if err != nil {
			// 记录错误但继续处理其他页面
			continue
		}

		for _, product := range pageResp.Products {
			productIDs = append(productIDs, int64(product.ID))
		}
	}

	return productIDs, nil
}

// 重构后的更新商品信息方法示例
func (c *Client) UpdateProductInfoV2(updateReq UpdateProductInfoReq) error {
	SPC_CDS := uuid.New().String()
	updateReq.Cookies += "SPC_CDS=" + SPC_CDS + ";"

	params := url.Values{
		"SPC_CDS":          {SPC_CDS},
		"SPC_CDS_VER":      {"2"},
		"cnsc_shop_id":     {updateReq.ShopID},
		"cbsc_shop_region": {updateReq.Region},
	}

	req := &UpdateProductInfoRequest{
		ProductID: updateReq.ProductId,
		ProductInfo: ProductInfo{
			EnableModelLevelDts: false,
		},
		IsDraft: false,
	}

	if updateReq.DaysToShip != 0 {
		preOrderInfo := PreOrderInfo{
			PreOrder:   true,
			DaysToShip: updateReq.DaysToShip,
		}
		if updateReq.DaysToShip == 2 {
			preOrderInfo.PreOrder = false
		}
		req.ProductInfo.PreOrderInfo = preOrderInfo
	} else {
		req.ProductInfo.Unlisted = updateReq.ProductStatus.Unlisted
	}

	rm := NewRequestManager(c)
	apiURL := APIPathUpdateProductInfo + "?" + params.Encode()

	// 使用新的请求管理器
	_, err := DoRequestWithCommonResponse[UpdateProductInfoData](
		rm,
		context.Background(),
		HTTPMethodPost,
		apiURL,
		req,
		updateReq.Cookies,
	)

	return err
}
