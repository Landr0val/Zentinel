package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zentinel/internal/application/usecase"
	"zentinel/internal/domain/service"
	"zentinel/internal/infrastructure/driven/ai"
	"zentinel/internal/infrastructure/driven/blockchain"
	"zentinel/internal/infrastructure/driven/persistence/postgres"
	zhttp "zentinel/internal/infrastructure/driving/http"
	"zentinel/internal/infrastructure/driving/http/handler"
	"zentinel/pkg/config"
	"zentinel/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	logger.Init(cfg.LogLevel)
	defer logger.Sync()
	logger.Info("Starting Zentinel API", zap.String("port", cfg.ServerPort))

	db, err := postgres.NewConnection(cfg)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	clientRepo := postgres.NewClientRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	alertRepo := postgres.NewAlertRepository(db)
	catalogueRepo := postgres.NewCatalogueRepository(db)

	aiService := ai.NewMockAnalyzer()

	// Initialize Blockchain Adapter
	var blockchainNotifier service.BlockchainNotifier
	blockchainNotifier, err = blockchain.NewPolygonAdapter()
	if err != nil {
		logger.Warn("Real Blockchain adapter failed to initialize, falling back to Mock adapter", zap.Error(err))
		blockchainNotifier = blockchain.NewMockBlockchainAdapter()
	} else {
		logger.Info("Blockchain adapter initialized successfully")
	}

	clientUseCase := usecase.NewClientUseCase(clientRepo, accountRepo, catalogueRepo)
	accountUseCase := usecase.NewAccountUseCase(accountRepo, clientRepo, catalogueRepo)
	alertUseCase := usecase.NewAlertUseCase(alertRepo, catalogueRepo, blockchainNotifier)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, accountRepo, clientRepo, alertRepo, catalogueRepo, aiService, blockchainNotifier)
	dashboardUseCase := usecase.NewDashboardUseCase(transactionRepo, alertRepo)

	healthHandler := handler.NewHealthHandler(db)
	clientHandler := handler.NewClientHandler(clientUseCase)
	accountHandler := handler.NewAccountHandler(accountUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	alertHandler := handler.NewAlertHandler(alertUseCase)
	dashboardHandler := handler.NewDashboardHandler(dashboardUseCase)

	router := zhttp.SetupRouter(
		healthHandler,
		clientHandler,
		accountHandler,
		transactionHandler,
		alertHandler,
		dashboardHandler,
	)

	srv := zhttp.NewServer(cfg, router)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", zap.Error(err))
			os.Exit(1)
		}
	}()

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
