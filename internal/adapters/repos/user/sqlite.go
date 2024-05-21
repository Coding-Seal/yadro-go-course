package user

import (
	"context"
	"database/sql"
	"errors"

	"yadro-go-course/internal/core/models"
	"yadro-go-course/internal/core/ports"
)

type SqliteRepo struct {
	db *sql.DB
}

var _ ports.UserRepo = (*SqliteRepo)(nil)

func NewSqliteRepo(db *sql.DB) *SqliteRepo {
	return &SqliteRepo{
		db: db,
	}
}

func (s *SqliteRepo) UserID(ctx context.Context, id int64) (models.User, error) {
	var user models.User

	row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE user_id = $1", id)

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.Join(ports.ErrNotFound, err)
		}

		return user, errors.Join(ports.ErrInternal, err)
	}

	return user, nil
}

func (s *SqliteRepo) UserLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE login = $1", login)

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.Join(ports.ErrNotFound, err)
		}

		return user, errors.Join(ports.ErrInternal, err)
	}

	return user, nil
}

func (s *SqliteRepo) RemoveUser(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", id)
	if err != nil {
		return errors.Join(err, ports.ErrInternal)
	}

	return nil
}

func (s *SqliteRepo) AddUser(ctx context.Context, user *models.User) error {
	row := s.db.QueryRowContext(ctx, "INSERT INTO users (login, password, is_admin) VALUES ($1, $2, $3) RETURNING user_id", user.Login, user.Password, user.IsAdmin)

	err := row.Scan(&user.ID)
	if err != nil {
		return errors.Join(err, ports.ErrInternal)
	}

	return nil
}
