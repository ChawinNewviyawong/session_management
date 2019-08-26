package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	gorillaSessions "github.com/gorilla/sessions"
)

var gorillastore = gorillaSessions.NewCookieStore([]byte("SESSION_KEY"))

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "My Session Test",
		})
	})

	// router.POST("/signin", func(c *gin.Context) {
	// 	user := Login{}
	// 	fmt.Print(user.Username)
	// 	if err := c.ShouldBindJSON(&user); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"status": http.StatusBadRequest,
	// 			"error":  err.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(200, gin.H{
	// 		"status":  http.StatusOK,
	// 		"message": user.Username + " " + user.Password,
	// 	})
	// })
	router.POST("/login", func(c *gin.Context) {
		user := Login{}
		// fmt.Println(user.Username)
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				// "status": http.StatusBadRequest,
				"error": err.Error(),
			})
			return
		}
		// session := sessions.Default(c)
		// session.Set("sessionid", "user.Username")
		// // session.Save()
		// ssid := session.Get("sessionid")

		now := time.Now()
		sec := now.Unix()
		body := user.Username
		fmt.Println(user.Username)
		sid, err := createSession(body + string(sec))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		fmt.Println(sid)
		proflie := Profile{
			Username:    "gear",
			Address:     "empiretower",
			Email:       "gear@email.com",
			CompanyName: "ice",
		}
		c.JSON(http.StatusOK, gin.H{
			// "message":   "user.Username",
			"profile":   proflie,
			"sessionid": sid,
		})
	})

	// router.GET("/gorillalogin", func(c *gin.Context) {
	// 	gorillasession, err := gorillastore.Get(c.Request, "sessionid")
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{
	// 			"error": http.StatusText(http.StatusUnauthorized),
	// 		})
	// 		return
	// 	}

	// 	gorillasession.Values["username"] = "username"
	// 	gorillasession.Save(c.Request, c.Writer)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": gorillasession.Values["username"],
	// 	})
	// })

	router.POST("/getAllCars", func(c *gin.Context) {
		profile := Profile{}
		if err := c.ShouldBindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				// "status": http.StatusBadRequest,
				"error": err.Error(),
			})
			return
		}

		value, err := getValue(profile.SId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": http.StatusText(http.StatusUnauthorized),
			})
			return
		}
		fmt.Println(value)

		// session := sessions.Default(c)
		// sessionidserver := session.Get("sessionid")
		// var sessionidclient string
		// if err := c.ShouldBindJSON(&sessionidclient); err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": http.StatusText(http.StatusUnauthorized),
		// 	})
		// } else if sessionidserver != sessionidclient {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": http.StatusText(http.StatusUnauthorized),
		// 	})
		// } else {

		// call ice's api
		// requestBody, _ := json.Marshal(profile)
		var request RequestAllCars
		request.Profile = profile
		// jsonDataAsByte, _ := json.Marshal(profile)
		// err = json.Unmarshal(jsonDataAsByte, &request)
		requestAsByte, _ := json.Marshal(request)
		fmt.Println(requestAsByte)
		// response, err := http.Post("http://3.16.217.238:8080/api/v1/queryAll", "application/json", bytes.NewBuffer(requestBody))
		req, err := http.NewRequest("POST", "http://3.16.217.238:8080/api/v1/queryAll", bytes.NewBuffer(requestAsByte))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		var cars Cars
		json.NewDecoder(res.Body).Decode(&cars)
		fmt.Println(cars)
		if err != nil || res.StatusCode != 201 {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"cars": res.Body,
			})
		}
		// }

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

type Cars struct {
	Cars []Car
}

type Car struct {
	Make   string
	Model  string
	Colour string
	Owner  string
}

type Login struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

type RequestAllCars struct {
	Profile Profile `json: profile`
}

type Profile struct {
	Username    string `json: "username"`
	Address     string `json: "address"`
	Email       string `json: "email"`
	CompanyName string `json: "companyName"`
	SId         string `json: "sid"`
}

type Sid struct {
	SId string `json: "sid"`
}

type Message struct {
	Message   string `json: "message"`
	SessionId string `json: "ssid"`
}
