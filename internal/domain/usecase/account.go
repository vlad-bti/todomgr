package usecase

import (
	"context"

	"testcode/test3/internal/domain/entity"
	"testcode/test3/pkg/logger"
)

type AccountStorage interface {
	Create(ctx context.Context, dto entity.Account) error
	Get(ctx context.Context, accountID uint) (*entity.Account, error)
	GetByName(ctx context.Context, name string) (*entity.Account, error)
	GetAll(ctx context.Context) ([]entity.Account, error)
	Delete(ctx context.Context, accountID uint) error
}

type accountUsecase struct {
	storage AccountStorage
	log     *logger.Logger
}

func NewAccountUsecase(log *logger.Logger, storage AccountStorage) *accountUsecase {
	return &accountUsecase{
		storage: storage,
		log:     log,
	}
}

func (r *accountUsecase) CreateAccount(ctx context.Context, dto entity.Account) error {
	if err := r.storage.Create(ctx, dto); err != nil {
		r.log.Error("AccountUsecase - CreateAccount - r.storage.Create: %v; Name=%v, Password=%v, AccountType=%v",
			err,
			dto.Name,
			dto.Password,
			dto.AccountType,
		)
		return err
	}
	return nil
}

func (r *accountUsecase) GetAccount(ctx context.Context, accountID uint) (*entity.Account, error) {
	ret, err := r.storage.Get(ctx, accountID)
	if err != nil {
		r.log.Error("AccountUsecase - GetAccount - r.storage.Get: %v; accountID=%v", err, accountID)
		return nil, err
	}
	return ret, nil
}

func (r *accountUsecase) GetAccountByName(ctx context.Context, name string) (*entity.Account, error) {
	ret, err := r.storage.GetByName(ctx, name)
	if err != nil {
		r.log.Error("AccountUsecase - GetAccountByName - r.storage.GetByName: %v; name=%v", err, name)
		return nil, err
	}
	return ret, nil
}

func (r *accountUsecase) GetAccountAll(ctx context.Context) ([]entity.Account, error) {
	ret, err := r.storage.GetAll(ctx)
	if err != nil {
		r.log.Error("AccountUsecase - GetAccountAll - r.storage.GetAll: %v", err)
		return nil, err
	}
	return ret, nil
}

func (r *accountUsecase) DeleteAccount(ctx context.Context, accountID uint) error {
	if err := r.storage.Delete(ctx, accountID); err != nil {
		r.log.Error("AccountUsecase - DeleteAccount - r.storage.Delete: %v", err)
		return err
	}
	return nil
}
