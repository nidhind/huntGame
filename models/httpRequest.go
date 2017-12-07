// Typesmodels in HTTP Requests

package models

// For user self registeration
type SignUpReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	EmailId   string `json:"emailId"`
	Password  string `json:"password"`
}

// For accesstoken request
type GenAccessToken struct {
	Email    string `json:"emailId"`
	Password string `json:"password"`
}
