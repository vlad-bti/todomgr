// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"testcode/test3/config"
	"testcode/test3/internal/adapters/db/mysql"
	"testcode/test3/internal/adapters/db/session"
	"testcode/test3/internal/adapters/notification/telegram"
	v1 "testcode/test3/internal/controller/http/v1"
	"testcode/test3/internal/domain/usecase"
	"testcode/test3/pkg/httpserver"
	"testcode/test3/pkg/logger"
	mysqlpool "testcode/test3/pkg/mysql"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	// Repository
	db, err := mysqlpool.New(cfg.MySql.URL, mysqlpool.MaxPoolSize(cfg.MySql.PoolMax))
	if err != nil {
		log.Fatal("app - Run - mysql.New: %v", err)
	}
	defer db.Close()

	accountStorage := mysql.NewAccountStorage(db)
	todoStorage := mysql.NewTodoStorage(db)
	sessionStorage := session.NewSessionStorage()

	// Notification
	telegramNotification := telegram.NewTelegramNotification(log)

	// Use case
	accountUsecase := usecase.NewAccountUsecase(log, accountStorage)
	todoUsecase := usecase.NewTodoUsecase(log, todoStorage, telegramNotification)
	sessionUsecase := usecase.NewSessionUsecase(sessionStorage)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, log, accountUsecase, todoUsecase, sessionUsecase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: %v", s.String())
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown: %v", err)
	}
}
