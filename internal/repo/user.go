package repo

import (
	"context"

	"github.com/thatmatin/subserv/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Exists(ctx context.Context, ID uint) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Exists(ctx context.Context, ID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", ID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
