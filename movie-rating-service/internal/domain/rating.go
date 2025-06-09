package domain

import (
	"gorm.io/gorm"
	"movie-rating-service/internal/application/models/response"
)

type Rating struct {
	gorm.Model
	UserID  uint    `json:"user_id" gorm:"index:,unique,composite:uni_user_movie"`
	MovieID uint    `json:"movie_id" gorm:"index:,unique,composite:uni_user_movie"`
	Score   float64 `json:"score"`
	Review  string  `json:"review"`

	Movie Movie `json:"-" gorm:"foreignKey:MovieID"`
	User  User  `json:"-" gorm:"foreignKey:UserID"`
}

func (r *Rating) CreateRatingResponse() *response.CreateRating {
	return &response.CreateRating{
		ID: r.ID,
	}
}

func (r *Rating) UpdateMovieResponse() *response.UpdateRating {
	return &response.UpdateRating{
		ID: r.ID,
	}
}

func (r *Rating) GetRatingByUserIdResponse() *response.Ratings {
	return &response.Ratings{
		RatedMovie: response.RatedMovie{
			Title:       r.Movie.Title,
			Description: r.Movie.Description,
			Genre:       r.Movie.Genre,
			Director:    r.Movie.Director,
			Year:        r.Movie.Year,
			Rating:      r.Movie.Rating,
		},
		Rating: response.Rating{
			Score:  r.Score,
			Review: r.Review,
		},
	}
}
