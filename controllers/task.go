package controllers

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"webserver/models"
	"webserver/utils"

	htmldiff "github.com/documize/html-diff"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/builder"
	"github.com/pkg/errors"
)

type Task struct{}

func (x *Task) URLPatterns() []utils.Route {
	return []utils.Route{
		{Method: http.MethodGet, Path: "/task/list", ResourceFunc: x.All},
		{Method: http.MethodGet, Path: "/task/ilist", ResourceFunc: x.Ileader},
		{Method: http.MethodGet, Path: "/task/icheck", ResourceFunc: x.Icheck},
		{Method: http.MethodGet, Path: "/task/itest", ResourceFunc: x.Itest},

		{Method: http.MethodGet, Path: "/task/listcon", ResourceFunc: x.ListCon},

		{Method: http.MethodGet, Path: "/task/project/:id", ResourceFunc: x.Project},
		{Method: http.MethodGet, Path: "/task/tag/:id", ResourceFunc: x.Tag},

		{Method: http.MethodPost, Path: "/task/list", ResourceFunc: x.List},
		{Method: http.MethodGet, Path: "/task/thread/:id", ResourceFunc: x.Get},
		{Method: http.MethodGet, Path: "/task/new/:id", ResourceFunc: x.NewTask},

		{Method: http.MethodPost, Path: "/task/save", ResourceFunc: x.Save},
		{Method: http.MethodPost, Path: "/task/batchupdate", ResourceFunc: x.BatchUpdate},
		{Method: http.MethodPost, Path: "/task/log/:id", ResourceFunc: x.GetLog},

		{Method: http.MethodPost, Path: "/task/related/:id", ResourceFunc: x.Related},
		{Method: http.MethodPost, Path: "/task/relatedto", ResourceFunc: x.RelatedTo},

		{Method: http.MethodPost, Path: "/task/uploadImage", ResourceFunc: x.UploadImage},

		{Method: http.MethodGet, Path: "/task/logdiff/:id", ResourceFunc: x.LogDiff},
	}
}

type listForm struct {
	models.Task
	AuthorDept int64  `json:"author_dept"`
	Page       int    `json:"page"`
	OrderBy    string `json:"order_by"`
	MainTask   int    `json:"main_task"`
}

func (x *listForm) getOrderBy() string {
	switch x.OrderBy {
	case "news":
		return "task.id desc"
	case "created_at":
		return "task.id desc"
	case "tag":
		return "tag.paixu desc, status, priority desc, level desc, task.id desc"
	case "end_at":
		return "end_at desc"
	default:
		return "status, priority desc, level desc, task.id desc"
	}
}

const (
	TaskListPageSize = 15
)

func (x *Task) All(ctx *gin.Context) {
	con := &listForm{
		Page:    1,
		Task:    models.Task{},
		OrderBy: "tag",
	}
	x.list(ctx, con)
}

func (x *Task) Ileader(ctx *gin.Context) {
	authUser := GetAuthUser(ctx)
	con := &listForm{
		Page: 1,
		Task: models.Task{
			LeaderDept: authUser.Department,
			Leader:     authUser.Id,
		},
	}
	x.list(ctx, con)
}

func (x *Task) Icheck(ctx *gin.Context) {
	authUser := GetAuthUser(ctx)
	con := &listForm{
		Page: 1,
		Task: models.Task{
			CheckerDept: authUser.Department,
			Checker:     authUser.Id,
		},
	}
	x.list(ctx, con)
}

func (x *Task) Itest(ctx *gin.Context) {
	authUser := GetAuthUser(ctx)
	con := &listForm{
		Page: 1,
		Task: models.Task{
			TesterDept: authUser.Department,
			Tester:     authUser.Id,
		},
	}
	x.list(ctx, con)
}

func (x *Task) Tag(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")

	tag := &models.Tag{}
	models.DB.Id(id).Get(tag)
	con := &listForm{
		Page: 1,
		Task: models.Task{
			Tag:     tag.Id,
			Project: tag.ProjectId,
		},
	}
	x.list(ctx, con)
}

func (x *Task) Project(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	con := &listForm{
		Page: 1,
		Task: models.Task{
			Project: int(id),
		},
	}
	x.list(ctx, con)
}

func (x *Task) ListCon(ctx *gin.Context) {
	con := &listForm{
		Page: 1,
		Task: models.Task{
			Tag:        GetQuery(ctx, "tag", 0),
			LeaderDept: GetQuery(ctx, "department", 0),
			Leader:     GetQuery(ctx, "leader", 0),
			Author:     GetQuery(ctx, "author", 0),
		},
	}

	x.list(ctx, con)
}

func (x *Task) List(ctx *gin.Context) {
	con := &listForm{
		OrderBy: "tag",
	}
	if !utils.ShouldBindJSON(ctx, con) {
		return
	}

	x.list(ctx, con)
}

func (x *Task) list(ctx *gin.Context, con *listForm) {
	where := builder.NewCond()

	if con.MainTask > 0 {
		where = where.And(builder.Eq{"pid": 0})
	}
	if con.Project > 0 {
		where = where.And(builder.Eq{"project": con.Project})
	}
	if con.Tag > 0 {
		where = where.And(builder.Eq{"tag": con.Tag})
	}

	if con.Author > 0 {
		where = where.And(builder.Eq{"author": con.Author})
	} else if con.AuthorDept > 0 {
		var authers []int64
		models.DB.Table("user").Where("department = ?", con.AuthorDept).Cols("id").Find(&authers)
		if len(authers) > 0 {
			where = where.And(builder.In("author", authers))
		}
	}
	if con.Leader > 0 {
		where = where.And(builder.Eq{"leader": con.Leader})
	} else if con.LeaderDept > 0 {
		where = where.And(builder.Eq{"leader_dept": con.LeaderDept})
	}
	if con.Checker > 0 {
		where = where.And(builder.Eq{"checker": con.Checker})
	} else if con.CheckerDept > 0 {
		where = where.And(builder.Eq{"checker_dept": con.CheckerDept})
	}
	if con.Tester > 0 {
		where = where.And(builder.Eq{"tester": con.Tester})
	} else if con.TesterDept > 0 {
		where = where.And(builder.Eq{"tester_dept": con.TesterDept})
	}

	if con.Caty > 0 {
		where = where.And(builder.Eq{"caty": con.Caty})
	}
	if con.Status > 0 {
		where = where.And(builder.Eq{"status": con.Status})
	} else {
		if con.OrderBy == "news" || con.OrderBy == "end_at" { //最新任务，排除已经完成
			where = where.And(builder.Lt{"status": 900})
		}
	}
	if con.Priority > 0 {
		where = where.And(builder.Eq{"priority": con.Priority})
	}
	if con.Level > 0 {
		where = where.And(builder.Eq{"level": con.Level})
	}
	if con.Title != "" {
		where = where.And(builder.Like{"title", con.Title})
	}

	sqlWhere, args, err := builder.ToSQL(where)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	total, err := models.DB.Where(sqlWhere, args...).Count(new(models.Task))
	if err != nil {
		utils.Error(ctx, err)
		return
	}
	pagination := &utils.Pagination{Page: int64(con.Page), Size: TaskListPageSize, Total: total}
	pagination.FormatPage()

	offset := pagination.GetOffset()

	result := make([]*models.Task, 0)
	err = models.DB.Table("task").
		Join("LEFT", "tag", "task.tag = tag.id").
		Where(sqlWhere, args...).OrderBy(con.getOrderBy()).Limit(TaskListPageSize, int(offset)).Find(&result)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	templateData := ctx.MustGet("templateData").(map[string]any)
	templateData["tasks"] = result
	templateData["pagination"] = pagination
	templateData["filter"] = con

	x.addCommonData(templateData)

	if utils.IsAjax(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"#task-table": utils.GetRenderedTemplateContent(ctx, "task-list.html"),
		})
	} else {
		orderedTags := make([]*models.Tag, 0)
		models.DB.OrderBy("paixu desc, id").Find(&orderedTags)
		templateData["orderedTags"] = orderedTags
		orderedUsers := make([]*models.User, 0)
		models.DB.OrderBy("team, id").Find(&orderedUsers)
		templateData["orderedUsers"] = orderedUsers

		utils.HTML(ctx, "task.html", nil)
	}
}

func (x *Task) addCommonData(userData map[string]any) {
	userData["tags"] = getTags()
	userData["users"] = getUsers()
	userData["projects"] = getProjects()

	userData["props"] = gin.H{
		"caty":       models.CatyDict,
		"status":     models.StatusDict,
		"priority":   models.PriorityDict,
		"department": models.DepartmentDict,
	}
}

func (x *Task) NewTask(ctx *gin.Context) {
	taskId := GetParamInt(ctx, "id")
	task := &models.Task{}
	if taskId > 0 {
		has, err := models.DB.Where("id = ?", taskId).Get(task)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONErrMsg(ctx, "没有该任务")
			return
		}
	} else {
		authUser := GetAuthUser(ctx)
		task.Leader = authUser.Id
		task.LeaderDept = authUser.Department
		task.Checker = authUser.Id
		task.CheckerDept = authUser.Department
		task.Tester = authUser.Id
		task.TesterDept = authUser.Department
		task.Priority = 99
	}

	templateData := ctx.MustGet("templateData").(map[string]any)
	templateData["task"] = task

	orderedTags := make([]*models.Tag, 0)
	models.DB.OrderBy("paixu desc, id").Find(&orderedTags)
	templateData["orderedTags"] = orderedTags
	orderedUsers := make([]*models.User, 0)
	models.DB.OrderBy("team, id").Find(&orderedUsers)
	templateData["orderedUsers"] = orderedUsers

	x.addCommonData(templateData)

	utils.HTML(ctx, "task-add.html", nil)
}

func (x *Task) Get(ctx *gin.Context) {
	taskId := GetParamInt(ctx, "id")
	if taskId == 0 {
		utils.JSONErrMsg(ctx, "参数错误")
		return
	}
	task := &models.Task{}
	has, err := models.DB.Where("id = ?", taskId).Get(task)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if !has {
		utils.JSONErrMsg(ctx, "没有该任务")
		return
	}

	searchId := task.Id
	if task.Pid > 0 {
		searchId = task.Pid
	}

	result := make([]*models.Task, 0)
	err = models.DB.Where("pid = ? or id = ?", searchId, searchId).OrderBy("pid, status, priority desc, level desc, id desc").Find(&result)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	templateData := ctx.MustGet("templateData").(map[string]any)
	templateData["task"] = task
	templateData["tasks"] = result
	templateData["pagination"] = &utils.Pagination{Page: 1, Total: int64(len(result)), Size: 10000}

	orderedTags := make([]*models.Tag, 0)
	models.DB.OrderBy("paixu desc, id").Find(&orderedTags)
	templateData["orderedTags"] = orderedTags
	orderedUsers := make([]*models.User, 0)
	models.DB.OrderBy("team, id").Find(&orderedUsers)
	templateData["orderedUsers"] = orderedUsers

	x.addCommonData(templateData)

	comments := make([]*models.Comment, 0)
	models.DB.Where("task_id = ?", task.Id).Find(&comments)
	templateData["comments"] = comments

	utils.HTML(ctx, "thread.html", templateData)
}

type FormRelated struct {
	TaskId    int   `json:"task_id"`
	RelatedId int64 `json:"related_id"`
}

func (x *Task) RelatedTo(ctx *gin.Context) {
	form := &FormRelated{}
	if !utils.ShouldBindJSON(ctx, form) {
		return
	}

	task := &models.Task{}
	has, err := models.DB.ID(form.TaskId).Get(task)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if !has {
		utils.JSONErrMsg(ctx, fmt.Sprintf("未找到任务:%d", form.TaskId))
		return
	}

	related := &models.Task{Id: form.RelatedId}
	has, err = models.DB.Get(related)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if !has {
		utils.JSONErrMsg(ctx, fmt.Sprintf("未找到任务:%d", form.RelatedId))
		return
	}

	if task.Pid == 0 {
		total, err := models.DB.Where("pid = ?", task.Id).Count(new(models.Task))
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if total > 0 {
			utils.JSONError(ctx, errors.New("该任务下有子任务，不能修改关联关系"))
			return
		}
	}

	if task.Pid == form.RelatedId {
		utils.JSONError(ctx, errors.New("没有修改关联关系"))
		return
	}

	if form.RelatedId == task.Id {
		utils.JSONError(ctx, errors.New("不能关联自己"))
		return
	}

	//修正父任务关联，如果父任务有父任务，则修正为父任务的父任务
	if related.Pid > 0 {
		form.RelatedId = related.Pid
	}

	update := map[string]interface{}{
		"pid": form.RelatedId,
	}
	_, err = models.DB.Table("task").Where("id = ?", task.Id).Update(update)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"toast": "关联成功",
	})
}

func (x *Task) Related(ctx *gin.Context) {
	taskId := GetParamInt(ctx, "id")
	if taskId == 0 {
		utils.JSONErrMsg(ctx, "参数错误")
		return
	}

	task := &models.Task{Id: int64(taskId)}
	has, err := models.DB.Get(task)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if !has {
		utils.JSONErrMsg(ctx, "任务不存在")
		return
	}

	var searchId int64
	if task.Pid > 0 {
		searchId = task.Pid
	} else {
		searchId = task.Id
	}

	tasks := make([]*models.Task, 0)
	err = models.DB.Where("pid = ? or id = ?", searchId, searchId).OrderBy("pid, status, priority desc, level desc, id desc").Find(&tasks)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	if len(tasks) == 1 {

		ctx.JSON(http.StatusOK, gin.H{
			"#related-tasks": "",
		})

		return
	}

	h := gin.H{
		"tasks":      tasks,
		"pagination": &utils.Pagination{Page: 1, Total: int64(len(tasks)), Size: 10000},

		"tags":     getTags(),
		"users":    getUsers(),
		"projects": getProjects(),

		"props": map[string]any{
			"caty":       models.CatyDict,
			"status":     models.StatusDict,
			"priority":   models.PriorityDict,
			"department": models.DepartmentDict,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"#related-tasks": utils.RenderedTemplateContent(ctx, "task-related.html", h),
	})
}

func (x *Task) checkUser(ctx *gin.Context, id *int64, dept *int64, text string) bool {
	if *id > 0 {
		leader, err := models.GetUserByID(*id)
		if err != nil {
			utils.JSONError(ctx, err)
			return false
		}
		if leader == nil {
			utils.JSONErrMsg(ctx, text+"不存在")
			return false
		}

		if leader.Department == 0 {
			utils.JSONErrMsg(ctx, text+"的部门未设置")
			return false
		}

		*dept = leader.Department
	}
	return true
}

func (x *Task) checkUpdate(ctx *gin.Context, update *models.Task) bool {
	if !x.checkUser(ctx, &update.Leader, &update.LeaderDept, "负责人") {
		return false
	}

	if !x.checkUser(ctx, &update.Checker, &update.CheckerDept, "审核人") {
		return false
	}

	if !x.checkUser(ctx, &update.Tester, &update.TesterDept, "测试人") {
		return false
	}

	if update.Caty > 0 {
		if _, ok := models.CatyDict[update.Caty]; !ok {
			utils.JSONErrMsg(ctx, "分类不存在")
			return false
		}
	}

	if update.Status > 0 {
		if _, ok := models.StatusDict[update.Status]; !ok {
			utils.JSONErrMsg(ctx, "状态不存在")
			return false
		}
	}

	if update.Priority > 0 {
		if _, ok := models.PriorityDict[update.Priority]; !ok {
			utils.JSONErrMsg(ctx, "优先级不存在")
			return false
		}
	}

	if update.Project > 0 {
		project := &models.Project{Id: update.Project}
		has, err := models.DB.Exist(project)
		if err != nil {
			utils.JSONError(ctx, err)
			return false
		}
		if !has {
			utils.JSONErrMsg(ctx, "项目不存在")
			return false
		}
	}

	if update.Tag > 0 {
		tag := &models.Tag{Id: update.Tag}
		has, err := models.DB.Exist(tag)
		if err != nil {
			utils.JSONError(ctx, err)
			return false
		}
		if !has {
			utils.JSONErrMsg(ctx, "版本不存在")
			return false
		}
	}

	return true
}

type TaskLock struct {
	sync.Mutex
	LastTime int64
}

// 在文件顶部添加全局变量声明
var taskLocks sync.Map

func init() {
	// 定时清理 任务锁
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			taskLocks.Range(func(key, value interface{}) bool {
				taskLock := value.(*TaskLock)
				if time.Now().Unix()-taskLock.LastTime > 60 {
					taskLocks.Delete(key)
				}
				return true
			})
		}
	}()
}

func (x *Task) GetLock(id int64) *TaskLock {
	if v, ok := taskLocks.Load(id); ok {
		return v.(*TaskLock)
	}

	taskLock := &TaskLock{}
	taskLocks.Store(id, taskLock)
	return taskLock
}

func (x *Task) Save(ctx *gin.Context) {
	update := &models.Task{}
	if !utils.ShouldBindJSON(ctx, update) {
		return
	}

	if update.Id > 0 {
		lock := x.GetLock(update.Id)
		isLock := lock.TryLock()
		if !isLock {
			utils.JSONErrMsg(ctx, "任务正在被其他用户修改")
			return
		}
		defer lock.Unlock()

		lock.LastTime = time.Now().Unix()
	}

	authUser := ctx.MustGet("authUser").(*models.User)

	if update.Id == 0 {
		update.Author = authUser.Id
		if update.Pid > 0 {
			parentTask := &models.Task{}
			has, err := models.DB.Where("id = ?", update.Pid).Get(parentTask)
			if err != nil {
				utils.JSONError(ctx, err)
				return
			}
			if !has {
				utils.JSONErrMsg(ctx, "父任务不存在")
				return
			}
			if parentTask.Pid > 0 {
				update.Pid = parentTask.Pid //修复父任务关联
			}
		}
	}
	update.Editor = authUser.Id

	if !x.checkUpdate(ctx, update) {
		return
	}

	task := &models.Task{}
	if update.Id > 0 {
		has, err := models.DB.Where("id = ?", update.Id).Get(task)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONErrMsg(ctx, "任务不存在")
			return
		}

		tag := &models.Tag{}
		has, err = models.DB.Where("id = ?", task.Tag).Get(tag)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONErrMsg(ctx, "版本不存在")
			return
		}

		if update.Lockn != task.Lockn+1 {
			utils.JSONErrMsg(ctx, "任务已经被其他人抢先一步修改，请刷新查看")
			return
		}
	}

	if update.Id > 0 {
		_, err := models.DB.Id(update.Id).Update(update)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}

		go x.log(authUser, task, update)
		go x.feishuTaskUpdate(task, update)
	} else {
		_, err := models.DB.InsertOne(update)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}

		go x.feishuTaskNew(update)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": update.Id,
	})
}

func (x *Task) log(authUser *models.User, task, newtask *models.Task) {
	defer utils.LogRecover()

	logs := make([]*models.TaskLog, 0)

	if newtask.Title != "" {
		logs = x.logContent(authUser, logs, task.Id, task.Title, newtask.Title, "title")
	}

	if newtask.Content != "" {
		logs = x.logContent(authUser, logs, task.Id, task.Content, newtask.Content, "content")
	}

	if newtask.Status > 0 {
		name := models.GetPropsName(task.Status, models.StatusDict)
		nameTo := models.GetPropsName(newtask.Status, models.StatusDict)
		logs = x.logContent(authUser, logs, task.Id, name, nameTo, "status")
	}
	if newtask.Priority > 0 {
		name := models.GetPropsName(task.Priority, models.PriorityDict)
		nameTo := models.GetPropsName(newtask.Priority, models.PriorityDict)
		logs = x.logContent(authUser, logs, task.Id, name, nameTo, "priority")
	}

	if newtask.Leader > 0 {
		logs = x.logUser(authUser, logs, task.Id, task.Leader, newtask.Leader, "leader")
	}
	if newtask.Tester > 0 {
		logs = x.logUser(authUser, logs, task.Id, task.Tester, newtask.Tester, "tester")
	}
	if newtask.Checker > 0 {
		logs = x.logUser(authUser, logs, task.Id, task.Checker, newtask.Checker, "checker")
	}

	if newtask.Project > 0 {
		logs = x.logProject(authUser, logs, task.Id, task.Project, newtask.Project)
	}
	if newtask.Tag > 0 {
		logs = x.logTag(authUser, logs, task.Id, task.Tag, newtask.Tag)
	}

	models.DB.Insert(&logs)
}

func (x *Task) logTag(authUser *models.User, logs []*models.TaskLog, taskId int64, value, valueto int64) []*models.TaskLog {
	if value == valueto {
		return logs
	}

	tag := &models.Tag{Id: value}
	if has, err := models.DB.Get(tag); err != nil || !has {
		return logs
	}
	tagTo := &models.Tag{Id: valueto}
	if has, err := models.DB.Get(tagTo); err != nil || !has {
		return logs
	}
	return x.logContent(authUser, logs, taskId, tag.Name, tagTo.Name, "tag")
}

func (x *Task) logProject(authUser *models.User, logs []*models.TaskLog, taskId int64, value, valueto int) []*models.TaskLog {
	if value == valueto {
		return logs
	}
	project := &models.Project{Id: value}
	if has, err := models.DB.Get(project); err != nil || !has {
		return logs
	}
	projectTo := &models.Project{Id: valueto}
	if has, err := models.DB.Get(projectTo); err != nil || !has {
		return logs
	}
	return x.logContent(authUser, logs, taskId, project.Name, projectTo.Name, "project")
}

func (x *Task) logUser(authUser *models.User, logs []*models.TaskLog, taskId int64, value, valueto int64, field string) []*models.TaskLog {
	if value == valueto {
		return logs
	}

	user, _ := models.GetUserByID(value)
	userTo, _ := models.GetUserByID(valueto)
	if user == nil || userTo == nil {
		return logs
	}

	return x.logContent(authUser, logs, taskId, user.Name, userTo.Name, field)
}

var logColName = map[string]string{
	"title":    "标题",
	"content":  "内容",
	"caty":     "分类",
	"priority": "优先级",
	"status":   "状态",
	"project":  "项目",
	"tag":      "版本",
	"leader":   "负责人",
	"checker":  "验收人",
	"tester":   "测试人",
}

func (x *Task) logContent(authUser *models.User, logs []*models.TaskLog, taskId int64, value, valueto, field string) []*models.TaskLog {
	if value == valueto {
		return logs
	}

	log := &models.TaskLog{
		TaskId:   taskId,
		Operator: authUser.Name,
		Col:      logColName[field],
		Value:    value,
		ValueTo:  valueto,
	}

	return append(logs, log)
}

type BatchUpdateForm struct {
	Ids    []int        `form:"ids"`
	Update *models.Task `form:"update"`
}

func (x *Task) BatchUpdate(ctx *gin.Context) {
	form := &BatchUpdateForm{}
	// 关键修改：初始化 Update 字段
	form.Update = &models.Task{}
	form.Ids = []int{}

	if !utils.ShouldBindJSON(ctx, form) {
		return
	}

	if !x.checkUpdate(ctx, form.Update) {
		return
	}

	if len(form.Ids) == 0 {
		utils.JSONErrMsg(ctx, "请选择要修改的任务")
		return
	}

	tasks := []*models.Task{}
	err := models.DB.In("id", form.Ids).Find(&tasks)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if len(tasks) == 0 {
		utils.JSONError(ctx, errors.New("没有找到相关任务"))
		return
	}

	_, err = models.DB.In("id", form.Ids).Update(form.Update)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	authUser := GetAuthUser(ctx)
	for _, task := range tasks {
		go x.log(authUser, task, form.Update)
	}

	utils.JSONMsg(ctx, "更新成功")
}

func (x *Task) GetLog(ctx *gin.Context) {
	taskId := GetParamInt(ctx, "id")

	results := make([]*models.TaskLog, 0)
	err := models.DB.Where("task_id = ?", taskId).Find(&results)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	task := &models.Task{}
	models.DB.Id(taskId).Get(task)

	utils.Dialog(ctx, "task-log.html", gin.H{
		"task": task,
		"logs": results,
	})
}

func (x *Task) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("upload")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}

	dir, monthdir, err := utils.GetImageSaveDir()

	// 检查文件目录是否存在
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	}

	// 重命名文件
	newFileName := RenameFileWithAutoIncrement(file.Filename)
	dst := fmt.Sprintf("%s/%s", dir, newFileName)

	// 保存文件
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"path": utils.AppConfig.Server.Imgurl + monthdir + "/" + newFileName,
	})
}

func (x *Task) LogDiff(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	log := &models.TaskLog{}
	has, err := models.DB.Id(id).Get(log)
	if err != nil {
		utils.Error(ctx, err)
		return
	}
	if !has {
		utils.ErrorMsg(ctx, "没有找到相关日志")
		return
	}

	task := &models.Task{}
	models.DB.Id(log.TaskId).Get(task)

	utils.HTML(ctx, "log-diff.html", gin.H{
		"log":      log,
		"task":     task,
		"authUser": GetAuthUser(ctx),
	})
}

func (x *Task) logDiff(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")

	log := &models.TaskLog{}
	has, err := models.DB.Id(id).Get(log)
	if err != nil {
		utils.Error(ctx, err)
		return
	}
	if !has {
		utils.ErrorMsg(ctx, "没有找到相关日志")
		return
	}

	cfg := &htmldiff.Config{
		Granularity:  5,
		InsertedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: palegreen;"}},
		DeletedSpan:  []htmldiff.Attribute{{Key: "style", Val: "background-color: lightpink;"}},
		ReplacedSpan: []htmldiff.Attribute{{Key: "style", Val: "background-color: lightskyblue;"}},
		CleanTags:    []string{""},
	}
	res, err := cfg.HTMLdiff([]string{log.Value, log.ValueTo})
	if err != nil {
		utils.Error(ctx, err)
		return
	}
	mergedHTML := res[0]

	task := &models.Task{}
	models.DB.Id(log.TaskId).Get(task)

	utils.HTML(ctx, "log-diff.html", gin.H{
		"mergedHTML": mergedHTML,
		"task":       task,
	})
}

func (x *Task) feishuTaskUpdate(task *models.Task, update *models.Task) {
	defer utils.LogRecover()

	taskTitle := task.Title
	if update.Title != "" {
		taskTitle = update.Title
	}

	title := fmt.Sprintf("\\n### [#%d %s](%s/task/thread/%d)", update.Id, taskTitle, utils.AppConfig.Server.Domain, update.Id)
	//修改了任务负责人
	if update.Leader == 0 || task.Leader != update.Leader {
		return
	}

	oldLeader, _ := models.GetUserByID(task.Leader)
	if oldLeader == nil {
		return
	}
	newLeader, _ := models.GetUserByID(update.Leader)
	if newLeader == nil {
		return
	}

	content := fmt.Sprintf("<font color='red'>任务被重新分配给 %s</font>，请做好交接工作。", newLeader.Name)
	content += title
	x.feishuTixing(oldLeader.Account, content)

	content = "<font color='green'>您有一条新任务等待接收</font>"
	priority := task.Priority
	if update.Priority > 0 {
		priority = update.Priority
	}
	if priority == 99999 {
		content += "\\n**<font color='red'>此任务优先级为加急</font>**"
	}
	content += title

	x.feishuTixing(newLeader.Account, content)
}

// 新任务
func (x *Task) feishuTaskNew(update *models.Task) {
	defer utils.LogRecover()

	leader, _ := models.GetUserByID(update.Leader)
	if leader == nil {
		return
	}

	title := fmt.Sprintf("\\n### [#%d %s](%s/task/thread/%d)", update.Id, update.Title, utils.AppConfig.Server.Domain, update.Id)
	content := "<font color='green'>您有一条新任务等待接收</font>"
	priority := update.Priority
	if priority == 99999 {
		content += "\\n**<font color='red'>此任务优先级为加急</font>**"
	}
	content += title
	x.feishuTixing(leader.Account, content)
}

func (x *Task) feishuTixing(user, content string) {
	if utils.AppConfig.Feishu.ClientID == "" {
		return
	}

	feishuMsgTmp := `
{
    "schema": "2.0",
    "config": {
        "update_multi": true,
        "style": {
            "text_size": {
                "normal_v2": {
                    "default": "normal",
                    "pc": "normal",
                    "mobile": "heading"
                }
            }
        }
    },
    "body": {
        "direction": "vertical",
        "padding": "12px 12px 12px 12px",
        "elements": [
            {
                "tag": "markdown",
                "content": "%s",
                "text_align": "left",
                "text_size": "normal_v2",
                "margin": "0px 0px 0px 0px"
            }
        ]
    }
}
	`

	content = fmt.Sprintf(feishuMsgTmp, content)
	utils.FeishuAddMsgToChannel(user, content)
}
