package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/go-redis/redis"
)

func (opt *operation) createConnection() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ping, err := client.Ping().Result()
	fmt.Println(ping, err)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error connection to Redis " + err.Error()
		go Logger("FATAL", ACTOR, "GO_API_SERVER", "", "createConnection", message, "", opt.Channel)

	}
	return client
}

func (opt *operation) setValue(key string, value string) (bool, string) {
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "setValue", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "GO_API_SERVER", "", "setValue", "key="+key+" value="+value, "", opt.Channel)
	client := opt.createConnection()
	err := client.Set(key, value, 0).Err()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error on Set Value to Redis " + err.Error()
		go Logger("FATAL", ACTOR, "GO_API_SERVER", "", "setValue", message, "", opt.Channel)

		return false, message
	}
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "setValue", "Success", "", opt.Channel)
	return true, ""
}

func (opt *operation) getValue(key string) ([]byte, error) {
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "getValue", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "GO_API_SERVER", "", "getValue", "key="+key, "", opt.Channel)
	client := opt.createConnection()
	value, err := client.Get(key).Result()
	valueAsByte := []byte(value)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error on Get Value to Redis " + err.Error()
		go Logger("FATAL", ACTOR, "GO_API_SERVER", "", "getValue", message, "", opt.Channel)
		return nil, err
	}
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "getValue", "Success", "", opt.Channel)
	return valueAsByte, err
}

func (opt *operation) delValue(key string) (bool, error) {
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "delValue", "Request Function", "", opt.Channel)
	go Logger("DEBUG", ACTOR, "GO_API_SERVER", "", "delValue", "key="+key, "", opt.Channel)
	client := opt.createConnection()
	err := client.Del(key).Err()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		message := "[" + file + "][" + strconv.Itoa(line) + "] : Error to Del Value on Redis " + err.Error()
		go Logger("FATAL", ACTOR, "GO_API_SERVER", "", "delValue", message, "", opt.Channel)

		return false, err
	}
	go Logger("INFO", ACTOR, "GO_API_SERVER", "", "delValue", "Success", "", opt.Channel)
	return true, nil
}
