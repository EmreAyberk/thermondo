package request

type CreateRating struct {
	MovieID uint    `json:"-" validate:"required"`
	UserID  uint    `json:"-" validate:"required"`
	Score   float64 `json:"score" validate:"required,gte=0,lte=5"`
	Review  string  `json:"review"`
}

type UpdateRating struct {
	MovieID uint    `json:"-" validate:"required"`
	UserID  uint    `json:"-" validate:"required"`
	Score   float64 `json:"score" validate:"required,gte=0,lte=5"`
	Review  string  `json:"review"`
}

type DeleteRating struct {
	MovieID uint `json:"-" validate:"required"`
	UserID  uint `json:"-" validate:"required"`
}

type GetUserRatings struct {
	UserID uint `param:"id" validate:"required"`
}
