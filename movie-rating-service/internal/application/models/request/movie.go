package request

type CreateMovie struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Director    string `json:"director" validate:"required"`
	Year        int    `json:"year"`
}

type UpdateMovie struct {
	ID          uint   `param:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Director    string `json:"director" validate:"required"`
	Year        int    `json:"year"`
}
type DeleteMovie struct {
	ID uint `param:"id"`
}

type GetMovie struct {
	ID uint `param:"id" validate:"required"`
}
