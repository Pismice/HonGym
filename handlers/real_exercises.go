package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"gin-app/misc"
)

func RealExercises(r *gin.RouterGroup, db *gorm.DB) {
	r.PATCH("/exercises/:id/done", func(c *gin.Context) {
		id := c.Param("id")
		var realExercise misc.RealExercise
		db.Preload("Template").First(&realExercise, id)

		reps := c.PostForm("reps")
		weight := c.PostForm("weight")

		var realSet misc.RealSet
		realSet.Reps, _ = strconv.Atoi(reps)
		realSet.Weight, _ = strconv.Atoi(weight)
		realSet.CorrespondingExercise = realExercise
		realSet.CorrespondingExerciseID = int(realExercise.ID)
		db.Save(&realSet)

		var count int64
		_ = db.Model(&misc.RealSet{}).Where("corresponding_exercise_id = ?", realExercise.ID).Count(&count).Error

		c.HTML(http.StatusOK, "doing_exercise.html", gin.H{"realExercise": realExercise, "count": count})
	})

	r.POST("/exercises/:id/start", func(c *gin.Context) {
		id := c.Param("id")
		var realExercise misc.RealExercise
		db.Preload("Template").First(&realExercise, id)

		var count int64
		_ = db.Model(&misc.RealSet{}).Where("corresponding_exercise_id = ?", realExercise.ID).Count(&count).Error

		c.HTML(http.StatusOK, "doing_exercise.html", gin.H{"realExercise": realExercise, "count": count})
	})

	r.POST("/exercises/:id/finish", func(c *gin.Context) {
		id := c.Param("id")
		var realExercise misc.RealExercise
		db.Preload("Template").First(&realExercise, id)

		realExercise.Finished = true
		db.Save(&realExercise)

		c.Redirect(http.StatusFound, "/protected/sessions/"+strconv.Itoa(realExercise.CorrespondingSeanceID)+"/display")
	})
}
