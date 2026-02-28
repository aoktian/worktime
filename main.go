package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	// 防止时区错误
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"webserver/controllers"
	"webserver/middlewares"
	"webserver/models"
	"webserver/templatefuncs"
	"webserver/utils"

	"github.com/sirupsen/logrus"
)

type cmd struct {
	desc string
	run  func()
}

var cmds = map[string]cmd{
	"start":         {"启动服务", startWebServer},
	"addUser":       {"添加用户", addUser},
	"setAdmin":      {"设置管理员", setAdmin},
	"resetPassword": {"重置密码", resetPassword},
	"createTables":  {"初始化数据库", createTables},
}

func help() {
	fmt.Println("输入正确的命令索引")
	for name, cmd := range cmds {
		fmt.Printf("%s: ./worktime %s\n", cmd.desc, name)
	}
}

func main() {
	utils.InitConf()

	models.ConnectDatabase()

	cmdindex := ""
	if len(os.Args) > 1 {
		cmdindex = os.Args[1]
	}

	if cmdindex == "" {
		help()
		return
	}

	cmd, ok := cmds[cmdindex]
	if !ok {
		help()
		return
	}

	cmd.run()
}

func createTables() {
	models.CreateTables()
}

func addUser() {
	u := &models.User{}

	for {
		fmt.Print("请输入登录帐号：")
		fmt.Scanln(&u.Account)

		has, err := models.DB.Where("account = ?", u.Account).Get(new(models.User))
		if err != nil {
			panic(err)
		}
		if has {
			fmt.Println("帐号已存在")
		} else {
			break
		}
	}

	fmt.Print("请输入登录密码：")
	fmt.Scanln(&u.Password)

	fmt.Print("请输入姓名或者昵称：")
	fmt.Scanln(&u.Name)

	fmt.Println("请输入部门编号：")
	for _, d := range models.DepartmentList {
		fmt.Printf("%d:%s\n", d.Id, d.Name)
	}
	fmt.Scanln(&u.Department)

	fmt.Print("是否是管理员(y/n)：")
	var isAdmin string
	fmt.Scanln(&isAdmin)
	if isAdmin == "y" {
		u.IsAdmin = true
	}

	err := u.SaveUser()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("用户添加成功")

	fmt.Print("是否继续添加用户(y/n)：")
	var isAdd string
	fmt.Scanln(&isAdd)
	if isAdd == "y" {
		addUser()
	} else {
		fmt.Println("已退出")
	}
}

func setAdmin() {
	id := 0

	fmt.Print("请输入ID：")
	for {
		fmt.Scanln(&id)
		has := models.CheckUserExist(id)
		if !has {
			fmt.Println("用户不存在，请重新输入：")
			continue
		} else {
			break
		}
	}

	fmt.Println("是否管理员？y/n")
	var isAdmin string
	fmt.Scanln(&isAdmin)

	is := false
	if isAdmin == "y" {
		is = true
	}

	update := &models.User{IsAdmin: is}
	affected, err := models.DB.ID(id).MustCols("is_admin").Update(update)

	if err != nil {
		panic(err)
	}

	if affected > 0 {
		fmt.Println("修改成功")
	} else {
		fmt.Println("未修改")
	}
}

func resetPassword() {
	id := 0

	fmt.Print("请输入ID：")
	for {
		fmt.Scanln(&id)
		has := models.CheckUserExist(id)
		if !has {
			fmt.Println("用户不存在，请重新输入：")
			continue
		} else {
			break
		}
	}

	password := utils.RandNumeric(8)

	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	update := &models.User{Password: string(b)}
	affected, err := models.DB.ID(id).Update(update)

	if err != nil {
		panic(err)
	}
	if affected > 0 {
		fmt.Println("修改成功，新密码: ", password)
	} else {
		fmt.Println("未修改")
	}
}

func startWebServer() {
	log := utils.GetLogger()
	gin.DefaultWriter = log.Out

	r := gin.Default()
	utils.SetGlobalRouter(r) // 设置全局 router

	// 使用自定义 Recovery 中间件以确保 panic 被记录
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"client": c.ClientIP(),
			"stack":  string(debug.Stack()),
		}).Errorf("Panic recovered: %v", recovered)

		utils.ErrorMsg(c, fmt.Sprintf("%v", recovered))
		c.Abort()
	}))

	templatefuncs.SetTemplateFuncs(r) //要在load 之前调用
	r.LoadHTMLGlob("templates/*")

	r.Use(middlewares.EnableCookieSession())

	// 前端项目静态资源
	r.Static("/static", "./static")
	r.StaticFile("/worktime.svg", "./static/worktime.svg")

	// 上传的图片
	r.Static("/upload", utils.AppConfig.Server.UploadDir)

	// 使用中间件处理跨域问题
	// r.Use(middlewares.CORSMiddleware())

	defaultController := &controllers.Task{}
	r.Handle(http.MethodGet, "/", middlewares.AuthSessionMiddleware(), defaultController.Ileader)

	r.Use(middlewares.ErrorHandler())

	public := r.Group("/auth")
	{
		register(public, &controllers.Auth{})
		register(public, &controllers.FeishuAuth{})
	}

	protected := r.Group("/")
	{
		protected.Use(middlewares.AuthSessionMiddleware())

		register(protected, &controllers.Help{})
		register(protected, &controllers.Tag{})
		register(protected, &controllers.User{})
		register(protected, defaultController)
		register(protected, &controllers.Comment{})
		register(protected, &controllers.Props{})
		register(protected, &controllers.Project{})
	}

	// 处理404请求并重定向到index.html
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "404.html", nil)
	})

	gin.SetMode(utils.AppConfig.Server.GIN_MODE)

	r.Run("0.0.0.0:" + utils.AppConfig.Server.Port)
}

func register(group *gin.RouterGroup, router utils.Router) {
	for _, route := range router.URLPatterns() {
		group.Handle(route.Method, route.Path, route.ResourceFunc)
	}
}
