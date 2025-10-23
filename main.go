package main

import (
	config "asset-pulse-api/configs"
	"asset-pulse-api/handler"
	dbRepo "asset-pulse-api/repositories/database"
	"asset-pulse-api/services/ai"
	"asset-pulse-api/usecase"
	mygorm "asset-pulse-api/utils/gorm"
	"asset-pulse-api/utils/logger"

	"github.com/caarlos0/env/v11"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")
	config := &config.Config{}
	if err := env.Parse(config); err != nil {
		panic(err)
	}

	logger.New(logger.Options{
		AppEnv:  config.AppEnv,
		AppName: config.AppName,
	})

	dbCredentials := mygorm.PostgresCredential{
		Host:        config.DbHost,
		Port:        config.DbPort,
		DBName:      config.DbName,
		Username:    config.DbUsername,
		Password:    config.DbPassword,
		SSLMode:     config.DbSSLMode,
		SSLRootCert: config.DbSSLRootCert,
		TimeZone:    config.DbTimeZone,
	}

	db := mygorm.NewPostgres(dbCredentials.Parse(), &gorm.Config{})

	dbRepo := dbRepo.New(db)

	// Initialize AI service (using mock for now)
	aiService := ai.NewMockAIService()

	uc := usecase.New(usecase.UsecaseOptions{
		DBRepo:    dbRepo,
		AIService: aiService,
	})

	newHandler := handler.NewHandler(handler.HandlerOptions{
		Usecase:   uc,
		AIService: aiService,
		JWTSecret: config.JWTSecret,
	})

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	apiRouter := handler.New(newHandler)

	router.Any("/api/*any", gin.WrapH(apiRouter))
	err := router.Run(":" + config.AppPort)
	if err != nil {
		panic(err)
	}
}
