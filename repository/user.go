package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lab-3/app/model"
	"lab-3/repository/connect"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
}

type userRepository struct {
	*connect.Postgres
}

func NewUserRepository(pg *connect.Postgres) UserRepository {
	return &userRepository{
		pg,
	}
}

func (u *userRepository) collectRow(row pgx.Row) (*model.User, error) {
	var user model.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err == pgx.ErrNoRows {
		return nil, ErrNoRows
	}
	errCode := ErrorCode(err)
	if errCode == ForeignKeyViolation {
		return nil, ErrForeignKeyViolation
	}
	if errCode == UniqueViolation {
		return nil, ErrUniqueViolation
	}

	return &user, err
}

func (u *userRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `select u.id,u.login, u.password
				from "user" u
				where login = $1`

	row := u.Pool.QueryRow(ctx, query, login)
	return u.collectRow(row)
}

func (u *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `insert into "user" (login,password) values ($1,$2) returning *`

	row := u.Pool.QueryRow(ctx, query, user.Login, user.Password)
	return u.collectRow(row)
}

func (u *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `select *
				from "user" u
				where id = $1`

	row := u.Pool.QueryRow(ctx, query, id)
	return u.collectRow(row)
}
