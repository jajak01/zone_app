package main

import (
	"net/http"
	"zone-app/database" // Ensure this matches your go.mod name
	"zone-app/models"

	"github.com/gin-gonic/gin"
)


type RegisterRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ProfileRequest struct {
	UserID       uint    `json:"user_id" binding:"required"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	BirthDate    string  `json:"birth_date"`
	BaseLocation string  `json:"base_location"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
}

type StartActivityRequest struct {
	UserID       uint   `json:"user_id" binding:"required"`
	TypeActivity string `json:"type_activity" binding:"required"` // e.g., "Running", "Cycling"
}

type LogDetailRequest struct {
	IDActivity    uint   `json:"id_activity" binding:"required"`
	UserID        uint   `json:"user_id" binding:"required"`
	LogGisSpatial string `json:"log_gis_spatial" binding:"required"` // The GPS string or JSON
}

func main() {
	// 1. Initialize DB Connection
	database.Connect()

	r := gin.Default()

	// 2. Registration Route
	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest

		// Validate JSON input
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// --- DATABASE TRANSACTION START ---
		// We use a transaction to ensure both tables are updated or none at all
		tx := database.DB.Begin()

		// Step A: Create base User
		newUser := models.User{
			UserName: req.UserName,
		}
		if err := tx.Create(&newUser).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create User profile"})
			return
		}

		// Step B: Create User Auth using the ID from Step A
		newAuth := models.UserAuth{
			UserID: newUser.UserID, // GORM automatically populates this after Create
			Email:  req.Email,
			Pass:   req.Password,
		}
		if err := tx.Create(&newAuth).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email might already exist"})
			return
		}

		tx.Commit()
		// --- DATABASE TRANSACTION END ---

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"user_id": newUser.UserID,
			"message": "Registration complete!",
		})
	})

	r.POST("/login", func(c *gin.Context) {
    	var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var auth models.UserAuth
		// Search for the user by email
		result := database.DB.Where("\"Email\" = ?", req.Email).First(&auth)

		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Compare passwords (Plain text for now as we agreed)
		if auth.Pass != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login successful",
			"user_id": auth.UserID,
			})
	})

	r.POST("/profile/update", func(c *gin.Context) {
		var req ProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Prepare the record
		info := models.UserInfo{
			UserID:       req.UserID,
			Email:        req.Email,
			Phone:        req.Phone,
			BirthDate:    req.BirthDate,
			BaseLocation: req.BaseLocation,
			Height:       req.Height,
			Weight:       req.Weight,
		}

		// "Save" in GORM performs an Upsert: 
		// It updates if the User_ID exists, or creates if it doesn't.
		if err := database.DB.Save(&info).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Profile updated successfully",
			})
	})

		// STEP 1: Start the session
	r.POST("/activity/start", func(c *gin.Context) {
		var req StartActivityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		activity := models.UserActivity{
			UserID:       req.UserID,
			TypeActivity: req.TypeActivity,
		}

		if err := database.DB.Create(&activity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start activity"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"id_activity": activity.IDActivity,
		})
	})

	// STEP 2: Log GPS coordinates (this would be called frequently by Flutter)
	r.POST("/activity/log", func(c *gin.Context) {
		var req LogDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		detail := models.ActivityDetail{
			IDActivity:    req.IDActivity,
			UserID:        req.UserID,
			LogGisSpatial: req.LogGisSpatial,
		}

		if err := database.DB.Create(&detail).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "logged"})
	})
	// Start server on port 8080
	r.Run(":8080")
}