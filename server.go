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
)

var UUID_LOGS string
var ACTOR string

var best = Profile{
	Username:    "best",
	Address:     "1234",
	Email:       "best@email.com",
	CompanyName: "ice_company",
	TxId:        "001",
}
var gear = Profile{
	Username:    "gear",
	Address:     "1234",
	Email:       "gear@email.com",
	CompanyName: "ice_company",
	TxId:        "001",
}

var SERVER_PORT = "3000"
var SERVER_HOST = "0.0.0.0"

type operation struct {
	Channel chan string
}

func main() {
	app := setupRouter()
	app.Run(SERVER_HOST + ":" + SERVER_PORT)
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

func (opt *operation) LoggerPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			go Logger("DEBUG", ACTOR, "sample_server", "POST", "LoggerPayload", "payload="+buf.String(), "", opt.Channel)

		} else if c.Request.Method == http.MethodPut {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			go Logger("DEBUG", ACTOR, "sample_server", "PUT", "LoggerPayload", "payload="+buf.String(), "", opt.Channel)

		} else {
			c.Next()
		}
	}
}

func (opt *operation) InitializeChannel() {
	opt.Channel = make(chan string)
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	//middleware

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.

	//
	opt := operation{}
	opt.InitializeChannel()

	app.Use(opt.LoggerPayload())

	log.SetFlags(log.Lshortfile)
	go Logger("INFO", ACTOR, "sample_server", "", "setupRouter", "Start API Server "+SERVER_HOST+":"+SERVER_PORT, "", opt.Channel)

	//app router
	app.GET("/", opt.FirstPage)
	app.POST("/login", opt.Login)
	app.POST("/getAllCars", opt.GetAllCars)

	// app.Use(logger.Setgo Logger() )

	return app
}

/*
	########################################################################################################
	#################################################### MODEL #############################################
	########################################################################################################
*/

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
	TxId        string `json: "txid"`
}

type Sid struct {
	SId string `json: "sid"`
}

type Message struct {
	Message   string `json: "message"`
	SessionId string `json: "ssid"`
}

/*
	########################################################################################################
	######################################### ROUTER&CONTROLLER ############################################
	########################################################################################################
*/
/*
	USER
*/
func (opt *operation) FirstPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "My Session Test",
	})
}

func (opt *operation) Login(c *gin.Context) {
	go Logger("INFO", "", "sample_server", "POST", "Login", "Request Function", "", opt.Channel)
	go Logger("DEBUG", "", "sample_server", "POST", "Login", "path="+c.Request.RequestURI, "", opt.Channel)

	//decode payload request
	user := Login{}
	if err := c.ShouldBindJSON(&user); err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", "", "sample_server", "POST", "Login", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	//generate uuid is index for logs store
	UUID_LOGSAsBytes, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	UUID_LOGS = string(UUID_LOGSAsBytes)

	ACTOR = user.Username
	go Logger("INFO", ACTOR, "sample_server", "POST", "Login", "ACTOR : "+ACTOR+"UUID_LOGS : "+UUID_LOGS, "", opt.Channel)

	var profile Profile
	if user.Username == "best" {
		profile = best
	} else if user.Username == "gear" {
		profile = gear
	}
	now := time.Now()
	sec := now.Unix()
	body := profile.Username
	sid, errmessage := opt.createSession(body+string(sec), profile)
	if errmessage != "" {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + errmessage
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Login", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	// hardcode
	// proflie := Profile{
	// 	Username:    "gear",
	// 	Address:     "empiretower",
	// 	Email:       "gear@email.com",
	// 	CompanyName: "ice_company",
	// 	TxId:        "001",
	// }

	go Logger("INFO", ACTOR, "sample_server", "POST", "Login", "Request Success:", strconv.Itoa(http.StatusOK), opt.Channel)

	c.JSON(http.StatusOK, gin.H{
		"profile":   profile,
		"sessionid": sid,
	})
	return
}

/*
	CARS
*/
func (opt *operation) GetAllCars(c *gin.Context) {
	go Logger("INFO", "", "sample_server", "POST", "GetAllCars", "Request Function", "", opt.Channel)
	go Logger("DEBUG", "", "sample_server", "POST", "GetAllCars", "path="+c.Request.RequestURI, "", opt.Channel)

	profile := Profile{}

	//decode payload request
	if err := c.ShouldBindJSON(&profile); err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	valueAsByte, err := opt.getValue(profile.SId)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Unauthorized " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusUnauthorized), opt.Channel)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		return
	}
	err = json.Unmarshal(valueAsByte, &profile)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	var request RequestAllCars
	request.Profile = profile
	requestAsByte, _ := json.Marshal(request)
	url := "http://3.16.217.238:8080/api/v1/queryAll"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestAsByte))
	go Logger("DEBUG", ACTOR, "sample_server", "POST", "GetAllCars", `http.NewRequest url:`+url, strconv.Itoa(http.StatusBadRequest), opt.Channel)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	fmt.Println(response)
	if err != nil || res.StatusCode != 201 {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : ServiceUnavailable " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": message,
		})
		return
	}
	go Logger("INFO", ACTOR, "sample_server", "POST", "GetAllCars", "Request Success:", strconv.Itoa(http.StatusOK), opt.Channel)

	c.JSON(http.StatusOK, gin.H{
		"cars": response,
	})
	return

}
