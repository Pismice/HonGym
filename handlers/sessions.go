package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"

	"gin-app/misc"
)

func Sessions(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/sessions", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)
		var sessions []misc.Seance
		db.Model(&misc.Seance{}).Preload("Exercises").Where("owner_id = ?", user.ID).Find(&sessions)
		c.HTML(http.StatusOK, "manage_sessions.html", gin.H{"sessions": sessions})
	})

	r.GET("/sessions/:id", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var session misc.Seance
		db.Preload("Exercises").First(&session, id)

		var exercises []misc.Exercise
		db.Where("owner_id = ?", user.ID).Find(&exercises)

		var exercises_selected []misc.Exercise
		var exercises_not_selected []misc.Exercise

		for _, exercise := range exercises {
			found := false
			for _, sessionExercise := range session.Exercises {
				if exercise.ID == sessionExercise.ID {
					found = true
					break
				}
			}
			if found {
				exercises_selected = append(exercises_selected, exercise)
			} else {
				exercises_not_selected = append(exercises_not_selected, exercise)
			}
		}

		c.HTML(http.StatusOK, "modify_session.html", gin.H{"session": session, "exercises_selected": exercises_selected, "exercises_not_selected": exercises_not_selected})
	})

	r.GET("/creation_session", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)
		var exercises []misc.Exercise
		db.Where("owner_id = ?", user.ID).Find(&exercises)
		c.HTML(http.StatusOK, "creation_session.html", gin.H{"exercises": exercises})
	})

	r.PATCH("/sessions/:id", func(c *gin.Context) {
		var request struct {
			Name               string `form:"name" json:"name" binding:"required"`
			Selected_exercises string `form:"selected-exercises-input" json:"selected-exercises-input"`
		}

		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parsing"})
			return
		}
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		println(request.Selected_exercises)

		strArr := strings.Split(request.Selected_exercises, ",")
		var exercisesId []int
		for _, str := range strArr {
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error converting:", str, err)
				continue
			}
			exercisesId = append(exercisesId, num)
		}

		var exercises []misc.Exercise
		for _, id := range exercisesId {
			var exercise misc.Exercise
			db.First(&exercise, id)
			exercises = append(exercises, exercise)
			println("Found exercise:", exercise.Name)
		}

		for _, exercise := range exercises {
			println("Exercise:", exercise.Name)
		}

		id := c.Param("id")
		var session misc.Seance
		db.First(&session, id)
		session.Name = request.Name
		session.Exercises = exercises
		db.Save(&session)
		db.Model(&session).Association("Exercises").Replace(exercises)
		for _, e := range session.Exercises {
			println("Deleting exercise:", e.Name)
		}
		c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Session modified"})
	})

	// TODO not used !!!!
	r.PATCH("/sessions/:id/start", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var template misc.Seance
		db.First(&template, id)

		var realSession misc.RealSeance
		realSession.Owner = user
		realSession.Template = template

		db.Create(&realSession)
		c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "New session started"})
	})

	r.POST("/sessions", func(c *gin.Context) {
		var request struct {
			Name               string `form:"name" json:"name" binding:"required"`
			Selected_exercises string `form:"selected_exercises" json:"selected_exercises"`
		}

		// Bind the request to the struct
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing name"})
			return
		}
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		strArr := strings.Split(request.Selected_exercises, ",")
		var exercisesId []int
		for _, str := range strArr {
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error converting:", str, err)
				continue
			}
			exercisesId = append(exercisesId, num)
		}

		var exercises []misc.Exercise
		for _, id := range exercisesId {
			var exercise misc.Exercise
			db.First(&exercise, id)
			exercises = append(exercises, exercise)
		}

		// Create the exercise for the user
		var idk = db.Create(&misc.Seance{Name: request.Name, Owner: user, Exercises: exercises})
		if idk.Error != nil {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": false, "message": idk.Error})
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Session modified"})
		}
	})

}
