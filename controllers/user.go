package controllers

import (
	"fmt"
	"net/http"
	"os"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type User struct {
}

func (x *User) URLPatterns() []utils.Route {
	return []utils.Route{
		{Path: "/user/list", ResourceFunc: x.List, Method: http.MethodGet},

		{Path: "/user/modify/:id", ResourceFunc: x.Modify, Method: http.MethodPost},
		{Path: "/user/save", ResourceFunc: x.Save, Method: http.MethodPost},
		{Path: "/user/deleteform/:id", ResourceFunc: x.DeleteForm, Method: http.MethodPost},
		{Path: "/user/delete", ResourceFunc: x.Delete, Method: http.MethodPost},

		{Path: "/permission/modify", ResourceFunc: x.modifyPermission, Method: http.MethodPost},
	}
}

func getUsers() map[int64]*models.User {
	results := make(map[int64]*models.User, 0)
	models.DB.OrderBy("team").Find(&results)
	return results
}

func (x *User) List(ctx *gin.Context) {
	results := make([]*models.User, 0)
	err := models.DB.OrderBy("department, team, id").Find(&results)
	if err != nil {
		utils.JSONError(ctx, err)
		return
	}

	h := ctx.MustGet("templateData").(map[string]any)
	h["users"] = results
	h["departments"] = models.DepartmentDict
	h["psGroups"] = utils.GetPathPermission().UserGroups

	utils.HTML(ctx, "users.html", nil)
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
	if !authUser.IsAdmin {
		utils.JSONErrMsg(ctx, "只有管理员有权限")
		return
	}

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

	_, err := models.DB.ID(user.Id).Update(user)
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
	if !authUser.IsAdmin {
		utils.JSONError(ctx, errors.New("不是管理员不能删除"))
		return
	}

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
