package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danielMensah/faceit-challenge/internal/api"
	"github.com/danielMensah/faceit-challenge/internal/handler"
	"github.com/danielMensah/faceit-challenge/internal/repository/mongo/old"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	router := echo.New()
	router.HideBanner = true

	router.GET("/_healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	swagger, err := api.GetSwagger()
	if err != nil {
		logrus.WithError(err).Fatal("failed to get swagger")
	}

	repo, err := old.New("mongodb://mongo:27017", "faceitchallenge")
	if err != nil {
		logrus.WithError(err).Fatal("failed to create repository")
	}
	defer repo.Close()

	userManagement := handler.NewUserManagement(repo)

	apiGroup := router.Group("", middleware.OapiRequestValidator(swagger))
	api.RegisterHandlersWithBaseURL(apiGroup, userManagement, "/api/v1")

	go func() {
		if err = router.Start("0.0.0.0:8000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Error("starting server")
		}
	}()

	quitGracefully(router)
}

func quitGracefully(router *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("shutting down server")
	}
}
