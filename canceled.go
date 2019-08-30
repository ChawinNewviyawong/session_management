package main

type Canceled struct {
	car     Car
	profile Profile
}

func (canceled Canceled) Submit(car Car, profile Profile) {

}

func (canceled Canceled) Approve(car Car, profile Profile) {

}

func (canceled Canceled) Reject(car Car, profile Profile) {

}

func (canceled Canceled) Cancel(car Car, profile Profile) {

}

func (canceled Canceled) Complete(car Car, profile Profile) {

}
