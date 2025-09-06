package model

import (
	"time"

	"github.com/donghui12/shopee_tool_base/consts"
)

// Account 存储虾皮账号信息
type Account struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	AccountId    int64     `json:"account_id" gorm:"primarykey"`
	MerchantName string    `json:"merchant_name" gorm:"type:varchar(255)"`  // shopee 账户
	Username     string    `json:"username" gorm:"type:varchar(255)"`       // shopee 账户
	Password     string    `json:"password" gorm:"type:varchar(255)"`       // shopee 密码
	Phone        string    `json:"phone" gorm:"type:varchar(255)"`          // 手机号
	Email        string    `json:"email" gorm:"type:varchar(255)"`          // 邮箱
	MachineCode  string    `json:"machine_code" gorm:"type:varchar(255)"`   // 机器码
	ActiveCode   string    `json:"active_code" gorm:"type:varchar(255)"`    // 激活码
	ExpiredAt    string    `json:"expired_at" gorm:"type:varchar(255)"`     // 过期时间
	Cookies      string    `json:"cookies" gorm:"type:text"`                // cookies
	Session      string    `json:"session" gorm:"type:text"`                // session
	Status       int       `json:"status" gorm:"type:int;default:1"`        // 状态：1=有效 0=无效
	ShopId       string    `json:"shop_id" gorm:"type:string"`              // 授权店铺 id
	AccessToken  string    `json:"access_token" gorm:"type:access_token"`   // token
	RefreshToken string    `json:"refresh_token" gorm:"type:refresh_token"` // 刷新 token
	CreatedAt    time.Time `json:"created_at" gorm:"index"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a *Account) TableName() string {
	return consts.AccountTable
}
