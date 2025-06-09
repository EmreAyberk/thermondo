package repository

import (
	"context"
	"gorm.io/gorm"
	"movie-rating-service/internal/domain"
	"time"
)

type ratingRepository struct {
	DB *gorm.DB
}

type RatingRepository interface {
	Create(ctx context.Context, rating domain.Rating) (*domain.Rating, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Rating, error)
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{DB: db}
}

func (r *ratingRepository) Create(ctx context.Context, rating domain.Rating) (*domain.Rating, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := r.DB.WithContext(ctxWithTimeout).Create(&rating).Error; err != nil {
		return nil, err
	}
	return &rating, nil
}
func (r *ratingRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Rating, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var ratings []domain.Rating
	if err := r.DB.WithContext(ctxWithTimeout).Preload("Movie").Where("user_id = ?", userID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}
