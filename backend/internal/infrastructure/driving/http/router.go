package http

import (
	"zentinel/internal/infrastructure/driving/http/handler"
	"zentinel/internal/infrastructure/driving/http/middleware"

	_ "zentinel/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	jwtSecret string,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	clientHandler *handler.ClientHandler,
	accountHandler *handler.AccountHandler,
	transactionHandler *handler.TransactionHandler,
	alertHandler *handler.AlertHandler,
	dashboardHandler *handler.DashboardHandler,
) *gin.Engine {
	r := gin.New()

	// Middlewares
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.SecurityMiddleware())
	r.Use(middleware.RateLimitMiddleware(10, 20))

	// Health Check
	r.GET("/health", healthHandler.Check)

	// Metrics (Placeholder)
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "metrics endpoint"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth
	auth := r.Group("/auth")
	{
		// auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// API v1
	v1 := r.Group("/api/v1")
	v1.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		// Clients
		clients := v1.Group("/clients")
		{
			clients.POST("", clientHandler.Create)
			clients.GET("", clientHandler.List)
			clients.GET("/:id", clientHandler.Get)
			clients.PUT("/:id", clientHandler.Update)
			clients.GET("/:id/accounts", clientHandler.GetAccounts)
		}

		// Accounts
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", accountHandler.Create)
			accounts.GET("", accountHandler.List)
			accounts.GET("/:id", accountHandler.Get)
			accounts.PUT("/:id", accountHandler.Update)
		}

		// Transactions
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", transactionHandler.Create)
			transactions.GET("", transactionHandler.List)
			transactions.GET("/:id", transactionHandler.Get)
			transactions.POST("/:id/analyze", transactionHandler.Analyze)
		}

		// Alerts
		alerts := v1.Group("/alerts")
		{
			alerts.POST("", alertHandler.Create)
			alerts.GET("", alertHandler.List)
			alerts.GET("/:id", alertHandler.Get)
			alerts.PATCH("/:id/status", alertHandler.UpdateStatus)
		}

		// Dashboard
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/stats", dashboardHandler.GetStats)
		}
	}

	return r
}
