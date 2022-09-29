package mysql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"testcode/test3/internal/domain/entity"
	"testcode/test3/pkg/mysql"
)

type todoStorage struct {
	baseStorage
}

func NewTodoStorage(db *mysql.Mysql) *todoStorage {
	return &todoStorage{
		baseStorage{db},
	}
}

func (r *todoStorage) Create(ctx context.Context, dto entity.Todo) error {
	sql, args, err := r.db.Builder.
		Insert("todo").
		Columns("owner_id, name, `desc`, status").
		Values(dto.OwnerId, dto.Name, dto.Desc, entity.TodoStatusDefault).
		ToSql()
	if err != nil {
		return fmt.Errorf("TodoStorage - Create - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TodoStorage - Create - r.Exec: %w", err)
	}

	return nil
}

func (r *todoStorage) Get(ctx context.Context, todoID uint) (*entity.Todo, error) {
	sql, args, err := r.db.Builder.
		Select("id, owner_id, name, `desc`, status").
		From("todo").
		Where(sq.Eq{"id": todoID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TodoStorage - Get - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TodoStorage - Get - r.Query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		e := entity.Todo{}
		err = rows.Scan(&e.Id, &e.OwnerId, &e.Name, &e.Desc, &e.Status)
		if err != nil {
			return nil, fmt.Errorf("TodoStorage - Get - rows.Scan: %w", err)
		}
		return &e, nil
	}
	return nil, nil
}

func (r *todoStorage) GetAll(ctx context.Context) ([]entity.Todo, error) {
	sql, _, err := r.db.Builder.
		Select("id, owner_id, name, `desc`, status").
		From("todo").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TodoStorage - GetAll - r.Builder: %w", err)
	}

	rows, err := r.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TodoStorage - GetAll - r.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Todo, 0, _defaultEntityCap)
	for rows.Next() {
		e := entity.Todo{}
		err = rows.Scan(&e.Id, &e.OwnerId, &e.Name, &e.Desc, &e.Status)
		if err != nil {
			return nil, fmt.Errorf("TodoStorage - GetAll - rows.Scan: %w", err)
		}
		entities = append(entities, e)
	}
	return entities, nil
}

func (r *todoStorage) Update(ctx context.Context, dto entity.Todo) error {
	builder := r.db.Builder.Update("todo")
	if dto.Name != "" {
		builder = builder.Set("name", dto.Name)
	}
	if dto.Desc != "" {
		builder = builder.Set("`desc`", dto.Desc)
	}
	if dto.Status > 0 {
		builder = builder.Set("status", dto.Status)
	}
	sql, args, err := builder.Where(sq.Eq{"id": dto.Id}).ToSql()
	if err != nil {
		return fmt.Errorf("TodoStorage - Update - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TodoStorage - Update - r.Exec: %w", err)
	}

	return nil
}

func (r *todoStorage) Delete(ctx context.Context, todoID uint) error {
	sql, args, err := r.db.Builder.
		Delete("todo").
		Where(sq.Eq{"id": todoID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("TodoStorage - Delete - r.Builder: %w", err)
	}

	_, err = r.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TodoStorage - Delete - r.Exec: %w", err)
	}
	return nil
}
