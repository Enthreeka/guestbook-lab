package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lab-3/app/model"
	"lab-3/repository/connect"
	"time"
)

type GuestbookRepository interface {
	GetAllByListID(ctx context.Context, listID int) ([]model.Guestbook, error)
	Create(ctx context.Context, guestbook *model.Guestbook) (*model.Guestbook, error)
	Delete(ctx context.Context, id int) error
}

type guestbookRepo struct {
	*connect.Postgres
}

func NewGuestbookRepository(pg *connect.Postgres) GuestbookRepository {
	return &guestbookRepo{
		pg,
	}
}

func (l *guestbookRepo) collectRow(row pgx.Row) (*model.Guestbook, error) {
	var guestbook model.Guestbook
	var t time.Time
	err := row.Scan(&guestbook.ID, &guestbook.Message, &t, &guestbook.UserID, &guestbook.ListID)
	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	nt := t.Format("15:04 02.01.2006")
	guestbook.CreatedAt = nt

	return &guestbook, err
}

func (l *guestbookRepo) collectRows(rows pgx.Rows) ([]model.Guestbook, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Guestbook, error) {
		guestbook, err := l.collectRow(row)
		return *guestbook, err
	})
}

func (g *guestbookRepo) GetAllByListID(ctx context.Context, listID int) ([]model.Guestbook, error) {
	query := `select g.id, g.message, g.created_at, g.user_id, g.list_id, u.login, l.name  from guestbook g
				join "user" u on g.user_id = u.id
				join list l on g.list_id = l.id
				where g.list_id = $1
				order by g.created_at desc`

	rows, err := g.Pool.Query(ctx, query, listID)
	if err != nil {
		return nil, err
	}

	data := make([]model.Guestbook, 0, 50)
	for rows.Next() {
		var g model.Guestbook
		var t time.Time

		err := rows.Scan(&g.ID, &g.Message, &t, &g.UserID, &g.ListID, &g.UserName, &g.ListName)
		if err != nil {
			if err == pgx.ErrNoRows {
				return []model.Guestbook{}, nil
			}
			return nil, err
		}
		nt := t.Format("15:04 02.01.2006")
		g.CreatedAt = nt

		data = append(data, g)
	}

	return data, nil
}

func (g *guestbookRepo) Create(ctx context.Context, guestbook *model.Guestbook) (*model.Guestbook, error) {
	query := `insert into guestbook (message,user_id,list_id) values ($1,$2,$3) returning *`

	row := g.Pool.QueryRow(ctx, query, guestbook.Message, guestbook.UserID, guestbook.ListID)
	return g.collectRow(row)
}

func (g *guestbookRepo) Delete(ctx context.Context, id int) error {
	query := `delete from guestbook where id = $1`

	_, err := g.Pool.Exec(ctx, query, id)
	return err
}
