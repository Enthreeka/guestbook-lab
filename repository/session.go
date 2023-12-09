package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lab-3/app/model"
	"lab-3/repository/connect"
)

type SessionRepository interface {
	GetToken(ctx context.Context, token string) (*model.Session, error)
	CreateToken(ctx context.Context, token *model.Session) error
	UpdateToken(ctx context.Context, userID int, token string) (*model.Session, error)
}

type sessionRepo struct {
	*connect.Postgres
}

func NewSessionRepository(pg *connect.Postgres) SessionRepository {
	return &sessionRepo{
		pg,
	}
}

func (l *sessionRepo) collectRow(row pgx.Row) (*model.Session, error) {
	var session model.Session
	err := row.Scan(&session.ID, &session.Token, &session.UserID)
	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}

	return &session, err
}

func (l *sessionRepo) GetToken(ctx context.Context, token string) (*model.Session, error) {
	query := `SELECT id,token,user_id FROM session WHERE token = $1`

	row := l.Pool.QueryRow(ctx, query, token)
	return l.collectRow(row)
}

func (l *sessionRepo) CreateToken(ctx context.Context, token *model.Session) error {
	query := `INSERT INTO session (token, user_id) VALUES ($1,$2)`

	_, err := l.Pool.Exec(ctx, query, token.Token, token.UserID)
	return err
}

func (l *sessionRepo) UpdateToken(ctx context.Context, userID int, token string) (*model.Session, error) {
	query := `UPDATE session SET token = $1 WHERE user_id = $2 returning *`

	row := l.Pool.QueryRow(ctx, query, token, userID)

	return l.collectRow(row)
}
