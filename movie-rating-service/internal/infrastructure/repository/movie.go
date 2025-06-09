package repository

import (
	"context"
	"gorm.io/gorm"
	"movie-rating-service/internal/domain"
	"time"
)

type movieRepository struct {
	DB *gorm.DB
}

type MovieRepository interface {
	Create(ctx context.Context, movie domain.Movie, tx ...*gorm.DB) (*domain.Movie, error)
	Get(ctx context.Context, id uint) (*domain.Movie, error)
	List(ctx context.Context) ([]domain.Movie, error)
	AddRating(ctx context.Context, movieID uint, score float64, tx ...*gorm.DB) error
	UpdateRating(ctx context.Context, movieID uint, oldScore, newScore float64, tx ...*gorm.DB) error
	DeleteRating(ctx context.Context, movieID uint, score float64, tx ...*gorm.DB) error
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{DB: db}
}

func (r *movieRepository) Create(ctx context.Context, movie domain.Movie, tx ...*gorm.DB) (*domain.Movie, error) {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := db.WithContext(ctxWithTimeout).Create(&movie).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) Get(ctx context.Context, id uint) (*domain.Movie, error) {
	db := r.DB

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	movie := domain.Movie{}
	return &movie, db.WithContext(ctxWithTimeout).Where("id=?", id).First(&movie).Error
}

func (r *movieRepository) List(ctx context.Context) ([]domain.Movie, error) {
	db := r.DB

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var movies []domain.Movie
	err := db.WithContext(ctxWithTimeout).Find(&movies).Error
	return movies, err
}

func (r *movieRepository) AddRating(ctx context.Context, movieID uint, score float64, tx ...*gorm.DB) error {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return db.WithContext(ctxWithTimeout).Model(&domain.Movie{}).Where("id=?", movieID).Updates(map[string]interface{}{
		"rating":       gorm.Expr("((rating * rating_count) + ? ) / GREATEST(rating_count + 1, 1)", score),
		"rating_count": gorm.Expr("rating_count + 1"),
	}).Error
}

func (r *movieRepository) UpdateRating(ctx context.Context, movieID uint, oldScore, newScore float64, tx ...*gorm.DB) error {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return db.WithContext(ctxWithTimeout).Model(&domain.Movie{}).Where("id=?", movieID).Updates(map[string]interface{}{
		"rating": gorm.Expr("((rating * rating_count) - ? + ? ) / GREATEST(rating_count, 1)", oldScore, newScore),
	}).Error
}

func (r *movieRepository) DeleteRating(ctx context.Context, movieID uint, score float64, tx ...*gorm.DB) error {
	db := r.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return db.WithContext(ctxWithTimeout).Model(&domain.Movie{}).Where("id=?", movieID).Updates(map[string]interface{}{
		"rating":       gorm.Expr("((rating * rating_count) - ? ) / GREATEST(rating_count - 1, 1)", score),
		"rating_count": gorm.Expr("rating_count - 1"),
	}).Error
}
