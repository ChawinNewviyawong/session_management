package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (opt *operation) queryProfile(login Login) (io.ReadCloser, error) {
	url := "http://623d1140.ngrok.io/api" + "/getUser/" + login.Username
	// fmt.Println("URL:>", url)
	// fmt.Println("payload:>", payload)

	//Json to byteArray
	// loginAsBytes, err := json.Marshal(login)
	// if err != nil {
	// 	fmt.Println("Marshal is error" + err.Error())
	// 	return nil, err
	// }

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	if err != nil {
		panic(err)
		return nil, err
	}
	defer res.Body.Close()

	return res.Body, nil
}

func queryPermission(function string) (io.ReadCloser, error) {
	url := "http://209.97.167.162:9000/api" + "/..."

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
