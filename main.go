package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gin-app/handlers"
	"gin-app/middlewares"
	"gin-app/misc"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var prod = false

func main() {
	_ = godotenv.Load()
	tursoToken := os.Getenv("TOKEN")
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.html")

	r.Static("/assets", "./assets")

	var db *gorm.DB
	var err error

	if prod {
		url := "libsql://hongym-pismice.turso.io?authToken=" + tursoToken
		db, err = gorm.Open(sqlite.New(sqlite.Config{
			DriverName: "libsql",
			DSN:        url,
		}), &gorm.Config{})

	} else {
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		return
	}

	db.AutoMigrate(&misc.Workout{})
	db.AutoMigrate(&misc.User{})
	db.AutoMigrate(&misc.Exercise{})
	db.AutoMigrate(&misc.Seance{})
	db.AutoMigrate(&misc.RealWorkout{}, &misc.RealSeance{}, &misc.RealExercise{}, &misc.RealSet{})

	protected := r.Group("/protected")
	protected.Use(middlewares.AuthMiddleware(db))

	protected.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{})
	})

	protected.GET("/home", func(c *gin.Context) {
		// Iterate over all the Workouts and return the first one that is active for the connected user
		var user misc.User
		sessionID, _ := c.Cookie("session_id")
		db.Where("session_id = ?", sessionID).First(&user)
		var activeRealWorkout misc.RealWorkout
		res := db.Preload("Template").Where("owner_id = ? AND active = ?", user.ID, true).First(&activeRealWorkout)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				sessionID, _ := c.Cookie("session_id")
				var user misc.User
				db.Where("session_id = ?", sessionID).First(&user)
				var workouts []misc.Workout
				db.Model(&misc.Workout{}).Preload("Seances").Where("owner_id = ?", user.ID).Find(&workouts)
				c.HTML(http.StatusOK, "choose_workout_to_activate.html", gin.H{"workouts": workouts})
			} else {
				c.JSON(200, gin.H{"message": "Unexpected error"})
			}
		} else {
			// Get the real seances for the active real workout
			var realSeances []misc.RealSeance
			db.Preload("Template").Preload("Template.Exercises").Where("corresponding_workout_id = ?", activeRealWorkout.ID).Find(&realSeances)
			println(len(realSeances))
			c.HTML(http.StatusOK, "choose_session_to_start.html", gin.H{"workout": activeRealWorkout, "sessions": realSeances})
		}
	})

	handlers.Sessions(protected, db)
	handlers.RealSessions(protected, db)
	handlers.Exercises(protected, db)
	handlers.RealExercises(protected, db)
	handlers.Workouts(protected, db)
	handlers.RealWorkouts(protected, db)
	handlers.Auth(&r.RouterGroup, db)

	protected.GET("/stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "stats.html", gin.H{})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	fmt.Println("APP SHOULD HAVE STARTED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	r.Run(":8080")
	//err = r.RunTLS(":8080", "cert.pem", "key.pem")
}
