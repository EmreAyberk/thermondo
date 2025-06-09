package response

type CreateRating struct {
	ID uint `json:"id"`
}

type GetUserRatings struct {
	Ratings []Ratings `json:"ratings"`
}

type Ratings struct {
	RatedMovie RatedMovie `json:"rated_movie"`
	Rating     Rating     `json:"rating"`
}
type RatedMovie struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Director    string  `json:"director"`
	Year        int     `json:"year"`
	Rating      float64 `json:"rating"`
}
type Rating struct {
	Score  float64 `json:"score"`
	Review string  `json:"review"`
}
