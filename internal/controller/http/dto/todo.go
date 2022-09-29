package dto

type CreateTodoRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type GetTodoRequest struct {
	Id uint `json:"id" form:"id" binding:"required"`
}

type UpdateTodoRequest struct {
	Id     uint   `json:"id" binding:"required"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status uint   `json:"status" binding:"omitempty,oneof=1 2"`
}

type DeleteTodoRequest struct {
	Id uint `json:"id" binding:"required"`
}
