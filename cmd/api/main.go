package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danielMensah/user-management/internal/api"
	"github.com/danielMensah/user-management/internal/config"
	"github.com/danielMensah/user-management/internal/handler"
	mongoRepo "github.com/danielMensah/user-management/internal/repository/mongo"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	router := echo.New()
	router.HideBanner = true

	router.GET("/_healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	swagger, err := api.GetSwagger()
	if err != nil {
		logrus.WithError(err).Fatal("failed to get swagger")
	}

	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to mongo")
	}
	defer func() {
		if err = conn.Disconnect(context.Background()); err != nil {
			logrus.WithError(err).Error("disconnecting from database")
		}
	}()

	repo := mongoRepo.New(conn.Database(cfg.MongoDB))
	handlers := handler.New(repo)

	apiGroup := router.Group("", middleware.OapiRequestValidator(swagger))
	api.RegisterHandlersWithBaseURL(apiGroup, handlers, "/api/v1")

	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort)
		if err = router.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
