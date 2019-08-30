package main

type Transition interface {
	Submit(car Car, profile Profile) (string, error)
	Approve(car Car, profile Profile) (string, error)
	Reject(car Car, profile Profile) (string, error)
	Cancel(car Car, profile Profile) (string, error)
	Complete(car Car, profile Profile) (string, error)
}
