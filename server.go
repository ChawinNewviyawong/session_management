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
	TaxId:       "001",
	Role:        "admin",
}
var gear = Profile{
	Username:    "gear",
	Address:     "1234",
	Email:       "gear@email.com",
	CompanyName: "ice_company",
	TaxId:       "001",
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
	app.POST("/getAllCars", opt.GetAllCars2)
	app.POST("/addCar", opt.AddCar)

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

type RequestAddCar struct {
	Profile Profile `json: "profile"`
	Key     string  `json: "key"`
	Make    string  `json: "make"`
	Model   string  `json: "model"`
	Colour  string  `json: "colour"`
	Owner   string  `json: "owner"`
}

type ResponseAddCar struct {
	Code    int64  `json: "code"`
	Message string `json: "message"`
}

type Session struct {
	Profile Profile `json: "profile"`
	SID     string  `json: "sessionId"`
	UUID    string  `json: "uuid"`
}

type Profile struct {
	Username    string `json: "username"`
	Password    string `json: "password"`
	Address     string `json: "address"`
	Email       string `json: "email"`
	CompanyName string `json: "companyName"`
	SId         string `json: "sid"`
	TaxId       string `json: "txid"`
	Role        string `json: "role"`
	UuId        string `json: "uuid"`
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
	return
}

func (opt *operation) Login(c *gin.Context) {
	go Logger("INFO", "", "sample_server", "POST", "Login", "Request Function", "", opt.Channel)
	go Logger("DEBUG", "", "sample_server", "POST", "Login", "path="+c.Request.RequestURI, "", opt.Channel)

	//decode payload request
	login := Login{}
	if err := c.ShouldBindJSON(&login); err != nil {
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

	// get Profile from DB
	res, err := opt.quireProfile(login)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", "", "sample_server", "POST", "Login", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	var profile Profile
	json.NewDecoder(res).Decode(&profile)
	if login.Username == "best" {
		profile = best
		profile.UuId = string(UUID_LOGSAsBytes)
	} else if login.Username == "gear" {
		profile = gear
		profile.UuId = string(UUID_LOGSAsBytes)
	}
	now := time.Now()
	sec := now.Unix()
	sid, errmessage := opt.createSession(profile.Username+strconv.FormatInt(sec, 10), profile)
	if errmessage != "" {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + errmessage
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Login", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	var session Session
	session.Profile = profile
	session.SID = sid
	session.UUID = string(UUID_LOGSAsBytes)
	opt.setUuidAndActor(c, sid)
	go Logger("INFO", ACTOR, "sample_server", "POST", "Login", "ACTOR : "+ACTOR+"UUID_LOGS : "+UUID_LOGS, "", opt.Channel)

	// get value from sid key
	// valueAsByte, err := opt.getValue(session.SID)
	// err = json.Unmarshal(valueAsByte, &profile)
	// if err != nil {
	// 	_, file, line, _ := runtime.Caller(1)
	// 	message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
	// 	go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": message,
	// 	})
	// 	return
	// }

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

func (opt *operation) Logout2(c *gin.Context) {
	var session Session
	// bind json request from client
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// set uuid and actor for log
	opt.setUuidAndActor(c, session.SID)

	// get Profile from Session
	profileAsByte, err := opt.getValue(session.SID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var profile Profile
	err = json.Unmarshal(profileAsByte, &profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err = opt.deleteSession(profile.Username, session.SID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
	return

}

func (opt *operation) Logout(c *gin.Context) {
	user := Profile{}
	// fmt.Println(user.Username)
	if err := c.ShouldBindJSON(&user); err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", "", "sample_server", "POST", "Logout", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	opt.setUuidAndActor(c, user.SId)
	go Logger("INFO", ACTOR, "sample_server", "POST", "Logout", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "POST", "Logout", "path="+c.Request.RequestURI, "", opt.Channel)
	valueAsByte, err := opt.getValue(user.SId)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Unauthorized " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Logout", message, strconv.Itoa(http.StatusUnauthorized), opt.Channel)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		return
	}
	err = json.Unmarshal(valueAsByte, &user)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Logout", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	err = opt.deleteSession(user.Username, user.SId)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "Logout", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	go Logger("INFO", ACTOR, "sample_server", "POST", "Logout", "Request Success:", strconv.Itoa(http.StatusOK), opt.Channel)
	c.JSON(http.StatusOK, gin.H{
		"sessionid": nil,
	})
	return
}

func (opt *operation) GetAllCars2(c *gin.Context) {
	function := "query"
	var session Session

	// bind json request from client
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// quire Permission from Database
	res, err := opt.quirePermission(function)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var permission string
	json.NewDecoder(res).Decode(&permission)

	checkedPermission := opt.checkedPermission(permission, session.Profile.Role)
	fmt.Println(checkedPermission)
	if !checkedPermission {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	var request RequestAllCars
	request.Profile = session.Profile
	requestAsByte, _ := json.Marshal(request)
	url := "http://3.16.217.238:8080/api/v1/queryAll"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestAsByte))
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

	if err != nil || res.StatusCode != 201 {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cars": response,
	})
	return
}

func (opt *operation) GetAllCars(c *gin.Context) {
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
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Unauthorized000 " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusUnauthorized), opt.Channel)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		return
	}
	opt.setUuidAndActor(c, profile.SId)
	go Logger("INFO", ACTOR, "sample_server", "POST", "GetAllCars", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "POST", "GetAllCars", "path="+c.Request.RequestURI, "", opt.Channel)
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

func (opt *operation) AddCar(c *gin.Context) {

	// var buf bytes.Buffer
	// newBody := &MyReadCloser{c.Request.Body, &buf}
	// c.Request.Body = newBody
	// fmt.Println(c.Request.Body)

	// car := Car{}
	// if err := c.ShouldBindJSON(&car); err != nil {
	// 	// err.Error() conv to string
	// 	_, file, line, _ := runtime.Caller(1)
	// 	message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
	// 	go Logger("ERROR", "", "sample_server", "POST", "AddCar", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": message,
	// 	})
	// 	return
	// }

	user := RequestAddCar{}
	// fmt.Println(user.Username)
	if err := c.ShouldBindJSON(&user); err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest000 " + err.Error()
		go Logger("ERROR", "", "sample_server", "POST", "AddCar", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	fmt.Println(user)
	opt.setUuidAndActor(c, user.Profile.SId)
	go Logger("INFO", ACTOR, "sample_server", "POST", "AddCar", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "POST", "AddCar", "path="+c.Request.RequestURI, "", opt.Channel)

	valueAsByte, err := opt.getValue(user.Profile.SId)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Unauthorized " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "AddCar", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		return
	}
	// fmt.Println(string(valueAsByte))
	profile := Profile{}
	err = json.Unmarshal(valueAsByte, &profile)
	user.Profile = profile
	if profile.Role != "admin" {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Unauthorized "
		go Logger("ERROR", ACTOR, "sample_server", "POST", "AddCar", message, strconv.Itoa(http.StatusUnauthorized), opt.Channel)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		return
	}
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "AddCar", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	var request RequestAddCar
	request.Profile = user.Profile
	request.Key = user.Key
	request.Make = user.Make
	request.Model = user.Model
	request.Colour = user.Colour
	request.Owner = user.Owner
	requestAsByte, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "http://3.16.217.238:8080/api/v1/createCar", bytes.NewBuffer(requestAsByte))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	var response ResponseAddCar
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "AddCar", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	if err != nil || res.StatusCode != 201 {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : ServiceUnavailable " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "GetAllCars", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": message,
		})
		return
	}

	go Logger("INFO", ACTOR, "sample_server", "POST", "AddCar", "Request Success:", strconv.Itoa(http.StatusOK), opt.Channel)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	return
}

func (opt *operation) checkedPermission(permission string, role string) bool {
	permissionAsInt, _ := strconv.ParseInt(permission, 10, 64)
	fmt.Println(permissionAsInt)
	userRoleAsInt, _ := strconv.ParseInt(role, 10, 64)
	fmt.Println(userRoleAsInt)
	checkedPermission := permissionAsInt & userRoleAsInt
	if permissionAsInt != checkedPermission {
		return false
	}
	return true
}

func (opt *operation) setUuidAndActor(c *gin.Context, sid string) {
	go Logger("INFO", ACTOR, "sample_server", "POST", "setUuidAndActor", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "POST", "setUuidAndActor", "path="+c.Request.RequestURI, "", opt.Channel)

	valueAsByte, err := opt.getValue(sid)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "", "setUuidAndActor", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}
	// fmt.Println(string(valueAsByte))
	profile := Profile{}
	err = json.Unmarshal(valueAsByte, &profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "", "setUuidAndActor", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
		return
	}

	go Logger("INFO", profile.Username, "sample_server", "POST", "Login", "ACTOR : "+profile.Username+"UUID_LOGS : "+profile.UuId, "", opt.Channel)
	UUID_LOGS = profile.UuId
	ACTOR = profile.Username
	return
}

// *note*
// admin 1,2,3,4  30
// approver 2,3,4  28
// user 1,2  6
//
// add 2
// quire 4
// approve 8
// cancel 16
