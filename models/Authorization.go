package models

import "time"

type Authorization struct {
	Id int
	Person *Person `orm:"rel(fk)"`
	Granted time.Time `orm:"auto_now_add;type(datetime)"`
	Target string
}