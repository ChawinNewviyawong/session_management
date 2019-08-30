package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

func (opt *operation) createSession(body string, profile Profile) (string, string) {
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "createSession", "Request Function", "", opt.Channel)
	bytestring := sha256.Sum256([]byte(body))
	sid := hex.EncodeToString(bytestring[:])
	fmt.Println("SHA256 String is ", sid)
	go Logger("DEBUG", ACTOR, "GO_API_SERVER", "", "createSession", "hash="+sid, "", opt.Channel)
	var jsonData []byte
	jsonData, err := json.Marshal(profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Marshal " + err.Error()
		go Logger("ERROR", ACTOR, "GO_API_SERVER", "", "createSession", message, "", opt.Channel)
	}

	if status, err := opt.setValue(sid, string(jsonData)); status != true || err != "" {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Set Value to Redis Fail " + err
		go Logger("ERROR", ACTOR, "GO_API_SERVER", "", "createSession", message, "", opt.Channel)
		return "", message
	}

	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "createSession", "Success", "", opt.Channel)
	return sid, ""
}

func (opt *operation) deleteSession(username string, sid string) error {
	profile := Profile{}
	valueAsByte, err := opt.getValue(sid)
	if err != nil {
		return err
	}
	err = json.Unmarshal(valueAsByte, &profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "GO_API_SERVER", "POST", "deleteSession", message, strconv.Itoa(http.StatusBadRequest), opt.Channel)
		return err
	}
	if profile.Username == username {
		if status, err := opt.delValue(sid); status != true || err != nil {
			return err
		}
	}

	return nil
}
