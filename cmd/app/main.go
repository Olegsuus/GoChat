package main

import (
	"context"
	"github.com/Olegsuus/Auth/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Olegsuus/Auth/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Start(); err != nil && err != http.ErrServerClosed {
			application.Logger.Error("Ошибка работы сервера", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}

	log.Println("Сервер успешно завершил работу")
}
