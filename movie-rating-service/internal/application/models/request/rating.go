package request

type CreateRating struct {
	MovieID uint    `param:"movie_id" validate:"required"`
	UserID  uint    `validate:"required"`
	Score   float64 `json:"score" validate:"required"`
	Review  string  `json:"review"`
}

type GetUserRatings struct {
	UserID uint `param:"id" validate:"required"`
}
