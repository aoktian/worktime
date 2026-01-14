package controllers

import (
	"fmt"
	"net/http"
	"time"
	"webserver/models"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Tag struct {
}

func (x *Tag) URLPatterns() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/manage/tag/:project_id/:page", ResourceFunc: x.InProject},
		{Method: http.MethodPost, Path: "/tag/modify/:id", ResourceFunc: x.Modify},
		{Method: http.MethodPost, Path: "/tag/save", ResourceFunc: x.Save},
		{Method: http.MethodPost, Path: "/tag/delete/:id", ResourceFunc: x.Delete},
		{Method: http.MethodPost, Path: "/tag/stats/:id", ResourceFunc: x.Stats},
		{Method: http.MethodPost, Path: "/tag/gantt/:id", ResourceFunc: x.Gantt},
	}
}

func getTags() map[int]*models.Tag {
	list := make([]*models.Tag, 0)
	models.DB.Find(&list)
	result := make(map[int]*models.Tag)
	for _, u := range list {
		result[u.Id] = u
	}
	return result
}

func (x *Tag) InProject(ctx *gin.Context) {
	project_id := GetParamInt(ctx, "project_id")
	page := GetParamInt(ctx, "page")

	project := &models.Project{}
	models.DB.Id(project_id).Get(project)

	pagination := &Pagination{Page: int(page), Size: 15}
	total, err := models.DB.Where("project_id = ?", project_id).Count(new(models.Tag))
	if err != nil {
		JSONError(ctx, err)
		return
	}
	pagination.Total = int(total)

	results := make([]*models.Tag, 0)
	sql := "select * from tag where project_id = ? order by paixu desc, id desc limit ? offset ?"
	err = models.DB.SQL(sql, project_id, pagination.Size, pagination.GetOffset()).Find(&results)
	if err != nil {
		JSONError(ctx, err)
		return
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["tags"] = results
	h["page"] = pagination
	h["project"] = project

	HTML(ctx, "tags.html", nil)
}

func (x *Tag) Modify(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")

	tag := &models.Tag{}
	if id > 0 {
		models.DB.Id(id).Get(tag)
	} else {
		tag.StartAt = time.Now().Unix()
		tag.EndAt = time.Now().Unix()
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["tag"] = tag

	Dialog(ctx, "tag-edit.html", nil)
}

func (x *Tag) Save(ctx *gin.Context) {
	update := &models.Tag{}
	if !ShouldBindJSON(ctx, update) {
		return
	}

	var err = errors.New("")
	if update.Id > 0 {
		_, err = models.DB.ID(update.Id).MustCols("is_lock").Update(update)
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

func (x *Tag) Delete(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	has := HasTask(ctx, "tag", id)
	if has {
		return
	}
	Delete(ctx, id, &models.Tag{})
}

type TaskStats struct {
	Xid    int
	Count  int
	Status int
}

type TaskStatsResult struct {
	Id    int         `json:"id"`
	Col   string      `json:"col"`
	Name  string      `json:"name"`
	Value map[int]int `json:"value"`
}

func (x *TaskStatsResult) New(id int, col, name string) *TaskStatsResult {
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
		JSONError(ctx, err)
		return
	}

	xTaskStatsResult := &TaskStatsResult{}
	r := make([]*TaskStatsResult, 0)

	total := xTaskStatsResult.New(int(tagId), "tag", "总数")
	r = append(r, total)

	departmentStats := make(map[int]*TaskStatsResult)
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
		JSONError(ctx, err)
		return
	}

	users, err := models.GetUsers()
	if err != nil {
		JSONError(ctx, err)
		return
	}

	leaderStats := make(map[int]*TaskStatsResult)
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

	Dialog(ctx, "tag-stats.html", nil)
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

var ganttPrioritys = map[int]string{
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
		Error(ctx, err)
		return
	}

	counts := make(map[int]int)
	for _, v := range list {
		counts[v.Leader]++
	}
	users := getUsers()
	userTrue := make(map[int]bool)

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
