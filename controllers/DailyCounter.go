package controllers

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

// 全局计数器结构体
type DailyCounter struct {
	counter int
	date    string
	mu      sync.Mutex
}

var dailyCounter = &DailyCounter{
	counter: 0,
	date:    time.Now().Format("20060102"),
}

func (dc *DailyCounter) GetNextID() int {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	// 如果日期变化，重置计数器
	if dc.date != time.Now().Format("20060102") {
		dc.counter = 0
		dc.date = time.Now().Format("20060102")
	}

	dc.counter++
	return dc.counter
}

// RenameFileWithAutoIncrement 返回新文件名，格式为 {date}_{id}{ext}
func RenameFileWithAutoIncrement(originalName string) string {
	ext := filepath.Ext(originalName) // 获取扩展名

	today := time.Now().Format("20060102") // 当前日期
	id := dailyCounter.GetNextID()         // 获取当日递增ID

	return fmt.Sprintf("%s_%d%s", today, id, ext)
}
