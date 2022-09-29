package usecase

import (
	"context"

	"testcode/test3/internal/domain/entity"
	"testcode/test3/pkg/logger"
)

type TodoStorage interface {
	Create(ctx context.Context, dto entity.Todo) error
	Get(ctx context.Context, todoID uint) (*entity.Todo, error)
	GetAll(ctx context.Context) ([]entity.Todo, error)
	Update(ctx context.Context, dto entity.Todo) error
	Delete(ctx context.Context, todoID uint) error
}

type Notification interface {
	Send(msg interface{})
}

type todoUsecase struct {
	storage      TodoStorage
	notification Notification
	log          *logger.Logger
}

func NewTodoUsecase(log *logger.Logger, storage TodoStorage, notification Notification) *todoUsecase {
	return &todoUsecase{
		storage:      storage,
		notification: notification,
		log:          log,
	}
}

func (r *todoUsecase) CreateTodo(ctx context.Context, dto entity.Todo) error {
	if err := r.storage.Create(ctx, dto); err != nil {
		r.log.Error("TodoUsecase - CreateTodo - r.storage.Create: %v; OwnerId=%v, Name=%v, Desc=%v, Status=%v",
			err,
			dto.OwnerId,
			dto.Name,
			dto.Desc,
			dto.Status,
		)
		return err
	}
	r.notification.Send(dto)
	return nil
}

func (r *todoUsecase) GetTodo(ctx context.Context, todoID uint) (*entity.Todo, error) {
	ret, err := r.storage.Get(ctx, todoID)
	if err != nil {
		r.log.Error("TodoUsecase - GetTodo - r.storage.Get: %v; todoID=%v", err, todoID)
		return nil, err
	}
	return ret, nil
}

func (r *todoUsecase) GetTodoAll(ctx context.Context) ([]entity.Todo, error) {
	ret, err := r.storage.GetAll(ctx)
	if err != nil {
		r.log.Error("TodoUsecase - GetTodoAll - r.storage.GetAll: %v", err)
		return nil, err
	}
	return ret, nil
}

func (r *todoUsecase) UpdateTodo(ctx context.Context, dto entity.Todo) error {
	if err := r.storage.Update(ctx, dto); err != nil {
		r.log.Error("TodoUsecase - UpdateTodo - r.storage.Update: %v; ID=%v, Name=%v, Desc=%v, Status=%v",
			err,
			dto.Id,
			dto.Name,
			dto.Desc,
			dto.Status,
		)
		return err
	}
	return nil
}

func (r *todoUsecase) DeleteTodo(ctx context.Context, todoID uint) error {
	if err := r.storage.Delete(ctx, todoID); err != nil {
		r.log.Error("TodoUsecase - DeleteTodo - r.storage.Delete: %v", err)
		return err
	}
	return nil
}
