package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"watchAlert/alert/process"
	"watchAlert/internal/cache"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/repo"
)

type AlertSimulator struct {
	ctx *ctx.Context
}

func NewAlertSimulator(ctx *ctx.Context) *AlertSimulator {
	return &AlertSimulator{
		ctx: ctx,
	}
}

// MockAlertConfig 模拟告警配置
type MockAlertConfig struct {
	RuleName         string                 `json:"rule_name"`
	Severity         string                 `json:"severity"`
	Labels           map[string]interface{} `json:"labels"`
	AutoCreateTicket bool                   `json:"auto_create_ticket"`
	AutoRecover      bool                   `json:"auto_recover"`
	RecoverAfter     time.Duration          `json:"recover_after"`
	Duration         time.Duration          `json:"duration"`
	TenantId         string                 `json:"tenant_id"`
	FaultCenterId    string                 `json:"fault_center_id"`
}

// CreateMockAlert 创建模拟告警
func (s *AlertSimulator) CreateMockAlert(config MockAlertConfig) (*models.AlertCurEvent, error) {
	if s.ctx == nil {
		return nil, fmt.Errorf("AlertSimulator未初始化，需要Context")
	}

	now := time.Now()
	fingerprint := fmt.Sprintf("mock-fp-%d", now.UnixNano())

	// 如果没有指定故障中心ID，使用默认值
	faultCenterId := config.FaultCenterId
	if faultCenterId == "" {
		faultCenterId = "default"
	}

	// 验证故障中心是否存在（可选，用于调试）
	logc.Infof(s.ctx.Ctx, "使用故障中心: %s", faultCenterId)

	// 如果没有指定租户ID，使用默认值
	tenantId := config.TenantId
	if tenantId == "" {
		tenantId = "demo-tenant-001"
	}

	// 创建告警事件，使用系统标准流程
	event := &models.AlertCurEvent{
		TenantId:             tenantId,
		EventId:              fmt.Sprintf("mock-%d", now.UnixNano()),
		RuleId:               "mock-rule-" + strconv.FormatInt(now.Unix(), 10),
		RuleGroupId:          "mock-group",
		RuleName:             config.RuleName,
		DatasourceType:       "Prometheus",
		DatasourceId:         "mock-datasource",
		Fingerprint:          fingerprint,
		Severity:             config.Severity,
		Status:               models.StatePreAlert, // 初始状态为预告警
		Labels:               config.Labels,
		EvalInterval:         60,
		ForDuration:          int64(config.Duration.Seconds()), // 持续时间
		RepeatNoticeInterval: 300,
		EffectiveTime:        models.EffectiveTime{},
		FaultCenterId:        faultCenterId,
		FirstTriggerTime:     now.Unix(),
		LastEvalTime:         now.Unix(),
	}

	// 使用系统的标准告警处理流程
	process.PushEventToFaultCenter(s.ctx, event)

	logc.Infof(s.ctx.Ctx, "成功创建模拟告警: %s (规则: %s, 严重程度: %s, 故障中心: %s)",
		event.EventId, event.RuleName, event.Severity, faultCenterId)

	// 模拟定期评估，触发状态转换
	go func() {
		// 等待持续时间后，再次推送告警以触发状态转换
		time.Sleep(config.Duration + 1*time.Second)

		// 更新评估时间
		event.LastEvalTime = time.Now().Unix()

		// 再次推送，这将触发状态转换从 pre_alert 到 alerting
		process.PushEventToFaultCenter(s.ctx, event)

		logc.Infof(s.ctx.Ctx, "模拟告警状态转换已触发: %s (pre_alert -> alerting)", event.EventId)

		// 获取故障中心配置，尝试发送通知
		time.Sleep(1 * time.Second) // 等待状态更新完成
		faultCenter, err := s.ctx.DB.FaultCenter().Get(tenantId, faultCenterId, "")
		if err != nil {
			logc.Errorf(s.ctx.Ctx, "获取故障中心配置失败: %v", err)
		} else if len(faultCenter.NoticeIds) > 0 {
			// 获取最新的告警事件（包含更新后的状态）
			cacheEvent, cacheErr := s.ctx.Redis.Alert().GetEventFromCache(tenantId, faultCenterId, fingerprint)
			if cacheErr == nil && cacheEvent.GetEventStatus() == models.StateAlerting {
				// 发送告警通知（使用第一个通知对象）
				noticeId := faultCenter.NoticeIds[0]
				err = process.HandleAlert(s.ctx, "alarm", faultCenter, noticeId, []*models.AlertCurEvent{&cacheEvent})
				if err != nil {
					logc.Errorf(s.ctx.Ctx, "发送模拟告警通知失败: %v", err)
				} else {
					logc.Infof(s.ctx.Ctx, "模拟告警通知已发送: %s, 通知对象: %s", event.EventId, noticeId)
				}
			}
		}

		// 如果配置了自动恢复，则在一段时间后触发恢复
		if config.AutoRecover {
			time.Sleep(config.RecoverAfter)

			// 标记为已恢复
			event.IsRecovered = true
			event.LastEvalTime = time.Now().Unix()

			// 推送恢复事件
			process.PushEventToFaultCenter(s.ctx, event)

			logc.Infof(s.ctx.Ctx, "模拟告警已恢复: %s", event.EventId)
		}
	}()

	return event, nil
}

// RecoverAlert 恢复告警
func (s *AlertSimulator) RecoverAlert(eventId string) error {
	if s.ctx == nil {
		return fmt.Errorf("AlertSimulator未初始化，需要Context")
	}

	// 从数据库或缓存中查找告警事件
	// 这里我们通过数据库查找
	var event models.AlertCurEvent
	err := s.ctx.DB.DB().Where("event_id = ?", eventId).First(&event).Error
	if err != nil {
		return fmt.Errorf("查找告警事件失败: %v", err)
	}

	// 创建恢复事件
	recoverEvent := &models.AlertCurEvent{
		TenantId:             event.TenantId,
		EventId:              event.EventId,
		RuleId:               event.RuleId,
		RuleGroupId:          event.RuleGroupId,
		RuleName:             event.RuleName,
		DatasourceType:       event.DatasourceType,
		DatasourceId:         event.DatasourceId,
		Fingerprint:          event.Fingerprint,
		Severity:             event.Severity,
		Status:               models.StateAlerting, // 设置为告警状态，让系统处理恢复逻辑
		Labels:               event.Labels,
		EvalInterval:         event.EvalInterval,
		ForDuration:          event.ForDuration,
		RepeatNoticeInterval: event.RepeatNoticeInterval,
		EffectiveTime:        event.EffectiveTime,
		FaultCenterId:        event.FaultCenterId,
		FirstTriggerTime:     event.FirstTriggerTime,
		LastEvalTime:         time.Now().Unix(),
		IsRecovered:          true, // 标记为已恢复
	}

	// 使用系统流程处理恢复
	process.PushEventToFaultCenter(s.ctx, recoverEvent)

	logc.Infof(s.ctx.Ctx, "告警恢复请求已提交: %s", eventId)
	return nil
}

// GetMockAlerts 获取模拟告警列表
func (s *AlertSimulator) GetMockAlerts(tenantId string) ([]models.AlertCurEvent, error) {
	if s.ctx == nil {
		return nil, fmt.Errorf("AlertSimulator未初始化，需要Context")
	}

	var events []models.AlertCurEvent
	query := s.ctx.DB.DB().Where("rule_id LIKE ?", "mock-rule-%")
	if tenantId != "" {
		query = query.Where("tenant_id = ?", tenantId)
	}

	if err := query.Order("first_trigger_time DESC").Limit(50).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("查询模拟告警失败: %v", err)
	}

	return events, nil
}

// CleanupMockAlerts 清理模拟告警
func (s *AlertSimulator) CleanupMockAlerts(tenantId string) error {
	if s.ctx == nil {
		return fmt.Errorf("AlertSimulator未初始化，需要Context")
	}

	// 清理数据库中的告警事件
	query := s.ctx.DB.DB().Where("rule_id LIKE ?", "mock-rule-%")
	if tenantId != "" {
		query = query.Where("tenant_id = ?", tenantId)
	}

	if err := query.Delete(&models.AlertCurEvent{}).Error; err != nil {
		return fmt.Errorf("清理模拟告警失败: %v", err)
	}

	// 清理历史记录
	historyQuery := s.ctx.DB.DB().Where("rule_id LIKE ?", "mock-rule-%")
	if tenantId != "" {
		historyQuery = historyQuery.Where("tenant_id = ?", tenantId)
	}

	if err := historyQuery.Delete(&models.AlertHisEvent{}).Error; err != nil {
		logc.Errorf(s.ctx.Ctx, "清理告警历史失败: %v", err)
	}

	// 清理相关工单
	ticketQuery := s.ctx.DB.DB().Where("created_by = ?", "alert-simulator")
	if tenantId != "" {
		ticketQuery = ticketQuery.Where("tenant_id = ?", tenantId)
	}

	if err := ticketQuery.Delete(&models.Ticket{}).Error; err != nil {
		logc.Errorf(s.ctx.Ctx, "清理工单失败: %v", err)
	}

	// 清理Redis缓存 - 清理故障中心的模拟告警
	// 如果指定了租户ID，只清理该租户的数据，否则清理所有租户
	tenantsToClean := []string{}
	if tenantId != "" {
		tenantsToClean = append(tenantsToClean, tenantId)
	} else {
		tenantsToClean = []string{"default", "demo-tenant-001"}
	}
	faultCenters := []string{"default", "mock-fault-center"}

	for _, tid := range tenantsToClean {
		for _, fcId := range faultCenters {
			key := fmt.Sprintf("w8t:%s:faultCenter:%s.events", tid, fcId)
			// 获取hash中的所有字段
			fields, err := s.ctx.Redis.Redis().HKeys(key).Result()
			if err != nil {
				continue
			}
			// 删除所有mock-fp开头的字段
			for _, field := range fields {
				if strings.HasPrefix(field, "mock-fp-") {
					s.ctx.Redis.Redis().HDel(key, field)
					logc.Infof(s.ctx.Ctx, "已清理Redis告警: %s (键: %s)", field, key)
				}
			}
		}
	}

	logc.Infof(s.ctx.Ctx, "已清理模拟告警数据")
	return nil
}

// NewAlertSimulatorWithDB 创建带数据库连接的模拟器（用于CLI工具）
func NewAlertSimulatorWithDB() *AlertSimulator {
	// 创建Context
	dbRepo := repo.NewRepoEntry()
	cacheEntry := cache.NewEntryCache()

	context := ctx.NewContext(nil, dbRepo, cacheEntry)

	return &AlertSimulator{
		ctx: context,
	}
}
