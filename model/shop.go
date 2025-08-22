package model

import (
	"time"

	"github.com/shopee_tool_base/consts"
	"github.com/shopee_tool_base/global"
)

type Shop struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey"`
	ShopID    string    `json:"shop_id" gorm:"column:shop_id;size:64;uniqueIndex;not null"`
	Region    string    `json:"region" gorm:"column:region;size:16;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (s *Shop) TableName() string {
	return consts.ShopTable
}

func (s *Shop) GetShops() ([]Shop, error) {
	var shops []Shop
	err := global.DB.Find(&shops).Error
	return shops, err
}

func (s *Shop) GetShopByID(shopID string) error {
	return global.DB.Where("shop_id = ?", shopID).First(&s).Error
}

func (s *Shop) SaveOrUpdateShop() error {
	return global.DB.Where("shop_id = ?", s.ShopID).
		Assign(s).FirstOrCreate(s).Error
}

func (s *Shop) BatchSaveOrUpdateShops(shops []Shop) error {
	for _, shop := range shops {
		if err := shop.SaveOrUpdateShop(); err != nil {
			return err
		}
	}
	return nil
}
