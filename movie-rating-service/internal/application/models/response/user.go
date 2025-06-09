package response

type CreateUser struct {
	ID uint `json:"id"`
}

type GetUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	IsAdmin  bool   `json:"is_admin"`
}

type Login struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
