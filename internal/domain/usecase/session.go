package usecase

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"testcode/test3/internal/domain/entity"
)

type SessionStorage interface {
	Get(key string) (entity.Account, bool)
	Set(key string, account entity.Account)
	Delete(key string)
}

func sha256hash(ori string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(ori))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type sessionUsecase struct {
	storage SessionStorage
}

func NewSessionUsecase(storage SessionStorage) *sessionUsecase {
	return &sessionUsecase{
		storage: storage,
	}
}

func (r *sessionUsecase) Get(key string) (entity.Account, bool) {
	return r.storage.Get(key)
}

func (r *sessionUsecase) Create(account entity.Account) string {
	id := uuid.New()
	key := sha256hash(id.String())
	r.storage.Set(key, account)
	return key
}

func (r *sessionUsecase) Delete(key string) {
	r.storage.Delete(key)
}
