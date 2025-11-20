package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var AppConfig Config

type Config struct {
	Server struct {
		Domain    string
		UploadDir string `yaml:"upload_dir"`
		Imgurl    string

		Port     string
		GIN_MODE string `yaml:"GIN_MODE"`
	} `yaml:"server"`

	Jwt struct {
		Secret    string
		ExpiresIn int `yaml:"expires_in"`
	} `yaml:"jwt"`

	Database struct {
		Host    string
		Port    int
		User    string
		Pwd     string
		Dbname  string
		ShowSql bool `yaml:"show_sql"`
	} `yaml:"mysql"`

	Feishu struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret" json:"-"`
		RedirectURI  string `yaml:"redirect_uri"`
	} `yaml:"feishu"`
}

func InitConf() {
	data, err := os.ReadFile("conf/application.yaml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		panic(err)
	}
	log.Println("AppConfig:", AppConfig)

	// 检查 UploadDir 是否有读写权限
	uploadDirInfo, err := os.Stat(AppConfig.Server.UploadDir)
	if err != nil {
		panic(err)
	}
	// 检查权限位
	mode := uploadDirInfo.Mode()
	if mode&0200 == 0 {
		panic("UploadDir has no write permission")
	}
	if mode&0400 == 0 {
		panic("UploadDir has no read permission")
	}
}
