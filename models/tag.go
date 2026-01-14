package models

type Tag struct {
	Id        int    `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	Name      string `json:"name" xorm:"notnull comment('版本名称')"`
	ProjectId int    `json:"project_id" xorm:"notnull comment('项目id')"`
	Paixu     int    `json:"paixu"`

	StartAt int64 `json:"start_at" xorm:"comment('开始时间')"`
	EndAt   int64 `json:"end_at" xorm:"comment('结束时间')"`

	CreatedAt int64 `json:"created_at" xorm:"comment('创建时间') created"`
	UpdatedAt int64 `json:"updated_at" xorm:"comment('更新时间') updated"`
}

func (x *Tag) GetId() int {
	return x.Id
}
