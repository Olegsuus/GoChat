// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"github.com/Olegsuus/GoChat/internal/handlers/chat"
	messageHandlers "github.com/Olegsuus/GoChat/internal/handlers/message"
	"github.com/Olegsuus/GoChat/internal/handlers/routers"
	HandlerUser "github.com/Olegsuus/GoChat/internal/handlers/user"
	"github.com/Olegsuus/GoChat/internal/tokens/jwt"

	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Olegsuus/GoChat/internal/config"
	serviceChat "github.com/Olegsuus/GoChat/internal/services/chat"
	serviceMessage "github.com/Olegsuus/GoChat/internal/services/message"
	services "github.com/Olegsuus/GoChat/internal/services/user"
	storageChat "github.com/Olegsuus/GoChat/internal/storage/chat"
	storageMessage "github.com/Olegsuus/GoChat/internal/storage/message"
	db "github.com/Olegsuus/GoChat/internal/storage/mongo"
	storage "github.com/Olegsuus/GoChat/internal/storage/user"
)

type App struct {
	Server      *http.Server
	Logger      *slog.Logger
	UserHandler *HandlerUser.UserHandler
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

	chatHandler := handlers.RegisterChatHandler(chatSvc, messageSvc)
	messageHandler := messageHandlers.RegisterMessageHandlers(messageSvc)

	userStorage := storage.RegisterStorage(mongoStorage)

	serviceUser := services.RegisterServices(userStorage, logger)

	tokenExpiry, err := time.ParseDuration(cfg.JWT.Expiry)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить срок действия токена: %w", err)
	}

	tokenManager := jwt.NewJWTManager(cfg.JWT.Secret, tokenExpiry)

	userHandler := HandlerUser.RegisterHandlers(serviceUser, tokenManager, cfg)

	// Настройка маршрутов с передачей userService
	router := routers.SetupRoutes(userHandler, tokenManager, chatHandler, messageHandler, serviceUser)

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
