package main

import (
	//"fmt"
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io"
	"net/http"
)

type Workout struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Username   string
	Password   string
	Session_id string
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var user User
		if err := db.Where("session_id = ?", sessionID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Workout{})
	db.AutoMigrate(&User{})
	db.Create(&Workout{Name: "Squatos"})

	protected := r.Group("/protected")
	protected.Use(AuthMiddleware(db))

	protected.GET("/", func(c *gin.Context) {
		tmpl, _ := template.New("").ParseFiles("templates/base.html", "templates/home.html")
		_ = tmpl.ExecuteTemplate(c.Writer, "base", nil)

	})

	protected.GET("/workouts", func(c *gin.Context) {
		tmpl, _ := template.New("").ParseFiles("templates/base.html", "templates/manage_workouts.html")
		_ = tmpl.ExecuteTemplate(c.Writer, "base", nil)
	})

	protected.GET("/sessions", func(c *gin.Context) {
		tmpl, _ := template.New("").ParseFiles("templates/base.html", "templates/manage_sessions.html")
		_ = tmpl.ExecuteTemplate(c.Writer, "base", nil)
	})

	protected.GET("/stats", func(c *gin.Context) {
		tmpl, _ := template.New("").ParseFiles("templates/base.html", "templates/stats.html")
		_ = tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

		tmpl, _ := template.New("").ParseFiles("templates/base.html", "templates/home.html")
		_ = tmpl.ExecuteTemplate(c.Writer, "base", nil)
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

		tmpl, _ := template.New("").ParseFiles("templates/base.html", "templates/home.html")
		_ = tmpl.ExecuteTemplate(c.Writer, "base", nil)

	})

	r.Run(":8080")
}

func GenerateSessionID(length int) (string, error) {
	// Create a byte slice with the given length
	bytes := make([]byte, length)

	// Fill the slice with random bytes
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	// Convert the byte slice to a hexadecimal string
	return hex.EncodeToString(bytes), nil
}
