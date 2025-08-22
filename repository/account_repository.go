package repository

import (
	"github.com/donghui12/shopee_tool_base/global"
	"github.com/donghui12/shopee_tool_base/model"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{db: global.DB}
}

// GetAccountByUsername 根据用户名获取账号
func (r *AccountRepository) GetAccountByUsername(username string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("username = ?", username).First(&account).Error
	return &account, err
}

// GetAccountByID 根据ID获取账号
func (r *AccountRepository) GetAccountByID(id uint) (*model.Account, error) {
	var account model.Account
	err := r.db.First(&account, id).Error
	return &account, err
}

// GetAccounts 获取所有账号
func (r *AccountRepository) GetAccounts() ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Find(&accounts).Error
	return accounts, err
}

// GetAllActiveAccounts 获取所有有效的账号
func (r *AccountRepository) GetAllActiveAccounts() ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Where("status = ?", 1).Find(&accounts).Error
	return accounts, err
}

// UpdateCookies 更新账号的cookies
func (r *AccountRepository) UpdateCookies(id uint, cookies string) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Update("cookies", cookies).Error
}

// UpdateSession 更新账号的session
func (r *AccountRepository) UpdateSession(id uint, session string) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Update("session", session).Error
}

// UpdateStatus 更新账号状态
func (r *AccountRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Update("status", status).Error
}

// SaveOrUpdateAccount 保存或更新账号
func (r *AccountRepository) SaveOrUpdateAccount(account *model.Account) error {
	return r.db.Where("username = ?", account.Username).
		Assign(account).FirstOrCreate(account).Error
}

// BatchSaveOrUpdateAccounts 批量保存或更新账号
func (r *AccountRepository) BatchSaveOrUpdateAccounts(accounts []model.Account) error {
	for i := range accounts {
		if err := r.SaveOrUpdateAccount(&accounts[i]); err != nil {
			return err
		}
	}
	return nil
}