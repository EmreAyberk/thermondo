package seeder

import (
	"gorm.io/gorm"
	"log/slog"
	"movie-rating-service/internal/domain"
)

type seeder struct {
	db *gorm.DB
}

type Seeder interface {
	Seed() error
}

func NewSeeder(db *gorm.DB) Seeder {
	return &seeder{db: db}
}

func (s *seeder) Seed() error {
	users := []domain.User{
		{
			Username: "alice",
			Password: "$2a$10$HuR2tFF3FV2Suyv9/rMX/Om6ylABQFADP15d2WuHsYGyFMqwFh38W",
			Name:     "Alice",
			Surname:  "Wonder",
			Email:    "alice@mail.com",
			Phone:    "1234567890",
			Address:  "123 Elm St",
			IsAdmin:  false,
		},
		{
			Username: "bob",
			Password: "$2a$10$HuR2tFF3FV2Suyv9/rMX/Om6ylABQFADP15d2WuHsYGyFMqwFh38W",
			Name:     "Bob",
			Surname:  "Builder",
			Email:    "bob@mail.com",
			Phone:    "2345678901",
			Address:  "456 Oak St",
			IsAdmin:  false,
		},
		{
			Username: "carol",
			Password: "$2a$10$HuR2tFF3FV2Suyv9/rMX/Om6ylABQFADP15d2WuHsYGyFMqwFh38W",
			Name:     "Carol",
			Surname:  "Smith",
			Email:    "carol@mail.com",
			Phone:    "3456789012",
			Address:  "789 Pine St",
			IsAdmin:  true,
		},
	}

	movies := []domain.Movie{
		{
			Title:       "Inception",
			Description: "A mind-bending thriller",
			Genre:       "Sci-Fi",
			Director:    "Christopher Nolan",
			Year:        2010,
		},
		{
			Title:       "Interstellar",
			Description: "Space exploration epic",
			Genre:       "Sci-Fi",
			Director:    "Christopher Nolan",
			Year:        2014,
		},
		{
			Title:       "The Godfather",
			Description: "Mafia family story",
			Genre:       "Crime",
			Director:    "Francis Ford Coppola",
			Year:        1972,
		},
	}

	ratings := []domain.Rating{
		{
			UserID:  1,
			MovieID: 1,
			Score:   4.5,
			Review:  "Amazing movie!",
		},
		{
			UserID:  1,
			MovieID: 2,
			Score:   5.0,
			Review:  "Loved the visuals",
		},
		{
			UserID:  2,
			MovieID: 1,
			Score:   4.0,
			Review:  "Great concept",
		},
		{
			UserID:  2,
			MovieID: 3,
			Score:   5.0,
			Review:  "Classic masterpiece",
		},
		{
			UserID:  3,
			MovieID: 2,
			Score:   4.8,
			Review:  "Emotional and thought-provoking",
		},
	}

	slog.Info("Seeder start...")
	s.db.Save(&users)
	s.db.Save(&movies)
	s.db.Save(&ratings)
	slog.Info("Seeder end...")

	return nil
}
