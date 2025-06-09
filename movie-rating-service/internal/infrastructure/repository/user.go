package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"movie-rating-service/internal/domain"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	GetByID(ctx context.Context, userID uint) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := r.DB.WithContext(ctxWithTimeout).Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, userID uint) (*domain.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	user := domain.User{}
	return &user, r.DB.WithContext(ctxWithTimeout).Where("id=?", userID).First(&user).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	user := domain.User{}
	return &user, r.DB.WithContext(ctxWithTimeout).Where("username = ?", username).First(&user).Error
}
