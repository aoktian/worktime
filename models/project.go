package models

type Project struct {
	Id   int    `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	Name string `json:"name" xorm:"notnull comment('项目名称')"`

	Paixu      int  `json:"paixu"`
	IsArchived bool `json:"is_archived"`

	StartAt int64 `json:"start_at" xorm:"comment('开始时间')"`
	EndAt   int64 `json:"end_at" xorm:"comment('结束时间')"`

	CreatedAt int64 `json:"created_at" xorm:"notnull comment('创建时间') created"`
	UpdatedAt int64 `json:"updated_at" xorm:"notnull comment('更新时间') updated"`
}

func (p *Project) GetId() int { return p.Id }
