package domain

import (
	"gorm.io/gorm"
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
