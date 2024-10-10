package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gin-app/handlers"
	"gin-app/middlewares"
	"gin-app/misc"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Static("/assets", "./assets")

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&misc.Workout{})
	db.AutoMigrate(&misc.User{})
	db.AutoMigrate(&misc.Exercise{})
	db.AutoMigrate(&misc.Seance{})
	db.AutoMigrate(&misc.RealWorkout{}, &misc.RealSeance{}, &misc.RealExercise{}, &misc.RealSet{})
	db.Create(&misc.Workout{Name: "Squatos"})

	protected := r.Group("/protected")
	protected.Use(middlewares.AuthMiddleware(db))

	protected.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{})
	})

	protected.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	handlers.Sessions(protected, db)
	handlers.Exercises(protected, db)
	handlers.Workouts(protected, db)
	handlers.Auth(&r.RouterGroup, db)

	protected.GET("/stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "stats.html", gin.H{})
	})

	r.GET("/", func(c *gin.Context) {
		var workouts []misc.Workout
		var users []misc.User
		db.Find(&workouts)
		db.Find(&users)
		c.HTML(http.StatusOK, "index.html", gin.H{"workouts": workouts, "users": users})
	})

	r.Run(":8080")
}
