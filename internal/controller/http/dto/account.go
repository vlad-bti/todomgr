package dto

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAccountRequest struct {
	Name        string `json:"name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	AccountType uint   `json:"type" binding:"required,oneof=1 2"`
}

type GetAccountRequest struct {
	Id uint `json:"id" binding:"required"`
}

type DeleteAccountRequest struct {
	Id uint `json:"id" binding:"required"`
}
