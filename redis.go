package main

import (
	"fmt"

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

func setValue(key string, value Profile) (bool, error) {
	client := createConnection()
	err := client.Set(key, value, 0).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func getValue(key string) (string, error) {
	client := createConnection()
	value, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, err
}

func delValue(key string) (bool, error) {
	client := createConnection()
	err := client.Del(key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
