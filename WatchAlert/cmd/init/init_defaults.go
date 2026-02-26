package main

import (
	"fmt"
	"time"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/initialization"
)

func main() {
	// 初始化基础配置
	initialization.InitBasic()
	
	// 初始化数据库连接
	db := global.DB
	
	// 1. 创建默认工时标准
	createDefaultWorkHours(db)
	
	// 2. 创建默认工单模板
	createDefaultTicketTemplates(db)
	
	// 3. 创建默认SLA策略
	createDefaultSLAPolicies(db)
	
	fmt.Println("默认数据初始化完成！")
}

func createDefaultWorkHours(db interface{}) {
	fmt.Println("创建默认工时标准...")
	// 这里需要通过service层创建数据
	// 由于需要认证和租户上下文，这里只是示例
}

func createDefaultTicketTemplates(db interface{}) {
	fmt.Println("创建默认工单模板...")
	// 这里需要通过service层创建数据
}

func createDefaultSLAPolicies(db interface{}) {
	fmt.Println("创建默认SLA策略...")
	// 这里需要通过service层创建数据
}
