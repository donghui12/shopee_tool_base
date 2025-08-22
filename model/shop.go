package model

import (
	"time"

	"github.com/donghui12/shopee_tool_base/consts"
)

type Shop struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey"`
	AccountId int64     `json:"account_id" gorm:"column:account_id;size:64"`
	ShopID    string    `json:"shop_id" gorm:"column:shop_id;size:64;uniqueIndex;not null"`
	Region    string    `json:"region" gorm:"column:region;size:16;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (s *Shop) TableName() string {
	return consts.ShopTable
}
