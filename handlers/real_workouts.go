package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"gin-app/misc"
)

func RealWorkouts(r *gin.RouterGroup, db *gorm.DB) {
	// J ai du mettre GET parce que sinon 404 sur "PATCH" /home :(
	r.GET("/workouts/:id/finish", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var realWorkout misc.RealWorkout
		db.Preload("Template").First(&realWorkout, id)

		if realWorkout.OwnerID != int(user.ID) {
			c.Abort()
		}

		realWorkout.Finished = true
		realWorkout.Active = false
		db.Save(&realWorkout)

		c.Redirect(http.StatusFound, "/protected/home")
	})

	r.POST("/workouts/:id/activate", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var template misc.Workout
		db.Preload("Seances").First(&template, id)

		if template.OwnerID != int(user.ID) {
			c.Abort()
		}

		var realWorkout misc.RealWorkout
		realWorkout.Active = true
		realWorkout.Owner = user
		realWorkout.CurrentWeek = 1
		realWorkout.Template = template

		db.Save(&realWorkout)

		for _, templateSeance := range template.Seances {
			var realSeance misc.RealSeance
			realSeance.Owner = user
			realSeance.Template = templateSeance
			realSeance.Week = realWorkout.CurrentWeek
			realSeance.CorrespondingWorkout = realWorkout
			realSeance.CorrespondingWorkoutID = int(realWorkout.ID)
			realSeance.Active = false
			db.Save(&realSeance)
		}

		c.Redirect(http.StatusFound, "/protected/home")
	})
}
