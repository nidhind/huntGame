package models

type LoginRes struct{
	Code string `json:"code"`
	Message string `json:"message"`
  Status string `json:"status"`
  Payload LoginResExtra `json:"payload"`
}

type LoginResExtra struct{
  AccessToken string `json:"access_token"`
}
