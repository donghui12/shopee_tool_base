package model

import (
	"time"

	"github.com/donghui12/shopee_tool_base/consts"
)

type ParentAccount struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Username    string    `gorm:"type:varchar(255)"`
	Password    string    `gorm:"type:varchar(255)"`
	Phone       string    `gorm:"type:varchar(255)"`
	MachineCode string    `gorm:"type:varchar(255)"`
	ActiveCode  string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
}

func (a *ParentAccount) TableName() string {
	return consts.ParentAccountTable
}

// Shopee 子账号表
type ShopeeAccount struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	ShopId          string    `json:"shop_id" gorm:"type:varchar(255)"`
	AccessToken     string    `json:"access_token" gorm:"type:varchar(255)"`
	RefreshToken    string    `json:"refresh_token" gorm:"type:varchar(255)"`
	ParentAccountId int64     `gorm:"not null;index"` // 外键，关联 ParentAccount.ID
	ExpiredAt       string    `json:"expired_at" gorm:"type:varchar(255)"`
	Status          string    `gorm:"type:varchar(50)"` // active / inactive 等
	CreatedAt       time.Time `gorm:"index"`
	UpdatedAt       time.Time
}

func (a *ShopeeAccount) TableName() string {
	return consts.ShopeeAccountTable
}
