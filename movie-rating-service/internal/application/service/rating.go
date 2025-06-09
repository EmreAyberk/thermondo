package service

import (
	"context"
	"fmt"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/domain"
	"movie-rating-service/internal/infrastructure/repository"
)

type RatingService interface {
	Create(ctx context.Context, req request.CreateRating) (*response.CreateRating, error)
	GetRatingsByUserID(ctx context.Context, req request.GetUserRatings) (*response.GetUserRatings, error)
}

type ratingService struct {
	ratingRepository repository.RatingRepository
	movieRepository  repository.MovieRepository
}

func NewRatingService(ratingRepository repository.RatingRepository, movieRepository repository.MovieRepository) RatingService {
	return &ratingService{ratingRepository: ratingRepository, movieRepository: movieRepository}
}

func (s *ratingService) Create(ctx context.Context, req request.CreateRating) (*response.CreateRating, error) {
	//TODO:add transaction

	rating, err := s.ratingRepository.Create(ctx, domain.Rating{
		UserID:  req.UserID,
		MovieID: req.MovieID,
		Score:   req.Score,
		Review:  req.Review,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to rate movie: %w", err)
	}

	err = s.movieRepository.UpdateRating(ctx, req.MovieID, req.Score)
	if err != nil {
		return nil, fmt.Errorf("failed to update rating: %w", err)
	}

	return rating.RateMovieResponse(), nil
}

func (s *ratingService) GetRatingsByUserID(ctx context.Context, req request.GetUserRatings) (*response.GetUserRatings, error) {
	userRatings, err := s.ratingRepository.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting user's ratings: %w", err)
	}

	resp := &response.GetUserRatings{}
	resp.Ratings = make([]response.Ratings, len(userRatings))
	for i, userRating := range userRatings {
		resp.Ratings[i] = *userRating.GetRatingByUserIdResponse()
	}

	return resp, nil
}
