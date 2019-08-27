package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gorillaSessions "github.com/gorilla/sessions"
)

var gorillastore = gorillaSessions.NewCookieStore([]byte("SESSION_KEY"))

var UUIR_LOGS string
var ACTOR = "robot_test"

func main() {

	UUIR_LOGS, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ACTOR : %s\n", ACTOR)
	fmt.Printf("UUIR_LOGS : %s\n", UUIR_LOGS)

	app := setupRouter()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// router.GET("/show", func(c *gin.Context) {
	// 	gorillasession, err := gorillastore.Get(c.Request, "sessionid")
	// 	if err != nil {
	// 		c.JSON(http.StatusUnauthorized, gin.H{
	// 			"error": http.StatusText(http.StatusUnauthorized),
	// 		})
	// 		return
	// 	}
	// 	fmt.Println(gorillasession.Values["username"])
	// 	if username, _ := gorillasession.Values["username"]; username == "username" {
	// 		c.JSON(http.StatusUnauthorized, gin.H{
	// 			"error": http.StatusText(http.StatusUnauthorized),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": gorillasession.Values["username"],
	// 	})

	// })
	app.Run(":3000")
}

/*
	########################################################################################################
	############################################## MIDELWARE ###############################################
	########################################################################################################
*/
type MyReadCloser struct {
	rc io.ReadCloser
	w  io.Writer
}

func (rc *MyReadCloser) Read(p []byte) (n int, err error) {
	n, err = rc.rc.Read(p)
	if n > 0 {
		if n, err := rc.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
}

func (rc *MyReadCloser) Close() error {
	return rc.rc.Close()
}

func (h *CustomerHandler) LoggerPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			go Logger("DEBUG", ACTOR, "sample_server", "POST", "LoggerPayload", "payload="+buf.String(), "", h.Channel)

		} else if c.Request.Method == http.MethodPut {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			go Logger("DEBUG", ACTOR, "sample_server", "PUT", "LoggerPayload", "payload="+buf.String(), "", h.Channel)

		} else {
			c.Next()
		}
	}
}

func (h *CustomerHandler) InitializeChannel() {
	h.Channel = make(chan string)
	return
}

/*
	########################################################################################################
	############################################## GIN FRANWORK ############################################
	########################################################################################################
*/
func setupRouter() *gin.Engine {

	//log fomat json

	// debug := flag.Bool("debug", true, "sets log level to debug")

	// flag.Parse()

	//เพื่อสร้าง Engine instance ของ Gin
	//มี middleware Logger และ Recovery ติดตั้งมาให้
	app := gin.Default()
	//เหมือน gin.Default() ; Full
	// app := gin.New()

	//middleware

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.

	//mysql
	system := CustomerHandler{}
	// system.InitializeMYSQL()

	system.InitializeChannel()

	app.Use(system.LoggerPayload())

	// result := make(chan string)

	log.SetFlags(log.Lshortfile)
	go Logger("INFO", ACTOR, "sample_server", "", "setupRouter", "Start API Server localhost:8080", "", system.Channel)

	//app router
	app.GET("/", system.FirstPage)
	app.POST("/login", system.Login)
	app.POST("/getAllCars", system.GetAllCars)

	// app.Use(logger.Setgo Logger() )

	return app
}

/*
	########################################################################################################
	######################################### ROUTER&CONTROLLER ############################################
	########################################################################################################
*/

func (h *CustomerHandler) FirstPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "My Session Test",
	})
}

func (h *CustomerHandler) Login(c *gin.Context) {
	go Logger("INFO", ACTOR, "sample_server", "POST", "Login", "Request Function", "", h.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "POST", "Login", "path="+c.Request.RequestURI, "", h.Channel)
	user := Login{}
	// fmt.Println(user.Username)
	if err := c.ShouldBindJSON(&user); err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Login", message, strconv.Itoa(http.StatusBadRequest), h.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	now := time.Now()
	sec := now.Unix()
	body := user.Username
	sid, err := h.createSession(body + string(sec))
	if err != "" {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Login", message, strconv.Itoa(http.StatusBadRequest), h.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	// hardcode
	proflie := Profile{
		Username:    "gear",
		Address:     "empiretower",
		Email:       "gear@email.com",
		CompanyName: "ice",
	}

	go Logger("INFO", ACTOR, "sample_server", "POST", "Login", "Request Success:", strconv.Itoa(http.StatusOK), h.Channel)

	c.JSON(http.StatusOK, gin.H{
		"profile":   proflie,
		"sessionid": sid,
	})
	return
}

func (h *CustomerHandler) GetAllCars(c *gin.Context) {
	profile := Profile{}
	if err := c.ShouldBindJSON(&profile); err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), h.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			// "status": http.StatusBadRequest,
			"error": message,
		})
		return
	}

	valueAsByte, err := h.getValue(profile.SId)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Unauthorized " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusUnauthorized), h.Channel)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		return
	}
	err = json.Unmarshal(valueAsByte, &profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), h.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	var request RequestAllCars
	request.Profile = profile
	requestAsByte, _ := json.Marshal(request)
	// response, err := http.Post("http://3.16.217.238:8080/api/v1/queryAll", "application/json", bytes.NewBuffer(requestBody))
	req, err := http.NewRequest("POST", "http://3.16.217.238:8080/api/v1/queryAll", bytes.NewBuffer(requestAsByte))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), h.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	fmt.Println(response)
	if err != nil || res.StatusCode != 201 {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : ServiceUnavailable " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), h.Channel)

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": message,
		})
		return
	}
	go Logger("INFO", ACTOR, "sample_server", "POST", "GetAllCars", "Request Success:", strconv.Itoa(http.StatusOK), h.Channel)

	c.JSON(http.StatusOK, gin.H{
		"cars": response,
	})
	return

}

type CustomerHandler struct {
	Channel chan string
}

type Response struct {
	Code    int64  `json: "code"`
	Message []Cars `json: "message"`
}

type Cars struct {
	Key    string `json: "Key"`
	Record Car    `json: "Record"`
}

type Car struct {
	Make   string `json: "make"`
	Model  string `json: "model"`
	Colour string `json: "colour"`
	Owner  string `json: "owner"`
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
