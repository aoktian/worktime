package controllers

import (
	"net/http"
	"webserver/middlewares"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
}

func (x *Auth) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/test", ResourceFunc: x.Test},
		{Method: http.MethodGet, Path: "/home", ResourceFunc: x.Home},
		{Method: http.MethodPost, Path: "/login", ResourceFunc: x.Login},
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

// /welcome/login 的请求体
type ReqLogin struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (x *Auth) Login(c *gin.Context) {
	if utils.AppConfig.Feishu.ClientID != "" {
		utils.JSONErrMsg(c, "已配置飞书授权登录，请使用飞书登录")
		return
	}

	var req ReqLogin
	if err := c.ShouldBind(&req); err != nil {
		utils.JSONError(c, err)
		return
	}

	u := &models.User{}

	has, err := models.DB.Where("account = ?", req.Account).Get(u)
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	if !has {
		utils.JSONError(c, errors.New("user not found"))
		return
	}

	err = models.VerifyPassword(req.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
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
