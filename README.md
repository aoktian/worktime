# WORKTIME 项目管理和问题跟踪工具

简单、好用、高效。
广泛使用、备受好评的项目管理和问题跟踪工具，提供可定制的工作流程和广泛的集成能力。

WORKTIME 是一个现代化的项目管理和问题跟踪工具，旨在帮助团队更好地管理和追踪工作任务。系统采用前后端分离架构，提供直观的用户界面和强大的后端支持。

## 编译
windows 和 linux 的可执行文件都可以在windows上生成
```shell
go build -o worktime.exe

SET GOPROXY=https://goproxy.cn
set GOARCH=amd64
set GOOS=linux

go build -o worktime

pause
```

## 开发编译参数

开发时候，可以不把静态文件打包到可执行文件。
方便调试静态文件 。

```shell
go build -tags="dev" -o worktime.exe
```

## 服务器后台运行
```shell
nohup ./worktime start > /dev/null 2>&1 &
```

## 管理脚本
```shell
chmod +x service.sh
./service.sh start # 启动
./service.sh stop # 停止
./service.sh restart # 重启
./service.sh status # 状态
```


## 技术栈

### 后端
- 语言：Go
- 配置管理：YAML

### 前端
- 语言：HTML + CSS + JavaScript
- 库：JQuery + Bootstrap


## 核心功能

### 任务管理
- 任务创建和编辑
- 任务状态跟踪
- 任务优先级管理
- 任务分类系统
- 任务日志记录

### 用户管理
- 用户认证和授权
- 部门管理
- 角色分配（负责人、审核人、测试人员等）
- 集成飞书授权登录

### 项目管理
- 项目创建和管理
- 版本标签管理
- 项目进度追踪

### 系统特性
- 多级任务支持（父子任务关系）
- 完整的任务生命周期管理
- 详细的操作日志
- 灵活的配置系统

## 运行环境

```
.
├── worktime                  # linux运行文件
├── worktime.exe              # windows运行文件
└── conf/                   # 配置文件，自行修改
    ├── application.yaml    # 主配置文件
    ├── caty.yaml           # 分类
    ├── department.yaml     # 部门
    ├── priority.yaml       # 优先级
    └── status.yaml         # 状态
```

## 视频教程

1. [介绍](https://www.bilibili.com/video/BV13sjZzgEEa/?share_source=copy_web&vd_source=42766cd92882fca8c755bb74903c2aa8)
2. [安装部署]( https://www.bilibili.com/video/BV1StjZzTEzY/?share_source=copy_web&vd_source=42766cd92882fca8c755bb74903c2aa8)

3. 接入飞书授权登录

4. 最佳实践 - 版本管理
5. 最佳实践 - 离职人员的处理
6. 最佳实践 - 上传图片的存储
7. 最佳实践 - 负责人的使用方法
8. 最佳实践 - 验收人的使用方法
9. 最佳实践 - 品控人的使用方法
10. 最佳实践 - 管理者的使用方法


## 许可证

请查看 [LICENSE](LICENSE) 文件了解详细信息。