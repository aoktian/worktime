package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	// 防止时区错误
	_ "time/tzdata"

	"github.com/gin-gonic/gin"

	"webserver/controllers"
	"webserver/middlewares"
	"webserver/models"
	"webserver/utils"

	"github.com/sirupsen/logrus"
)

type cmd struct {
	desc string
	run  func()
}

var cmds = map[string]cmd{
	"start":  {"启动服务", startWebServer},
	"initDb": {"初始化数据库", initDb},
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

func initDb() {
	models.CreateTables()
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

	loadTemplatesAndStatic(r)

	r.Use(middlewares.EnableCookieSession())

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
