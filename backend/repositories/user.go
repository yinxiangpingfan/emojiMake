package repositories

import (
	"emoji-maker-backend/models"

	"xorm.io/xorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *models.User) error
	FindByPhone(phone string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
	Update(user *models.User) error
}

// xormUserRepository 用户仓库实现
type xormUserRepository struct {
	engine *xorm.Engine
}

// NewXormUserRepository 创建用户仓库实例
func NewXormUserRepository(engine *xorm.Engine) UserRepository {
	return &xormUserRepository{engine: engine}
}

// Create 创建用户
func (r *xormUserRepository) Create(user *models.User) error {
	_, err := r.engine.Insert(user)
	return err
}

// FindByPhone 根据手机号查找用户
func (r *xormUserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	has, err := r.engine.Where("phone = ?", phone).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil // 用户不存在
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *xormUserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	has, err := r.engine.ID(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil // 用户不存在
	}
	return &user, nil
}

// Update 更新用户信息
func (r *xormUserRepository) Update(user *models.User) error {
	_, err := r.engine.ID(user.ID).Update(user)
	return err
}
