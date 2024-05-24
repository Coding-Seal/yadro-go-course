package ports

import (
	"context"

	"yadro-go-course/internal/core/models"
)

type UserRepo interface {
	UserID(ctx context.Context, id int64) (models.User, error)
	UserLogin(ctx context.Context, login string) (models.User, error)
	RemoveUser(ctx context.Context, id int64) error
	AddUser(ctx context.Context, user *models.User) error
}
