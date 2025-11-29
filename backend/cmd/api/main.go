package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zentinel/internal/application/usecase"
	"zentinel/internal/infrastructure/driven/ai"
	"zentinel/internal/infrastructure/driven/persistence/postgres"
	zhttp "zentinel/internal/infrastructure/driving/http"
	"zentinel/internal/infrastructure/driving/http/handler"
	"zentinel/pkg/config"
	"zentinel/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// 2. Initialize Logger
	logger.Init(cfg.LogLevel)
	defer logger.Sync()
	logger.Info("Starting Zentinel API", zap.String("port", cfg.ServerPort))

	// 3. Database Connection
	db, err := postgres.NewConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	// 4. Initialize Driven Adapters (Repositories & External Services)
	clientRepo := postgres.NewClientRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	alertRepo := postgres.NewAlertRepository(db)

	aiService := ai.NewMockAnalyzer()

	// 5. Initialize Application Layer (Use Cases)
	clientUseCase := usecase.NewClientUseCase(clientRepo, accountRepo)
	accountUseCase := usecase.NewAccountUseCase(accountRepo, clientRepo)
	alertUseCase := usecase.NewAlertUseCase(alertRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, accountRepo, clientRepo, alertRepo, aiService)
	dashboardUseCase := usecase.NewDashboardUseCase(transactionRepo, alertRepo)

	// 6. Initialize Driving Adapters (HTTP Handlers)
	healthHandler := handler.NewHealthHandler(db)
	clientHandler := handler.NewClientHandler(clientUseCase)
	accountHandler := handler.NewAccountHandler(accountUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	alertHandler := handler.NewAlertHandler(alertUseCase)
	dashboardHandler := handler.NewDashboardHandler(dashboardUseCase)

	// 7. Setup Router
	router := zhttp.SetupRouter(
		healthHandler,
		clientHandler,
		accountHandler,
		transactionHandler,
		alertHandler,
		dashboardHandler,
	)

	// 8. Start Server
	srv := zhttp.NewServer(cfg, router)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 9. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
