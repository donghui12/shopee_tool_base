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

// UpdateStatus 更新账号状态
func (r *AccountRepository) UpdateAccountId(id uint, accountId int64) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Update("account_id", accountId).Error
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

func (s *AccountRepository) UpdateMachineCode(username, machineCode string) error {
	// 更新数据库中的机器码
	result := s.db.Model(&model.Account{}).
		Where("username = ?", username).
		Update("machine_code", machineCode)
	return result.Error
}

func (s *AccountRepository) GetMachineCode(username, machineCode string) error {
	var selectMachineCode string
	result := s.db.Model(&model.Account{}).
		Where("username = ? AND machine_code = ?", username, machineCode).
		Select("machine_code").Scan(&selectMachineCode)
	return result.Error
}

func (s *AccountRepository) GetActiveCodeByActiveCode(activeCode string) error {
	var selectActiveCode string
	result := s.db.Model(&model.Account{}).
		Where("active_code = ?", activeCode).
		Select("active_code").Scan(&selectActiveCode)
	return result.Error
}

func (s *AccountRepository) UpdateActiveCode(username, activeCode string) error {
	result := s.db.Model(&model.Account{}).
		Where("username = ?", username).
		Update("active_code", activeCode)
	return result.Error
}

func (s *AccountRepository) VerifyActiveCode(activeCode string) (string, error) {
	var selectActiveCode string
	result := s.db.Model(&model.Account{}).
		Where("active_code = ?", activeCode).
		Select("active_code").Scan(&selectActiveCode)
	return selectActiveCode, result.Error
}

func (s *AccountRepository) GetActiveCode(username string) (string, error) {
	var activeCode string
	result := s.db.Model(&model.Account{}).
		Where("username = ?", username).
		Select("active_code").Scan(&activeCode)
	return activeCode, result.Error
}

// GetActiveCodeByUsernames 根据用户名获取激活码
func (s *AccountRepository) GetActiveCodeByUsernames(usernames []string) ([]string, error) {
	var activeCodes []string
	result := s.db.Model(&model.Account{}).
		Where("username in ?", usernames).
		Select("active_code").Scan(&activeCodes)
	return activeCodes, result.Error
}

// QueryAccount 根据用户名或激活码获取用户列表，并支持分页
func (s *AccountRepository) QueryAccount(username, code []string, page, pageSize int) ([]model.Account, error) {
	var accounts []model.Account
	query := s.db.Model(&model.Account{})

	// 条件构造
	if len(username) > 0 {
		query = query.Where("username in (?)", username)
	} else if len(code) > 0 {
		query = query.Where("active_code in (?)", code)
	}
	// 分页处理
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// 执行查询
	err := query.Find(&accounts).Error
	return accounts, err
}

func (s *AccountRepository) GetCookies(username string) (string, error) {
	account, err := s.GetAccountByUsername(username)
	if err != nil {
		return "", err
	}
	return account.Cookies, nil
}

func (s *AccountRepository) GetSession(username string) (string, error) {
	account, err := s.GetAccountByUsername(username)
	if err != nil {
		return "", err
	}
	return account.Session, nil
}

func (s *AccountRepository) GetSessionByUsernameAndPassword(username, password string) (string, error) {
	var session string
	result := s.db.Model(&model.Account{}).
		Where("username = ? AND password = ? AND updated_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)", username, password).
		Select("session").Scan(&session)
	return session, result.Error
}
