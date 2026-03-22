// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"log"
	"main/storage"
	"main/storage/migrations"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"main/internal/api/grpc"
	ht "main/internal/api/http"
	config "main/internal/config"
	merger "main/internal/core/merger"
	"main/internal/core/validator"
	logger "main/internal/logger"
	repo "main/internal/repo"
	watcher "main/internal/watcher"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var cfg *config.Config

func init() {
	cfg = config.MustLoad()
}

// Start запускает HTTP и gRPC серверы
func Start() {
	// Инициализируем логгеры
	apilogger := logger.GetAPILogger(
		fmt.Sprintf("%s/%s", cfg.Logger.LogsDir, cfg.Logger.APIlp),
	)
	dblogger := logger.GetDBLogger(
		fmt.Sprintf("%s/%s", cfg.Logger.LogsDir, cfg.Logger.DBlp),
	)

	// Подключаемся к БД
	db := checkUpDB(dblogger)

	// Инициализируем компоненты ядра
	configRepo := repo.NewConfigRepository(db)
	watcherManager := watcher.NewWatcherManager()
	merger := merger.NewMerger()

	// Инициализируем валидатор (если есть схема)
	var val *validator.Validator = validator.NewValidator()
	if cfg.Validator.SchemaPath != "" {
		var err error
		validator, err = validator.NewValidator(cfg.Validator.SchemaPath)
		if err != nil {
			apilogger.Printf("Warning: Failed to load validator schema: %v", err)
		}
	}

	// Создаем HTTP сервер
	httpServer := ht.NewHTTPServer(
		cfg.HTTPServer.Host,
		cfg.HTTPServer.Port,
		apilogger,
	)

	// Создаем gRPC сервер
	grpcServer := grpc.NewGRPCServer(
		configRepo,
		watcherManager,
	)

	// Контекст для graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Запускаем HTTP сервер
	go func() {
		apilogger.Printf("HTTP Server starting on %s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
		log.Printf("HTTP server listening on %s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Запускаем gRPC сервер
	go func() {
		apilogger.Printf("gRPC Server starting on %s:%d", cfg.GRPCServer.Host, cfg.GRPCServer.Port)
		log.Printf("gRPC server listening on %s:%d", cfg.GRPCServer.Host, cfg.GRPCServer.Port)

		if err := grpcServer.Start(cfg.GRPCServer.Host, cfg.GRPCServer.Port); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down servers...")
	apilogger.Println("Received shutdown signal, stopping servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer shutdownCancel()

	httpShutdownDone := make(chan struct{})
	go func() {
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		close(httpShutdownDone)
	}()

	grpcShutdownDone := make(chan struct{})
	go func() {
		grpcServer.GracefulStop(cfg.GRPCServer.ShutdownTimeout)
		close(grpcShutdownDone)
	}()

	select {
	case <-httpShutdownDone:
		log.Println("HTTP server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("HTTP server shutdown timeout, forcing close...")
		if err := httpServer.Close(); err != nil {
			log.Printf("HTTP server force close error: %v", err)
		}
	}

	select {
	case <-grpcShutdownDone:
		log.Println("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("gRPC server shutdown timeout, forcing stop...")
		grpcServer.Stop()
	}

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	log.Println("Config Service stopped")
}

func checkUpDB(logger gormlogger.Interface) *gorm.DB {
	db := storage.Psql(cfg.Db.Role, cfg.Db.Pass, cfg.Db.Name, cfg.Db.Host, cfg.Db.Port, logger)
	migrations.RunMigration(db)

	return db
}
