package app

import (
	"context"
	"fmt"
	handlersChat "github.com/Olegsuus/Auth/internal/handlers/chat"
	handlersMessage "github.com/Olegsuus/Auth/internal/handlers/message"
	"github.com/Olegsuus/Auth/internal/handlers/routers"
	"github.com/Olegsuus/Auth/internal/handlers/user"
	"github.com/Olegsuus/Auth/internal/tokens/jwt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Olegsuus/Auth/internal/config"
	serviceChat "github.com/Olegsuus/Auth/internal/services/chat"
	serviceMessage "github.com/Olegsuus/Auth/internal/services/message"
	services "github.com/Olegsuus/Auth/internal/services/user"
	storageChat "github.com/Olegsuus/Auth/internal/storage/chat"
	storageMessage "github.com/Olegsuus/Auth/internal/storage/message"
	db "github.com/Olegsuus/Auth/internal/storage/mongo"
	storage "github.com/Olegsuus/Auth/internal/storage/user"
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

	chatStore := storageChat.RegisterStorageChat(mongoStorage)
	messageStore := storageMessage.RegisterStorageMessage(mongoStorage)

	chatSvc := serviceChat.RegisterChatService(chatStore, logger)
	messageSvc := serviceMessage.RegisterServiceMessage(messageStore, logger)

	chatHandler := handlersChat.RegisterChatHandler(chatSvc, messageSvc)
	messageHandler := handlersMessage.RegisterMessageHandlers(messageSvc)

	userStorage := storage.RegisterStorage(mongoStorage)

	serviceUser := services.RegisterServices(userStorage, logger)

	tokenExpiry, err := time.ParseDuration(cfg.JWT.Expiry)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить срок действия токена: %w", err)
	}

	tokenManager := jwt.NewJWTManager(cfg.JWT.Secret, tokenExpiry)

	userHandler := handlers.RegisterHandlers(serviceUser, tokenManager, cfg)

	router := routers.SetupRoutes(userHandler, tokenManager, chatHandler, messageHandler)

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
