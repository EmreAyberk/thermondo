package repository

import (
	"context"
	"gorm.io/gorm"
	"movie-rating-service/internal/domain"
	"sync"
	"time"
)

type cachedMovieRepository struct {
	movieRepository MovieRepository
	idCache         map[uint]cacheItem
	mu              sync.RWMutex
	ttl             time.Duration
}
type cacheItem struct {
	data      *domain.Movie
	expiresAt time.Time
}

func NewCachedMovieRepository(movieRepository MovieRepository, ttl time.Duration) MovieRepository {
	return &cachedMovieRepository{
		movieRepository: movieRepository,
		idCache:         make(map[uint]cacheItem),
		ttl:             ttl,
	}
}

func (c *cachedMovieRepository) Get(ctx context.Context, id uint) (*domain.Movie, error) {
	now := time.Now()

	c.mu.RLock()
	item, ok := c.idCache[id]
	if ok && item.expiresAt.After(now) {
		c.mu.RUnlock()
		return item.data, nil
	}
	c.mu.RUnlock()

	movie, err := c.movieRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.idCache[id] = cacheItem{data: movie, expiresAt: now.Add(c.ttl)}
	c.mu.Unlock()

	return movie, nil
}

func (c *cachedMovieRepository) Create(ctx context.Context, movie domain.Movie, tx ...*gorm.DB) (*domain.Movie, error) {
	return c.movieRepository.Create(ctx, movie)
}

func (c *cachedMovieRepository) List(ctx context.Context) ([]domain.Movie, error) {
	return c.movieRepository.List(ctx)
}

func (c *cachedMovieRepository) AddRating(ctx context.Context, movieID uint, score float64, tx ...*gorm.DB) error {
	err := c.movieRepository.AddRating(ctx, movieID, score, tx...)
	if err != nil {
		return err
	}

	movie, err := c.movieRepository.Get(ctx, movieID)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.idCache[movieID] = cacheItem{data: movie, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()

	return nil
}

func (c *cachedMovieRepository) UpdateRating(ctx context.Context, movieID uint, oldScore, newScore float64, tx ...*gorm.DB) error {
	err := c.movieRepository.UpdateRating(ctx, movieID, oldScore, newScore, tx...)
	if err != nil {
		return err
	}

	movie, err := c.movieRepository.Get(ctx, movieID)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.idCache[movieID] = cacheItem{data: movie, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()

	return nil
}

func (c *cachedMovieRepository) DeleteRating(ctx context.Context, movieID uint, score float64, tx ...*gorm.DB) error {
	err := c.movieRepository.DeleteRating(ctx, movieID, score, tx...)
	if err != nil {
		return err
	}

	movie, err := c.movieRepository.Get(ctx, movieID)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.idCache[movieID] = cacheItem{data: movie, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()

	return nil
}
