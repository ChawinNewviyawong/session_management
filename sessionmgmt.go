package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
)

func (opt *operation) createSession(body string, profile Profile) (string, string) {
	go Logger("INFO", ACTOR, "sample_server", "", "createSession", "Request Function", "", opt.Channel)
	bytestring := sha256.Sum256([]byte(body))
	sid := hex.EncodeToString(bytestring[:])
	fmt.Println("SHA256 String is ", sid)
	go Logger("DEBUG", ACTOR, "sample_server", "", "createSession", "hash="+sid, "", opt.Channel)
	var jsonData []byte
	jsonData, err := json.Marshal(profile)
	if err != nil {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Marshal " + err.Error()
		go Logger("ERROR", ACTOR, "sample_server", "", "createSession", message, "", opt.Channel)
	}

	if status, err := opt.setValue(sid, string(jsonData)); status != true || err != "" {
		// err.Error() conv to string
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Set Value to Redis Fail " + err
		go Logger("ERROR", ACTOR, "sample_server", "", "createSession", message, "", opt.Channel)
		return "", message
	}

	go Logger("INFO", ACTOR, "sample_server", "", "createSession", "Success", "", opt.Channel)
	return sid, ""
}
