package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/internal/cache"
	"watchAlert/internal/ctx"
	"watchAlert/internal/repo"
	"watchAlert/internal/services"
)

type debugController struct{}

var DebugController = new(debugController)

func (debugController debugController) API(gin *gin.RouterGroup) {
	debug := gin.Group("debug")
	{
		debug.POST("createMockAlert", debugController.CreateMockAlert)
		debug.POST("recoverMockAlert", debugController.RecoverMockAlert)
		debug.POST("cleanupMockAlerts", debugController.CleanupMockAlerts)
		debug.GET("testStatus", debugController.TestStatus)
	}
}

func (debugController debugController) TestStatus(c *gin.Context) {
	logc.Info(c.Request.Context(), "测试接口连接成功")
	now := time.Now()
	c.JSON(200, gin.H{
		"message":   "WatchAlert调试接口运行正常",
		"database":  "已连接并包含新字段",
		"timestamp": gin.H{"unix": now.Unix(), "str": now.String()},
	})
}

// CreateMockAlert 创建模拟告警
func (debugController debugController) CreateMockAlert(c *gin.Context) {
	var req struct {
		RuleName         string                 `json:"rule_name"`
		Severity         string                 `json:"severity"`
		Labels           map[string]interface{} `json:"labels"`
		AutoCreateTicket bool                   `json:"auto_create_ticket"`
		AutoRecover      bool                   `json:"auto_recover"`
		RecoverAfter     int64                  `json:"recover_after"`
		Duration         int64                  `json:"duration"`
		TenantId         string                 `json:"tenant_id"`
		FaultCenterId    string                 `json:"fault_center_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logc.Errorf(c.Request.Context(), "参数解析失败: %v", err)
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "参数解析失败",
		})
		return
	}

	logc.Infof(c.Request.Context(), "接收到的参数 - RuleName: %s (长度: %d)", req.RuleName, len(req.RuleName))

	// 设置默认值
	if req.RuleName == "" {
		req.RuleName = "CPU使用率过高"
	}
	if req.Severity == "" {
		req.Severity = "Critical"
	}
	if req.TenantId == "" {
		req.TenantId = "demo-tenant-001"
	}
	if req.FaultCenterId == "" {
		req.FaultCenterId = "default"
	}
	if req.Duration == 0 {
		req.Duration = 30
	}
	if req.RecoverAfter == 0 {
		req.RecoverAfter = 60
	}

	// 创建模拟器
	dbRepo := repo.NewRepoEntry()
	cacheEntry := cache.NewEntryCache()
	simCtx := ctx.NewContext(c.Request.Context(), dbRepo, cacheEntry)
	simulator := services.NewAlertSimulator(simCtx)

	if simulator == nil {
		logc.Error(c.Request.Context(), "模拟器初始化失败")
		c.JSON(500, gin.H{"error": "模拟器初始化失败", "message": "服务错误"})
		return
	}

	// 创建模拟告警配置
	config := services.MockAlertConfig{
		RuleName:         req.RuleName,
		Severity:         req.Severity,
		Labels:           req.Labels,
		AutoCreateTicket: req.AutoCreateTicket,
		AutoRecover:      req.AutoRecover,
		RecoverAfter:     time.Duration(req.RecoverAfter) * time.Second,
		Duration:         time.Duration(req.Duration) * time.Second,
		TenantId:         req.TenantId,
		FaultCenterId:    req.FaultCenterId,
	}

	// 创建模拟告警
	event, err := simulator.CreateMockAlert(config)
	if err != nil {
		logc.Errorf(c.Request.Context(), "创建模拟告警失败: %v", err)
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "创建模拟告警失败",
		})
		return
	}

	logc.Infof(c.Request.Context(), "成功创建模拟告警: %s", event.EventId)
	logc.Infof(c.Request.Context(), "告警规则名称: %s", event.RuleName)

	c.JSON(200, gin.H{
		"success": true,
		"message": "模拟告警创建成功",
		"alert": gin.H{
			"event_id":  event.EventId,
			"rule_name": event.RuleName,
			"severity":  event.Severity,
			"status":    event.Status,
		},
	})
}

// RecoverMockAlert 恢复模拟告警
func (debugController debugController) RecoverMockAlert(c *gin.Context) {
	var req struct {
		EventId string `json:"event_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logc.Errorf(c.Request.Context(), "参数解析失败: %v", err)
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "参数解析失败",
		})
		return
	}

	// 创建模拟器
	dbRepo := repo.NewRepoEntry()
	cacheEntry := cache.NewEntryCache()
	simCtx := ctx.NewContext(c.Request.Context(), dbRepo, cacheEntry)
	simulator := services.NewAlertSimulator(simCtx)

	if simulator == nil {
		logc.Error(c.Request.Context(), "模拟器初始化失败")
		c.JSON(500, gin.H{"error": "模拟器初始化失败", "message": "服务错误"})
		return
	}

	// 恢复告警
	err := simulator.RecoverAlert(req.EventId)
	if err != nil {
		logc.Errorf(c.Request.Context(), "恢复告警失败: %v", err)
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "恢复告警失败",
		})
		return
	}

	logc.Infof(c.Request.Context(), "成功恢复告警: %s", req.EventId)

	c.JSON(200, gin.H{
		"success":  true,
		"message":  "告警已恢复",
		"event_id": req.EventId,
	})
}

// CleanupMockAlerts 清理模拟告警
func (debugController debugController) CleanupMockAlerts(c *gin.Context) {
	var req struct {
		TenantId string `json:"tenant_id"`
	}

	c.ShouldBindJSON(&req)

	// 创建模拟器
	dbRepo := repo.NewRepoEntry()
	cacheEntry := cache.NewEntryCache()
	simCtx := ctx.NewContext(c.Request.Context(), dbRepo, cacheEntry)
	simulator := services.NewAlertSimulator(simCtx)

	if simulator == nil {
		logc.Error(c.Request.Context(), "模拟器初始化失败")
		c.JSON(500, gin.H{"error": "模拟器初始化失败", "message": "服务错误"})
		return
	}

	// 清理数据
	err := simulator.CleanupMockAlerts(req.TenantId)
	if err != nil {
		logc.Errorf(c.Request.Context(), "清理数据失败: %v", err)
		c.JSON(500, gin.H{
			"error":   err.Error(),
			"message": "清理数据失败",
		})
		return
	}

	logc.Info(c.Request.Context(), "已清理模拟告警数据")

	c.JSON(200, gin.H{
		"success": true,
		"message": "数据清理完成",
	})
}
