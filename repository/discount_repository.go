package repository

import (
	"github.com/donghui12/shopee_tool_base/global"
	"github.com/donghui12/shopee_tool_base/model"
	"gorm.io/gorm"
)

type DiscountRepository struct {
	db *gorm.DB
}

func NewDiscountRepository() *DiscountRepository {
	return &DiscountRepository{db: global.DB}
}

// GetDiscounts 根据店铺ID获取折扣列表
func (r *DiscountRepository) GetDiscounts(shopID string) ([]model.Discount, error) {
	var discounts []model.Discount
	err := r.db.Where("shop_id = ?", shopID).Find(&discounts).Error
	return discounts, err
}

// GetDiscountsByShopID 根据店铺ID获取有效的折扣列表
func (r *DiscountRepository) GetDiscountsByShopID(shopID string) ([]model.Discount, error) {
	var discounts []model.Discount
	err := r.db.Where("shop_id = ? ", shopID, 1).Find(&discounts).Error
	return discounts, err
}

// GetDiscountsByShopID 根据店铺ID获取有效的折扣列表
func (r *DiscountRepository) GetDiscountsByShopIDList(shopIDList []string) ([]model.Discount, error) {
	var discounts []model.Discount
	err := r.db.Where("shop_id in ?", shopIDList).Find(&discounts).Error
	return discounts, err
}

// GetAllActiveDiscounts 获取所有有效的折扣
func (r *DiscountRepository) GetAllActiveDiscounts() ([]model.Discount, error) {
	var discounts []model.Discount
	err := r.db.Where("status = ?", 1).Find(&discounts).Error
	return discounts, err
}

// SaveOrUpdateDiscount 保存或更新折扣
func (r *DiscountRepository) SaveOrUpdateDiscount(discount *model.Discount) error {
	return r.db.Where("discount_id = ?", discount.DiscountID).
		Assign(discount).FirstOrCreate(discount).Error
}

// BatchSaveOrUpdateDiscounts 批量保存或更新折扣
func (r *DiscountRepository) BatchSaveOrUpdateDiscounts(discounts []model.Discount) error {
	for i := range discounts {
		if err := r.SaveOrUpdateDiscount(&discounts[i]); err != nil {
			return err
		}
	}
	return nil
}
