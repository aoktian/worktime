package models

type Task struct {
	Id      int64  `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	Title   string `json:"title" xorm:"notnull comment('任务标题')"`
	Content string `json:"content" xorm:"text comment('任务内容')"`

	Author int64 `json:"author" xorm:"notnull comment('创建人id')"`
	Editor int64 `json:"editor" xorm:"notnull comment('编辑人id')"`

	CreatedAt int64 `json:"created_at" xorm:"notnull comment('创建时间') created"`
	UpdatedAt int64 `json:"updated_at" xorm:"notnull comment('更新时间') updated"`

	Project int   `json:"project" xorm:"notnull comment('项目id')"`
	Tag     int64 `json:"tag" xorm:"notnull comment('版本id')"`

	Leader      int64 `json:"leader" xorm:"notnull comment('负责人id')"`
	LeaderDept  int64 `json:"leader_dept" xorm:"notnull comment('负责人部门id')"`
	Checker     int64 `json:"checker" xorm:"notnull comment('审核人员id')"`
	CheckerDept int64 `json:"checker_dept" xorm:"notnull comment('审核人员部门id')"`
	Tester      int64 `json:"tester" xorm:"notnull comment('测试人员id')"`
	TesterDept  int64 `json:"tester_dept" xorm:"notnull comment('测试人员部门id')"`

	Caty     int64 `json:"caty" xorm:"notnull comment('分类')"`
	Status   int64 `json:"status" xorm:"notnull comment('状态')"`
	Priority int64 `json:"priority" xorm:"notnull comment('优先级')"`
	Level    int64 `json:"level" xorm:"notnull comment('级别')"`

	Pid int64 `json:"pid" xorm:"notnull default 0 comment('父任务id')"`

	Lockn int `json:"lockn" xorm:"notnull default 0 comment('锁号')"`

	StartAt  int64 `json:"start_at" xorm:"comment('开始时间')"`
	EndAt    int64 `json:"end_at" xorm:"comment('结束时间')"`
	ActualAt int64 `json:"actual_at" xorm:"comment('实际完成时间')"`
}

type TaskLog struct {
	Id     int64 `json:"id" xorm:"pk autoincr notnull comment('主键')"`
	TaskId int64 `json:"task_id" xorm:"notnull comment('任务id') index"`

	Operator string `json:"operator" xorm:"notnull comment('操作人')"`

	Col     string `json:"col" xorm:"notnull comment('字段名')"`
	Value   string `json:"value" xorm:"text notnull comment('值')"`
	ValueTo string `json:"value_to" xorm:"text notnull comment('值')"`

	CreatedAt int64 `json:"created_at" xorm:"notnull comment('创建时间') created"`
}
