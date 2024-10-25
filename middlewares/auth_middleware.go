package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"gin-app/misc"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized because no session_id cookie"})
			fmt.Println("Unauthorized because no session_id cookie")
			//c.Abort()
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		var user misc.User
		if err := db.Where("session_id = ?", sessionID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized because no corresponding for the session_id = " + sessionID})
			fmt.Println("Unauthorized because no corresponding for the session_id = " + sessionID)
			//c.Abort()
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		c.Next()
	}
}
