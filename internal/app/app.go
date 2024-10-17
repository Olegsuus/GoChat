package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Olegsuus/Auth/internal/config"
	handlers "github.com/Olegsuus/Auth/internal/handlers/handlers"
	routers "github.com/Olegsuus/Auth/internal/handlers/routers"
	services "github.com/Olegsuus/Auth/internal/services/user"
	db "github.com/Olegsuus/Auth/internal/storage/mongo"
	storage "github.com/Olegsuus/Auth/internal/storage/user"
	"github.com/Olegsuus/Auth/internal/tokens"
)

type App struct {
	Server      *http.Server
	Logger      *slog.Logger
	UserHandler *handlers.UserHandler
}

func NewApp(cfg *config.Config) (*App, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	mongoStorage, err := db.NewMongoStorage(cfg)
	if err != nil {
		logger.Error("Ошибка подключения к БД", slog.String("error", err.Error()))
		return nil, err
	}

	userStorage := storage.RegisterStorage(mongoStorage)

	serviceUser := services.RegisterServices(userStorage, logger)

	tokenExpiry, err := time.ParseDuration(cfg.JWT.Expiry)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить срок действия токена: %w", err)
	}

	tokenManager := tokens.NewJWTManager(cfg.JWT.Secret, tokenExpiry)

	userHandler := handlers.RegisterHandlers(serviceUser, tokenManager)

	router := routers.SetupRoutes(userHandler)

	addr := ":" + cfg.Server.Port

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		Server:      server,
		Logger:      logger,
		UserHandler: userHandler,
	}, nil
}

func (app *App) Start() error {
	app.Logger.Info("Сервер запущен на порту " + app.Server.Addr)
	return app.Server.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	app.Logger.Info("Завершаем работу сервера...")
	return app.Server.Shutdown(ctx)
}
