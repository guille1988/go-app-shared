package dtos

type UserLoggedIn struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
