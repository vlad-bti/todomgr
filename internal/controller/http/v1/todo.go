package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"testcode/test3/internal/controller/http/dto"
	"testcode/test3/internal/domain/entity"
	"testcode/test3/pkg/logger"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, dto entity.Account) error
	GetAccount(ctx context.Context, accountID uint) (*entity.Account, error)
	GetAccountByName(ctx context.Context, name string) (*entity.Account, error)
	GetAccountAll(ctx context.Context) ([]entity.Account, error)
	DeleteAccount(ctx context.Context, accountID uint) error
}

type TodoUsecase interface {
	CreateTodo(ctx context.Context, dto entity.Todo) error
	GetTodo(ctx context.Context, todoID uint) (*entity.Todo, error)
	GetTodoAll(ctx context.Context) ([]entity.Todo, error)
	UpdateTodo(ctx context.Context, dto entity.Todo) error
	DeleteTodo(ctx context.Context, todoID uint) error
}

type SessionUsecase interface {
	Get(key string) (entity.Account, bool)
	Create(account entity.Account) string
	Delete(key string)
}

type todoHandler struct {
	accountUsecase AccountUsecase
	todoUsecase    TodoUsecase
	sessionUsecase SessionUsecase
	log            *logger.Logger
}

func (r *todoHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - Login: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	account, err := r.accountUsecase.GetAccountByName(c.Request.Context(), req.Name)
	if err != nil {
		r.log.Error("http - v1 - Login: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}
	if account == nil || account.Password != req.Password {
		err = errors.New("account not found")
		r.log.Error("http - v1 - Login: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeUnauthenticated, err.Error()))
		return
	}

	resp := r.sessionUsecase.Create(*account)
	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok", LoginResponse{resp}))
}

func (r *todoHandler) Logout(c *gin.Context) {
	token := c.Request.Header.Get(HeaderAuthKey)
	if token != "" {
		r.sessionUsecase.Delete(token)
	}
	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok"))
}

func (r *todoHandler) CreateAccount(c *gin.Context) {
	account := c.MustGet(UserKey).(entity.Account)
	if account.AccountType != entity.AccountTypeAdmin {
		err := errors.New("No access")
		r.log.Error("http - v1 - CreateAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeNoAccess, err.Error()))
		return
	}

	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - CreateAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	acc := entity.Account{
		Name:        req.Name,
		Password:    req.Password,
		AccountType: entity.AccountType(req.AccountType),
	}
	if err := r.accountUsecase.CreateAccount(c.Request.Context(), acc); err != nil {
		r.log.Error("http - v1 - CreateAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok"))
}

func (r *todoHandler) GetAccount(c *gin.Context) {
	var req dto.GetAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - GetAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	resp, err := r.accountUsecase.GetAccount(c.Request.Context(), req.Id)
	if err != nil {
		r.log.Error("http - v1 - GetAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok", resp))
}

func (r *todoHandler) GetAccounts(c *gin.Context) {
	resp, err := r.accountUsecase.GetAccountAll(c.Request.Context())
	if err != nil {
		r.log.Error("http - v1 - GetAccountAll: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok", resp))
}

func (r *todoHandler) DeleteAccount(c *gin.Context) {
	account := c.MustGet(UserKey).(entity.Account)
	if account.AccountType != entity.AccountTypeAdmin {
		err := errors.New("No access")
		r.log.Error("http - v1 - DeleteAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeNoAccess, err.Error()))
		return
	}

	var req dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - DeleteAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	err := r.accountUsecase.DeleteAccount(c.Request.Context(), req.Id)
	if err != nil {
		r.log.Error("http - v1 - DeleteAccount: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok"))
}

func (r *todoHandler) CreateTodo(c *gin.Context) {
	account := c.MustGet(UserKey).(entity.Account)
	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - CreateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	todo := entity.Todo{
		OwnerId: account.Id,
		Name:    req.Name,
		Desc:    req.Desc,
	}
	if err := r.todoUsecase.CreateTodo(c.Request.Context(), todo); err != nil {
		r.log.Error("http - v1 - CreateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok"))
}

func (r *todoHandler) GetTodo(c *gin.Context) {
	var req dto.GetTodoRequest
	if err := c.ShouldBind(&req); err != nil {
		r.log.Error("http - v1 - GetTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	resp, err := r.todoUsecase.GetTodo(c.Request.Context(), req.Id)
	if err != nil {
		r.log.Error("http - v1 - GetTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok", resp))
}

func (r *todoHandler) GetTodos(c *gin.Context) {
	resp, err := r.todoUsecase.GetTodoAll(c.Request.Context())
	if err != nil {
		r.log.Error("http - v1 - GetTodoAll: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok", resp))
}

func (r *todoHandler) UpdateTodo(c *gin.Context) {
	account := c.MustGet(UserKey).(entity.Account)

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - UpdateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	resp, err := r.todoUsecase.GetTodo(c.Request.Context(), req.Id)
	if err != nil {
		r.log.Error("http - v1 - UpdateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}
	if resp == nil {
		err = errors.New("not found")
		r.log.Error("http - v1 - UpdateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}
	if account.AccountType != entity.AccountTypeAdmin && account.Id != resp.OwnerId {
		err = errors.New("No access")
		r.log.Error("http - v1 - UpdateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeNoAccess, err.Error()))
		return
	}

	todo := entity.Todo{
		Id:     req.Id,
		Name:   req.Name,
		Desc:   req.Desc,
		Status: entity.TodoStatus(req.Status),
	}

	if err = r.todoUsecase.UpdateTodo(c.Request.Context(), todo); err != nil {
		r.log.Error("http - v1 - UpdateTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok"))
}

func (r *todoHandler) DeleteTodo(c *gin.Context) {
	account := c.MustGet(UserKey).(entity.Account)

	var req dto.DeleteTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("http - v1 - DeleteTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInvalidArgument, err.Error()))
		return
	}

	resp, err := r.todoUsecase.GetTodo(c.Request.Context(), req.Id)
	if err != nil {
		r.log.Error("http - v1 - DeleteTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}
	if resp == nil {
		err = errors.New("not found")
		r.log.Error("http - v1 - DeleteTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}
	if account.AccountType != entity.AccountTypeAdmin && account.Id != resp.OwnerId {
		err = errors.New("No access")
		r.log.Error("http - v1 - DeleteTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeNoAccess, err.Error()))
		return
	}

	if err = r.todoUsecase.DeleteTodo(c.Request.Context(), req.Id); err != nil {
		r.log.Error("http - v1 - DeleteTodo: %v", err)
		c.JSON(http.StatusOK, NewResp(ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, NewResp(ErrCodeNone, "ok"))
}
