package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"gin-app/misc"
)

func RealSessions(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/sessions/:id/start", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var realSeance misc.RealSeance
		db.Preload("Template").Preload("Template.Exercises").First(&realSeance, id)

		realSeance.Active = true

		for _, templateExercise := range realSeance.Template.Exercises {
			var realExercise misc.RealExercise
			realExercise.Template = templateExercise
			realExercise.Owner = user
			realExercise.CorrespondingSeance = realSeance // TODO vrm besoin de cette ligne et en dessous ?
			realExercise.CorrespondingSeanceID = int(realSeance.ID)
			db.Save(&realExercise)
		}

		var realExercises []misc.RealExercise
		db.Preload("Template").Where("corresponding_seance_id = ?", realSeance.ID).Find(&realExercises)

		println(len(realExercises))

		c.HTML(http.StatusOK, "choose_exercise_to_start.html", gin.H{"exercises": realExercises})
	})
}
