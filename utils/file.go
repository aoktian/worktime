package utils

import (
	"fmt"
	"os"
	"time"
)

// 检查文件目录是否存在
func CheckDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetImageSaveDir() (string, string, error) {
	now := time.Now()
	monthdir := fmt.Sprintf("%d%02d", now.Year(), now.Month())

	dir := AppConfig.Server.UploadDir + "/" + monthdir

	// 检查文件目录是否存在
	return dir, monthdir, CheckDir(dir)
}
