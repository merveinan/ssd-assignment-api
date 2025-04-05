package handlers

import (
	"net/http"
	"ssd-assignment-api/services"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description This endpoint registers a new user with username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User info"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Bad request"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save user to database (This part should ideally include password hashing)
	// For now, we're skipping user storage.

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Log in an existing user
// @Description This endpoint allows an existing user to log in using their username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User login info"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad request"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials (This part should ideally include password comparison)
	// For now, assuming credentials are correct.

	token, err := services.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
