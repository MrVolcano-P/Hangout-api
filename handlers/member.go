package handlers

type Member struct {
	ID       uint `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
