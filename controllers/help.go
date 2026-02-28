package controllers

import (
	"net/http"
	"webserver/utils"

	"github.com/gin-gonic/gin"
)

type Help struct {
}

func (it *Help) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/help/:page", ResourceFunc: it.Index},
	}
}

func (it *Help) Index(ctx *gin.Context) {
	page := ctx.Param("page")
	utils.HTML(ctx, "help_"+page+".html", nil)
}
