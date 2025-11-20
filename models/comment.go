package models

type Comment struct {
	Id      int    `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	Content string `json:"content" xorm:"text comment('内容')"`

	Author string `json:"author" xorm:"notnull comment('作者')"`
	Editor string `json:"editor" xorm:"notnull comment('编辑')"`

	CreatedAt int64 `json:"created_at" xorm:"notnull comment('创建时间') created"`
	UpdatedAt int64 `json:"updated_at" xorm:"notnull comment('更新时间') updated"`

	TaskId int `json:"task_id" xorm:"notnull comment('任务id')"`
}

func (x *Comment) GetId() int { return x.Id }

type CommentLog struct {
	Id        int    `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	CommentId int    `json:"comment_id" xorm:"notnull comment('日志id') index"`
	Operator  string `json:"operator" xorm:"notnull comment('操作人')"`
	CreatedAt int64  `json:"created_at" xorm:"notnull comment('创建时间') created"`

	Content   string `json:"content" xorm:"text comment('内容')"`
	ContentTo string `json:"content_to" xorm:"text comment('内容')"`
}
