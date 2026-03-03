package controllers

import (
	"net/http"
	"webserver/middlewares"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Auth struct {
}

func (x *Auth) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/test", ResourceFunc: x.Test},
		{Method: http.MethodGet, Path: "/home", ResourceFunc: x.Home},
		{Method: http.MethodPost, Path: "/login", ResourceFunc: x.Login},
		{Method: http.MethodPost, Path: "/showregister", ResourceFunc: x.Register},
		{Method: http.MethodPost, Path: "/register", ResourceFunc: x.register},
		{Method: http.MethodGet, Path: "/logout", ResourceFunc: x.Logout},
	}
}

func (x *Auth) Test(c *gin.Context) {
	name := models.CatyDict[0].Name
	c.String(http.StatusOK, name)
}

func (x *Auth) Home(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

func (x *Auth) Register(c *gin.Context) {
	utils.Dialog(c, "user-register.html", gin.H{
		"departments": models.DepartmentDict,
	})
}

func (x *Auth) register(c *gin.Context) {
	params := &models.User{}
	if !utils.ShouldBindJSON(c, params) {
		return
	}

	if params.Name == "" {
		utils.JSONErrMsg(c, "姓名不能为空")
		return
	}
	if params.Department == 0 {
		utils.JSONErrMsg(c, "部门不能为空")
		return
	}
	if params.Password == "" {
		utils.JSONErrMsg(c, "密码不能为空")
		return
	}

	hashedPassword, err := params.GenerateFromPassword()
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	user := &models.User{
		Account:    params.Account,
		Name:       params.Name,
		Department: params.Department,
		Nick:       params.Nick,
		Password:   string(hashedPassword),
	}

	first := &models.User{}
	has, _ := models.DB.Limit(1).Get(first)
	if !has {
		user.IsAdmin = true
		user.Ps = 95
	}

	_, err = models.DB.InsertOne(user)
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	utils.JSONMsg(c, "创建成功，请通知管理员开通权限。")
}

func (x *Auth) Login(c *gin.Context) {
	if utils.AppConfig.Feishu.ClientID != "" {
		utils.JSONErrMsg(c, "已配置飞书授权登录，请使用飞书登录")
		return
	}

	params := &models.User{}
	if !utils.ShouldBindJSON(c, params) {
		return
	}

	u := &models.User{}

	has, err := models.DB.Where("account = ?", params.Account).Get(u)
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	if !has {
		utils.JSONError(c, errors.New("user not found"))
		return
	}

	err = models.VerifyPassword(params.Password, u.Password)
	if err != nil {
		utils.JSONError(c, err)
		return
	}

	middlewares.SaveSession(c, u.Id)

	utils.Redirect(c, "/task/ilist", 0)
}

func (x *Auth) Logout(ctx *gin.Context) {
	middlewares.ClearSession(ctx)
	ctx.Redirect(http.StatusMovedPermanently, "/")
}
