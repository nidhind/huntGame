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

//For adding new puzzles
type PuzzleReq struct {
		Level        int    `json:"level,string"`
		Image        string `json:"image"`
		Clue         string `json:"clue"`
		SolutionHash string `json:"solutionHash"`
}

//For submitting answers
type AnswerReq struct {
	Answer string `json:"answer"`
}
