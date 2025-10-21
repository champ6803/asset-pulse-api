package handler

import (
	"asset-pulse-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase usecase.Usecase
}

type HandlerOptions struct {
	Usecase usecase.Usecase
}

func New(handler *Handler) *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// User routes
	api.GET("/users", handler.GetUsers)

	return router
}

func NewHandler(options HandlerOptions) *Handler {
	return &Handler{
		useCase: options.Usecase,
	}
}

