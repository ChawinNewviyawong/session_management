package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	gorillaSessions "github.com/gorilla/sessions"
)

var gorillastore = gorillaSessions.NewCookieStore([]byte("SESSION_KEY"))

func main() {
	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
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
		// session.Save()
		ssid := session.Get("sessionid")
		c.JSON(http.StatusOK, gin.H{
			"message":   "user.Username",
			"sessionid": ssid,
		})
	})

	router.GET("/gorillalogin", func(c *gin.Context) {
		gorillasession, err := gorillastore.Get(c.Request, "sessionid")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		gorillasession.Values["username"] = "username"
		gorillasession.Save(c.Request, c.Writer)
		c.JSON(http.StatusOK, gin.H{
			"message": gorillasession.Values["username"],
		})
	})

	router.POST("/getAllCars", func(c *gin.Context) {
		session := sessions.Default(c)
		sessionidserver := session.Get("sessionid")
		var sessionidclient string
		if err := c.ShouldBindJSON(&sessionidclient); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
		} else if sessionidserver != sessionidclient {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
		} else {
			// call ice's api
			requestBody, _ := json.Marshal(map[string]string{
				"profile": sessionidclient})
			response, err := http.Post("", "application/json", bytes.NewBuffer(requestBody))
			if err != nil || response.StatusCode != http.StatusOK {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": http.StatusText(http.StatusServiceUnavailable),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"cars": response.Body,
				})
			}
		}

	})
	router.GET("/show", func(c *gin.Context) {
		gorillasession, err := gorillastore.Get(c.Request, "sessionid")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
			return
		}
		fmt.Println(gorillasession.Values["username"])
		if username, _ := gorillasession.Values["username"]; username == "username" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": gorillasession.Values["username"],
		})
		// session := sessions.Default(c)
		// ssid := session.Get("sessionid")
		// if ssid == "user.Username" {
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"sessionid": ssid,
		// 		"message":   "my session",
		// 	})
		// } else {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": http.StatusText(http.StatusUnauthorized),
		// 	})
		// }
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
