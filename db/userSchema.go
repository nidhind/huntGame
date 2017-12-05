package db

type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Rank        int    `json:"rank"`
	AccessLevel string `json:"accessLevel"`
	AccessToken string `json:"accessToken"`
}
