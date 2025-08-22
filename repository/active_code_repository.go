package repository

import (
	"github.com/donghui12/shopee_tool_base/global"
	"github.com/donghui12/shopee_tool_base/model"
	"gorm.io/gorm"
)

type ActiveCodeRepository struct {
	db *gorm.DB
}

func NewActiveCodeRepository() *ActiveCodeRepository {
	return &ActiveCodeRepository{db: global.DB}
}

// CreateActiveCode 创建激活码
func (r *ActiveCodeRepository) CreateActiveCode(activeCode *model.ActiveCode) error {
	return r.db.Create(activeCode).Error
}

// GetActiveCodeByCode 根据激活码获取
func (r *ActiveCodeRepository) GetActiveCodeByCode(code string) (*model.ActiveCode, error) {
	var activeCode model.ActiveCode
	err := r.db.Where("code = ?", code).First(&activeCode).Error
	return &activeCode, err
}

// GetAllActiveCodes 获取所有激活码
func (r *ActiveCodeRepository) GetAllActiveCodes() ([]model.ActiveCode, error) {
	var activeCodes []model.ActiveCode
	err := r.db.Find(&activeCodes).Error
	return activeCodes, err
}

// DeleteActiveCode 删除激活码
func (r *ActiveCodeRepository) DeleteActiveCode(id uint) error {
	return r.db.Delete(&model.ActiveCode{}, id).Error
}