package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"gin-app/misc"
)

func RealExercises(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/exercises/:id/start", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var realExercise misc.RealExercise
		db.Preload("Template").First(&realExercise, id)

		c.HTML(http.StatusOK, "doing_exercise.html", gin.H{"realExercise": realExercise})
	})
}
