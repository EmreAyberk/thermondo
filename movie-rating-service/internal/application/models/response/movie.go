package response

type CreateMovie struct {
	ID uint `json:"id"`
}

type GetMovie struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Director    string  `json:"director"`
	Year        int     `json:"year"`
	Rating      float64 `json:"rating"`
	RatingCount int64   `json:"rating_count"`
}
