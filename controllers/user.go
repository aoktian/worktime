package controllers

import (
	"fmt"
	"net/http"
	"os"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/builder"
	"github.com/pkg/errors"
)

type User struct {
}

func (x *User) URLPatterns() []utils.Route {
	return []utils.Route{
		{Path: "/user/list", ResourceFunc: x.index, Method: http.MethodGet},
		{Path: "/user/search", ResourceFunc: x.search, Method: http.MethodPost},
		{Path: "/user/selector", ResourceFunc: x.selector, Method: http.MethodPost},

		{Path: "/user/modify/:id", ResourceFunc: x.Modify, Method: http.MethodPost},
		{Path: "/user/save", ResourceFunc: x.Save, Method: http.MethodPost},
		{Path: "/user/deleteform/:id", ResourceFunc: x.DeleteForm, Method: http.MethodPost},
		{Path: "/user/delete", ResourceFunc: x.Delete, Method: http.MethodPost},

		{Path: "/user/showimodify", ResourceFunc: x.Imodify, Method: http.MethodPost},
		{Path: "/user/imodify", ResourceFunc: x.imodify, Method: http.MethodPost},

		{Path: "/permission/modify", ResourceFunc: x.modifyPermission, Method: http.MethodPost},
	}
}

func getUsers() map[int64]*models.User {
	results := make(map[int64]*models.User, 0)
	models.DB.OrderBy("team").Find(&results)
	return results
}

func (x *User) selector(ctx *gin.Context) {
	params := &selectorReq{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}

	pageSize := 20
	offset := (params.Page - 1) * pageSize

	results := make([]*models.User, 0)
	if params.Keyword == "" {
		models.DB.Where("department = ? and is_leave = 0", params.ProjectId).
			OrderBy("team desc").
			Limit(pageSize, offset).Find(&results)
	} else {
		models.DB.Where("department = ? and name like ?", params.ProjectId, "%"+params.Keyword+"%").
			OrderBy("team desc").
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

type userSearch struct {
	models.User
	Page        int64 `json:"page"`
	SearchState int   `json:"search_state"`
	CreatedAt0  int64 `json:"created_at0"`
	CreatedAt9  int64 `json:"created_at9"`
}

func (x *User) index(ctx *gin.Context) {
	x.list(ctx, &userSearch{})
	utils.HTML(ctx, "users.html", nil)
}

func (x *User) search(ctx *gin.Context) {
	params := &userSearch{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}
	x.list(ctx, params)
	ctx.JSON(http.StatusOK, gin.H{
		"#data-table": utils.GetRenderedTemplateContent(ctx, "users-content.html"),
	})
}

func (x *User) list(ctx *gin.Context, con *userSearch) {
	where := builder.NewCond()
	if con.Department > 0 {
		where = where.And(builder.Eq{"department": con.Department})
	}
	if con.Name != "" {
		where = where.And(builder.Like{"name", con.Name})
	}
	if con.Nick != "" {
		where = where.And(builder.Like{"nick", con.Nick})
	}
	if con.CreatedAt0 > 0 {
		where = where.And(builder.Gte{"created_at": con.CreatedAt0})
	}
	if con.CreatedAt9 > 0 {
		where = where.And(builder.Lte{"created_at": con.CreatedAt9})
	}
	if con.SearchState > -1 {
		where = where.And(builder.Eq{"is_leave": con.SearchState})
	}

	sqlWhere, args, err := builder.ToSQL(where)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	total, err := models.DB.Where(sqlWhere, args...).Count(new(models.User))
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	results := make([]*models.User, 0)

	var pageSize int64 = 12
	pagination := &utils.Pagination{Page: con.Page, Size: pageSize, Total: total}
	pagination.FormatPage()

	offset := pagination.GetOffset()
	err = models.DB.Where(sqlWhere, args...).
		OrderBy("is_leave, department, team, id").
		Limit(int(pageSize), int(offset)).
		Find(&results)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["users"] = results
	h["departments"] = models.DepartmentDict
	h["psGroups"] = utils.GetPathPermission().UserGroups
	h["pagination"] = pagination
}

func (x *User) Imodify(ctx *gin.Context) {
	user := GetAuthUser(ctx)
	h := ctx.MustGet("templateData").(map[string]any)
	h["user"] = user
	h["departments"] = models.DepartmentDict
	h["psGroups"] = utils.GetPathPermission().UserGroups

	ctx.JSON(http.StatusOK, gin.H{
		"dialog": utils.GetRenderedTemplateContent(ctx, "user-imodify.html"),
	})
}

func (x *User) imodify(c *gin.Context) {
	params := &models.User{}
	if !utils.ShouldBindJSON(c, params) {
		return
	}

	authUser := GetAuthUser(c)

	user := &models.User{
		Nick: params.Nick,
	}

	if params.Password != "" {
		hashedPassword, err := params.GenerateFromPassword()
		if err != nil {
			utils.JSONError(c, err)
			return
		}
		user.Password = string(hashedPassword)
	}

	_, err := models.DB.ID(authUser.Id).Update(user)
	if err != nil {
		utils.JSONError(c, err)
		return
	}
	utils.JSONMsg(c, "修改成功。")
}

func (x *User) Modify(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	user := &models.User{}
	if id > 0 {
		has, err := models.DB.ID(id).Get(user)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONError(ctx, errors.New("用户不存在"))
			return
		}
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["user"] = user
	h["departments"] = models.DepartmentDict
	h["psGroups"] = utils.GetPathPermission().UserGroups

	ctx.JSON(http.StatusOK, gin.H{
		"dialog": utils.GetRenderedTemplateContent(ctx, "user-edit.html"),
	})
}

func (x *User) Save(ctx *gin.Context) {
	user := &models.User{}
	if !utils.ShouldBindJSON(ctx, user) {
		return
	}

	// 保持上下的沟通，故意这么设计
	authUser := ctx.MustGet("authUser").(*models.User)

	if user.Department == 0 {
		utils.JSONErrMsg(ctx, "部门不能为空")
		return
	}

	//重置密码
	if user.Password != "" {
		hashedPassword, err := user.GenerateFromPassword()
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		user.Password = string(hashedPassword)
	}

	if user.Id == 0 {
		_, err := models.DB.InsertOne(user)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id": user.Id,
		})
		return
	}

	target, _ := models.GetUserByID(user.Id)
	if target == nil {
		utils.JSONErrMsg(ctx, "用户不存在")
		return
	}

	if target.Ps != user.Ps || target.Name != user.Name {
		if !utils.GetPathPermission().HasPermission(authUser.Ps, "/permission/modify") {
			utils.JSONErrMsg(ctx, "没有权限")
			return
		}
	}

	session := models.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		utils.JSONError(ctx, err)
		return
	}

	if target.Department != user.Department {
		//部门发生了变化
		_, err := session.Exec("UPDATE task SET leader_dept = ? WHERE leader = ?", user.Department, user.Id)
		if err != nil {
			session.Rollback()
			utils.JSONError(ctx, err)
			return
		}
	}

	_, err := models.DB.ID(user.Id).MustCols("is_leave").Update(user)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": user.Id,
	})
}

type DelForm struct {
	Id      int64 `json:"id"`
	Reciver int   `json:"reciver"`
}

func (x *User) DeleteForm(ctx *gin.Context) {
	id := GetParamInt(ctx, "id")
	user := &models.User{}
	if id > 0 {
		has, err := models.DB.ID(id).Get(user)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}
		if !has {
			utils.JSONError(ctx, errors.New("用户不存在"))
			return
		}
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["user"] = user
	h["departments"] = models.DepartmentDict

	users := getUsers()
	delete(users, user.Id)
	h["users"] = users

	utils.Dialog(ctx, "user-delete.html", nil)
}

func (x *User) Delete(ctx *gin.Context) {
	authUser := GetAuthUser(ctx)

	form := &DelForm{}
	if !utils.ShouldBindJSON(ctx, form) {
		return
	}

	if authUser.Id == form.Id {
		utils.JSONError(ctx, errors.New("不能删除自己"))
		return
	}

	if form.Reciver == 0 {
		utils.JSONError(ctx, errors.New("请选择接收人"))
		return
	}

	user := &models.User{}
	has, err := models.DB.ID(form.Id).Get(user)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if !has {
		utils.JSONError(ctx, errors.New("用户不存在"))
		return
	}

	receiver := &models.User{}
	has, err = models.DB.ID(form.Reciver).Get(receiver)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}
	if !has {
		utils.JSONError(ctx, errors.New("接收人不存在"))
		return
	}

	if user.Id == receiver.Id {
		utils.JSONError(ctx, errors.New("接收人不能是自己"))
		return
	}

	session := models.DB.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		utils.JSONError(ctx, err)
		return
	}

	// 辅助函数用于更新任务表中的字段
	updateTaskField := func(field string) error {
		_, err := session.Exec(fmt.Sprintf("UPDATE task SET %s = ? WHERE %s = ?", field, field), receiver.Id, user.Id)
		return err
	}

	fieldsToUpdate := []string{"author", "editor"}
	for _, field := range fieldsToUpdate {
		if err := updateTaskField(field); err != nil {
			session.Rollback()
			utils.JSONError(ctx, err)
			return
		}
	}

	// 单独处理 leader 和 department 字段
	_, err = session.Exec("UPDATE task SET leader = ?, leader_dept = ? WHERE leader = ?", receiver.Id, receiver.Department, user.Id)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}
	_, err = session.Exec("UPDATE task SET checker = ?, checker_dept = ? WHERE checker = ?", receiver.Id, receiver.Department, user.Id)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}
	_, err = session.Exec("UPDATE task SET tester = ?, tester_dept = ? WHERE tester = ?", receiver.Id, receiver.Department, user.Id)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}

	_, err = session.Exec("UPDATE comment SET author = ? WHERE author = ?", receiver.Id, user.Id)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}
	_, err = session.Exec("UPDATE comment SET editor = ? WHERE editor = ?", receiver.Id, user.Id)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}

	_, err = session.ID(user.Id).Delete(user)
	if err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		utils.JSONError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": user.Id,
	})
}

func (x *User) modifyPermission(ctx *gin.Context) {
	params := &models.User{}
	if !utils.ShouldBindJSON(ctx, params) {
		return
	}

	if params.Name == "" {
		data, err := os.ReadFile("./conf/permission.yaml")
		if err != nil {
			utils.Error(ctx, err)
			return
		}
		utils.Dialog(ctx, "permission.html", gin.H{
			"data": string(data),
		})
		return
	}

	//写入到文件
	err := os.WriteFile("./conf/permission.yaml", []byte(params.Name), 0644)
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	newPermission, err := utils.LoadPermissionConfig("./conf/permission.yaml")
	if err != nil {
		utils.Error(ctx, err)
		return
	}

	utils.SetPathPermission(newPermission)

	utils.JSONMsg(ctx, "修改成功")
}
