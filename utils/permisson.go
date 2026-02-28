package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// PermissionConfig 权限配置结构
type PermissionConfig struct {
	UserGroups map[int]UserGroup `yaml:"user_groups"`
}

var pathPermission *PermissionConfig

// UserGroup 用户组结构
type UserGroup struct {
	Id            int      `yaml:"id"`
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	IncludeGroups []int    `yaml:"include_groups"`
	Permissions   []string `yaml:"permissions"`
	WhiteList     []string `yaml:"white_list"`
	BlackList     []string `yaml:"black_list"`
}

// GetUserPermissions 获取指定用户组的权限
func (pc *PermissionConfig) GetUserPermissions(groupID int) []string {
	group, exists := pc.UserGroups[groupID]
	if !exists {
		return []string{}
	}
	return group.Permissions
}

// HasPermission 检查用户组是否有访问特定路径的权限
func (pc *PermissionConfig) HasPermission(groupID int, path string) bool {
	group, exists := pc.UserGroups[groupID]
	if !exists {
		return false
	}

	// 处理白名单和黑名单
	if len(group.WhiteList) > 0 {
		for _, whitePath := range group.WhiteList {
			if pc.matchPath(whitePath, path) {
				return true
			}
		}
		return false
	}
	if len(group.BlackList) > 0 {
		for _, blackPath := range group.BlackList {
			if pc.matchPath(blackPath, path) {
				return false
			}
		}
	}

	for _, permission := range group.Permissions {
		if pc.matchPath(permission, path) {
			return true
		}
	}

	for _, includeGroup := range group.IncludeGroups {
		if pc.HasPermission(includeGroup, path) {
			return true
		}
	}

	return false
}

// matchPath 匹配路径（支持通配符）
func (pc *PermissionConfig) matchPath(pattern, path string) bool {
	// 简单的通配符匹配实现
	if pattern == "**/*" {
		return true // 匹配所有路径
	}

	// 处理 /api/* 类型的模式
	if len(pattern) >= 2 && pattern[len(pattern)-2:] == "/*" {
		prefix := pattern[:len(pattern)-2] // 去掉 "/*"
		if path == prefix || len(path) > len(prefix) && path[:len(prefix)] == prefix && path[len(prefix)] == '/' {
			return true
		}
	}

	// 精确匹配
	if pattern == path {
		return true
	}

	return false
}

// LoadPermissionConfig 从YAML文件加载权限配置
func LoadPermissionConfig(filePath string) (*PermissionConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config PermissionConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %v", err)
	}

	return &config, nil
}

func init() {
	var err error
	pathPermission, err = LoadPermissionConfig("./conf/permission.yaml")
	if err != nil {
		panic(err)
	}
}

func GetPathPermission() *PermissionConfig { return pathPermission }
func SetPathPermission(p *PermissionConfig) {
	pathPermission = p
}
