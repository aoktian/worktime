package utils

import (
	"net/http"
	"strconv"

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

func ApiResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": data,
	})
}

func ApiError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  err.Error(),
	})
}

func ApiErrorMsg(c *gin.Context, err string) {
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  err,
	})
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
			"dialog": GetRenderedTemplateContent(ctx, tpl),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"dialog": RenderedTemplateContent(ctx, tpl, userData),
		})
	}
}

func DomHtml(ctx *gin.Context, dom, tpl string, userData gin.H) {
	ctx.JSON(http.StatusOK, gin.H{
		dom: RenderedTemplateContent(ctx, tpl, userData),
	})
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
