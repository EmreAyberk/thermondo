package service

import (
	"context"
	"fmt"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/domain"
	"movie-rating-service/internal/infrastructure/db"
	"movie-rating-service/internal/infrastructure/repository"
)

type RatingService interface {
	Create(ctx context.Context, req request.CreateRating) (*response.CreateRating, error)
	GetRatingsByUserID(ctx context.Context, req request.GetUserRatings) (*response.GetUserRatings, error)
	Update(ctx context.Context, req request.UpdateRating) (*response.UpdateRating, error)
	Delete(ctx context.Context, req request.DeleteRating) error
}

type ratingService struct {
	ratingRepository repository.RatingRepository
	movieRepository  repository.MovieRepository
}

func NewRatingService(ratingRepository repository.RatingRepository, movieRepository repository.MovieRepository) RatingService {
	return &ratingService{ratingRepository: ratingRepository, movieRepository: movieRepository}
}

func (s *ratingService) Create(ctx context.Context, req request.CreateRating) (*response.CreateRating, error) {
	tx := db.BeginTransaction()

	rating, err := s.ratingRepository.Create(ctx, domain.Rating{
		UserID:  req.UserID,
		MovieID: req.MovieID,
		Score:   req.Score,
		Review:  req.Review,
	}, tx)
	if err != nil {
		err = tx.Rollback().Error
		if err != nil {
			return nil, fmt.Errorf("failed to rollback rate movie: %w", err)
		}
		return nil, fmt.Errorf("failed to rate movie: %w", err)
	}

	err = s.movieRepository.AddRating(ctx, req.MovieID, req.Score, tx)
	if err != nil {
		err = tx.Rollback().Error
		if err != nil {
			return nil, fmt.Errorf("failed to rollback add rating: %w", err)
		}
		return nil, fmt.Errorf("failed to add rating: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return rating.CreateRatingResponse(), nil
}

func (s *ratingService) GetRatingsByUserID(ctx context.Context, req request.GetUserRatings) (*response.GetUserRatings, error) {
	userRatings, err := s.ratingRepository.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's ratings: %w", err)
	}

	resp := &response.GetUserRatings{}
	resp.Ratings = make([]response.Ratings, len(userRatings))
	for i, userRating := range userRatings {
		resp.Ratings[i] = *userRating.GetRatingByUserIdResponse()
	}

	return resp, nil
}

func (s *ratingService) Update(ctx context.Context, req request.UpdateRating) (*response.UpdateRating, error) {
	rating, err := s.ratingRepository.GetByUserIDAndMovieID(ctx, req.UserID, req.MovieID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's rating on the selected movie: %w", err)
	}

	tx := db.BeginTransaction()

	err = s.ratingRepository.Update(ctx, domain.Rating{
		UserID:  req.UserID,
		MovieID: req.MovieID,
		Score:   req.Score,
		Review:  req.Review,
	}, tx)
	if err != nil {
		err = tx.Rollback().Error
		if err != nil {
			return nil, fmt.Errorf("ailed to rollback rate movie: %w", err)
		}
		return nil, fmt.Errorf("failed to rate movie: %w", err)
	}

	err = s.movieRepository.UpdateRating(ctx, req.MovieID, rating.Score, req.Score, tx)
	if err != nil {
		err = tx.Rollback().Error
		if err != nil {
			return nil, fmt.Errorf("failed to rollback update rating: %w", err)
		}
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return rating.UpdateMovieResponse(), nil
}

func (s *ratingService) Delete(ctx context.Context, req request.DeleteRating) error {
	rating, err := s.ratingRepository.GetByUserIDAndMovieID(ctx, req.UserID, req.MovieID)
	if err != nil {
		return fmt.Errorf("failed to get user's rating on the selected movie: %w", err)
	}

	tx := db.BeginTransaction()

	err = s.ratingRepository.Delete(ctx, domain.Rating{
		UserID:  req.UserID,
		MovieID: req.MovieID,
	}, tx)

	if err != nil {
		err = tx.Rollback().Error
		if err != nil {
			return fmt.Errorf("ailed to rollback rate movie: %w", err)
		}
		return fmt.Errorf("failed to rate movie: %w", err)
	}

	err = s.movieRepository.DeleteRating(ctx, req.MovieID, rating.Score, tx)
	if err != nil {
		err = tx.Rollback().Error
		if err != nil {
			return fmt.Errorf("failed to rollback update rating: %w", err)
		}
		return fmt.Errorf("failed to update rating: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
