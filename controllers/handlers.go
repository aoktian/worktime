package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
)

func GetAuthUser(ctx *gin.Context) *models.User {
	return ctx.MustGet("authUser").(*models.User)
}

func GetQuery(ctx *gin.Context, name string, defaultValue int64) int64 {
	intParam := ctx.DefaultQuery(name, fmt.Sprintf("%d", defaultValue))
	value, err := strconv.ParseInt(intParam, 10, 64)
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

func GetParamInt64(ctx *gin.Context, name string) int64 {
	idParam := ctx.Param(name)
	if idParam == "" {
		return 0
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func HasTask(ctx *gin.Context, taskProp string, id uint64) bool {
	t := &models.Task{}
	has, err := models.DB.Where(taskProp+" = ?", id).Get(t)
	if err != nil {
		utils.JSONError(ctx, err)
		return true
	}

	if has {
		utils.JSONMsg(ctx, "尚有任务关联此属性，不能删除。")
		return true
	}
	return false
}

func Delete(ctx *gin.Context, id uint64, beans interface{}) {
	_, err := models.DB.ID(id).Delete(beans)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
