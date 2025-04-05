package handlers

import (
	"fmt"
	"log"
	"net/http"
	"ssd-assignment-api/models" // models paketini import et
	"ssd-assignment-api/services"

	"github.com/gin-gonic/gin"
)

// GetAllConfigs godoc
// @Summary Get all configurations
// @Description Retrieves all configurations stored in memory
// @Tags configuration
// @Accept json
// @Produce json
// @Success 200 {array} models.Config "List of configurations"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /api/configuration/all [get]
func GetAllConfigs(service *services.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		configs, err := service.GetAllConfigs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, configs)
	}
}

// GetConfigByID godoc
// @Summary Get configuration by ID
// @Description Retrieves a specific configuration by its ID
// @Tags configuration
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 200 {object} models.Config
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/configuration/{id} [get]
func GetConfigByID(service *services.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		config, err := service.GetConfigByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Config not found"})
			return
		}
		c.JSON(http.StatusOK, config)
	}
}

// AddConfig godoc
// @Summary Add a new configuration
// @Description Adds a new configuration to the system
// @Tags configuration
// @Accept json
// @Produce json
// @Param config body models.Config true "Configuration"
// @Success 201 {object} models.Config
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/configuration [post]
func AddConfig(service *services.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Received POST request to /api/configuration") // Added log statement

		var config models.Config
		if err := c.ShouldBindJSON(&config); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("JSON binding error: %s", err.Error())})
			return
		}

		log.Println("Successfully bound JSON:", config) // Log the received config

		// Attempt to add the configuration
		if err := service.AddConfig(config); err != nil {
			log.Println("Error adding config:", err)
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("Failed to add config: %s", err.Error())})
			return
		}

		log.Println("Config added successfully:", config)
		c.JSON(http.StatusCreated, config)
	}
}

// UpdateConfig godoc
// @Summary Update an existing configuration
// @Description Updates an existing configuration by ID
// @Tags configuration
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param config body models.Config true "Updated configuration"
// @Success 200 {object} models.Config
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/configuration/{id} [put]
func UpdateConfig(service *services.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var config models.Config
		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}
		if err := service.UpdateConfig(id, config); err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Config not found"})
			return
		}
		c.JSON(http.StatusOK, config)
	}
}

// DeleteConfig godoc
// @Summary Delete a configuration
// @Description Deletes a specific configuration by ID
// @Tags configuration
// @Param id path string true "Configuration ID"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/configuration/{id} [delete]
func DeleteConfig(service *services.ConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		// Veritabanında id'nin var olup olmadığını kontrol et
		if err := service.DeleteConfig(id); err != nil {
			// Eğer ID bulunamazsa, 404 döndür
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Config not found"})
			return
		}
		// Silme başarılı ise
		c.JSON(http.StatusOK, models.MessageResponse{Message: "Config deleted"})
	}
}
