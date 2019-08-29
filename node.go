package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func (opt *operation) queryAllCar(request RequestAllCars) (io.ReadCloser, error) {
	requestAsByte, _ := json.Marshal(request)
	url := "http://3.16.217.238:8080/api/v1/queryAll"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestAsByte))
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
