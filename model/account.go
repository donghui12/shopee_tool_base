package model

import (
	"time"

	"github.com/shopee_tool_base/consts"
	"github.com/shopee_tool_base/global"
)

// Account 存储虾皮账号信息
type Account struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	AccountId   int64     `json:"account_id" gorm:"primarykey"`
	Username    string    `json:"username" gorm:"type:varchar(255)"`     // shopee 账户
	Password    string    `json:"password" gorm:"type:varchar(255)"`     // shopee 密码
	Phone       string    `json:"phone" gorm:"type:varchar(255)"`        // 手机号
	MachineCode string    `json:"machine_code" gorm:"type:varchar(255)"` // 机器码
	ActiveCode  string    `json:"active_code" gorm:"type:varchar(255)"`  // 激活码
	ExpiredAt   string    `json:"expired_at" gorm:"type:varchar(255)"`   // 过期时间
	Cookies     string    `json:"cookies" gorm:"type:text"`              // cookies
	Session     string    `json:"session" gorm:"type:text"`              // session
	Status      int       `json:"status" gorm:"type:int;default:1"`      // 状态：1=有效 0=无效
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Account) TableName() string {
	return consts.AccountTable
}

func (a *Account) GetAccountByUsername(username string) error {
	return global.DB.Where("username = ?", username).First(&a).Error
}

func (a *Account) GetAccountByID(id uint) error {
	return global.DB.First(&a, id).Error
}

func (a *Account) GetAccounts() ([]Account, error) {
	var accounts []Account
	err := global.DB.Find(&accounts).Error
	return accounts, err
}

func (a *Account) UpdateCookies(newCookies string) error {
	a.Cookies = newCookies
	a.UpdatedAt = time.Now()
	return global.DB.Save(a).Error
}

func (a *Account) UpdateSession(newSession string) error {
	a.Session = newSession
	a.UpdatedAt = time.Now()
	return global.DB.Save(a).Error
}

func (a *Account) UpdateStatus(status int) error {
	a.Status = status
	a.UpdatedAt = time.Now()
	return global.DB.Save(a).Error
}

func (a *Account) SaveOrUpdateAccount() error {
	return global.DB.Where("username = ?", a.Username).
		Assign(a).FirstOrCreate(a).Error
}

func (a *Account) BatchSaveOrUpdateAccounts(accounts []Account) error {
	for _, account := range accounts {
		if err := account.SaveOrUpdateAccount(); err != nil {
			return err
		}
	}
	return nil
}

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
