package templatefuncs

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func SetTemplateFuncs(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"mod":           mod,
		"json":          jsonString,
		"html":          html,
		"percent":       percent,
		"intDiv":        intDiv,
		"pages":         pagesHtml,
		"dict":          dict,
		"formatDate":    formatDate,
		"unixToDay":     unixToDay,
		"SecondsToTime": secondsToTime,
	})
}
