package routers

import (
	"github.com/astaxie/beego"
	"github.com/wfplhatch/membership/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/checkin",&controllers.CheckinController{})
	beego.Router("/signup", &controllers.SignupController{})
}
