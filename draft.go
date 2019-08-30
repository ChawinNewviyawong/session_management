package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Draft struct {
	car     Car
	profile Profile
}

func (draft Draft) Submit(car Car, profile Profile) (string, error) {
	function := "submit"
	// query Permission from Database
	permission, err := queryPermission(function)
	if err != nil {
		return "", err
	}
	var permissionAsString string
	json.NewDecoder(permission).Decode(&permissionAsString)

	// check permission
	checkedPermission := checkedPermission(permissionAsString, profile.Role)
	fmt.Println(checkedPermission)
	if !checkedPermission {
		err := errors.New(string(http.StatusUnauthorized))
		return "", err
	}

	// check state
	if car.State != "draft" {
		err := errors.New("State is not DRAFT, State is " + car.State)
		return "", err
	}

	price, err := strconv.Atoi(car.Price)
	if err != nil {
		return "", nil
	}
	var message string
	if price < 50000 {
		// connect node
		message = "update CAR KEY: " + car.Key + " STATE: " + car.State + " >> COMPLETED"
		car.State = "completed"
	} else if price > 50000 {
		// connect node
		message = "update CAR KEY: " + car.Key + " STATE: " + car.State + " >> AWAITING"
		car.State = "awaiting"
	}

	return message, nil
}

func (draft Draft) Approve(car Car, profile Profile) (string, error) {
	err := errors.New("DRAFT can't APPROVE")
	return "", err
}

func (draft Draft) Reject(car Car, profile Profile) (string, error) {
	err := errors.New("DRAFT can't REJECT")
	return "", err
}

func (draft Draft) Cancel(car Car, profile Profile) (string, error) {
	function := "cancel"
	// query Permission from Database
	permission, err := queryPermission(function)
	if err != nil {
		return "", err
	}
	var permissionAsString string
	json.NewDecoder(permission).Decode(&permissionAsString)

	// check permission
	checkedPermission := checkedPermission(permissionAsString, profile.Role)
	fmt.Println(checkedPermission)
	if !checkedPermission {
		err := errors.New(string(http.StatusUnauthorized))
		return "", err
	}

	// connect node
	message := "update CAR KEY: " + car.Key + " STATE: " + car.State + " >> AWAITING"
	return message, nil
}

func (draft Draft) Complete(car Car, profile Profile) (string, error) {
	err := errors.New("DRAFT can't COMPLETE")
	return "", err
}
