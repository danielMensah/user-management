package repository

import (
	"context"

	"github.com/danielMensah/faceit-challenge/internal/api"
)

// UserRepository represents the user repository contract
type UserRepository interface {
	GetUsers(ctx context.Context, params api.GetUsersParams) (*[]api.User, error)
	CreateUser(ctx context.Context, user *api.User) (string, error)
	UpdateUser(ctx context.Context, id string, data *api.UserUpdateData) (*api.User, error)
	DeleteUser(ctx context.Context, id string) error
	Close(ctx context.Context) error
}
