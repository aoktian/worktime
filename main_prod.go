//go:build !dev
// +build !dev

package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"webserver/templatefuncs"

	"github.com/gin-gonic/gin"
)

// 静态资源
//
//go:embed templates/* static/*
var assetsFS embed.FS

func loadTemplatesAndStatic(r *gin.Engine) {
	templ := template.New("").Funcs(templatefuncs.FuncMap)
	templ = template.Must(templ.ParseFS(assetsFS, "templates/*.html"))
	r.SetHTMLTemplate(templ)

	// 示例中的代码不行
	// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/README.md
	// example: /public/assets/images/example.png
	// r.StaticFS("/public", http.FS(f))

	// example: /static/images/example.png
	fp, _ := fs.Sub(assetsFS, "static")
	r.StaticFS("/static", http.FS(fp))
}
