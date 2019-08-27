package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/go-redis/redis"
)

func createConnection() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ping, err := client.Ping().Result()
	fmt.Println(ping, err)
	return client
}

func (h *CustomerHandler) setValue(key string, value string) (bool, string) {
	go Logger("INFO", ACTOR, "sample_server", "", "setValue", "Request Function", "", h.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "", "setValue", "key="+key+" value="+value, "", h.Channel)
	client := createConnection()
	err := client.Set(key, value, 0).Err()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error on Set Value to Redis " + err.Error()
		go Logger("FATAL", ACTOR, "sample_server", "", "setValue", message, "", h.Channel)

		return false, message
	}
	go Logger("INFO", ACTOR, "sample_server", "", "setValue", "Success", "", h.Channel)
	return true, ""
}

func (h *CustomerHandler) getValue(key string) ([]byte, error) {
	go Logger("INFO", ACTOR, "sample_server", "", "getValue", "Request Function", "", h.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "", "getValue", "key="+key, "", h.Channel)
	client := createConnection()
	value, err := client.Get(key).Result()
	valueAsByte := []byte(value)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error on Get Value to Redis " + err.Error()
		go Logger("FATAL", ACTOR, "sample_server", "", "getValue", message, "", h.Channel)
		return nil, err
	}
	go Logger("INFO", ACTOR, "sample_server", "", "getValue", "Success", "", h.Channel)
	return valueAsByte, err
}

func (h *CustomerHandler) delValue(key string) (bool, error) {
	go Logger("INFO", ACTOR, "sample_server", "", "delValue", "Request Function", "", h.Channel)
	go Logger("DEBUG", ACTOR, "sample_server", "", "delValue", "key="+key, "", h.Channel)
	client := createConnection()
	err := client.Del(key).Err()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error to Del Value on Redis " + err.Error()
		go Logger("FATAL", ACTOR, "sample_server", "", "delValue", message, "", h.Channel)

		return false, err
	}
	go Logger("INFO", ACTOR, "sample_server", "", "delValue", "Success", "", h.Channel)
	return true, nil
}
