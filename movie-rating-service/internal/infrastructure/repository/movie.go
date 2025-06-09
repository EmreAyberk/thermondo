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
	Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error)
	Get(ctx context.Context, id uint) (*domain.Movie, error)
	List(ctx context.Context) ([]domain.Movie, error)
	UpdateRating(ctx context.Context, movieID uint, score float64) error
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{DB: db}
}

func (r *movieRepository) Create(ctx context.Context, movie domain.Movie) (*domain.Movie, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := r.DB.WithContext(ctxWithTimeout).Create(&movie).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) Get(ctx context.Context, id uint) (*domain.Movie, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	movie := domain.Movie{}
	return &movie, r.DB.WithContext(ctxWithTimeout).Where("id=?", id).First(&movie).Error
}

func (r *movieRepository) List(ctx context.Context) ([]domain.Movie, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var movies []domain.Movie
	err := r.DB.WithContext(ctxWithTimeout).Find(&movies).Error
	return movies, err
}

func (r *movieRepository) UpdateRating(ctx context.Context, movieID uint, score float64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return r.DB.WithContext(ctxWithTimeout).Model(&domain.Movie{}).Where("id=?", movieID).Updates(map[string]interface{}{
		"rating":       gorm.Expr("((rating * rating_count) + ? ) / (rating_count +1)", score),
		"rating_count": gorm.Expr("rating_count + 1"),
	}).Error
}
