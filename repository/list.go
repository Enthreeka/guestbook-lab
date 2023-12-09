package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lab-3/app/model"
	"lab-3/repository/connect"
)

type ListRepository interface {
	GetByName(ctx context.Context, name string) (*model.List, error)
	GetByID(ctx context.Context, id int) (*model.List, error)
	Create(ctx context.Context, list *model.List) (*model.List, error)
	GetAll(ctx context.Context) ([]model.List, error)
	DeleteByID(сtx context.Context, id int) error
}

type listRepo struct {
	*connect.Postgres
}

func NewListRepository(pg *connect.Postgres) ListRepository {
	return &listRepo{
		pg,
	}
}

func (l *listRepo) collectRow(row pgx.Row) (*model.List, error) {
	var list model.List
	err := row.Scan(&list.ID, &list.Name)
	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	return &list, err
}

func (l *listRepo) collectRows(rows pgx.Rows) ([]model.List, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.List, error) {
		list, err := l.collectRow(row)
		return *list, err
	})
}

func (l *listRepo) GetAll(ctx context.Context) ([]model.List, error) {
	query := `select * from list`

	rows, err := l.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return l.collectRows(rows)
}

func (l *listRepo) GetByName(ctx context.Context, name string) (*model.List, error) {
	query := `select * from list where name = $1`

	row := l.Pool.QueryRow(ctx, query, name)
	return l.collectRow(row)
}

func (l *listRepo) GetByID(ctx context.Context, id int) (*model.List, error) {
	query := `select * from list where id = $1`

	row := l.Pool.QueryRow(ctx, query, id)
	return l.collectRow(row)
}

func (l *listRepo) Create(ctx context.Context, list *model.List) (*model.List, error) {
	query := `insert into list (name) values ($1) returning *`

	row := l.Pool.QueryRow(ctx, query, list.Name)
	return l.collectRow(row)
}
func (l *listRepo) DeleteByID(сtx context.Context, id int) error {
	query := `delete from list where id = $1`

	_, err := l.Pool.Exec(сtx, query, id)
	return err
}
