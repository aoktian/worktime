package models

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// HasId 接口定义了一个 GetId 方法，用于获取 Id 字段的值
type HasId interface {
	GetId() int
}

type Props struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (p *Props) GetId() int {
	return p.Id
}

var CatyList = []*Props{}
var StatusList = []*Props{}
var PriorityList = []*Props{}
var DepartmentList = []*Props{}

var CatyDict = map[int]*Props{}
var StatusDict = map[int]*Props{}
var PriorityDict = map[int]*Props{}
var DepartmentDict = map[int]*Props{}

func (props *Props) InitFromConf() {
	props.loadYAML("conf/caty.yaml", &CatyList)
	props.loadYAML("conf/status.yaml", &StatusList)
	props.loadYAML("conf/priority.yaml", &PriorityList)
	props.loadYAML("conf/department.yaml", &DepartmentList)

	props.fillDict(CatyList, &CatyDict)
	props.fillDict(StatusList, &StatusDict)
	props.fillDict(PriorityList, &PriorityDict)
	props.fillDict(DepartmentList, &DepartmentDict)
}

func (props *Props) fillDict(list []*Props, dict *map[int]*Props) {
	for _, v := range list {
		(*dict)[v.Id] = v
	}
}

func (props *Props) loadYAML(filePath string, target interface{}) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read file %s: %v", filePath, err))
	}
	if err := yaml.Unmarshal(data, target); err != nil {
		panic(fmt.Sprintf("failed to unmarshal YAML from file %s: %v", filePath, err))
	}
}

func (props *Props) GetName(id int, dict *map[int]*Props) string {
	props, ok := (*dict)[id]
	if ok {
		return props.Name
	}
	return ""
}
