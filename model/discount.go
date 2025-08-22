package model

import (
	"time"

	"github.com/shopee_tool_base/consts"
	"github.com/shopee_tool_base/global"
)

type Discount struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey"`
	Name       string    `json:"name" gorm:"column:name;size:255"`
	ShopID     string    `json:"shop_id" gorm:"column:shop_id;size:64;not null;index"`
	DiscountID int64     `json:"discount_id" gorm:"column:discount_id;uniqueIndex;not null"`
	Status     int       `json:"status" gorm:"column:status;not null;default:1"`
	StartTime  string    `json:"start_time" gorm:"column:start_time"`
	EndTime    string    `json:"end_time" gorm:"column:end_time"`
	ItemCount  int       `json:"item_count" gorm:"column:item_count"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (d *Discount) TableName() string {
	return consts.DiscountTable
}

func (d *Discount) GetDiscounts(shopID string) ([]Discount, error) {
	var discounts []Discount
	err := global.DB.Where("shop_id = ?", shopID).
		Find(&discounts).Error
	return discounts, err
}

func (d *Discount) SaveOrUpdateDiscount(discount Discount) error {
	err := global.DB.Where("discount_id = ?", d.DiscountID).
		Assign(discount).FirstOrCreate(&discount).Error
	return err
}
func (d *Discount) BatchSaveOrUpdateDiscounts(discounts []Discount) error {
	for i := range discounts {
		if err := d.SaveOrUpdateDiscount(discounts[i]); err != nil {
			return err
		}
	}
	return nil
}
