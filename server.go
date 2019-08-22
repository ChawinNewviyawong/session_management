package main

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var sessionManager *scs.SessionManager

func main() {
	router := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "My Session Test",
		})
	})
	router.POST("/signin", func(c *gin.Context) {
		user := Login{}
		fmt.Print(user.Username)
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  http.StatusOK,
			"message": user.Username + " " + user.Password,
		})
	})
	router.GET("/login", func(c *gin.Context) {
		// user := Login{}
		// fmt.Println(user.Username)
		// if err := c.ShouldBindJSON(&user); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"status": http.StatusBadRequest,
		// 		"error":  err.Error(),
		// 	})
		// 	return
		// }
		session := sessions.Default(c)
		session.Set("sessionid", "user.Username")
		session.Save()
		ssid := session.Get("sessionid")
		c.JSON(http.StatusOK, gin.H{
			"message":   "user.Username",
			"sessionid": ssid,
		})
	})
	router.GET("/show", func(c *gin.Context) {
		session := sessions.Default(c)
		ssid := session.Get("sessionid")
		if ssid == "user.Username" {
			c.JSON(http.StatusOK, gin.H{
				"sessionid": ssid,
				"message":   "my session",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
		}
	})
	router.Run(":3000")
}

// func createSession(c *gin.Context, username string) {
// 	// user := Login{}
// 	// if err := c.ShouldBindJSON(&user); err != nil {
// 	// 	return err.Error()
// 	// }
// 	try {
// 		sessionManager.Put(c, "memory", username+" "+"test")
// 	} catch {
// 		err := "create session fail"
// 		return err;
// 	}

// }

type Login struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

type Message struct {
	Message   string `json: "message"`
	SessionId string `json: "ssid"`
}
