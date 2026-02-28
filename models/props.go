package models

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// HasId 接口定义了一个 GetId 方法，用于获取 Id 字段的值
type IProps interface {
	GetId() int64
	GetName() string
}

type Props struct {
	Id        int64  `json:"id" yaml:"id"`
	Name      string `json:"name" yaml:"name"`
	Alias     string `json:"alias" yaml:"alias"`
	ShortName string `json:"short_name" yaml:"short_name"`
}

func (x *Props) GetId() int64    { return x.Id }
func (x *Props) GetName() string { return x.Name }

var CatyList = []*Props{}
var StatusList = []*Props{}
var PriorityList = []*Props{}
var DepartmentList = []*Props{}

var CatyDict = map[int64]*Props{}
var StatusDict = map[int64]*Props{}
var PriorityDict = map[int64]*Props{}
var DepartmentDict = map[int64]*Props{}

func init() {
	loadPropsYAML("conf/caty.yaml", &CatyList)
	loadPropsYAML("conf/status.yaml", &StatusList)
	sss, _ := json.Marshal(StatusList)
	fmt.Println(string(sss))
	loadPropsYAML("conf/priority.yaml", &PriorityList)
	loadPropsYAML("conf/department.yaml", &DepartmentList)

	propsFillDict(CatyList, CatyDict)
	propsFillDict(StatusList, StatusDict)
	propsFillDict(PriorityList, PriorityDict)
	propsFillDict(DepartmentList, DepartmentDict)
}

// 泛型方法实现
func propsFillDict[T IProps](list []T, dict map[int64]T) {
	for _, v := range list {
		id := v.GetId()
		dict[id] = v
	}
}

func loadPropsYAML(filePath string, target interface{}) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read file %s: %v", filePath, err))
	}
	if err := yaml.Unmarshal(data, target); err != nil {
		panic(fmt.Sprintf("failed to unmarshal YAML from file %s: %v", filePath, err))
	}
}

func GetPropsName[T IProps](id int64, dict map[int64]T) string {
	item, ok := dict[id]
	if ok {
		return item.GetName()
	}
	return ""
}
