package repository

import (
	"time"

	"gorm.io/gorm"

	"showcase-backend/model"
)

// UserRepo 用户数据访问。
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateLastLogin(id uint) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login_at", time.Now()).Error
}
