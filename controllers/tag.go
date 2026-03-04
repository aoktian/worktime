package controllers

import (
	"fmt"
	"net/http"
	"time"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/builder"
	"github.com/pkg/errors"
)

type Tag struct {
}

func (x *Tag) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/project/tag/:project_id", ResourceFunc: x.InProject},
		{Method: http.MethodPost, Path: "/tag/search", ResourceFunc: x.search},
		{Method: http.MethodPost, Path: "/tag/selector", ResourceFunc: x.selector},
		{Method: http.MethodPost, Path: "/tag/modify", ResourceFunc: x.Modify},
		{Method: http.MethodPost, Path: "/tag/save", ResourceFunc: x.Save},
		{Method: http.MethodPost, Path: "/tag/delete/:id", ResourceFunc: x.Delete},
		{Method: http.MethodPost, Path: "/tag/stats/:id", ResourceFunc: x.Stats},
		{Method: http.MethodPost, Path: "/tag/gantt/:id", ResourceFunc: x.Gantt},
	}
}

func getTags() map[int64]*models.Tag {
	results := make(map[int64]*models.Tag, 0)
	models.DB.Find(&results)
	return results
}

type selectorReq struct {
	Keyword   string `json:"keyword"`
	ProjectId int64  `json:"project_id"`
	Page      int    `json:"page"`
}

func (x *Tag) selector(ctx *gin.Context) {
	params := &selectorReq{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}

	pageSize := 20

	results := make([]*models.Tag, 0)
	offset := (params.Page - 1) * pageSize
	if params.Keyword == "" {
		models.DB.Where("project_id = ? and is_archived = 0", params.ProjectId).
			OrderBy("paixu desc").
			Limit(pageSize, offset).Find(&results)
	} else {
		models.DB.Where("project_id = ? and name like ?", params.ProjectId, "%"+params.Keyword+"%").
			OrderBy("paixu desc").
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

type tagSearch struct {
	models.Tag
	Page        int   `json:"page"`
	SearchState int   `json:"search_state"`
	CreatedAt0  int64 `json:"created_at0"`
	CreatedAt9  int64 `json:"created_at9"`
}

func (x *Tag) InProject(ctx *gin.Context) {
	project_id := GetParamInt64(ctx, "project_id")
	params := &tagSearch{
		Tag:  models.Tag{ProjectId: int(project_id)},
		Page: 1,
	}
	x.list(ctx, params)
	utils.HTML(ctx, "tags.html", nil)
}

func (x *Tag) search(ctx *gin.Context) {
	params := &tagSearch{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}
	x.list(ctx, params)
	ctx.JSON(http.StatusOK, gin.H{
		"#data-table": utils.GetRenderedTemplateContent(ctx, "tags-content.html"),
	})
}

type listTag struct {
	Id        int64  `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	Name      string `json:"name" xorm:"notnull unique(projectid_name) comment('版本名称')"`
	ProjectId int    `json:"project_id" xorm:"notnull unique(projectid_name) comment('项目id')"`

	Paixu      int  `json:"paixu"`
	IsArchived bool `json:"is_archived"`

	StartAt int64 `json:"start_at" xorm:"comment('开始时间')"`
	EndAt   int64 `json:"end_at" xorm:"comment('结束时间')"`

	CreatedAt int64 `json:"created_at" xorm:"comment('创建时间') created"`
	UpdatedAt int64 `json:"updated_at" xorm:"comment('更新时间') updated"`

	ProjectName string `json:"project_name"`
}

func (x *Tag) list(ctx *gin.Context, con *tagSearch) {
	where := builder.NewCond()
	if con.ProjectId > 0 {
		where = where.And(builder.Eq{"project_id": con.ProjectId})
	}
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
		where = where.And(builder.Eq{"tag.is_archived": con.SearchState})
	}

	sqlWhere, args, err := builder.ToSQL(where)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	project := &models.Project{}
	if con.ProjectId > 0 {
		models.DB.Id(con.ProjectId).Get(project)
	} else {
		project.Name = "全部"
	}

	var pageSize int64 = 12
	pagination := &utils.Pagination{Page: int64(con.Page), Size: pageSize}
	total, err := models.DB.Where(sqlWhere, args...).Count(new(models.Tag))
	if err != nil {
		utils.Error(ctx, err)
		return
	}
	pagination.Total = total

	results := make([]*listTag, 0)
	err = models.DB.Table("tag").Select("tag.*, project.name as project_name").
		Join("LEFT", "project", "project.id = tag.project_id").
		Where(sqlWhere, args...).
		OrderBy("paixu desc, tag.id desc").
		Limit(int(pagination.Size), int(pagination.GetOffset())).Find(&results)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["tags"] = results
	h["page"] = pagination
	h["project"] = project
}

func (x *Tag) Modify(ctx *gin.Context) {
	params := &models.Tag{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}

	project := &models.Project{}
	tag := &models.Tag{}
	if params.Id > 0 {
		models.DB.Id(params.Id).Get(tag)
		models.DB.Id(tag.ProjectId).Get(project)
	} else {
		tag.StartAt = time.Now().Unix()
		tag.EndAt = time.Now().Unix()
		if params.ProjectId > 0 {
			models.DB.Id(params.ProjectId).Get(project)
		}
	}
	if project.Id <= 0 {
		project.Name = "项目"
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["tag"] = tag
	h["project"] = project

	utils.Dialog(ctx, "tag-edit.html", nil)
}

func (x *Tag) Save(ctx *gin.Context) {
	update := &models.Tag{}
	if !utils.ShouldBindJSON(ctx, update) {
		return
	}

	if update.ProjectId <= 0 {
		utils.JSONErrMsg(ctx, "项目id不能为空")
		return
	}

	var err = errors.New("")
	if update.Id > 0 {
		_, err = models.DB.ID(update.Id).MustCols("is_archived").Update(update)
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

func (x *Tag) Delete(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	has := HasTask(ctx, "tag", id)
	if has {
		return
	}
	Delete(ctx, id, &models.Tag{})
}

type TaskStats struct {
	Xid    int64
	Count  int
	Status int
}

type TaskStatsResult struct {
	Id    int64       `json:"id"`
	Col   string      `json:"col"`
	Name  string      `json:"name"`
	Value map[int]int `json:"value"`
}

func (x *TaskStatsResult) New(id int64, col, name string) *TaskStatsResult {
	return &TaskStatsResult{
		Id:    id,
		Col:   col,
		Name:  name,
		Value: map[int]int{0: 0, -1: 0}, //总数 完成
	}
}

func (x *TaskStatsResult) StatusValue(status int, value int) {
	x.Value[status] = value
	x.Value[0] += value
	if status > 90 {
		x.Value[-1] += value
	}
}

func (x *Tag) Stats(ctx *gin.Context) {
	tagId := GetParamInt(ctx, "id")
	results := make([]*TaskStats, 0)

	sql := "select leader_dept as xid, status, count(id) as count from task where tag = ? group by leader_dept, status order by xid "
	err := models.DB.SQL(sql, tagId).Find(&results)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	xTaskStatsResult := &TaskStatsResult{}
	r := make([]*TaskStatsResult, 0)

	total := xTaskStatsResult.New(int64(tagId), "tag", "总数")
	r = append(r, total)

	departmentStats := make(map[int64]*TaskStatsResult)
	for _, v := range results {
		ele, ok := departmentStats[v.Xid]
		if !ok {
			ele = xTaskStatsResult.New(v.Xid, "department", models.DepartmentDict[v.Xid].Name)
			departmentStats[v.Xid] = ele
			r = append(r, ele)
		}
		ele.StatusValue(v.Status, v.Count)
	}

	results = make([]*TaskStats, 0)
	sql = "select leader as xid, status, count(id) as count, max(leader_dept) as leader_dept from task where tag = ? group by leader, status order by leader_dept "
	err = models.DB.SQL(sql, tagId).Find(&results)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	users, err := models.GetUsers()
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	leaderStats := make(map[int64]*TaskStatsResult)
	for _, v := range results {
		ele, ok := leaderStats[v.Xid]
		if !ok {
			ele = xTaskStatsResult.New(v.Xid, "leader", users[v.Xid].Name)
			leaderStats[v.Xid] = ele
			r = append(r, ele)
		}
		ele.StatusValue(v.Status, v.Count)
		total.StatusValue(v.Status, v.Count)
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["stats"] = r
	h["status"] = models.StatusList

	utils.Dialog(ctx, "tag-stats.html", nil)
}

type GanttValue struct {
	From        int64  `json:"from"`
	To          int64  `json:"to"`
	Label       string `json:"label"`
	Desc        string `json:"desc"`
	CustomClass string `json:"customClass"` //ganttBlue ganttRed ganttOrange ganttGreen
}

type Gantt struct {
	Name   string        `json:"name"`
	Desc   string        `json:"desc"`
	Values []*GanttValue `json:"values"`
}

var ganttPrioritys = map[int64]string{
	9:     "bg-secondary'>低",
	99:    "bg-primary'>中",
	999:   "bg-warning'>高",
	99999: "bg-danger'>急",
}

func (x *Tag) Gantt(ctx *gin.Context) {
	tagId := GetParamInt(ctx, "id")

	list := make([]models.Task, 0)
	err := models.DB.Where("tag = ?", tagId).OrderBy("leader").Find(&list)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	counts := make(map[int64]int)
	for _, v := range list {
		counts[v.Leader]++
	}
	users := getUsers()
	userTrue := make(map[int64]bool)

	now := time.Now().Unix() //当前时间

	source := make([]*Gantt, 0)

	for _, v := range list {
		gantt := &Gantt{
			Values: make([]*GanttValue, 1),
		}
		source = append(source, gantt)

		if !userTrue[v.Leader] {
			userTrue[v.Leader] = true
			gantt.Name = fmt.Sprintf("<span class='badge text-bg-secondary'>%d</span>%s", counts[v.Leader], users[v.Leader].Name)
		}

		gantt.Desc = fmt.Sprintf("<span class='badge %s</span> <a href='/task/thread/%d' target='_blank'>%s</a>", ganttPrioritys[v.Priority], v.Id, v.Title)

		customClass := "ganttOrange"
		if v.Status > 50 {
			customClass = "ganttGreen"
		} else {
			if v.EndAt > 0 && time.Now().Unix() < v.EndAt {
				customClass = "ganttRed" //超期
			}
		}

		from := v.StartAt
		if from == 0 {
			from = now
		}
		to := v.EndAt
		if to == 0 {
			to = now + 86400*7
		}

		gantt.Values[0] = &GanttValue{
			From:        from * 1000,
			To:          to * 1000,
			CustomClass: customClass,
			Label:       fmt.Sprintf("%.0f%s", float64(now-from)/float64(to-from)*100, "%"), //计算进度
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"source": source})
}
