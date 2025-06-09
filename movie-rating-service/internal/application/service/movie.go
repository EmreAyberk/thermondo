package service

import (
	"context"
	"fmt"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/domain"
	"movie-rating-service/internal/infrastructure/repository"
)

type MovieService interface {
	Create(ctx context.Context, req request.CreateMovie) (*response.CreateMovie, error)
	Get(ctx context.Context, req request.GetMovie) (*response.GetMovie, error)
}

type movieService struct {
	movieRepository repository.MovieRepository
}

func NewMovieService(movieRepository repository.MovieRepository) MovieService {
	return &movieService{
		movieRepository: movieRepository,
	}
}

func (s *movieService) Get(ctx context.Context, req request.GetMovie) (*response.GetMovie, error) {
	movie, err := s.movieRepository.Get(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("error occurred while finding movie: %w", err)
	}
	return movie.GetMovieResponse(), nil
}

func (s *movieService) Create(ctx context.Context, req request.CreateMovie) (*response.CreateMovie, error) {
	movie, err := s.movieRepository.Create(ctx, domain.Movie{
		Title:       req.Title,
		Description: req.Description,
		Genre:       req.Genre,
		Director:    req.Director,
		Year:        req.Year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}
	return movie.CreateMovieResponse(), nil
}
