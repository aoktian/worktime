-- worktime.comment definition

CREATE TABLE `comment` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `content` text COMMENT '内容',
  `author` varchar(255) NOT NULL,
  `editor` varchar(255) NOT NULL,
  `created_at` bigint NOT NULL COMMENT '创建时间',
  `updated_at` bigint NOT NULL COMMENT '更新时间',
  `task_id` int NOT NULL COMMENT '任务id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb3;


-- worktime.comment_log definition

CREATE TABLE `comment_log` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `comment_id` int DEFAULT NULL,
  `operator` varchar(255) NOT NULL COMMENT '操作人',
  `created_at` bigint NOT NULL COMMENT '创建时间',
  `content` text COMMENT '内容',
  `content_to` text COMMENT '内容',
  PRIMARY KEY (`id`),
  KEY `IDX_feed_log_feed_id` (`comment_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb3;


-- worktime.project definition

CREATE TABLE `project` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '项目名称',
  `created_at` bigint NOT NULL COMMENT '创建时间',
  `updated_at` bigint NOT NULL COMMENT '更新时间',
  `is_archived` tinyint NOT NULL DEFAULT '0',
  `paixu` int NOT NULL DEFAULT '0',
  `start_at` bigint NOT NULL DEFAULT '0',
  `end_at` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb3;


-- worktime.tag definition

CREATE TABLE `tag` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '标签名称',
  `project_id` int NOT NULL COMMENT '项目id',
  `start_at` bigint DEFAULT NULL COMMENT '开始时间',
  `end_at` bigint DEFAULT NULL COMMENT '结束时间',
  `created_at` bigint DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint DEFAULT NULL COMMENT '更新时间',
  `paixu` int NOT NULL DEFAULT '0',
  `is_archived` tinyint NOT NULL DEFAULT '0' COMMENT '是否归档',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb3;


-- worktime.task definition

CREATE TABLE `task` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(255) NOT NULL COMMENT '任务标题',
  `content` text COMMENT '任务内容',
  `author` int NOT NULL COMMENT '创建人id',
  `editor` int NOT NULL COMMENT '编辑人id',
  `created_at` bigint NOT NULL COMMENT '创建时间',
  `updated_at` bigint NOT NULL COMMENT '更新时间',
  `project` int NOT NULL COMMENT '项目id',
  `tag` int NOT NULL COMMENT '标签id',
  `leader` int NOT NULL COMMENT '负责人id',
  `checker` int NOT NULL COMMENT '审核人员id',
  `tester` int NOT NULL COMMENT '测试人员id',
  `caty` int NOT NULL COMMENT '分类',
  `status` int NOT NULL COMMENT '状态',
  `priority` int NOT NULL COMMENT '优先级',
  `level` int NOT NULL COMMENT '级别',
  `pid` int NOT NULL DEFAULT '0' COMMENT '父任务id',
  `lockn` int NOT NULL DEFAULT '0' COMMENT '锁号',
  `leader_dept` int NOT NULL COMMENT '负责人部门id',
  `checker_dept` int NOT NULL COMMENT '审核人员部门id',
  `tester_dept` int NOT NULL COMMENT '测试人员部门id',
  `start_at` bigint DEFAULT NULL COMMENT '开始时间',
  `end_at` bigint DEFAULT NULL COMMENT '结束时间',
  `actual_at` int NOT NULL DEFAULT '0',
  INDEX `IDX_task_project` (`project`),
  INDEX `IDX_task_tag` (`tag`),
  INDEX `IDX_task_author` (`author`),
  INDEX `IDX_task_leader` (`leader`),
  INDEX `IDX_task_checker` (`checker`),
  INDEX `IDX_task_tester` (`tester`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb3;


-- worktime.task_log definition

CREATE TABLE `task_log` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int NOT NULL COMMENT '任务id',
  `operator` varchar(255) NOT NULL COMMENT '操作人',
  `col` varchar(255) NOT NULL COMMENT '字段名',
  `value` text NOT NULL COMMENT '值',
  `value_to` text NOT NULL COMMENT '值',
  `created_at` bigint NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `IDX_task_log_task_id` (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb3;


-- worktime.`user` definition

CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account` varchar(255) NOT NULL COMMENT '账号',
  `name` varchar(255) NOT NULL COMMENT '姓名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `department` int NOT NULL COMMENT '部门',
  `team` int NOT NULL COMMENT '用户组',
  `is_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否管理员',
  `created_at` bigint NOT NULL COMMENT '创建时间',
  `updated_at` bigint NOT NULL COMMENT '更新时间',
  `ps` int NOT NULL DEFAULT '0',
  `nick` varchar(255) NOT NULL DEFAULT '',
  `is_leave` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_user_username` (`name`),
  UNIQUE KEY `UQE_user_account` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb3;
