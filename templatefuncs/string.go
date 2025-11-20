package templatefuncs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"webserver/controllers"
)

func html(v interface{}) template.HTML {
	return template.HTML(v.(string))
}

func jsonString(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

func pagesHtml(pagination *controllers.Pagination, pagef string) template.HTML {
	if pagination == nil {
		return ""
	}

	pageShow := 10
	offsetShow := 2
	pages := pagination.GetTotalPage()

	result := "<nav><ul class='pagination'>"

	var from int
	var to int
	if pages <= pageShow {
		from = 1
		to = pages
	} else {
		if pagination.Page <= offsetShow {
			from = 1
			to = pageShow
		} else if pagination.Page >= pages-offsetShow {
			from = pages - pageShow + 1
			to = pages
		} else {
			from = pagination.Page - offsetShow
			to = pagination.Page + offsetShow
		}
	}

	if pagination.Page-offsetShow > 1 && pages > pageShow {
		result += "<li class='page-item'><button class='page-link' onclick='" + pagef + "(1)'>首页</button></li>"
	}
	if pagination.Page > 1 {
		result += fmt.Sprintf("<li class='page-item'><button class='page-link' onclick='"+pagef+"(%d)'>&laquo;</button></li>", pagination.Page-1)
	}
	for i := from; i <= to; i++ {
		if i == pagination.Page {
			result += fmt.Sprintf("<li class='page-item active'><span class='page-link'>%d</span></li>", i)
		} else {
			result += fmt.Sprintf("<li class='page-item'><button class='page-link' onclick='"+pagef+"(%d)'>%d</button></li>", i, i)
		}
	}
	if pagination.Page < pages {
		result += fmt.Sprintf("<li class='page-item'><button class='page-link' onclick='"+pagef+"(%d)'>&raquo;</button></li>", pagination.Page+1)
	}
	if to < pages {
		result += fmt.Sprintf("<li class='page-item'><button class='page-link' onclick='"+pagef+"(%d)'>末页</button></li>", pages)
	}
	result += fmt.Sprintf("<li class='page-item disabled'><em class='page-link'>共 %d 条</em></li>", pagination.Total)
	result += "</ul></nav>"

	return template.HTML(result)
}
