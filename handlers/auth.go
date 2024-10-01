package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gin-app/misc"
)

func Auth(r *gin.RouterGroup, db *gorm.DB) {
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
		var existingUser misc.User
		if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		// Hash the password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

		// Create the cookie
		user := misc.User{Username: request.Username, Password: string(hashedPassword)}
		sessionID, err := misc.GenerateSessionID(32)
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
		var user misc.User
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
		sessionID, err := misc.GenerateSessionID(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session ID"})
			return
		}
		user.Session_id = sessionID
		db.Save(&user)
		c.SetCookie("session_id", sessionID, 3600, "/", "localhost", false, true)

		c.HTML(http.StatusOK, "base.html", gin.H{})
	})

}
