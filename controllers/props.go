package controllers

import (
	"net/http"
	"webserver/models"

	"github.com/gin-gonic/gin"
)

type Props struct{}

func (x *Props) URLPatterns() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/props/caty", ResourceFunc: x.Caty},
		{Method: http.MethodGet, Path: "/props/status", ResourceFunc: x.Status},
		{Method: http.MethodGet, Path: "/props/priority", ResourceFunc: x.Priority},
		{Method: http.MethodGet, Path: "/props/department", ResourceFunc: x.Department},
	}
}

func (x *Props) Caty(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": models.CatyList,
	})
}

func (x *Props) Status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": models.StatusList,
	})
}

func (x *Props) Priority(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": models.PriorityList,
	})
}

func (x *Props) Department(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": models.DepartmentList,
	})
}
