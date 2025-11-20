package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Help struct {
}

func (it *Help) URLPatterns() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/help/:page", ResourceFunc: it.Index},
	}
}

func (it *Help) Index(ctx *gin.Context) {
	page := ctx.Param("page")
	HTML(ctx, "help_"+page+".html", nil)
}
