package controllers

import (
	"net/http"
	"strconv"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
)

type Router interface {
	URLPatterns() []Route
}

type Route struct {
	Method       string                 //Method is one of the following: GET,PUT,POST,DELETE. required
	Path         string                 //Path contains a path pattern. required
	ResourceFunc gin.HandlerFunc        //the func this API calls. you must set this field or ResourceFunc, if you set both, ResourceFunc will be used
	FuncDesc     string                 //tells what this route is all about. Optional.
	Metadata     map[string]interface{} //Metadata adds or updates a key=value pair to api
}

func IsAjax(ctx *gin.Context) bool {
	if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	} else {
		return false
	}
}

func HTML(ctx *gin.Context, name string, userData gin.H) {
	if userData == nil {
		userData = ctx.MustGet("templateData").(map[string]any)
	}
	ctx.HTML(http.StatusOK, name, userData)
}

// 页面跳转
func Redirect(ctx *gin.Context, url string, delay int) {
	if IsAjax(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"redirect_url": url,
		})
	} else {
		ctx.HTML(http.StatusOK, "redirect.html", gin.H{
			"url":   url,
			"delay": 0,
		})
	}
}

func Error(ctx *gin.Context, err error) {
	ErrorMsg(ctx, err.Error())
}

func ErrorMsg(ctx *gin.Context, err string) {
	if IsAjax(ctx) {
		JSONErrMsg(ctx, err)
	} else {
		ctx.HTML(http.StatusOK, "error.html", gin.H{
			"msg": err,
		})
	}
}

func Dialog(ctx *gin.Context, tpl string, userData gin.H) {
	if userData == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"dialog": utils.GetRenderedTemplateContent(ctx, tpl),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"dialog": utils.RenderedTemplateContent(ctx, tpl, userData),
		})
	}
}

func JSONError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"assertAlert": err.Error(),
	})
}

func JSONErrMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"assertAlert": msg,
	})
}

func JSONMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"alert": msg,
	})
}

func ShouldBindJSON(ctx *gin.Context, beans interface{}) bool {
	if err := ctx.ShouldBindJSON(beans); err != nil {
		JSONError(ctx, err)
		return false
	}
	return true
}

func GetAuthUser(ctx *gin.Context) *models.User {
	return ctx.MustGet("authUser").(*models.User)
}

func GetQuery(ctx *gin.Context, name string, defaultValue int) int {
	intParam := ctx.DefaultQuery(name, strconv.Itoa(defaultValue))

	value, err := strconv.Atoi(intParam)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetParamInt(ctx *gin.Context, name string) uint64 {
	idParam := ctx.Param(name)
	if idParam == "" {
		return 0
	}
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func HasTask(ctx *gin.Context, taskProp string, id uint64) bool {
	t := &models.Task{}
	has, err := models.DB.Where(taskProp+" = ?", id).Get(t)
	if err != nil {
		JSONError(ctx, err)
		return true
	}

	if has {
		JSONMsg(ctx, "尚有任务关联此属性，不能删除。")
		return true
	}
	return false
}

func Delete(ctx *gin.Context, id uint64, beans interface{}) {
	_, err := models.DB.ID(id).Delete(beans)
	if err != nil {
		JSONError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
