package repository

import (
	"openvpn-admin-go/model"

	"gorm.io/gorm"
)

// UserFilter 用户列表过滤条件
type UserFilter struct {
	DepartmentID string
	UserID       string // 仅查询单个用户时使用
}

// UserRepository 用户数据访问接口
type UserRepository interface {
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByName(name string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	UpdateFields(id string, fields map[string]interface{}) error
	Delete(id string) error
	List(filter UserFilter) ([]model.User, error)
}

// GormUserRepository 基于 GORM 的用户仓库实现
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) FindByName(name string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) UpdateFields(id string, fields map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *GormUserRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *GormUserRepository) List(filter UserFilter) ([]model.User, error) {
	var users []model.User
	db := r.db.Order("created_at desc")
	if filter.DepartmentID != "" {
		db = db.Where("department_id = ?", filter.DepartmentID)
	}
	if filter.UserID != "" {
		db = db.Where("id = ?", filter.UserID)
	}
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
