// Types used in HTTP Requests

package models

type SignUpReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginReq struct{
	Email string `json:"emailId"`
	Password string `json:"password"`
}
