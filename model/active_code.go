package model

import (
	"time"

	"github.com/donghui12/shopee_tool_base/consts"
)

// ActiveCode 激活码模型
type ActiveCode struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code" gorm:"not null"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ac *ActiveCode) TableName() string {
	return consts.ActiveCodeTable
}
