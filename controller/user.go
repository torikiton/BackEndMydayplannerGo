package controller

import (
	"go-gorm/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/api")
	{
		routes.POST("/get_user", func(c *gin.Context) {
			GetUser(c, db)
		})
		routes.POST("/create_acc", func(c *gin.Context) {
			CreateUser(c, db)
		})
		routes.GET("/get_all_user", func(c *gin.Context) {
			GetAllUsers(c, db)
		})
		routes.DELETE("/account", func(c *gin.Context) {
			DeleteUser(c, db)
		})
	}
}

func GetUser(c *gin.Context, db *gorm.DB) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email is required."})
		return
	}

	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database query failed."})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context, db *gorm.DB) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		return
	}

	if user.HashedPassword != "" {
		user.HashedPassword = "-"
	} else {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
		user.HashedPassword = string(hashedPassword)
	}

	user.CreateAt = time.Now()
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to save data."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User account created successfully.", "user_id": user.UserID})
}

func GetAllUsers(c *gin.Context, db *gorm.DB) {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database query failed."})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No users found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "users": users})
}

func DeleteUser(c *gin.Context, db *gorm.DB) {
	email := c.PostForm("email")
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
