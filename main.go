package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.Static("/assets", "./assets")

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Workout{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Exercise{})
	db.AutoMigrate(&Seance{})
	db.Create(&Workout{Name: "Squatos"})

	protected := r.Group("/protected")
	protected.Use(AuthMiddleware(db))

	protected.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base.html", gin.H{})
	})

	protected.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	protected.GET("/workouts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "manage_workouts.html", gin.H{})
	})

	protected.GET("/sessions", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user User
		db.Where("session_id = ?", sessionID).First(&user)
		var sessions []Seance
		db.Model(&Seance{}).Preload("Exercises").Where("owner_id = ?", user.ID).Find(&sessions)
		c.HTML(http.StatusOK, "manage_sessions.html", gin.H{"sessions": sessions})
	})

	protected.GET("/creation_session", func(c *gin.Context) {
		sessionID, _ := c.Cookie("session_id")
		var user User
		db.Where("session_id = ?", sessionID).First(&user)
		var exercises []Exercise
		db.Where("owner_id = ?", user.ID).Find(&exercises)
		c.HTML(http.StatusOK, "creation_session.html", gin.H{"exercises": exercises})
	})

	protected.POST("/sessions", func(c *gin.Context) {
		var request struct {
			Name               string `form:"name" json:"name" binding:"required"`
			Selected_exercises string `form:"selected_exercises" json:"selected_exercises" binding:"required"`
		}

		// Bind the request to the struct
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing name"})
			return
		}

		strArr := strings.Split(request.Selected_exercises, ",")
		var exercisesId []int
		for _, str := range strArr {
			num, err := strconv.Atoi(str) // Convert string to int
			if err != nil {
				fmt.Println("Error converting:", str, err)
				continue
			}
			exercisesId = append(exercisesId, num) // Add the integer to the slice
		}

		sessionID, _ := c.Cookie("session_id")
		var user User
		db.Where("session_id = ?", sessionID).First(&user)

		var exercises []Exercise
		for _, id := range exercisesId {
			var exercise Exercise
			db.First(&exercise, id)
			exercises = append(exercises, exercise)
		}

		// Create the exercise for the user
		var idk = db.Create(&Seance{Name: request.Name, Owner: user, Exercises: exercises})
		if idk.Error != nil {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": false, "message": idk.Error})
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Session created"})
		}
	})

	protected.GET("/exercises", func(c *gin.Context) {
		// Get the connected user
		sessionID, _ := c.Cookie("session_id")
		var user User
		db.Where("session_id = ?", sessionID).First(&user)

		// Get the sessions of the user
		var exercises []Exercise
		db.Where("owner_id = ?", user.ID).Find(&exercises)
		c.HTML(http.StatusOK, "manage_exercises.html", gin.H{"exercises": exercises})
	})

	protected.GET("/creation_exercise", func(c *gin.Context) {
		c.HTML(http.StatusOK, "creation_exercise.html", gin.H{})
	})

	protected.POST("/exercises", func(c *gin.Context) {
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
		var user User
		db.Where("session_id = ?", sessionID).First(&user)

		// Create the exercise for the user
		var idk = db.Create(&Exercise{Name: request.Name, Owner: user})
		if idk.Error != nil {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": false, "message": idk.Error})
		} else {
			c.HTML(http.StatusOK, "result.html", gin.H{"success": true, "message": "Exercise created"})
		}
	})

	protected.GET("/stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "stats.html", gin.H{})
	})

	r.GET("/", func(c *gin.Context) {
		var workouts []Workout
		var users []User
		db.Find(&workouts)
		db.Find(&users)
		c.HTML(http.StatusOK, "index.html", gin.H{"workouts": workouts, "users": users})
	})

	r.GET("/loginregister", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	// POST request for Registration
	r.POST("/register", func(c *gin.Context) {
		var request struct {
			Username string `form:"username" json:"username" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}

		// Bind the request to the struct
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or password"})
			return
		}

		// Check if the username already exists
		var existingUser User
		if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		// Hash the password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

		// Create the cookie
		user := User{Username: request.Username, Password: string(hashedPassword)}
		sessionID, err := GenerateSessionID(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session ID"})
			return
		}
		user.Session_id = sessionID
		db.Save(&user)
		c.SetCookie("session_id", sessionID, 3600, "/", "localhost", false, true)

		// Create a new user
		db.Create(&user)

		c.HTML(http.StatusOK, "base.html", gin.H{})
	})

	// POST request for Login
	r.POST("/login", func(c *gin.Context) {
		var request struct {
			Username string `form:"username" json:"username" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}

		// Bind the request to the struct
		if err := c.Bind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or password"})
			return
		}

		// Find the user in the database
		var user User
		if err := db.Where("username = ?", request.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Compare the hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// If successful, return a success message (or start a session)
		sessionID, err := GenerateSessionID(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session ID"})
			return
		}
		user.Session_id = sessionID
		db.Save(&user)
		c.SetCookie("session_id", sessionID, 3600, "/", "localhost", false, true)

		c.HTML(http.StatusOK, "base.html", gin.H{})
	})

	r.Run(":8080")
}
