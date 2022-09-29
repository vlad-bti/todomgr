package v1

import (
	"github.com/gin-gonic/gin"
	"testcode/test3/pkg/logger"
)

func NewRouter(handler *gin.Engine, log *logger.Logger, accountUsecase AccountUsecase, todoUsecase TodoUsecase, sessionUsecase SessionUsecase) {
	r := &todoHandler{accountUsecase, todoUsecase, sessionUsecase, log}

	handler.Use(Auth(sessionUsecase, "/v1/login"))
	// Routers
	h := handler.Group("/v1")
	{
		h.POST("/login", r.Login)
		h.POST("/logout", r.Logout)

		h.GET("/accounts", r.GetAccounts)
		h.POST("/account", r.CreateAccount)
		h.DELETE("/account", r.DeleteAccount)

		h.GET("/todos", r.GetTodos)
		h.GET("/todo", r.GetTodo)
		h.POST("/todo", r.CreateTodo)
		h.PUT("/todo", r.UpdateTodo)
		h.DELETE("/todo", r.DeleteTodo)
	}
}
