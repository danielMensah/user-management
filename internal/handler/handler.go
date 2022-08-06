package handler

import (
	"net/http"

	"github.com/danielMensah/faceit-challenge/internal/api"
	"github.com/danielMensah/faceit-challenge/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo repository.UserRepository
}

func NewUserManagement(repo repository.UserRepository) *Service {
	return &Service{repo}
}

func (s *Service) GetUsers(ctx echo.Context, params api.GetUsersParams) error {
	users, err := s.repo.GetUsers(ctx.Request().Context(), params)
	if err != nil {
		logrus.WithError(err).Error("failed to get users")
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: "failed to get users"})
	}

	return ctx.JSON(http.StatusOK, api.GetUsersResponse{Users: users})
}

func (s *Service) CreateUser(ctx echo.Context) error {
	body := new(api.User)
	if err := ctx.Bind(body); err != nil {
		logrus.WithError(err).Error("failed to bind body")
		return ctx.JSON(http.StatusBadRequest, api.Error{Message: "failed to parse request body"})
	}

	id, err := s.repo.CreateUser(ctx.Request().Context(), body)
	if err != nil {
		logrus.WithError(err).Error("failed to create user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: "failed to create user"})
	}

	return ctx.JSON(http.StatusCreated, api.CreateUserResponse{Id: id})
}

func (s *Service) UpdateUser(ctx echo.Context, id api.IdPath) error {
	body := new(api.UserUpdateData)
	if err := ctx.Bind(body); err != nil {
		logrus.WithError(err).Error("failed to bind body")
		return ctx.JSON(http.StatusBadRequest, api.Error{Message: "failed to parse request body"})
	}

	user, err := s.repo.UpdateUser(ctx.Request().Context(), id, body)
	if err != nil {
		logrus.WithError(err).Error("failed to update user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: "failed to update user"})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (s *Service) DeleteUser(ctx echo.Context, id string) error {
	if err := s.repo.DeleteUser(ctx.Request().Context(), id); err != nil {
		logrus.WithError(err).Error("failed to delete user")
		return ctx.JSON(http.StatusInternalServerError, api.Error{Message: "failed to delete user"})
	}

	return ctx.JSON(http.StatusNoContent, nil)
}
