package shopee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/donghui12/shopee_tool_base/pkg/constant"
	"github.com/donghui12/shopee_tool_base/pkg/logger"
	"github.com/donghui12/shopee_tool_base/pkg/pool"
)

var (
	shopeeClient      *Client
	shopeeClientForTw *Client
	once              sync.Once
	mu                sync.RWMutex
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	cookies    []*http.Cookie
	userAgent  string
	retryTimes int
	retryDelay time.Duration
	timeout    time.Duration
}

type ClientOption func(*Client)

// WithBaseURL 设置基础URL
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithRetry 设置重试次数和延迟
func WithRetry(times int, delay time.Duration) ClientOption {
	return func(c *Client) {
		c.retryTimes = times
		c.retryDelay = delay
	}
}

// InitShopeeClient 创建新的客户端
func InitShopeeClient() {
	once.Do(func() {
		// 配置 HTTP 客户端的连接池
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}

		shopeeClientForTw = &Client{
			baseURL: BaseSellerURLForTw,
			httpClient: &http.Client{
				Timeout:   30 * time.Second,
				Transport: transport,
			},
			userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			retryTimes: 3,
			retryDelay: 2 * time.Second,
			timeout:    30 * time.Second,
		}

		shopeeClient = &Client{
			baseURL: BaseSellerURL,
			httpClient: &http.Client{
				Timeout:   30 * time.Second,
				Transport: transport,
			},
			userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			retryTimes: 3,
			retryDelay: 5 * time.Second,
			timeout:    30 * time.Second,
		}
	})
}

func GetShopeeClient() *Client {
	mu.RLock()
	defer mu.RUnlock()
	return shopeeClient
}

func GetTwShopeeClient() *Client {
	mu.RLock()
	defer mu.RUnlock()
	return shopeeClientForTw
}

// Login 登录
func (c *Client) Login(account, password, vcode, loginType string) (string, error) {
	if account == "" || password == "" {
		return "", fmt.Errorf("账号或密码不能为空")
	}

	cookieString := ""
	// 构建表单数据
	loginParam := LoginParam{
		PasswordHash: MD5Hash(password),
		Remember:     false,
		OtpType:      DefaultOptType,
		Vcode:        vcode,
	}

	if loginType == LoginTypeEmail {
		loginParam.OtpType = EmailOptType
	}

	if isPhone(account) {
		loginParam.SubaccountPhone = formatPhone(account)
	} else if VerifyEmailFormat(account) {
		loginParam.SubaccountEmail = account
	} else {
		loginParam.Subaccount = account
	}
	form := loginParam.ToFormValues()

	// 创建请求
	req, err := http.NewRequest(HTTPMethodPost, c.baseURL+APIPathLogin, strings.NewReader(form.Encode()))
	if err != nil {
		return cookieString, fmt.Errorf("create login request failed: %w", err)
	}

	// 设置表单请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// 执行请求
	resp, err := c.executeWithProxy(req)
	if err != nil {
		return cookieString, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cookieString, fmt.Errorf("read login response failed: %w", err)
	}

	// 解析响应
	var commonResp CommonResponse[json.RawMessage]
	if err := json.Unmarshal(body, &commonResp); err != nil {
		return cookieString, fmt.Errorf("parse login response failed: %w", err)
	}

	// 检查响应状态
	if commonResp.Code != ResponseCodeSuccess {
		return cookieString, fmt.Errorf("login failed: %s", commonResp.Message)
	}

	if commonResp.Message == "error_server" {
		return cookieString, fmt.Errorf("请联系管理员")
	}

	if commonResp.Message == "error_need_vcode" {
		return cookieString, fmt.Errorf("需要验证码")
	}
	if commonResp.Message == "error_invalid_vcode" {
		return cookieString, fmt.Errorf("验证码错误")
	}
	if commonResp.Message == "error_name_or_password_incorrect" {
		return cookieString, fmt.Errorf("账号或密码错误")
	}
	if commonResp.Message != "" {
		return cookieString, fmt.Errorf(commonResp.Message)
	}

	cookies := resp.Header["Set-Cookie"]
	// 将 cookie 转换为字符串
	for _, cookie := range cookies {
		if strings.Contains(cookie, "SPC_CDS") {
			continue
		}
		cookie = strings.Split(cookie, ";")[0]
		cookieString += cookie + "; "
	}

	return cookieString, nil
}

// GetMerchantShopListWithRegion 获取某一地区店铺列表
func (c *Client) GetMerchantShopListWithRegion(cookies, region string) ([]MerchantShop, error) {
	allShopList, err := c.GetMerchantShopList(cookies)
	if err != nil {
		return nil, err
	}
	var shopList []MerchantShop
	for _, shop := range allShopList {
		if shop.Region == region {
			shopList = append(shopList, shop)
		}
	}
	return shopList, nil
}

// GetMerchantShopList 获取全部地区店铺列表
func (c *Client) GetSession(cookies string) (AccountInfo, error) {
	var accountInfo AccountInfo
	if cookies == "" {
		return accountInfo, fmt.Errorf("cookies不能为空")
	}

	param := CommomParam{}
	url := APIPathGetSession + "?" + param.ToFormValues().Encode()

	getSessionResp := &GetSessionResp{}
	resp, err := c.doRequest(HTTPMethodGet, url, nil, cookies)
	if err != nil {
		return accountInfo, fmt.Errorf("get session failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return accountInfo, fmt.Errorf("read response body failed: %w", err)
	}
	err = json.Unmarshal(body, &getSessionResp)
	if err != nil {
		return accountInfo, fmt.Errorf("unmarshal merchant shop list response failed: %w", err)
	}
	if getSessionResp.Code != 0 || getSessionResp.Errcode != 0 {
		return accountInfo, fmt.Errorf("获取账户信息列表失败:%s", getSessionResp.Message)
	}
	return getSessionResp.AccountInfo, nil
}

// GetMerchantShopList 获取全部地区店铺列表
func (c *Client) GetMerchantShopList(cookies string) ([]MerchantShop, error) {
	if cookies == "" {
		return nil, fmt.Errorf("cookies不能为空")
	}

	merchantShopListResp := &MerchantShopListResponse{}
	resp, err := c.doRequest(HTTPMethodGet, APIPathGetMerchantShopList, nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("get merchant shop list failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	err = json.Unmarshal(body, &merchantShopListResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal merchant shop list response failed: %w", err)
	}
	if merchantShopListResp.Error != "" {
		return nil, fmt.Errorf("获取店铺列表失败:%s", merchantShopListResp.Error)
	}
	return merchantShopListResp.Data.Shops, nil
}

// GetProductList 获取商品列表
func (c *Client) GetProductList(cookies, shopID, region, listType string) ([]int64, error) {
	if cookies == "" || shopID == "" || region == "" {
		return nil, fmt.Errorf("参数不能为空: cookies=%s, shopID=%s, region=%s", cookies, shopID, region)
	}

	var productIDs []int64
	var productIDMap sync.Map
	var wg sync.WaitGroup

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
	}

	// 获取第一页以确定总数
	firstPageParams := copyURLValues(baseParams)
	firstPageParams.Set("page_number", "1")

	APIProductList := APIPathProductList + "?" + firstPageParams.Encode()
	resp, err := c.doRequest(HTTPMethodGet, APIProductList, nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("get first page failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read first page response failed: %w", err)
	}

	var firstPageResp ProductListResponse
	err = json.Unmarshal(body, &firstPageResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal first page response failed: %w", err)
	}

	pageSize := 48
	totalPages := (firstPageResp.Data.PageInfo.Total + pageSize - 1) / pageSize

	// 使用工作池处理剩余页面
	workerPool := pool.GetWorkerPool()

	for pageNumber := 1; pageNumber <= totalPages; pageNumber++ {
		wg.Add(1)
		currentPage := pageNumber

		task := pool.Task{
			Topic: constant.TopicProduct,
			Execute: func() error {
				defer wg.Done()

				// 为每个 goroutine 创建新的参数副本
				params := copyURLValues(baseParams)
				params.Set("page_number", strconv.Itoa(currentPage))
				apiURL := APIPathProductList + "?" + params.Encode()

				resp, err := c.doRequest(HTTPMethodGet, apiURL, nil, cookies)
				if err != nil {
					logger.Error("获取商品列表失败",
						zap.Int("page", currentPage),
						zap.Error(err),
					)
					return err
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Error("读取响应失败",
						zap.Int("page", currentPage),
						zap.Error(err),
					)
					return err
				}

				var pageResp ProductListResponse
				if err := json.Unmarshal(body, &pageResp); err != nil {
					logger.Error("解析响应失败",
						zap.Int("page", currentPage),
						zap.Error(err),
					)
					return err
				}

				// 处理商品数据
				for _, product := range pageResp.Data.Products {
					productIDMap.Store(int64(product.ID), true)
					for _, campaign := range product.PromotionDetail.OngoingCampaigns {
						productIDMap.Store(int64(campaign.ProductID), true)
					}
				}

				logger.Debug("成功处理页面",
					zap.Int("page", currentPage),
					zap.Int("products", len(pageResp.Data.Products)),
				)
				return nil
			},
		}

		if err := workerPool.Submit(task); err != nil {
			logger.Error("提交任务失败",
				zap.Int("page", currentPage),
				zap.Error(err),
			)
			return nil, err
		}
	}

	// 等待所有任务完成
	wg.Wait()

	// 收集结果
	productIDMap.Range(func(key, value interface{}) bool {
		productIDs = append(productIDs, key.(int64))
		return true
	})

	return productIDs, nil
}

// GetProductList 获取商品详细信息列表
func (c *Client) GetProductDetailList(cookies, shopID, region, listType string) ([]Product, error) {
	var ProductDetailList []Product
	var productIDMap sync.Map
	var wg sync.WaitGroup

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
	}

	// 获取第一页以确定总数
	firstPageParams := copyURLValues(baseParams)
	firstPageParams.Set("page_number", "1")

	APIProductList := APIPathProductList + "?" + firstPageParams.Encode()
	resp, err := c.doRequest(HTTPMethodGet, APIProductList, nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("get first page failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read first page response failed: %w", err)
	}

	var firstPageResp ProductListResponse
	err = json.Unmarshal(body, &firstPageResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal first page response failed: %w", err)
	}

	pageSize := 48
	totalPages := (firstPageResp.Data.PageInfo.Total + pageSize - 1) / pageSize

	// 使用工作池处理剩余页面
	workerPool := pool.GetWorkerPool()

	for pageNumber := 1; pageNumber <= totalPages; pageNumber++ {
		wg.Add(1)
		currentPage := pageNumber

		task := pool.Task{
			Topic: constant.TopicProduct,
			Execute: func() error {
				defer wg.Done()

				// 为每个 goroutine 创建新的参数副本
				params := copyURLValues(baseParams)
				params.Set("page_number", strconv.Itoa(currentPage))
				apiURL := APIPathProductList + "?" + params.Encode()

				resp, err := c.doRequest(HTTPMethodGet, apiURL, nil, cookies)
				if err != nil {
					logger.Error("获取商品列表失败",
						zap.Int("page", currentPage),
						zap.Error(err),
					)
					return err
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Error("读取响应失败",
						zap.Int("page", currentPage),
						zap.Error(err),
					)
					return err
				}

				var pageResp ProductListResponse
				if err := json.Unmarshal(body, &pageResp); err != nil {
					logger.Error("解析响应失败",
						zap.Int("page", currentPage),
						zap.Error(err),
					)
					return err
				}

				// 处理商品数据
				for _, product := range pageResp.Data.Products {
					productIDMap.Store(int64(product.ID), product)
				}

				logger.Debug("成功处理页面",
					zap.Int("page", currentPage),
					zap.Int("products", len(pageResp.Data.Products)),
				)
				return nil
			},
		}

		if err := workerPool.Submit(task); err != nil {
			logger.Error("提交任务失败",
				zap.Int("page", currentPage),
				zap.Error(err),
			)
			return nil, err
		}
	}

	// 等待所有任务完成
	wg.Wait()

	// 收集结果
	productIDMap.Range(func(key, value interface{}) bool {
		ProductDetailList = append(ProductDetailList, value.(Product))
		return true
	})

	return ProductDetailList, nil
}

// GetProductListWithDayToShip 获取带出货时间的商品列表
func (c *Client) GetProductListWithDayToShip(cookies, shopID, region, listType string, dayToShip int) ([]Product, error) {
	var ProductDetailList []Product
	var productIDMap sync.Map

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
		"page_size":        {"50"},
		"source":           {SourceAttributeTool},
		"version":          {"4.0.0"},
	}

	// 获取第一页以确定总数
	firstPageParams := copyURLValues(baseParams)
	firstPageParams.Set("page_number", "1")

	APIProductList := APIPathProductDetailList + "?" + firstPageParams.Encode()
	resp, err := c.doRequest(HTTPMethodGet, APIProductList, nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("get first page failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read first page response failed: %w", err)
	}

	var firstPageResp ProductDetailListResponse
	err = json.Unmarshal(body, &firstPageResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal first page response failed: %w", err)
	}

	pageSize := 50
	totalPages := (firstPageResp.Data.PageInfo.Total + pageSize - 1) / pageSize

	logger.Info("获取页面", zap.Int("总计页:", totalPages),
		zap.Int("总计:", firstPageResp.Data.PageInfo.Total))
	// 使用工作池处理剩余页面

	currentCursor := firstPageResp.Data.PageInfo.Cursor
	for pageNumber := 0; pageNumber <= totalPages; pageNumber++ {
		currentPage := pageNumber
		params := copyURLValues(baseParams)
		params.Set("cursor", currentCursor)
		apiURL := APIPathProductDetailList + "?" + params.Encode()

		resp, err := c.doRequest(HTTPMethodGet, apiURL, nil, cookies)
		if err != nil {
			logger.Error("获取商品列表失败",
				zap.Int("page", currentPage),
				zap.Error(err),
			)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("读取响应失败",
				zap.Int("page", currentPage),
				zap.Error(err),
			)
			continue
		}

		var pageResp ProductDetailListResponse
		if err := json.Unmarshal(body, &pageResp); err != nil {
			logger.Error("解析响应失败",
				zap.Int("page", currentPage),
				zap.Error(err),
			)
			continue
		}

		// 处理商品数据
		currentTodoProcessNumber := 0
		for _, product := range pageResp.Data.List {
			if product.DaysToShip == dayToShip {
				continue
			}
			currentTodoProcessNumber += 1
			productIDMap.Store(int64(product.ID), product)
		}
		logger.Info("页面处理成功", zap.Int("当前所在页:", pageNumber),
			zap.Int("当前页面总数:", len(pageResp.Data.List)),
			zap.Int("需要处理的商品:", currentTodoProcessNumber))
		currentCursor = pageResp.Data.PageInfo.Cursor
	}

	// 收集结果
	productIDMap.Range(func(key, value interface{}) bool {
		productDetail := value.(ProductDetail)
		product := Product{
			ID:        productDetail.ID,
			ModelList: productDetail.ModelList,
		}
		ProductDetailList = append(ProductDetailList, product)
		return true
	})

	return ProductDetailList, nil
}

// GetAccessTokenWithAreaTw 获取 TW shopee accessToken
func (c *Client) GetAccessTokenWithAreaTw(shopId, code, refreshToken string) (string, string, string, error) {
	var accessToken, newRefreshToken, path, expireTimeFormatted string
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	shopIdInt, _ := strconv.ParseInt(shopId, 10, 64)
	PartnerIdInt, _ := strconv.ParseInt(constant.PartnerId, 10, 64)
	req := &GetAccessTokenReq{
		ShopId:    shopIdInt,
		PartnerId: PartnerIdInt,
	}
	if code != "" {
		path = APIPathAuthTokenForTw
		req.Code = code
	} else {
		path = APIPathAccessTokenForTw
		req.RefreshToken = refreshToken
	}
	signature := GetPublicSign(constant.PartnerId, path, timestampStr)
	APIPathAccessTokenForTw := fmt.Sprintf("%s?partner_id=%s&timestamp=%s&sign=%s",
		path, constant.PartnerId, timestampStr, signature)
	resp, err := c.doRequest(HTTPMethodPost, APIPathAccessTokenForTw, req, accessToken)
	if err != nil {
		return accessToken, newRefreshToken, expireTimeFormatted, fmt.Errorf("get first page failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return accessToken, newRefreshToken, expireTimeFormatted, fmt.Errorf("read first page response failed: %w", err)
	}
	var currentResp TWGetAccessTokenResp
	err = json.Unmarshal(body, &currentResp)
	if err != nil {
		return accessToken, newRefreshToken, expireTimeFormatted, fmt.Errorf("unmarshal first page response failed: %w", err)
	}
	if currentResp.ErrorCode != "" {
		return accessToken, newRefreshToken, expireTimeFormatted, fmt.Errorf("get access_token error %s", currentResp.Message)
	}

	accessToken = currentResp.AccessToken
	newRefreshToken = currentResp.RefreshToken
	expireInInt := currentResp.ExpireIn
	// 将 expireInInt 转换为 time.Duration 类型（秒 * time.Second）
	expireInSecond := time.Duration(expireInInt) * time.Second
	currentTime := time.Now()
	// 计算到期时间
	expireTime := currentTime.Add(expireInSecond)

	// 格式化时间为 "yyyy-MM-dd HH:mm:ss"
	expireTimeFormatted = expireTime.Format("2006-01-02 15:04:05")
	return accessToken, newRefreshToken, expireTimeFormatted, nil
}

// GetProductListWithAreaTw 获取 tw 商品列表
func (c *Client) GetProductListWithAreaTw(accessToken, shopId string) ([]int64, error) {
	var productIDs []int64
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	signature := GetShopSign(constant.PartnerId,
		APIPathProductListForTw, timestampStr, accessToken, shopId)
	// 创建基础参数
	baseParams := url.Values{
		"partner_id":       {constant.PartnerId},
		"access_token":     {accessToken},
		"item_status":      {"NORMAL"},
		"timestamp":        {timestampStr},
		"shop_id":          {shopId},
		"page_size":        {constant.DefaultPageSize},
		"sign":             {signature},
		"update_time_from": {"1264143919"},
		"update_time_to":   {timestampStr},
	}

	// 是否还有下一页
	hasNextPage := true
	offset := int64(1)

	for hasNextPage {
		params := copyURLValues(baseParams)
		offsetStr := strconv.FormatInt(int64(offset), 10)
		params.Set("offset", offsetStr)

		APIProductList := APIPathProductListForTw + "?" + params.Encode()
		resp, err := c.doRequest(HTTPMethodGet, APIProductList, nil, accessToken)
		if err != nil {
			return nil, fmt.Errorf("get first page failed: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read first page response failed: %w", err)
		}
		var currentResp TWProductListResponse
		err = json.Unmarshal(body, &currentResp)
		if err != nil {
			hasNextPage = false
			return nil, fmt.Errorf("unmarshal first page response failed: %w", err)
		}

		if currentResp.Error != "" {
			return nil, fmt.Errorf("获取商品列表信息失败: %s", currentResp.Message)
		}

		for _, item := range currentResp.Response.Items {
			productIDs = append(productIDs, item.ItemId)
		}

		offset = currentResp.Response.NextOffset
		hasNextPage = currentResp.Response.HasNextPage

	}
	logger.Info("商品列表获取完成",
		zap.Int("total_products", len(productIDs)),
	)

	return productIDs, nil
}

// UpdateProductInfoWithAreaTwItem 更新商品信息请求
type UpdateProductInfoWithAreaTwItem struct {
	DaysToShip int  `json:"days_to_ship"`
	IsPreOrder bool `json:"is_pre_order"`
}

// UpdateProductInfoWithAreaTw 更新商品信息请求
type UpdateProductInfoWithAreaTw struct {
	ItemId   int64                           `json:"item_id"`
	PreOrder UpdateProductInfoWithAreaTwItem `json:"pre_order"`
}

// UpdateProductInfoWithAreaTw 更新 tw 商品
func (c *Client) UpdateProductInfoWithAreaTw(accessToken, shopId string, itemId int64, item UpdateProductInfoWithAreaTwItem) error {
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	signature := GetShopSign(constant.PartnerId,
		APIPathProductUpdateForTw, timestampStr, accessToken, shopId)
	shopIdInt, _ := strconv.ParseInt(shopId, 10, 64)

	// 创建基础参数
	baseParams := url.Values{
		"partner_id":   {constant.PartnerId},
		"timestamp":    {timestampStr},
		"access_token": {accessToken},
		"sign":         {signature},
	}

	req := UpdateProductInfoWithAreaTw{ItemId: itemId, PreOrder: item}
	params := copyURLValues(baseParams)

	APIProductUpdate := APIPathProductUpdateForTw + "?" + params.Encode()
	APIProductUpdate = fmt.Sprintf("%s&%s=%d", APIProductUpdate, "shop_id", shopIdInt)
	resp, err := c.doRequest(HTTPMethodPost, APIProductUpdate, req, accessToken)
	if err != nil {
		return fmt.Errorf("get first page failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("read first page response failed: %w", err)
	}
	var currentResp TWProductUpdateResponse
	err = json.Unmarshal(body, &currentResp)
	if err != nil {
		return fmt.Errorf("unmarshal first page response failed: %w", err)
	}
	if currentResp.Error != "" {
		return fmt.Errorf("更新商品:%d, 失败: %s", itemId, currentResp.Message)
	}
	if currentResp.Msg == RateLimitError {
		return fmt.Errorf(RateLimitError)
	}
	if item.DaysToShip == currentResp.Response.PreOrder.DaysToShip {
		logger.Info("商品更新成功", zap.Int64("product_id", itemId))
	} else {
		logger.Error("商品更新失败", zap.Int64("product_id", itemId), zap.Any("currentResp", currentResp))
	}

	return nil
}

// GetProductBaseInfoWithAreaTw 获取商品基础信息，主要是为了获取 出货时间，验证更新结果
func (c *Client) GetProductBaseInfoWithAreaTw(accessToken, shopId string, itemIdList []int64) ([]ProductBaseInfoWithAreaTwComplate, error) {
	var productInfos []ProductBaseInfoWithAreaTwComplate
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	signature := GetShopSign(constant.PartnerId,
		APIPathGetBaseProductInfo, timestampStr, accessToken, shopId)

	// 将 itemIDs 转换为逗号分隔的字符串
	itemIDStrs := make([]string, len(itemIdList))
	for i, id := range itemIdList {
		itemIDStrs[i] = strconv.FormatInt(id, 10)
	}
	itemIDList := strings.Join(itemIDStrs, ",")
	// 创建基础参数
	baseParams := url.Values{
		"partner_id":   {constant.PartnerId},
		"access_token": {accessToken},
		"timestamp":    {timestampStr},
		"shop_id":      {shopId},
		"sign":         {signature},
		"item_id_list": {itemIDList},
	}

	params := copyURLValues(baseParams)

	APIProductList := APIPathGetBaseProductInfo + "?" + params.Encode()
	resp, err := c.doRequest(HTTPMethodGet, APIProductList, nil, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get first page failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read first page response failed: %w", err)
	}
	var currentResp ProductBaseInfoWithAreaTwResp
	err = json.Unmarshal(body, &currentResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal first page response failed: %w", err)
	}

	if currentResp.Error != "" {
		return nil, fmt.Errorf("获取商品列表信息失败: %s", currentResp.Message)
	}

	for _, item := range currentResp.Response.ItemList {
		productInfos = append(productInfos, ProductBaseInfoWithAreaTwComplate{
			ItemId:     item.ItemId,
			DaysToShip: item.PreOrder.DaysToShip,
			IsPreOrder: item.PreOrder.IsPreOrder,
		})
	}

	return productInfos, nil
}

// copyURLValues 创建 url.Values 的深拷贝
func copyURLValues(src url.Values) url.Values {
	dst := make(url.Values, len(src))
	for k, vs := range src {
		dst[k] = append([]string(nil), vs...)
	}
	return dst
}

// UpdateProductInfoResponse 更新商品信息响应
type UpdateProductInfoResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
	Data        struct {
		ProductID int64 `json:"product_id"`
	} `json:"data"`
}

// BatchUpdateProductInfoRespItem 更新商品信息响应
type BatchUpdateProductInfoRespItem struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
	ID          int64  `json:"id"`
}

// BatchUpdateProductInfoResponse 批量更新商品信息响应
type BatchUpdateProductInfoResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
	Data        struct {
		Result []BatchUpdateProductInfoRespItem `json:"result"`
	} `json:"data"`
}

// UpdateProductInfo 更新商品信息
func (c *Client) UpdateProductInfo(updateProductInfoReq UpdateProductInfoReq) error {
	SPC_CDS := uuid.New().String()
	updateProductInfoReq.Cookies += "SPC_CDS=" + SPC_CDS + ";"

	updateProductInfoParams := url.Values{}
	updateProductInfoParams.Set("SPC_CDS", SPC_CDS)
	updateProductInfoParams.Set("SPC_CDS_VER", "2")
	updateProductInfoParams.Set("cnsc_shop_id", updateProductInfoReq.ShopID)
	updateProductInfoParams.Set("cbsc_shop_region", updateProductInfoReq.Region)

	req := &UpdateProductInfoRequest{
		ProductID: updateProductInfoReq.ProductId,
		ProductInfo: ProductInfo{
			EnableModelLevelDts: false,
		},
		IsDraft: false,
	}

	if updateProductInfoReq.DaysToShip != 0 {
		preOrderInfo := PreOrderInfo{
			PreOrder:   true,
			DaysToShip: updateProductInfoReq.DaysToShip,
		}
		if updateProductInfoReq.DaysToShip == 2 {
			preOrderInfo.PreOrder = false
		}
		req.ProductInfo.PreOrderInfo = preOrderInfo
	} else {
		req.ProductInfo.Unlisted = updateProductInfoReq.ProductStatus.Unlisted
	}

	APIUpdateProductInfo := APIPathUpdateProductInfo + "?" + updateProductInfoParams.Encode()
	resp, err := c.doRequest(HTTPMethodPost, APIUpdateProductInfo, req, updateProductInfoReq.Cookies)
	if err != nil {
		return fmt.Errorf("update product info failed, request error: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update product info failed, status code: %d, message: %s", resp.StatusCode, string(body))
	}
	var updateProductInfoResp UpdateProductInfoResponse
	err = json.Unmarshal(body, &updateProductInfoResp)
	if err != nil {
		return fmt.Errorf("unmarshal update product info response failed: %w", err)
	}
	if updateProductInfoResp.Code != ResponseCodeSuccess {
		return fmt.Errorf("update product info failed, message: %s", updateProductInfoResp.UserMessage)
	}

	return nil
}

// BatchUpdateProductInfoWithV3 使用 V3 接口批量更新商品信息
func (c *Client) BatchUpdateProductInfoWithV3(updateProductInfoReq UpdateProductInfoReq,
	shopIdList []int64, source, action string) ([]BatchUpdateProductInfoRespItem, error) {
	SPC_CDS := uuid.New().String()
	updateProductInfoReq.Cookies += "SPC_CDS=" + SPC_CDS + ";"

	updateProductInfoParams := url.Values{}
	updateProductInfoParams.Set("SPC_CDS", SPC_CDS)
	updateProductInfoParams.Set("SPC_CDS_VER", "2")
	updateProductInfoParams.Set("cnsc_shop_id", updateProductInfoReq.ShopID)
	updateProductInfoParams.Set("cbsc_shop_region", updateProductInfoReq.Region)
	updateProductInfoParams.Set("version", "3.1.0")
	updateProductInfoParams.Set("source", source)

	batchUpdateProductInfoReq := make([]*BatchUpdateProductInfoItem, 0, len(shopIdList))
	for _, shopId := range shopIdList {
		currentReq := &BatchUpdateProductInfoItem{
			ID: shopId,
		}
		if updateProductInfoReq.DaysToShip > 0 {
			currentReq.DaysToShip = updateProductInfoReq.DaysToShip
			currentReq.PreOrder = true
			if updateProductInfoReq.DaysToShip == 2 {
				currentReq.PreOrder = false
			}
		} else {
			// 默认是上架
			currentReq.Unlisted = UnlistedStatus
			if action == constant.ActionUnlisted {
				// action = 下架时，需要下架
				currentReq.Unlisted = ListedStatus
			}
		}
		batchUpdateProductInfoReq = append(batchUpdateProductInfoReq, currentReq)

	}

	APIUpdateProductInfo := APIPathBatchUpdateProductInfo + "?" + updateProductInfoParams.Encode()
	resp, err := c.doRequest(HTTPMethodPost, APIUpdateProductInfo, batchUpdateProductInfoReq, updateProductInfoReq.Cookies)
	if err != nil {
		return nil, fmt.Errorf("update product info failed, request error: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode == RateLimitCode {
		return nil, fmt.Errorf(RateLimitError)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update product info failed, status code: %d, message: %s", resp.StatusCode, string(body))
	}
	var updateProductInfoResp BatchUpdateProductInfoResponse
	err = json.Unmarshal(body, &updateProductInfoResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal update product info response failed: %w", err)
	}
	if updateProductInfoResp.Code != ResponseCodeSuccess {
		return updateProductInfoResp.Data.Result, fmt.Errorf("update product info failed, message: %s", updateProductInfoResp.UserMessage)
	}

	return updateProductInfoResp.Data.Result, nil
}

// BatchUpdateProductInfoWithFile 使用 excel 接口批量更新商品信息
func (c *Client) BatchUpdateProductInfoWithFile(updateProductInfoReq UpdateProductInfoReq, filename string) error {
	SPC_CDS := uuid.New().String()
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)

	updateProductInfoParams := url.Values{}
	updateProductInfoParams.Set("SPC_CDS", SPC_CDS)
	updateProductInfoParams.Set("SPC_CDS_VER", "2")
	updateProductInfoParams.Set("cnsc_shop_id", updateProductInfoReq.ShopID)
	updateProductInfoParams.Set("cbsc_shop_region", updateProductInfoReq.Region)
	updateProductInfoParams.Set("timestamp", timestampStr)

	APIUpdateProductInfo := APIPathBatchUpdateProductInfoWithFile + "?" + updateProductInfoParams.Encode()
	resp, err := c.doRequestWithFile(HTTPMethodPost, APIUpdateProductInfo, updateProductInfoReq.Cookies, "file", filename)
	if err != nil {
		return fmt.Errorf("update product info failed, request error: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode == RateLimitCode {
		return fmt.Errorf(RateLimitError)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update product info failed, status code: %d, message: %s", resp.StatusCode, string(body))
	}
	var updateProductInfoResp BatchUpdateProductInfoResponse
	err = json.Unmarshal(body, &updateProductInfoResp)
	if err != nil {
		return fmt.Errorf("unmarshal update product info response failed: %w", err)
	}
	if updateProductInfoResp.Code != ResponseCodeSuccess &&
		updateProductInfoResp.Code != ProcessCode {
		return fmt.Errorf("update product info failed, message: %+v\b", updateProductInfoResp)
	}

	return nil

}

// GetOrSetShop 获取或设置店铺
func (c *Client) SwitchMerchantShop(cookies, region, shopId string) error {
	// 构造 URL 参数
	param := CommomParam{
		ShopId: shopId,
		Region: region,
	}
	shopIdInt, _ := strconv.ParseInt(shopId, 10, 64)
	respBody := GetOrSetShopReq{
		ShopId: shopIdInt,
	}

	url := APIPathSwitchMerchantShop + "?" + param.ToFormValues().Encode()

	resp, err := c.doRequest(HTTPMethodPost, url, respBody, cookies)
	if err != nil {
		return fmt.Errorf("switch_merchant_shop failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("switch_merchant_shop request failed: %s", string(body))
	}

	var commonResp CommonResponse[json.RawMessage]
	err = json.Unmarshal(body, &commonResp)
	if err != nil {
		return fmt.Errorf("unmarshal switch_merchant_shop response failed: %w", err)
	}
	if commonResp.Message != SuccessMessage {
		return fmt.Errorf("switch_merchant_shop response failed: %s", commonResp.Message)
	}
	return nil

}

// GetOrSetShop 获取或设置店铺
func (c *Client) GetOrSetShop(cookies string) error {
	// 构建请求
	SPC_CDS := uuid.New().String()
	cookies += "SPC_CDS=" + SPC_CDS + ";"
	param := CommomParam{}

	APIGetOrSetShop := APIPathGetOrSetShop + "?" + param.ToFormValues().Encode()
	respBody := map[string]interface{}{}
	resp, err := c.doRequestWithProxy(HTTPMethodPost, APIGetOrSetShop, respBody, cookies)
	if err != nil {
		return fmt.Errorf("get or set shop request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get or set shop request failed: %s", string(body))
	}

	var commonResp CommonResponse[json.RawMessage]
	err = json.Unmarshal(body, &commonResp)
	if err != nil {
		return fmt.Errorf("unmarshal get or set shop response failed: %w", err)
	}
	if commonResp.Message != SuccessMessage {
		return fmt.Errorf("get or set shop response failed: %s", commonResp.Message)
	}
	return nil

}

// GetInactiveProducts 获取不活跃的商品信息
func (c *Client) GetInactiveProducts(cookies, shopId, region string, batch int) ([]int64, error) {
	var productIdList []int64
	productList, err := c.GetProductDetailList(cookies, shopId, region, ListTypeLive)
	if err != nil {
		logger.Info("获取商品详细信息失败", zap.Any("err", err))
		return productIdList, err
	}

	logger.Info("获取商品详细信息", zap.Any("total product", len(productList)))

	var inActiveProductList []Product

	// 首先获取：喜欢，销量，点击 为 0 的商品信息
	for _, product := range productList {
		if product.Statistics.IsActiveProduct() {
			inActiveProductList = append(inActiveProductList, product)
		}
	}

	// 根据创建时间进行排序，获取前 50 个商品Id
	// 根据创建时间进行排序（升序：越旧的在前）
	sort.Slice(inActiveProductList, func(i, j int) bool {
		return inActiveProductList[i].CreateTime < inActiveProductList[j].CreateTime
	})

	// 取前 50 个商品 ID
	limit := batch
	if len(inActiveProductList) < limit {
		limit = len(inActiveProductList)
	}
	for i := 0; i < limit; i++ {
		productIdList = append(productIdList, int64(inActiveProductList[i].ID))
	}

	return productIdList, nil
}

// ListedOrUnlistedProducts 上下架商品
func (c *Client) ListedOrUnlistedProducts(shopId, cookies, region string, listStatus bool, productIds []int64) int64 {
	var successfulNumber int64
	for _, productId := range productIds {
		time.Sleep(time.Second * 2)
		currentUpdateProductInfoReq := UpdateProductInfoReq{
			ProductId: productId,
			ShopID:    shopId,
			Region:    region,
			Cookies:   cookies,
			ProductStatus: ProductStatusInfo{
				Unlisted: listStatus,
			},
		}
		if err := c.UpdateProductInfo(currentUpdateProductInfoReq); err != nil {
			logger.Info("商品更新失败", zap.String("shopId", shopId),
				zap.Int64("productId", productId), zap.Any("error", err))
			continue
		}
		successfulNumber += 1
	}

	logger.Info("更新商品完成", zap.Int64("总计更新", successfulNumber))
	return successfulNumber
}

// get_discount_list
func (c *Client) GetDiscountList(cookies, shopId, region string) ([]Discount, error) {
	discountList := []Discount{}
	discountIdMap := make(map[int64]int)
	// now := time.Now()
	// // 取过去一个月的时间范围
	// periodFrom := now.AddDate(0, -1, 0).Unix() // 一个月前
	// periodTo := now.Unix()                     // 当前时间

	req := &getDiscountListRequest{
		DiscountType: 1,
		Offset:       0,
		Limit:        10,
		TimeStatus:   2,
	}

	// 构造 URL 参数
	param := CommomParam{
		ShopId: shopId,
		Region: region,
	}

	url := APIPathGetDiscountList + "?" + param.ToFormValues().Encode()

	for i := 0; i <= 100; i++ {
		resp, err := c.doRequestWithLocalProxy(HTTPMethodPost, url, req, cookies)
		if err != nil {
			return nil, fmt.Errorf("get discount list failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body failed: %w", err)
		}
		data, err := ParseCommonResponse[DiscountList](body)
		if err != nil {
			return nil, err
		}
		for _, discount := range data.Discounts {
			if _, ok := discountIdMap[discount.SellerDiscount.DiscountID]; ok {
				continue
			}
			discountIdMap[discount.SellerDiscount.DiscountID] = 1
			discountList = append(discountList, discount)
		}
		if len(discountList) >= data.TotalCount {
			break
		}
		req.Offset += 1
	}
	return discountList, nil
}

func (c *Client) GetDiscountItem(cookies, shopId, region string, discountId int64) ([]DiscountItemList, error) {
	data := &DiscountItemData{}
	var discountItemList []DiscountItemList
	totalItemCount := 0
	totalItemStatusCount := 0

	req := &getDiscountItemRequest{
		PromotionId: discountId,
		Offset:      0,
		Limit:       100,
	}

	// 构造 URL 参数
	param := CommomParam{shopId, region}

	url := APIPathGetDiscountItem + "?" + param.ToFormValues().Encode()
	for i := 0; i < 100; i++ {
		resp, err := c.doRequest(HTTPMethodPost, url, req, cookies)
		if err != nil {
			return discountItemList, fmt.Errorf("get discount list failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return discountItemList, fmt.Errorf("read response body failed: %w", err)
		}
		data, err = ParseCommonResponse[DiscountItemData](body)
		if err != nil {
			return discountItemList, fmt.Errorf("解析失败: %w", err)
		}
		discountItemStatusMap := make(map[int64]int)
		for _, discountItem := range data.ItemInfo {
			if discountItem.Status != 1 {
				totalItemStatusCount += 1
			}
			discountItemStatusMap[discountItem.ItemID] = discountItem.Status
		}

		for _, discountItem := range data.DiscountItemList {
			if status, ok := discountItemStatusMap[discountItem.ItemID]; ok && status == 1 {
				discountItemList = append(discountItemList, discountItem)
			}
		}
		totalItemCount += len(data.ItemInfo)
		if totalItemCount >= data.TotalCount {
			break
		}
		req.Offset += 100
	}
	return discountItemList, nil
}

// UpdateDiscountItem 更新折扣商品
func (c *Client) UpdateDiscountItem(cookies, shopId, region string, req UpdateDiscountItemRequest) (int, error) {
	data := &UpdateSellerDiscountItemsResp{}

	// 构造 URL 参数
	param := CommomParam{shopId, region}

	url := APIPathUpdateDiscountItem + "?" + param.ToFormValues().Encode()

	resp, err := c.doRequest(HTTPMethodPost, url, req, cookies)
	if err != nil {
		return 0, fmt.Errorf("get discount list failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response body failed: %w", err)
	}
	data, err = ParseCommonResponse[UpdateSellerDiscountItemsResp](body)
	if err != nil {
		return 0, fmt.Errorf("解析失败: %w", err)
	}
	if len(data.ErrorList) > 0 {
		return 0, fmt.Errorf("更新失败: %s", data.ErrorList[0].ErrorMessage)
	}
	return data.SuccessCount, nil
}

// DeleteProducts 删除商品
func (c *Client) DeleteProducts(shopId, cookies, region string, productIds []int64) (int64, error) {
	SPC_CDS := uuid.New().String()
	var successfulNumber int64

	var deleteProductReq DeleteProductReq
	deleteProductReq.ProductIdList = productIds

	deleteProductParams := url.Values{}
	deleteProductParams.Set("version", "3.1.0")
	deleteProductParams.Set("SPC_CDS", SPC_CDS)
	deleteProductParams.Set("SPC_CDS_VER", "2")
	deleteProductParams.Set("cnsc_shop_id", shopId)
	deleteProductParams.Set("cbsc_shop_region", region)

	APIUpdateProductInfo := APIPathDeleteProduct + "?" + deleteProductParams.Encode()
	resp, err := c.doRequest(HTTPMethodPost, APIUpdateProductInfo, deleteProductReq, cookies)
	if err != nil {
		return successfulNumber, fmt.Errorf("delete product info failed, request error: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return successfulNumber, fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return successfulNumber, fmt.Errorf("delete product info failed, status code: %d, message: %s", resp.StatusCode, string(body))
	}
	var updateProductInfoResp UpdateProductInfoResponse
	err = json.Unmarshal(body, &updateProductInfoResp)
	if err != nil {
		return successfulNumber, fmt.Errorf("unmarshal update product info response failed: %w", err)
	}
	if updateProductInfoResp.Code != ResponseCodeSuccess {
		return successfulNumber, fmt.Errorf("delete product info failed, message: %s", updateProductInfoResp.UserMessage)
	}

	successfulNumber = int64(len(productIds))

	return successfulNumber, nil
}

func (c *Client) doRequestWithProxy(method, path string, reqBody interface{}, cookies string) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if reqBody != nil {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置 JSON 请求头
	req.Header.Set("Cookie", cookies)
	c.setCommonHeaders(req)

	resp, err := c.executeWithLocalProxy(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) doRequestWithLocalProxy(method, path string, reqBody interface{}, cookies string) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if reqBody != nil {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置 JSON 请求头
	req.Header.Set("Cookie", cookies)
	c.setCommonHeaders(req)

	resp, err := c.executeWithLocalProxy(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}

	return resp, nil
}

// delete_discount
func (c *Client) DeleteDiscounts(cookies, shopId, region string, discountID []int64, action int) (int64, error) {
	var successfulNumber int64

	for _, promotionId := range discountID {
		// 构建请求体
		reqBody := DeleteDiscountReq{
			PromotionID: promotionId,
			Action:      action,
		}

		// 构建请求 URL
		param := CommomParam{ShopId: shopId, Region: region}
		apiURL := APIPathDeleteDiscount + "?" + param.ToFormValues().Encode()

		// 发起请求
		resp, err := c.doRequest(HTTPMethodPost, apiURL, reqBody, cookies)
		if err != nil {
			logger.Error("请求失败: ",
				zap.Any("折扣Id:", promotionId), zap.Error(err))
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logger.Error("读取响应失败: ",
				zap.Any("折扣Id:", promotionId), zap.Error(err))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			logger.Error("请求状态异常: ",
				zap.Any("折扣Id:", promotionId), zap.Error(err))
			continue
		}

		data, err := ParseCommonResponse[DeleteDiscountData](body)
		if err != nil {
			logger.Error("响应解析失败: ",
				zap.Any("折扣Id:", promotionId), zap.Error(err))
			continue
		}

		if len(data.ErrorList) > 0 {
			logger.Error("部分失败: ",
				zap.Any("折扣Id:", promotionId), zap.Error(err))
			continue
		}

		successfulNumber++
	}
	return successfulNumber, nil
}

// CopyCreateDiscounts 复制创建折扣
func (c *Client) CopyCreateDiscount(cookies, shopId, region string, discount Discount) (int64, error) {
	// 构建请求体
	currentReq := CreateDiscountReq{}
	currentReq.ConvertFromDiscount(discount)

	// 构建请求 URL
	param := CommomParam{
		ShopId: shopId,
		Region: region,
	}

	apiURL := APIPathCreateDiscount + "?" + param.ToFormValues().Encode()

	// 发起请求
	resp, err := c.doRequest(HTTPMethodPost, apiURL, currentReq, cookies)
	if err != nil {
		log.Printf("请求失败: discount=%d, err=%v",
			discount.SellerDiscount.DiscountID, err)
		return 0, fmt.Errorf("复制创建折扣: %d, 失败: %w", discount.SellerDiscount.DiscountID, err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("请求失败: discount=%d, err=%v",
			discount.SellerDiscount.DiscountID, err)
		return 0, fmt.Errorf("复制创建折扣: %d, 失败: %w", discount.SellerDiscount.DiscountID, err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("请求失败: discount=%d, status=%d, body=%s",
			discount.SellerDiscount.DiscountID, resp.StatusCode, string(body))
		return 0, fmt.Errorf("复制创建折扣: %d, 失败: %w", discount.SellerDiscount.DiscountID, err)
	}
	data, err := ParseCommonResponse[CreateDiscountData](body)
	if err != nil {
		log.Printf("请求失败: discount=%d, err=%v",
			discount.SellerDiscount.DiscountID, err)
		return 0, fmt.Errorf("复制创建折扣: %d, 失败: %w", discount.SellerDiscount.DiscountID, err)
	}

	return data.PromationId, nil
}

// doRequest 执行请求
func (c *Client) doRequest(method, path string, reqBody interface{}, cookies string) (*http.Response, error) {
	return c.doRequestWithContext(context.Background(), method, path, reqBody, cookies)
}

// doRequestWithContext 执行带上下文的请求
func (c *Client) doRequestWithContext(ctx context.Context, method, path string, reqBody interface{}, cookies string) (*http.Response, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	url := c.baseURL + path

	var bodyReader io.Reader
	if reqBody != nil {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置 JSON 请求头
	req.Header.Set("Cookie", cookies)
	c.setCommonHeaders(req)

	resp, err := c.executeWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) doRequestWithFile(method, path, cookies, fileFieldName, filePath string) (*http.Response, error) {
	url := c.baseURL + path

	// 创建 multipart/form-data 请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 如果有文件上传，添加文件字段
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("open file failed: %w", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fileFieldName, filepath.Base(filePath)) // fileFieldName 是表单中的文件字段名
		if err != nil {
			return nil, fmt.Errorf("create form field [%s]: %w", fileFieldName, err)
		}

		_, err = io.Copy(part, file) // 将文件内容写入表单
		if err != nil {
			return nil, fmt.Errorf("copy file to form failed: %w", err)
		}
	}

	// 关闭 writer，确保所有表单数据已写入
	writer.Close()

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置 Cookie
	req.Header.Set("Cookie", cookies)
	// 设置 Content-Type 头
	req.Header.Set("Content-Type", writer.FormDataContentType()) // 设置 multipart/form-data 的 content-type

	// 执行请求
	resp, err := c.executeWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}

	return resp, nil
}
