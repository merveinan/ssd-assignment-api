package main

import (
	"log"
	"ssd-assignment-api/handlers"
	"ssd-assignment-api/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files" // The correct package name here is "files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ssd-assignment-api/docs" // Necessary for Swagger documentation
)

// @title SSD Assignment API
// @version 1.0
// @description A Go-based API for managing configurations with JWT authentication
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Initialize the in-memory ConfigService
	// Specifying the file path

	configService, err := services.NewConfigService("config_files")
	if err != nil {
		log.Fatal("Error loading YAML: ", err)
	}
	specificService, err := services.NewSpecificConfigService("specific_configs")
	if err != nil {
		log.Fatal("Specific config service error: ", err)
	}

	// Set up the Gin router
	r := gin.Default()
	// Custom CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // frontend adresi
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Apply the CORS middleware
	r.Use(cors.New(corsConfig))
	// Apply the CORS middleware to all routes

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
	}

	// Configuration Routes
	configRoutes := r.Group("/api/configuration")
	configRoutes.Use(services.TokenAuthMiddleware())
	{
		configRoutes.GET("/all", handlers.GetAllConfigs(configService))
		configRoutes.GET("/:id", handlers.GetConfigByID(configService))
		configRoutes.POST("/", handlers.AddConfig(configService))
		configRoutes.PUT("/:id", handlers.UpdateConfig(configService))
		configRoutes.DELETE("/:id", handlers.DeleteConfig(configService))
	}

	// Specific Configuration Routes
	specificRoutes := r.Group("/api/specific")
	specificRoutes.Use(services.TokenAuthMiddleware())
	{
		specificRoutes.GET("/", handlers.GetSpecificConfigs(specificService))
		specificRoutes.GET("/all", handlers.GetAllSpecificConfigs(specificService))
		specificRoutes.GET("/:id", handlers.GetSpecificConfigByID(specificService))
		specificRoutes.POST("/", handlers.AddSpecificConfig(specificService))
		specificRoutes.PUT("/:id", handlers.UpdateSpecificConfig(specificService))
		specificRoutes.DELETE("/:id", handlers.DeleteSpecificConfig(specificService))
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	r.Run(":8000")
}
