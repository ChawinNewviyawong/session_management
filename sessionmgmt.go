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

func (h *CustomerHandler) createSession(body string, profile Profile) (string, string) {
	go Logger("INFO", ACTOR, "sample_server", "", "createSession", "Request Function", "", h.Channel)
	bytestring := sha256.Sum256([]byte(body))
	sid := hex.EncodeToString(bytestring[:])
	fmt.Println("SHA256 String is ", sid)
	go Logger("DEBUG", ACTOR, "sample_server", "", "createSession", "hash="+sid, "", h.Channel)
	var jsonData []byte
	jsonData, err := json.Marshal(profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Marshal " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "", "createSession", message, "", h.Channel)
	}

	if status, err := h.setValue(sid, string(jsonData)); status != true || err != "" {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Set Value to Redis Fail " + err
		go Logger("ERROR", ACTOR, "sample_server", "", "createSession", message, "", h.Channel)
		return "", message
	}

	go Logger("INFO", ACTOR, "sample_server", "", "createSession", "Success", "", h.Channel)
	return sid, ""
}

func (h *CustomerHandler) deleteSession(username string, sid string) error {
	profile := Profile{}
	valueAsByte, err := h.getValue(sid)
	if err != nil {
		return err
	}
	err = json.Unmarshal(valueAsByte, &profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : BadRequest " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "POST", "deleteSession", message, strconv.Itoa(http.StatusBadRequest), h.Channel)
		return err
	}
	if profile.Username == username {
		if status, err := h.delValue(sid); status != true || err != nil {
			return err
		}
	}

	return nil
}
