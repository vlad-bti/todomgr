package mysql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"testcode/test3/internal/domain/entity"
	"testcode/test3/pkg/mysql"
)

type accountStorage struct {
	baseStorage
}

func NewAccountStorage(db *mysql.Mysql) *accountStorage {
	return &accountStorage{
		baseStorage{db},
	}
}

func (r *accountStorage) Create(ctx context.Context, dto entity.Account) error {
	sql, args, err := r.db.Builder.
		Insert("account").
		Columns("name, password, account_type").
		Values(dto.Name, dto.Password, dto.AccountType).
		ToSql()
	if err != nil {
		return fmt.Errorf("AccountStorage - Create - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Create - r.Exec: %w", err)
	}

	return nil
}

func (r *accountStorage) Get(ctx context.Context, accountID uint) (*entity.Account, error) {
	sql, args, err := r.db.Builder.
		Select("id, name, password, account_type").
		From("account").
		Where(sq.Eq{"id": accountID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AccountStorage - Get - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AccountStorage - Get - r.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		e := entity.Account{}
		err = rows.Scan(&e.Id, &e.Name, &e.Password, &e.AccountType)
		if err != nil {
			return nil, fmt.Errorf("AccountStorage - Get - rows.Scan: %w", err)
		}
		return &e, nil
	}
	return nil, nil
}

func (r *accountStorage) GetByName(ctx context.Context, name string) (*entity.Account, error) {
	sql, args, err := r.db.Builder.
		Select("id, name, password, account_type").
		From("account").
		Where(sq.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AccountStorage - GetByName - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AccountStorage - GetByName - r.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		e := entity.Account{}
		err = rows.Scan(&e.Id, &e.Name, &e.Password, &e.AccountType)
		if err != nil {
			return nil, fmt.Errorf("AccountStorage - GetByName - rows.Scan: %w", err)
		}
		return &e, nil
	}
	return nil, nil
}

func (r *accountStorage) GetAll(ctx context.Context) ([]entity.Account, error) {
	sql, args, err := r.db.Builder.
		Select("id, name, password, account_type").
		From("account").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AccountStorage - GetAll - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AccountStorage - GetAll - r.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Account, 0, _defaultEntityCap)
	for rows.Next() {
		e := entity.Account{}
		err = rows.Scan(&e.Id, &e.Name, &e.Password, &e.AccountType)
		if err != nil {
			return nil, fmt.Errorf("AccountStorage - GetAll - rows.Scan: %w", err)
		}
		entities = append(entities, e)
	}
	return entities, nil
}

func (r *accountStorage) Delete(ctx context.Context, accountID uint) error {
	sql, args, err := r.db.Builder.
		Delete("account").
		Where(sq.Eq{"id": accountID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("AccountStorage - Delete - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountStorage - Delete - r.Exec: %w", err)
	}
	return nil
}
