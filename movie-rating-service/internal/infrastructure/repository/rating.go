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
	Create(ctx context.Context, rating domain.Rating, tx ...*gorm.DB) (*domain.Rating, error)
	GetByUserID(ctx context.Context, userID uint, tx ...*gorm.DB) ([]domain.Rating, error)
	GetByUserIDAndMovieID(ctx context.Context, userID, movieID uint, tx ...*gorm.DB) (*domain.Rating, error)
	Update(ctx context.Context, rating domain.Rating, tx ...*gorm.DB) error
	Delete(ctx context.Context, rating domain.Rating, tx ...*gorm.DB) error
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{DB: db}
}

func (r *ratingRepository) Create(ctx context.Context, rating domain.Rating, tx ...*gorm.DB) (*domain.Rating, error) {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := db.WithContext(ctxWithTimeout).Create(&rating).Error; err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepository) GetByUserID(ctx context.Context, userID uint, tx ...*gorm.DB) ([]domain.Rating, error) {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var ratings []domain.Rating
	if err := db.WithContext(ctxWithTimeout).Preload("Movie").Where("user_id = ?", userID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *ratingRepository) GetByUserIDAndMovieID(ctx context.Context, userID, movieID uint, tx ...*gorm.DB) (*domain.Rating, error) {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	ratings := domain.Rating{}
	return &ratings, db.WithContext(ctxWithTimeout).Preload("Movie").
		Where("user_id = ?", userID).
		Where("movie_id = ?", movieID).
		First(&ratings).Error
}

func (r *ratingRepository) Update(ctx context.Context, rating domain.Rating, tx ...*gorm.DB) error {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := db.WithContext(ctxWithTimeout).
		Model(&domain.Rating{}).
		Where("user_id = ?", rating.UserID).
		Where("movie_id = ?", rating.MovieID).
		Updates(rating).Error; err != nil {
		return err
	}

	return nil
}

func (r *ratingRepository) Delete(ctx context.Context, rating domain.Rating, tx ...*gorm.DB) error {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := db.WithContext(ctxWithTimeout).
		Where("user_id = ?", rating.UserID).
		Where("movie_id = ?", rating.MovieID).
		Delete(&domain.Rating{}).Error; err != nil {
		return err
	}

	return nil
}
