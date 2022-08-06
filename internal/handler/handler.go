package handler

import (
	"net/http"

	"github.com/danielMensah/user-management/internal/api"
	"github.com/danielMensah/user-management/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	errGetUsers   = "failed to get users"
	errParseBody  = "failed to parse request body"
	errCreateUser = "failed to create user"
	errUpdateUser = "failed to update user"
	errDeleteUser = "failed to delete user"
	errEncryptPwd = "failed to encrypt password"
)

// Handler represents handlers for user management
type Handler struct {
	repo repository.UserRepository
}

func (h *Handler) GetHealthz(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

// New creates a new user handler
func New(repo repository.UserRepository) *Handler {
	return &Handler{repo}
}

// GetUsers returns a list of users
func (h *Handler) GetUsers(ctx echo.Context, params api.GetUsersParams) error {
	users, err := h.repo.GetUsers(ctx.Request().Context(), params)
	if err != nil {
		logrus.WithError(err).Error(errGetUsers)
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: errGetUsers})
	}

	return ctx.JSON(http.StatusOK, api.GetUsersResponse{Users: users})
}

// CreateUser creates a new user
func (h *Handler) CreateUser(ctx echo.Context) error {
	var err error
	body := new(api.UserCreateData)
	if err = ctx.Bind(body); err != nil {
		logrus.WithError(err).Error(errParseBody)
		return ctx.JSON(http.StatusBadRequest, api.Error{Message: errParseBody})
	}

	if body.Password, err = encryptPassword(body.Password); err != nil {
		logrus.WithError(err).Error(errCreateUser)
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: errEncryptPwd})
	}

	id, err := h.repo.CreateUser(ctx.Request().Context(), body)
	if err != nil {
		logrus.WithError(err).Error(errCreateUser)
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: errCreateUser})
	}

	return ctx.JSON(http.StatusCreated, api.CreateUserResponse{Id: id})
}

// UpdateUser updates a user
func (h *Handler) UpdateUser(ctx echo.Context, id string) error {
	body := new(api.UserUpdateData)
	if err := ctx.Bind(body); err != nil {
		logrus.WithError(err).Error(errParseBody)
		return ctx.JSON(http.StatusBadRequest, api.Error{Message: errParseBody})
	}

	if body.Password != nil && *body.Password != "" {
		p, err := encryptPassword(*body.Password)
		if err != nil {
			logrus.WithError(err).Error(errEncryptPwd)
			return ctx.JSON(http.StatusInternalServerError, api.Error{Message: errEncryptPwd})
		}

		body.Password = &p
	}

	user, err := h.repo.UpdateUser(ctx.Request().Context(), id, body)
	if err != nil {
		logrus.WithError(err).Error(errUpdateUser)
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: errUpdateUser})
	}

	return ctx.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(ctx echo.Context, id string) error {
	if err := h.repo.DeleteUser(ctx.Request().Context(), id); err != nil {
		logrus.WithError(err).Error(errDeleteUser)
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: errDeleteUser})
	}

	return ctx.JSON(http.StatusNoContent, nil)
}
