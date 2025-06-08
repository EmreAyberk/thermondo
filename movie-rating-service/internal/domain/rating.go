package domain

import (
	"gorm.io/gorm"
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
