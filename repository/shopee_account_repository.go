package repository

import (
	"fmt"

	"github.com/donghui12/shopee_tool_base/global"
	"github.com/donghui12/shopee_tool_base/model"
	"gorm.io/gorm"
)

type ShopeeAccountRepository struct {
	db *gorm.DB
}

func NewShopeeAccountRepository() *ShopeeAccountRepository {
	return &ShopeeAccountRepository{db: global.DB}
}

// GetAccounts 获取所有账号
func (s *ShopeeAccountRepository) GetAuthShopeeAccounts() ([]model.ShopeeAccount, error) {
	var accounts []model.ShopeeAccount
	err := s.db.Where("shop_id != ''").Find(&accounts).Error
	return accounts, err
}

// UpdateToken 更新账号的 token, access_token, refresh_token
func (s *ShopeeAccountRepository) RefreshToken(id int64, accessToken, refreshToken, expireTime string) error {
	return s.db.Model(&model.ShopeeAccount{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expired_at":    expireTime,
		}).Error
}

// 创建 shopee tw 账户
func (s *ShopeeAccountRepository) CreateShopeeAccount(shopId, accessToken, refreshToken, expireAt string) error {
	// 创建账户
	account := model.ShopeeAccount{
		ShopId:       shopId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    expireAt,
	}
	// 保存到数据库, 如果账户已存在则更新
	var existingShopeeAccount model.ShopeeAccount
	result := s.db.Where("shop_id = ?", shopId).First(&existingShopeeAccount)
	if result.Error == nil {
		existingShopeeAccount.AccessToken = accessToken
		existingShopeeAccount.RefreshToken = refreshToken
		existingShopeeAccount.ExpiredAt = expireAt
		return s.db.Save(&existingShopeeAccount).Error
	}
	// 创建账户
	result = s.db.Create(&account)
	if result.Error != nil {
		return fmt.Errorf("创建账户失败: %w", result.Error)
	}
	return nil
}

// 判断是否授权
func (s *ShopeeAccountRepository) IsAuth(shopId string) (string, error) {
	var accessToken string
	var account model.ShopeeAccount
	result := s.db.Where("shop_id = ?", shopId).First(&account)
	if result.Error != nil {
		return accessToken, result.Error
	}
	if account.AccessToken == "" {
		return accessToken, fmt.Errorf("未授权")
	}
	return account.AccessToken, nil
}

// BindParentAccount 绑定父账户
func (s *ShopeeAccountRepository) BindParentAccount(parentAccountId int64, shopId string) error {
	// 保存到数据库, 如果账户不存在则退出报错
	var existingShopeeAccount model.ShopeeAccount
	result := s.db.Where("shop_id = ?", shopId).First(&existingShopeeAccount)
	if result.Error != nil {
		return result.Error
	}
	if existingShopeeAccount.ShopId != shopId {
		return fmt.Errorf("shop_id 不存在: %s,请授权", shopId)
	}
	return s.db.Model(&model.ShopeeAccount{}).
		Where("shop_id = ?", shopId).
		Updates(map[string]interface{}{
			"parent_account_id": parentAccountId,
		}).Error
}
