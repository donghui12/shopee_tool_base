# Shopee Client 优化说明

本文档介绍了 Shopee Client 的模块化优化，提供了更加清晰、可维护和高效的客户端实现。

## 优化内容

### 1. 错误处理模块 (`errors.go`)

#### 功能特性
- **统一错误类型**: 所有错误都包装为 `ShopeeError` 类型
- **错误分类**: 按照错误类型进行分类（认证、网络、验证等）
- **重试判断**: 自动判断错误是否可重试
- **错误链**: 支持错误包装和展开

#### 使用示例
```go
// 创建不同类型的错误
authErr := NewAuthError(401, "token expired", nil)
networkErr := NewNetworkError("connection failed", originalErr)
validationErr := NewValidationError("invalid parameters")

// 检查错误类型
if shopeeErr, ok := err.(*ShopeeError); ok {
    if shopeeErr.IsType(ErrTypeAuth) {
        // 处理认证错误
    }
    if shopeeErr.IsRetryable() {
        // 错误可重试
    }
}
```

### 2. 请求管理模块 (`request.go`)

#### 功能特性
- **统一请求接口**: 提供一致的请求方法
- **灵活配置**: 支持超时、重试、代理等配置
- **自动解析**: 自动处理JSON响应和通用响应格式
- **错误处理**: 集成错误处理逻辑

#### 使用示例
```go
// 创建请求管理器
client := NewDefaultClient()
rm := NewRequestManager(client)

// 基础请求
resp, err := rm.DoRequest(ctx, "GET", "/api/path", nil, cookies)

// 带响应解析的请求
body, err := rm.DoRequestWithResponse(ctx, "POST", "/api/path", reqData, cookies)

// 通用响应格式解析
data, err := rm.DoRequestWithCommonResponse[ResponseType](ctx, "GET", "/api/path", nil, cookies)

// 带配置选项的请求
data, err := rm.DoRequestWithCommonResponse[ResponseType](
    ctx, "POST", "/api/path", reqData, cookies,
    WithRequestTimeout(10*time.Second),
    WithRequestRetry(3, 2*time.Second),
    WithHeaders(map[string]string{"Custom": "value"}),
    WithProxy(true),
)
```

### 3. 配置管理模块 (`config.go`)

#### 功能特性
- **预定义配置**: 提供默认、台湾、高性能、代理等配置
- **工厂方法**: 便捷的客户端创建方法
- **配置应用**: 自动应用配置到客户端

#### 使用示例
```go
// 使用预定义配置
defaultClient := NewDefaultClient()
taiwanClient := NewTaiwanClient()
highPerfClient := NewHighPerformanceClient()
proxyClient := NewProxyClient()

// 自定义配置
config := DefaultConfig()
config.RetryTimes = 5
config.Timeout = 60 * time.Second
customClient := NewClientWithConfig(config)
```

### 4. 重构示例 (`client_v2.go`)

展示了如何使用新模块重构现有方法，包括：
- 登录方法重构
- 获取店铺列表重构
- 获取商品列表重构
- 更新商品信息重构

## 使用优势

### 1. 代码复用
- 消除了重复的错误处理逻辑
- 统一了请求处理流程
- 减少了样板代码

### 2. 错误处理
- 结构化的错误信息
- 自动的重试判断
- 更好的错误追踪

### 3. 可维护性
- 模块化设计
- 清晰的职责分离
- 易于扩展和修改

### 4. 性能优化
- 可配置的连接池
- 智能重试机制
- 代理轮换支持

## 迁移指南

### 旧代码
```go
resp, err := c.doRequest(HTTPMethodGet, APIPathGetMerchantShopList, nil, cookies)
if err != nil {
    return nil, fmt.Errorf("get merchant shop list failed: %w", err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    return nil, fmt.Errorf("read response body failed: %w", err)
}

var merchantShopListResp MerchantShopListResponse
err = json.Unmarshal(body, &merchantShopListResp)
if err != nil {
    return nil, fmt.Errorf("unmarshal merchant shop list response failed: %w", err)
}

if merchantShopListResp.Error != "" {
    return nil, fmt.Errorf("获取店铺列表失败:%s", merchantShopListResp.Error)
}
```

### 新代码
```go
rm := NewRequestManager(c)
data, err := rm.DoRequestWithCommonResponse[MerchantShopListData](
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
```

## 测试

运行测试确保优化正常工作：
```bash
go test ./client/shopee/
```

## 注意事项

1. **向后兼容**: 新模块与现有代码兼容，可以逐步迁移
2. **错误类型**: 所有新方法返回 `*ShopeeError` 类型
3. **配置管理**: 建议使用配置工厂方法创建客户端
4. **性能**: 新的请求管理器提供更好的性能和灵活性

## 后续优化建议

1. **缓存机制**: 添加响应缓存功能
2. **限流控制**: 实现客户端限流
3. **监控指标**: 添加请求指标和监控
4. **连接池**: 优化HTTP连接池配置
5. **异步支持**: 添加异步请求支持