package repository

import (
	"github.com/donghui12/shopee_tool_base/global"
	"github.com/donghui12/shopee_tool_base/model"
	"gorm.io/gorm"
)

type ShopRepository struct {
	db *gorm.DB
}

func NewShopRepository() *ShopRepository {
	return &ShopRepository{db: global.DB}
}

// GetShops 获取所有店铺
func (r *ShopRepository) GetShops() ([]model.Shop, error) {
	var shops []model.Shop
	err := r.db.Find(&shops).Error
	return shops, err
}

// GetShopByID 根据店铺ID获取店铺
func (r *ShopRepository) GetShopByID(shopID string) (*model.Shop, error) {
	var shop model.Shop
	err := r.db.Where("shop_id = ?", shopID).First(&shop).Error
	return &shop, err
}

// GetShopByID 根据店铺ID获取店铺
func (r *ShopRepository) GetShopByAccountID(accountID int64) ([]model.Shop, error) {
	var shops []model.Shop
	err := r.db.Where("account_id = ?", accountID).Find(&shops).Error
	return shops, err
}

// SaveOrUpdateShop 保存或更新店铺
func (r *ShopRepository) SaveOrUpdateShop(shop *model.Shop) error {
	return r.db.Where("shop_id = ?", shop.ShopID).
		Assign(shop).FirstOrCreate(shop).Error
}

// BatchSaveOrUpdateShops 批量保存或更新店铺
func (r *ShopRepository) BatchSaveOrUpdateShops(shops []model.Shop) error {
	for i := range shops {
		if err := r.SaveOrUpdateShop(&shops[i]); err != nil {
			return err
		}
	}
	return nil
}
