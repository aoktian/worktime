package controllers

import (
	"net/http"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/builder"
)

type Project struct {
}

func (x *Project) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/project/list", ResourceFunc: x.index},
		{Method: http.MethodPost, Path: "/project/search", ResourceFunc: x.search},
		{Method: http.MethodPost, Path: "/project/selector", ResourceFunc: x.selector},

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

func (x *Project) selector(ctx *gin.Context) {
	params := &selectorReq{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}

	pageSize := 20
	offset := (params.Page - 1) * pageSize

	results := make([]*models.Project, 0)
	if params.Keyword == "" {
		models.DB.Where("is_archived = 0").OrderBy("paixu desc").
			Limit(pageSize, offset).Find(&results)
	} else {
		models.DB.Where("name like ?", "%"+params.Keyword+"%").OrderBy("paixu desc").
			Limit(pageSize, offset).Find(&results)
	}

	more := false
	if len(results) == pageSize {
		more = true
	}
	ctx.JSON(http.StatusOK, gin.H{
		"results":    results,
		"pagination": gin.H{"more": more},
	})
}

type projectSearch struct {
	models.Project
	Page        int64 `json:"page"`
	CreatedAt0  int64 `json:"created_at0"`
	CreatedAt9  int64 `json:"created_at9"`
	SearchState int   `json:"search_state"`
}

func (x *Project) index(ctx *gin.Context) {
	x.list(ctx, &projectSearch{})
	utils.HTML(ctx, "projects.html", nil)
}

func (x *Project) search(ctx *gin.Context) {
	params := &projectSearch{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}
	x.list(ctx, params)
	ctx.JSON(http.StatusOK, gin.H{
		"#data-table": utils.GetRenderedTemplateContent(ctx, "projects-content.html"),
	})
}

func (x *Project) list(ctx *gin.Context, con *projectSearch) {
	where := builder.NewCond()

	if con.Name != "" {
		where = where.And(builder.Like{"name", con.Name})
	}
	if con.CreatedAt0 > 0 {
		where = where.And(builder.Gte{"created_at": con.CreatedAt0})
	}
	if con.CreatedAt9 > 0 {
		where = where.And(builder.Lte{"created_at": con.CreatedAt9})
	}
	if con.SearchState > -1 {
		where = where.And(builder.Eq{"is_archived": con.SearchState})
	}

	sqlWhere, args, err := builder.ToSQL(where)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	var pageSize int64 = 12
	pagination := &utils.Pagination{Page: int64(con.Page), Size: pageSize}
	total, err := models.DB.Where(sqlWhere, args...).Count(new(models.Project))
	if err != nil {
		utils.Error(ctx, err)
		return
	}
	pagination.Total = total

	results := make([]*models.Project, 0)
	models.DB.Where(sqlWhere, args...).
		OrderBy("paixu desc, id desc").
		Limit(int(pagination.Size), int(pagination.GetOffset())).Find(&results)

	h := ctx.MustGet("templateData").(map[string]any)
	h["projects"] = results
	h["page"] = pagination

}

func (x *Project) Modify(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")

	project := &models.Project{}
	if id > 0 {
		models.DB.Id(id).Get(project)
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["project"] = project

	utils.Dialog(ctx, "project-edit.html", nil)
}

func (x *Project) Save(ctx *gin.Context) {
	update := &models.Project{}
	if !utils.ShouldBindJSON(ctx, update) {
		return
	}

	var err error
	if update.Id > 0 {
		_, err = models.DB.ID(update.Id).AllCols().Update(update)
	} else {
		_, err = models.DB.InsertOne(update)
	}
	if err != nil {
		utils.JSONError(ctx, err)
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
