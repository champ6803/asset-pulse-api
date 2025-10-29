package handler

import (
	"asset-pulse-api/middleware"
	dbRepo "asset-pulse-api/repositories/database"
	"asset-pulse-api/services"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase"
	"asset-pulse-api/utils/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase                 usecase.Usecase
	aiService               ai.AIService
	catalogSearchService    *ai.CatalogSearchService
	jwtManager              *jwt.JWTManager
	authMiddleware          *middleware.AuthMiddleware
	dbRepo                  dbRepo.DatabaseRepository
	softwareGroupingService *services.SoftwareGroupingService
}

type HandlerOptions struct {
	Usecase                 usecase.Usecase
	AIService               ai.AIService
	CatalogSearchService    *ai.CatalogSearchService
	JWTSecret               string
	DBRepo                  dbRepo.DatabaseRepository
	SoftwareGroupingService *services.SoftwareGroupingService
}

func New(handler *Handler) *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
	})

	software := api.Group("/software")
	{
		software.GET("/licenses", handler.GetAllLicenses)
		software.POST("/licenses", handler.InsertLicense)
		software.DELETE("/licenses/:id", handler.DeleteLicense)
		software.GET("/licenses/grouped", handler.GetGroupedSoftware)
		software.POST("/licenses/grouped/generate", handler.GenerateGroupedSoftware)
	}

	// Authentication routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/logout", handler.Logout)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(handler.authMiddleware.RequireAuth())
	{
		// User routes
		protected.GET("/users", handler.GetUsers)
		protected.GET("/me", handler.GetCurrentUser)

		// License routes
		licenses := protected.Group("/licenses")
		{
			licenses.GET("", handler.GetLicenses)
			licenses.GET("/active", handler.GetActiveLicenses)
		}

		// Request routes
		requests := protected.Group("/requests")
		{
			requests.GET("/pending", handler.GetPendingRequests)
		}

		// AI routes
		ai := protected.Group("/ai")
		{
			ai.POST("/recommendations/jd-match", handler.GenerateJDRecommendations)
			ai.POST("/recommendations/new-hire", handler.GetNewHireRecommendations)
			ai.POST("/consolidation/memo", handler.GenerateConsolidationMemo)
			ai.POST("/similarity", handler.CalculateSoftwareSimilarity)
			ai.GET("/recommendations", handler.GetAIRecommendations) // New endpoint
			ai.POST("/catalog/search", handler.SearchCatalog)
		}

		// Recommendation routes
		recommendations := protected.Group("/recommendations")
		{
			// Placeholder routes - handlers not implemented yet
			recommendations.GET("/seat-optimization", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			recommendations.GET("/cross-sub-match", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
		}

		// Template routes
		templates := protected.Group("/templates")
		{
			// Placeholder routes - handlers not implemented yet
			templates.POST("/purchase", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
		}

		// Role-based routes
		employee := protected.Group("/employee")
		employee.Use(handler.authMiddleware.RequireRole("employee", "manager", "subsidiary-cto", "group-cto"))
		{
			employee.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			employee.GET("/licenses", handler.GetUserLicenses)
		}

		manager := protected.Group("/manager")
		manager.Use(handler.authMiddleware.RequireRole("manager", "subsidiary-cto", "group-cto"))
		{
			manager.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			manager.GET("/approvals", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
		}

		cto := protected.Group("/cto")
		cto.Use(handler.authMiddleware.RequireRole("subsidiary-cto", "group-cto"))
		{
			cto.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			cto.GET("/optimization", handler.GetSeatOptimization)
		}

		groupCto := protected.Group("/group-cto")
		groupCto.Use(handler.authMiddleware.RequireRole("group-cto"))
		{
			groupCto.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
			})
			groupCto.GET("/consolidation", handler.GetConsolidationOpportunities)
			groupCto.GET("/consolidation/:id", handler.GetConsolidationOpportunityByID)
		}
	}

	return router
}

func NewHandler(options HandlerOptions) *Handler {
	jwtManager := jwt.NewJWTManager(options.JWTSecret, 24*time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	return &Handler{
		useCase:                 options.Usecase,
		aiService:               options.AIService,
		catalogSearchService:    options.CatalogSearchService,
		jwtManager:              jwtManager,
		authMiddleware:          authMiddleware,
		dbRepo:                  options.DBRepo,
		softwareGroupingService: options.SoftwareGroupingService,
	}
}
