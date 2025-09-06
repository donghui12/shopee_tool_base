package repository

import (
	"errors"
	"fmt"

	"github.com/donghui12/shopee_tool_base/global"
	"github.com/donghui12/shopee_tool_base/model"
	"github.com/donghui12/shopee_tool_base/utils"
	"gorm.io/gorm"
)

type ParentAccountRepository struct {
	db *gorm.DB
}

func NewParentAccountRepository() *ParentAccountRepository {
	return &ParentAccountRepository{db: global.DB}
}

func (s *ParentAccountRepository) GetByID(parentAccountId int64) (model.ParentAccount, error) {
	var account model.ParentAccount
	result := s.db.Model(&model.ParentAccount{}).
		Where("id = ?", parentAccountId).
		First(&account)
	return account, result.Error
}

func (s *ParentAccountRepository) UpdateMachineCode(username, MachineCode string) error {
	// 更新数据库中的机器码
	result := s.db.Model(&model.ParentAccount{}).
		Where("username = ?", username).
		Update("machine_code", MachineCode)
	return result.Error
}

func (s *ParentAccountRepository) Register(username, password, phone string) error {
	// 查重
	var count int64
	s.db.Model(&model.ParentAccount{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}

	// md5 加密
	hashedPwd := utils.Encrypt(password)

	// 插库
	account := model.ParentAccount{
		Username: username,
		Password: hashedPwd,
		Phone:    phone,
	}
	if err := s.db.Create(&account).Error; err != nil {
		return fmt.Errorf("注册失败: %w", err)
	}

	return nil
}

// 1. 登录校验（用户名 + md5密码）
func (s *ParentAccountRepository) Login(username, password string) (model.ParentAccount, error) {
	encryptedPassword := utils.Encrypt(password)
	var account model.ParentAccount

	err := s.db.Where("username = ? AND password = ?",
		username, encryptedPassword).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return account, fmt.Errorf("用户名或密码错误")
	}
	return account, err
}

// 2. 绑定激活码
func (s *ParentAccountRepository) BindActivationCode(username, activationCode string) error {
	result := s.db.Model(&model.ParentAccount{}).
		Where("username = ?", username).
		Update("active_code", activationCode)
	if result.Error != nil {
		return fmt.Errorf("绑定激活码失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到对应的用户")
	}
	return nil
}

// 3. 解绑激活码
func (s *ParentAccountRepository) UnbindActivationCode(username string) error {
	result := s.db.Model(&model.ParentAccount{}).
		Where("username = ?", username).
		Update("active_code", "")
	if result.Error != nil {
		return fmt.Errorf("解绑激活码失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到对应的用户")
	}
	return nil
}

func (s *ParentAccountRepository) VerifyActiveCode(activeCode string) (string, error) {
	var selectActiveCode string
	result := s.db.Model(&model.ParentAccount{}).
		Where("active_code = ?", activeCode).
		Select("active_code").Scan(&selectActiveCode)
	return selectActiveCode, result.Error
}

func (s *ParentAccountRepository) GetActiveCodeByUsername(username string) (string, error) {
	var selectActiveCode string
	result := s.db.Model(&model.ParentAccount{}).
		Where("username = ?", username).
		Select("active_code").Scan(&selectActiveCode)
	return selectActiveCode, result.Error
}

func (s *ParentAccountRepository) GetActiveCodeByUsernames(usernames []string) ([]string, error) {
	var activeCodes []string
	result := s.db.Model(&model.ParentAccount{}).
		Where("username in ?", usernames).
		Select("active_code").Scan(&activeCodes)
	return activeCodes, result.Error
}
