package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xxvote/app/logic"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	//相关的路径 放在这里

	index := r.Group("")
	index.Use(checkUser)
	index.GET("/index", logic.Index)
	index.GET("/vote", logic.GetVoteInfo)
	index.POST("/vote", logic.DoVote)
	r.GET("/", logic.Index)

	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	r.GET("/logout", logic.Logout)
	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败！")
	}
}

func checkUser(context *gin.Context) {
	name, err := context.Cookie("name")
	if err != nil || name == "" {
		context.Redirect(http.StatusFound, "/login")
	}
	context.Next()
}
