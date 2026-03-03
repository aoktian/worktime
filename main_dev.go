//go:build dev
// +build dev

package main

import (
	"webserver/templatefuncs"

	"github.com/gin-gonic/gin"
)

func loadTemplatesAndStatic(r *gin.Engine) {
	r.SetFuncMap(templatefuncs.FuncMap) //要在load 之前调用
	r.LoadHTMLGlob("templates/*")

	// 前端项目静态资源
	r.Static("/static", "./static")
}
