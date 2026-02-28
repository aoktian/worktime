package controllers

import (
	"net/http"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
)

type Comment struct {
}

func (x *Comment) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodPost, Path: "/comment/save", ResourceFunc: x.Save},
		{Method: http.MethodPost, Path: "/comment/log/:id", ResourceFunc: x.GetLog},
	}
}

func (x *Comment) GetLog(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	results := make([]*models.CommentLog, 0)
	err := models.DB.Where("comment_id = ?", id).Find(&results)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	comment := &models.Comment{}
	models.DB.Where("id = ?", id).Get(comment)

	utils.Dialog(ctx, "comment-log.html", gin.H{
		"logs":    results,
		"comment": comment,
	})
}

func (x *Comment) Save(ctx *gin.Context) {
	update := &models.Comment{}
	if !utils.ShouldBindJSON(ctx, update) {
		return
	}

	authUser := ctx.MustGet("authUser").(*models.User)
	if update.Id == 0 {
		update.Author = authUser.Name
	}
	update.Editor = authUser.Name

	if update.TaskId > 0 {
		task := &models.Task{}
		has, err := models.DB.Where("id = ?", update.TaskId).Get(task)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONErrMsg(ctx, "任务不存在")
			return
		}
	}

	if update.Id > 0 {
		comment := &models.Comment{}
		has, err := models.DB.Where("id = ?", update.Id).Get(comment)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONErrMsg(ctx, "编辑的记录不存在")
			return
		}
		go x.log(authUser, comment, update)
	}

	isNew := update.Id == 0
	if update.Id == 0 {
		_, err := models.DB.Insert(update)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
	} else {
		_, err := models.DB.ID(update.Id).Update(update)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
	}

	comment := &models.Comment{}
	models.DB.Where("id = ?", update.Id).Get(comment)

	if isNew {
		h := ctx.MustGet("templateData").(map[string]any)
		h["comment"] = comment
		h["users"] = getUsers()

		ctx.JSON(http.StatusOK, gin.H{
			"content": utils.GetRenderedTemplateContent(ctx, "comment-one.html"),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (x *Comment) log(authUser *models.User, comment *models.Comment, newComment *models.Comment) {
	if comment.Content == newComment.Content {
		return
	}
	//记录的是上一个版本的内容
	log := &models.CommentLog{
		CommentId: comment.Id,
		Operator:  authUser.Name,
		Content:   comment.Content,
		ContentTo: newComment.Content,
	}
	models.DB.Insert(log)
}
