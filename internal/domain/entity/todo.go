package entity

type TodoStatus uint

const (
	_ TodoStatus = iota
	TodoStatusDefault
	TodoStatusDone
)

type Todo struct {
	Id      uint       `json:"id"`
	OwnerId uint       `json:"owner_id"`
	Name    string     `json:"name"`
	Desc    string     `json:"desc"`
	Status  TodoStatus `json:"status"`
}
