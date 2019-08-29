package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (opt *operation) quireProfile(login Login) (io.ReadCloser, error) {
	url := "http://ac5dc220.ngrok.io/api" + "/..."
	// fmt.Println("URL:>", url)
	// fmt.Println("payload:>", payload)

	//Json to byteArray
	loginAsBytes, err := json.Marshal(login)
	if err != nil {
		fmt.Println("Marshal is error" + err.Error())
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(loginAsBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		panic(err)
		return nil, err
	}
	defer res.Body.Close()

	return res.Body, nil
}

func (opt *operation) quirePermission(function string) (io.ReadCloser, error) {
	url := "http://ac5dc22.ngrok.io/api" + "/..."

	functionAsByte := []byte(function)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(functionAsByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer res.Body.Close()

	return res.Body, nil
}
