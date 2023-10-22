package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
	"xxvote/app/model"
	"xxvote/app/tools"
)

type User struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.tmpl", nil)
}

func DoLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(), //这里有风险
		})
		return
	}

	tools.Logger.Infof("user:%+v", user)

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}

	ret := model.GetUser(user.Name)
	if ret.Id < 1 || ret.Password != encryptV1(user.Password) {
		context.JSON(http.StatusOK, tools.UserErr)
		return
	}

	//context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	//context.SetCookie("Id", fmt.Sprint(ret.Id), 3600, "/", "", true, false)

	_ = model.SetSession(context, user.Name, ret.Id)

	context.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
	})
	return
}

func Logout(context *gin.Context) {
	//context.SetCookie("name", "", 3600, "/", "", true, false)
	//context.SetCookie("Id", "", 3600, "/", "", true, false)
	_ = model.FlushSession(context)
	context.Redirect(http.StatusFound, "/login")
}

// 新创建一个结构体
type CUser struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
}

func CreateUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
		return
	}
	fmt.Printf("user:%+v", user)

	//encrypt(user.Password)
	//encryptV1(user.Password)
	//encryptV2(user.Password)
	//return

	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}

	//校验密码
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同！", //这里有风险
		})
		return
	}

	nameLen := len(user.Name)
	password := len(user.Password)
	if nameLen > 16 || nameLen < 8 || password > 16 || password < 8 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "账号或密码大于8小于16",
		})
		return
	}

	//密码不能是纯数字 -》 数字+小写字母+大写字母
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.Password) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字", //这里有风险
		})
		return
	}

	//这里有一个巨大的BUG，并发安全！
	if oldUser := model.GetUser(user.Name); oldUser.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在！",
		})
		return
	}

	newUser := model.User{
		Name:        user.Name,
		Password:    encryptV1(user.Password),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := model.CreateUser(&newUser); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10007,
			Message: "新用户创建失败！", //这里有风险
		})
		return
	}

	context.JSON(http.StatusOK, tools.OK)
	return
}

// 最基础的版本
func encrypt(pwd string) string {

	hash := md5.New()
	hash.Write([]byte(pwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码：%s\n", hashString)

	return hashString
}

func encryptV1(pwd string) string {
	newPwd := pwd + "香香编程喵喵喵" //不能随便起，且不能暴露
	hash := md5.New()
	hash.Write([]byte(newPwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码：%s\n", hashString)

	return hashString
}

func encryptV2(pwd string) string {
	//基于Blowfish 实现加密。简单快速，但有安全风险
	//golang.org/x/crypto/ 中有大量的加密算法
	newPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("密码加密失败：", err)
		return ""
	}
	newPwdStr := string(newPwd)
	fmt.Printf("加密后的密码：%s\n", newPwdStr)
	return newPwdStr
}
