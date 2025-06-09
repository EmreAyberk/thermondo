package domain

import (
	"gorm.io/gorm"
	"movie-rating-service/internal/application/models/response"
)

type Movie struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Director    string  `json:"director"`
	Year        int     `json:"year"`
	Rating      float64 `json:"rating"`
	RatingCount int64   `json:"rating_count"`
}

func (m *Movie) GetMovieResponse() *response.GetMovie {
	return &response.GetMovie{
		Title:       m.Title,
		Description: m.Description,
		Genre:       m.Genre,
		Director:    m.Director,
		Year:        m.Year,
		Rating:      m.Rating,
		RatingCount: m.RatingCount,
	}
}

func (m *Movie) CreateMovieResponse() *response.CreateMovie {
	return &response.CreateMovie{
		ID: m.ID,
	}
}
