package controllers

import (
	"net/http"
	"webserver/models"

	"github.com/gin-gonic/gin"
)

type Project struct {
}

func (x *Project) URLPatterns() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/manage/projects", ResourceFunc: x.List},

		{Method: http.MethodPost, Path: "project/modify/:id", ResourceFunc: x.Modify},
		{Method: http.MethodPost, Path: "/project/save", ResourceFunc: x.Save},

		{Method: http.MethodPost, Path: "/project/delete/:id", ResourceFunc: x.Delete},
	}
}

func getProjects() map[int]*models.Project {
	list := make([]*models.Project, 0)
	models.DB.Find(&list)
	result := make(map[int]*models.Project)
	for _, u := range list {
		result[u.Id] = u
	}
	return result
}

func (x *Project) List(ctx *gin.Context) {
	results := make([]*models.Project, 0)
	models.DB.Find(&results)

	h := ctx.MustGet("templateData").(map[string]any)
	h["projects"] = results

	HTML(ctx, "projects.html", nil)
}

func (x *Project) Modify(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")

	project := &models.Project{}
	if id > 0 {
		models.DB.Id(id).Get(project)
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["project"] = project

	Dialog(ctx, "project-edit.html", nil)
}

func (x *Project) Save(ctx *gin.Context) {
	update := &models.Project{}
	if !ShouldBindJSON(ctx, update) {
		return
	}

	var err error
	if update.Id > 0 {
		_, err = models.DB.ID(update.Id).Update(update)
	} else {
		_, err = models.DB.InsertOne(update)
	}
	if err != nil {
		JSONError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": update.Id,
	})
}

func (x *Project) Delete(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	has := HasTask(ctx, "project", id)
	if has {
		return
	}
	Delete(ctx, id, new(models.Project))
}
