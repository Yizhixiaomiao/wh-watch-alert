package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"watchAlert/config"
	"watchAlert/internal/global"
	"watchAlert/internal/services"
)

func main() {
	// 初始化配置
	global.Config = config.InitConfig()

	simulator := services.NewAlertSimulatorWithDB()
	if simulator == nil {
		log.Fatal("无法初始化告警模拟器")
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		createMockAlert(simulator)
	case "list":
		listMockAlerts(simulator)
	case "recover":
		recoverAlert(simulator)
	case "cleanup":
		cleanupMockAlerts(simulator)
	case "demo":
		runDemo(simulator)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("WatchAlert 告警模拟器")
	fmt.Println("使用方法:")
	fmt.Println("  go run alert_simulator_cli.go create        - 创建模拟告警")
	fmt.Println("  go run alert_simulator_cli.go list          - 列出模拟告警")
	fmt.Println("  go run alert_simulator_cli.go recover <id>  - 恢复指定告警")
	fmt.Println("  go run alert_simulator_cli.go cleanup       - 清理所有模拟数据")
	fmt.Println("  go run alert_simulator_cli.go demo          - 运行演示场景")
}

func createMockAlert(simulator *services.AlertSimulator) {
	config := services.MockAlertConfig{
		RuleName: "CPU使用率过高",
		Severity: "Critical",
		Labels: map[string]interface{}{
			"instance": "web-server-01",
			"service":  "user-service",
			"env":      "production",
			"cpu":      "85%",
		},
		AutoCreateTicket: true,
		AutoRecover:      false,
		Duration:         30 * time.Second,
		TenantId:         "demo-tenant-001",
		FaultCenterId:    "mock-fault-center",
	}

	event, err := simulator.CreateMockAlert(config)
	if err != nil {
		log.Fatal("创建模拟告警失败:", err)
	}

	fmt.Printf("✅ 成功创建模拟告警:\n")
	fmt.Printf("   事件ID: %s\n", event.EventId)
	fmt.Printf("   规则名称: %s\n", event.RuleName)
	fmt.Printf("   严重程度: %s\n", event.Severity)
	fmt.Printf("   状态: %s\n", event.Status)
	fmt.Printf("   告警将通过系统标准流程处理，请等待状态转换和工单创建...\n")
}

func listMockAlerts(simulator *services.AlertSimulator) {
	alerts, err := simulator.GetMockAlerts("")
	if err != nil {
		log.Fatal("获取告警列表失败:", err)
	}

	if len(alerts) == 0 {
		fmt.Println("📭 没有找到模拟告警")
		return
	}

	fmt.Printf("📋 找到 %d 个模拟告警:\n\n", len(alerts))

	for i, alert := range alerts {
		fmt.Printf("[%d] %s\n", i+1, alert.RuleName)
		fmt.Printf("    事件ID: %s\n", alert.EventId)
		fmt.Printf("    严重程度: %s\n", alert.Severity)
		fmt.Printf("    状态: %s\n", alert.Status)
		fmt.Printf("    首次触发: %s\n", time.Unix(alert.FirstTriggerTime, 0).Format("2006-01-02 15:04:05"))

		// 打印标签
		if len(alert.Labels) > 0 {
			fmt.Printf("    标签: ")
			labels, _ := json.Marshal(alert.Labels)
			fmt.Printf("%s\n", string(labels))
		}
		fmt.Println()
	}
}

func recoverAlert(simulator *services.AlertSimulator) {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请提供告警事件ID")
		fmt.Println("使用方法: go run alert_simulator_cli.go recover <event_id>")
		return
	}

	eventId := os.Args[2]
	err := simulator.RecoverAlert(eventId)
	if err != nil {
		log.Fatal("恢复告警失败:", err)
	}

	fmt.Printf("✅ 告警已恢复: %s\n", eventId)
}

func cleanupMockAlerts(simulator *services.AlertSimulator) {
	err := simulator.CleanupMockAlerts("")
	if err != nil {
		log.Fatal("清理数据失败:", err)
	}

	fmt.Println("🧹 已清理所有模拟数据")
}

func runDemo(simulator *services.AlertSimulator) {
	fmt.Println("🎬 开始运行演示场景...")

	// 清理旧数据
	fmt.Println("🧹 清理旧的演示数据...")
	simulator.CleanupMockAlerts("demo-tenant-001")

	// 场景1: 创建多个不同严重程度的告警
	fmt.Println("\n🚨 场景1: 创建多个告警...")

	scenarios := []services.MockAlertConfig{
		{
			RuleName: "CPU使用率过高",
			Severity: "Critical",
			Labels: map[string]interface{}{
				"instance": "web-server-01",
				"service":  "user-service",
				"cpu":      "95%",
			},
			AutoCreateTicket: true,
			AutoRecover:      true,
			RecoverAfter:     10 * time.Second,
			Duration:         30 * time.Second,
			TenantId:         "demo-tenant-001",
			FaultCenterId:    "mock-fault-center",
		},
		{
			RuleName: "内存使用率过高",
			Severity: "Warning",
			Labels: map[string]interface{}{
				"instance": "db-server-01",
				"service":  "database",
				"memory":   "85%",
			},
			AutoCreateTicket: true,
			AutoRecover:      true,
			RecoverAfter:     15 * time.Second,
			Duration:         30 * time.Second,
			TenantId:         "demo-tenant-001",
			FaultCenterId:    "mock-fault-center",
		},
		{
			RuleName: "磁盘空间不足",
			Severity: "Info",
			Labels: map[string]interface{}{
				"instance": "storage-01",
				"service":  "storage",
				"disk":     "78%",
			},
			AutoCreateTicket: true,
			AutoRecover:      true,
			RecoverAfter:     20 * time.Second,
			Duration:         30 * time.Second,
			TenantId:         "demo-tenant-001",
			FaultCenterId:    "mock-fault-center",
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("   创建告警 %d: %s\n", i+1, scenario.RuleName)
		event, err := simulator.CreateMockAlert(scenario)
		if err != nil {
			log.Printf("❌ 创建告警失败: %v", err)
			continue
		}
		fmt.Printf("   ✅ 事件ID: %s\n", event.EventId)
		time.Sleep(2 * time.Second) // 间隔2秒
	}

	// 场景2: 等待一段时间查看状态
	fmt.Println("\n⏳ 等待告警自动恢复...")
	time.Sleep(25 * time.Second)

	// 场景3: 查看最终状态
	fmt.Println("\n📊 查看最终状态...")
	listMockAlerts(simulator)

	fmt.Println("\n🎉 演示完成！")
	fmt.Println("💡 提示:")
	fmt.Println("   - 告警通过系统标准流程处理，会自动创建工单")
	fmt.Println("   - 告警恢复后会同步更新工单状态")
	fmt.Println("   - 可以通过前端界面查看工单状态变化")
	fmt.Println("   - 工单和通知会通过正常的告警处理流程生成")
}
