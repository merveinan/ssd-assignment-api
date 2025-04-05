package handlers

import (
	"net/http"
	"ssd-assignment-api/models"
	"ssd-assignment-api/services"

	"github.com/gin-gonic/gin"
)

// GetAllSpecificConfigs godoc
// @Summary Get all specific configurations
// @Description Retrieves all specific configurations
// @Tags specific
// @Produce json
// @Success 200 {array} models.SpecificConfig
// @Failure 500 {object} models.ErrorResponse
// @Router /api/specific/all [get]
func GetAllSpecificConfigs(service *services.SpecificConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		configs, err := service.GetAllSpecificConfigs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, configs)
	}
}

// GetSpecificConfigByID godoc
// @Summary Get specific configuration by ID
// @Description Retrieves a specific configuration by its ID
// @Tags specific
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 200 {object} models.SpecificConfig
// @Failure 404 {object} models.ErrorResponse
// @Router /api/specific/{id} [get]
func GetSpecificConfigByID(service *services.SpecificConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		config, err := service.GetSpecificConfigByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Specific config not found"})
			return
		}
		c.JSON(http.StatusOK, config)
	}
}

// UpdateSpecificConfig godoc
// @Summary Update specific configuration
// @Description Updates an existing specific configuration
// @Tags specific
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param config body models.SpecificConfig true "Updated Configuration"
// @Success 200 {object} models.SpecificConfig
// @Failure 400 {object} models.ErrorResponse
// @Router /api/specific/{id} [put]
func UpdateSpecificConfig(service *services.SpecificConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var config models.SpecificConfig
		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		if err := service.UpdateSpecificConfig(id, config); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, config)
	}
}

// DeleteSpecificConfig godoc
// @Summary Delete specific configuration
// @Description Deletes a specific configuration by ID
// @Tags specific
// @Param id path string true "Configuration ID"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/specific/{id} [delete]
func DeleteSpecificConfig(service *services.SpecificConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeleteSpecificConfig(id); err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Specific config not found"})
			return
		}
		c.JSON(http.StatusOK, models.MessageResponse{Message: "Specific config deleted"})
	}
}

// GetSpecificConfigs godoc
// @Summary Get matching configurations
// @Description Get configuration IDs based on host, url or page
// @Tags specific
// @Produce json
// @Param host query string false "Target host"
// @Param url query string false "Target URL path"
// @Param page query string false "Target page name"
// @Success 200 {object} map[string][]string
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/specific [get]
func GetSpecificConfigs(service *services.SpecificConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Query("host")
		url := c.Query("url")
		page := c.Query("page")

		// Validate at least one parameter is provided
		if host == "" && url == "" && page == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "At least one parameter (host, url or page) is required",
			})
			return
		}

		configIDs, err := service.GetMatchingConfigs(host, url, page)
		if err != nil {
			if err.Error() == "no matching configurations found" {
				c.JSON(http.StatusNotFound, models.ErrorResponse{
					Error: "No configurations found for the specified criteria",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to retrieve configurations: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"config_ids": configIDs})
	}
}

// AddSpecificConfig godoc
// @Summary Add new specific configuration
// @Description Add a new specific configuration mapping
// @Tags specific
// @Accept json
// @Produce json
// @Param config body models.SpecificConfig true "Specific Configuration"
// @Success 201 {object} models.SpecificConfig
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/specific [post]
func AddSpecificConfig(service *services.SpecificConfigService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var config models.SpecificConfig

		// Bind and validate request body
		if err := c.ShouldBindJSON(&config); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid request format: " + err.Error(),
			})
			return
		}

		// Validate required fields
		if config.ID == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Configuration ID is required",
			})
			return
		}

		if len(config.DataSource.Pages) == 0 &&
			len(config.DataSource.URLs) == 0 &&
			len(config.DataSource.Hosts) == 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "At least one datasource mapping is required",
			})
			return
		}

		// Add configuration
		if err := service.AddSpecificConfig(config); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to add configuration: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, config)
	}
}
