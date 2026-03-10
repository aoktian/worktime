package models

import (
	"fmt"
	"webserver/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DB *xorm.Engine

func ConnectDatabase() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=%s",
		utils.AppConfig.Database.User, utils.AppConfig.Database.Pwd,
		utils.AppConfig.Database.Host, utils.AppConfig.Database.Port,
		utils.AppConfig.Database.Dbname, "Asia%2fShanghai",
	)
	fmt.Println(utils.AppConfig.Database)
	database, err := xorm.NewEngine("mysql", dataSource)

	if err != nil {
		panic(err)
	}

	if err := database.Ping(); err != nil {
		panic(err)
	}

	DB = database

	// 使用与 gin 相同的日志输出，确保 XORM 日志也写入文件和控制台
	multiWriter := utils.GetMultiWriter()
	xormLogger := xorm.NewSimpleLogger(multiWriter)
	xormLogger.ShowSQL(utils.AppConfig.Database.ShowSql)
	DB.ShowSQL(utils.AppConfig.Database.ShowSql)
	DB.SetLogger(xormLogger)
}

func CreateTables() {
	// 迁移数据表
	err := DB.Sync2(
		new(User),
		new(Project), new(Tag),
		new(Task), new(TaskLog), new(Comment), new(CommentLog),
	)
	if err != nil {
		panic(err)
	}

	DB.Exec("ALTER TABLE task AUTO_INCREMENT=10000")
	DB.Exec("ALTER TABLE user AUTO_INCREMENT=100")
	DB.Exec("ALTER TABLE project AUTO_INCREMENT=10")
	DB.Exec("ALTER TABLE tag AUTO_INCREMENT=100")
}
