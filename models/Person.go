package models

import (
	"time"
)

type Person struct {
	Id                int
	Name              string    `orm:"unique"`
	Birthdate         time.Time `orm:"null"`
	Created           time.Time `orm:"auto_now_add;type(datetime)"`
	Town              string
	LibraryCardNumber int64
	Volunteer		  bool
	Actions           []*Action `orm:"reverse(many)"`
	Authorizations    []*Authorization `orm:"reverse(many)"`
}
