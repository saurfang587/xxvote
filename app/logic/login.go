package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xxvote/app/model"
	"xxvote/app/tools"
)

type User struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.tmpl", nil)
}

func DoLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.Ecode{
			Message: err.Error(), //这里有风险
		})
	}

	ret := model.GetUser(user.Name)
	if ret.Id < 1 || ret.Password != user.Password {
		context.JSON(http.StatusOK, tools.Ecode{
			Message: "帐号密码错误！",
		})
	}

	context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	context.SetCookie("Id", fmt.Sprint(ret.Id), 3600, "/", "", true, false)

	context.JSON(http.StatusOK, tools.Ecode{
		Message: "登录成功",
	})
}

func Logout(context *gin.Context) {
	context.SetCookie("name", "", 3600, "/", "", true, false)
	context.SetCookie("Id", "", 3600, "/", "", true, false)
	context.Redirect(http.StatusFound, "/login")
}
