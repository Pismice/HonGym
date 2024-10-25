package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"gin-app/misc"
)

func RealSessions(r *gin.RouterGroup, db *gorm.DB) {
	// J ai du mettre GET parce que sinon 404 sur "PATCH" /home :(
	r.GET("/sessions/:id/finish", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var realSeance misc.RealSeance
		db.Preload("Template").First(&realSeance, id)

		if realSeance.OwnerID != int(user.ID) {
			c.Abort()
		}

		realSeance.Finished = true
		realSeance.Active = false
		db.Save(&realSeance)

		c.Redirect(http.StatusFound, "/protected/home")
	})

	r.GET("/sessions/:id/display", func(c *gin.Context) {
		id := c.Param("id")
		var realSeance misc.RealSeance
		db.Preload("Template").Preload("Template.Exercises").First(&realSeance, id)

		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		if realSeance.OwnerID != int(user.ID) {
			c.Abort()
		}

		var realExercises []misc.RealExercise
		db.Preload("Template").Where("corresponding_seance_id = ?", realSeance.ID).Find(&realExercises)

		c.HTML(http.StatusOK, "choose_exercise_to_start.html", gin.H{"exercises": realExercises, "realSeance": realSeance})
	})

	r.POST("/sessions/:id/start", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var realSeance misc.RealSeance
		db.Preload("Template").Preload("Template.Exercises").First(&realSeance, id)

		if realSeance.OwnerID != int(user.ID) {
			c.Abort()
		}

		if realSeance.Active != true {

			for _, templateExercise := range realSeance.Template.Exercises {
				var realExercise misc.RealExercise
				realExercise.Template = templateExercise
				realExercise.Owner = user
				realExercise.CorrespondingSeance = realSeance // TODO vrm besoin de cette ligne et en dessous ?
				realExercise.CorrespondingSeanceID = int(realSeance.ID)
				db.Save(&realExercise)
			}
		}

		realSeance.Active = true
		realSeance.Finished = false
		db.Save(&realSeance)

		var realExercises []misc.RealExercise
		db.Preload("Template").Where("corresponding_seance_id = ?", realSeance.ID).Find(&realExercises)

		c.HTML(http.StatusOK, "choose_exercise_to_start.html", gin.H{"exercises": realExercises, "realSeance": realSeance})
	})
}
