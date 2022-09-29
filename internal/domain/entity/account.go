package entity

type AccountType uint

const (
	_ AccountType = iota
	AccountTypeAdmin
	AccountTypeUser
)

type Account struct {
	Id          uint        `json:"id"`
	Name        string      `json:"name"`
	Password    string      `json:"password"`
	AccountType AccountType `json:"account_type"`
}
