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

func Workouts(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/workouts", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)
		var workouts []misc.Workout
		db.Model(&misc.Workout{}).Preload("Seances").Where("owner_id = ?", user.ID).Find(&workouts)
		c.HTML(http.StatusOK, "manage_workouts.html", gin.H{"workouts": workouts})
	})

	r.GET("/workouts/:id", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var workout misc.Workout
		db.Preload("Seances").First(&workout, id)

		if workout.OwnerID != int(user.ID) {
			c.Abort()
		}

		var sessions []misc.Seance
		db.Where("owner_id = ?", user.ID).Find(&sessions)

		var sessions_selected []misc.Seance
		var sessions_not_selected []misc.Seance

		for _, session := range sessions {
			found := false
			for _, workoutsSessions := range workout.Seances {
				if session.ID == workoutsSessions.ID {
					found = true
					break
				}
			}
			if found {
				sessions_selected = append(sessions_selected, session)
			} else {
				sessions_not_selected = append(sessions_not_selected, session)
			}
		}

		c.HTML(http.StatusOK, "modify_workout.html", gin.H{"workout": workout, "sessions_selected": sessions_selected, "sessions_not_selected": sessions_not_selected})
	})

	r.GET("/creation_workout", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)
		var sessions []misc.Seance
		db.Where("owner_id = ?", user.ID).Find(&sessions)
		c.HTML(http.StatusOK, "creation_workout.html", gin.H{"sessions": sessions})
	})

	r.PATCH("/workouts/:id", func(c *gin.Context) {
		var request struct {
			Name              string `form:"name" json:"name" binding:"required"`
			Selected_sessions string `form:"selected-sessions-input" json:"selected-sessions-input"`
		}

		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parsing"})
			return
		}
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		id := c.Param("id")
		var workout misc.Workout
		db.First(&workout, id)

		if workout.OwnerID != int(user.ID) {
			c.Abort()
		}

		strArr := strings.Split(request.Selected_sessions, ",")
		var sessionsId []int
		for _, str := range strArr {
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error converting:", str, err)
				continue
			}
			sessionsId = append(sessionsId, num)
		}

		var sessions []misc.Seance
		for _, id := range sessionsId {
			var session misc.Seance
			db.First(&session, id)
			sessions = append(sessions, session)
		}

		workout.Name = request.Name
		workout.Seances = sessions
		db.Save(&workout)
		db.Model(&workout).Association("Seances").Replace(sessions)
		c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Workout modified"})
	})

	r.POST("/workouts", func(c *gin.Context) {
		var request struct {
			Name              string `form:"name" json:"name" binding:"required"`
			Selected_sessions string `form:"selected_sessions" json:"selected_sessions"`
		}

		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind the request"})
			return
		}
		sessionID, _ := c.Cookie("session_id")
		var user misc.User
		db.Where("session_id = ?", sessionID).First(&user)

		strArr := strings.Split(request.Selected_sessions, ",")
		var sessionsId []int
		for _, str := range strArr {
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error converting:", str, err)
				continue
			}
			sessionsId = append(sessionsId, num)
		}

		var sessions []misc.Seance
		for _, id := range sessionsId {
			var session misc.Seance
			db.First(&session, id)
			sessions = append(sessions, session)
		}

		var idk = db.Create(&misc.Workout{Name: request.Name, Owner: user, Seances: sessions})
		if idk.Error != nil {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": false, "message": idk.Error})
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Workout modified"})
		}
	})
}
