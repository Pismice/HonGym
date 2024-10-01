package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"gin-app/misc"
)

func Exercises(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/exercises", func(c *gin.Context) {
		// Get the connected user
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		// Get the sessions of the user
		var exercises []misc.Exercise
		db.Where("owner_id = ?", user.ID).Find(&exercises)
		c.HTML(http.StatusOK, "manage_exercises.html", gin.H{"exercises": exercises})
	})

	r.GET("/creation_exercise", func(c *gin.Context) {
		c.HTML(http.StatusOK, "creation_exercise.html", gin.H{})
	})

	r.POST("/exercises", func(c *gin.Context) {
		var request struct {
			Name string `form:"name" json:"name" binding:"required"`
		}

		// Bind the request to the struct
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing name"})
			return
		}

		// Get the session id of the connected user
		sessionID, _ := c.Cookie("session_id")

		// Get the corresponding user
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		// Create the exercise for the user
		var idk = db.Create(&misc.Exercise{Name: request.Name, Owner: user})
		if idk.Error != nil {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": false, "message": idk.Error})
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Exercise created"})
		}
	})
}
