package models

import "time"

type Action struct {
	Id      int
	Person  *Person   `orm:"rel(fk)"`
	Time    time.Time `orm:"auto_now_add;type(datetime)"`
	Entered bool
}
