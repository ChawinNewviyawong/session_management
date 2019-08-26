package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func createSession(body string) (string, error) {
	bytestring := sha256.Sum256([]byte(body))
	sid := hex.EncodeToString(bytestring[:])
	fmt.Println("SHA256 String is ", sid)
	proflie := Profile{
		Username:    "gear",
		Address:     "empiretower",
		Email:       "gear@email.com",
		CompanyName: "ice",
	}

	if status, err := setValue(sid, proflie); status != true || err != nil {
		return "", err
	}
	return sid, nil
}

func deleteSession(username string, sid string) error {
	value, err := getValue(sid)
	if err != nil {
		return err
	}

	if value == username {
		if status, err := delValue(sid); status != true || err != nil {
			return err
		}
	}

	return nil
}
