package session

import "testcode/test3/internal/domain/entity"

type sessionStorage struct {
	m map[string]entity.Account
}

func NewSessionStorage() *sessionStorage {
	return &sessionStorage{m: make(map[string]entity.Account)}
}

func (r *sessionStorage) Get(key string) (entity.Account, bool) {
	if v, ok := r.m[key]; ok {
		return v, true
	}
	return entity.Account{}, false
}

func (r *sessionStorage) Set(key string, account entity.Account) {
	r.m[key] = account
}

func (r *sessionStorage) Delete(key string) {
	delete(r.m, key)
}
