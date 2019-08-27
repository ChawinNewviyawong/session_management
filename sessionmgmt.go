package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
)

func (h *CustomerHandler) createSession(body string) (string, string) {
	go Logger("INFO", ACTOR, "sample_server", "", "createSession", "Request Function", "", h.Channel)
	bytestring := sha256.Sum256([]byte(body))
	sid := hex.EncodeToString(bytestring[:])
	fmt.Println("SHA256 String is ", sid)
	go Logger("DEBUG", ACTOR, "sample_server", "", "createSession", "hash="+sid, "", h.Channel)
	proflie := Profile{
		Username:    "gear",
		Address:     "empiretower",
		Email:       "gear@email.com",
		CompanyName: "ice",
	}
	var jsonData []byte
	jsonData, err := json.Marshal(proflie)
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

// func deleteSession(username string, sid string) error {
// 	valueAsByte, err := getValue(sid)
// 	if err != nil {
// 		return err
// 	}

// 	// if valueAsByte == username {
// 	// 	if status, err := delValue(sid); status != true || err != nil {
// 	// 		return err
// 	// 	}
// 	// }

// 	return nil
// }
