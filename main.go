package main

import (
	"log"
	"ssd-assignment-api/handlers"
	"ssd-assignment-api/services"

	"github.com/gin-contrib/cors" // Adding the CORS package
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // The correct package name here is "files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ssd-assignment-api/docs" // Necessary for Swagger documentation
)

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
	r.Use(cors.Default()) // For testing purposes, it will accept all incoming requests

	// Configuration Routes
	configRoutes := r.Group("/api/configuration")
	{
		configRoutes.GET("/all", handlers.GetAllConfigs(configService))
		configRoutes.GET("/:id", handlers.GetConfigByID(configService))
		configRoutes.POST("/", handlers.AddConfig(configService))
		configRoutes.PUT("/:id", handlers.UpdateConfig(configService))
		configRoutes.DELETE("/:id", handlers.DeleteConfig(configService))
	}

	// Specific Configuration Routes
	specificRoutes := r.Group("/api/specific")
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
