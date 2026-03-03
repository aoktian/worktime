package templatefuncs

import (
	"html/template"
)

var FuncMap = template.FuncMap{
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
}
